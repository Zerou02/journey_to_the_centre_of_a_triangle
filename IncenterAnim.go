package main

import (
	"math"

	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type IncenterCornerAnim struct {
	dirP, cornerP glm.Vec2
	animX, animY  closedGL.Animation
	otherCornerP  glm.Vec2
}

func newCornerAnim(cornerP, p1, p2 glm.Vec2, centroid glm.Vec2) IncenterCornerAnim {
	var vec1 = p1.Sub(&cornerP)
	var vec2 = p2.Sub(&cornerP)
	var angle = AngleTo(vec1, vec2) / 2
	var rotated = Rotate(angle, vec2)
	var rotated2 = Rotate(2*math.Pi-angle, vec2)

	var norm = rotated.Normalized()
	var norm2 = rotated2.Normalized()

	var steps = distBetweenPoints(cornerP, centroid)
	var tenSteps = norm.Mul(steps)
	var tenSteps2 = norm2.Mul(steps)

	var newP = cornerP.Add(&tenSteps)
	var newP2 = cornerP.Add(&tenSteps2)
	var dist1 = distBetweenPoints(newP, centroid)
	var dist2 = distBetweenPoints(newP2, centroid)

	var p glm.Vec2

	if dist1 < dist2 {
		p = newP
	} else {
		p = newP2
	}
	var baseEq = CalcLinearEquation(cornerP, p)
	var oppositeEq = CalcLinearEquation(p1, p2)
	var crossing = CalcCrossingPoint(baseEq, oppositeEq)
	var vec = crossing.Sub(&cornerP)
	var scaled = vec.Mul(1.2)
	var newC = cornerP.Add(&scaled)

	var animDur float32 = 1
	return IncenterCornerAnim{
		cornerP:      cornerP,
		dirP:         newC,
		animX:        closedGL.NewAnimation(cornerP[0], newC[0], animDur, false, false),
		animY:        closedGL.NewAnimation(cornerP[1], newC[1], animDur, false, false),
		otherCornerP: p1,
	}
}

func (this *IncenterCornerAnim) process(delta float32, currState, allowedState int) {
	if currState == allowedState {
		this.animX.Process(delta)
		this.animY.Process(delta)
	}

}

// 0 -> nichts, 1 -> regulaer, 2-> post
func (this *IncenterCornerAnim) draw(ctx *closedGL.ClosedGLContext) {
	var p = glm.Vec2{this.animX.GetValue(), this.animY.GetValue()}
	var dist = distToLine(this.cornerP, this.otherCornerP, p)
	ctx.DrawCircle(CartesianToSSPoint(p, ctx.Window.Wh), glm.Vec4{0, 0.5, 0, 0}, glm.Vec4{0, 0.5, 0, 0}, 10, 3, 2)
	drawCartesianLine(this.cornerP, p, ctx, 2, glm.Vec4{0, 0.5, 0, 1})
	if !this.animX.IsFinished() {
		ctx.DrawCircle(CartesianToSSPoint(p, ctx.Window.Wh), glm.Vec4{0, 0.5, 0, 0}, glm.Vec4{0, 0.5, 0, 0}, dist, 3, 3)
	}
}

type IncenterAnim struct {
	tri         *Tri
	centroid    glm.Vec2
	cornerAnims [3]IncenterCornerAnim
	currState   int
	unitAnim    closedGL.Animation
}

func newIncenterAnim(tri *Tri) IncenterAnim {
	return IncenterAnim{tri: tri, currState: 0}
}

func (this *IncenterAnim) Process(delta float32) {
	for i, x := range this.cornerAnims {
		this.cornerAnims[i].process(delta, this.currState, i)
		if x.animX.IsFinished() {
			this.currState = i + 1
		}
	}
	if this.currState == 3 {
		this.unitAnim.Process(delta)
	}
}
func (this *IncenterAnim) init() {
	this.centroid = this.tri.calcCentroid()
	this.cornerAnims[0] = newCornerAnim(this.tri.Points[0], this.tri.Points[1], this.tri.Points[2], this.centroid)
	this.cornerAnims[1] = newCornerAnim(this.tri.Points[1], this.tri.Points[0], this.tri.Points[2], this.centroid)
	this.cornerAnims[2] = newCornerAnim(this.tri.Points[2], this.tri.Points[0], this.tri.Points[1], this.centroid)
	this.unitAnim = closedGL.NewAnimation(0, 1, 1, false, false)
}
func (this *IncenterAnim) Draw() {

	for i := 0; i < len(this.cornerAnims); i++ {
		if this.currState >= i {
			this.cornerAnims[i].draw(this.tri.Ctx)
		}
	}
	if this.currState == 3 {
		var incenter = this.tri.calcIncenter()

		var p = findLineCircleIntersectionPoint(this.tri.calcIncenter(), this.tri.Points[0], this.tri.Points[1])
		var len = incenter.Sub(&p)
		drawCartesianCircle(incenter, this.tri.Ctx, glm.Vec4{0, 1, 0, 0}, glm.Vec4{0, 0.5, 0, 1}, 3, this.unitAnim.GetValue()*len.Len(), 3)

	}
}

func (this *IncenterAnim) IsFinished() bool {
	return this.unitAnim.IsFinished()
}
