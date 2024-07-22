package main

import (
	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type IncenterCornerAnim struct {
	dirP, cornerP glm.Vec2
	animX, animY  closedGL.Animation
	otherCornerP  glm.Vec2
}

func newCornerAnim(cornerP, p1, p2 glm.Vec2, tri *Tri, animDur float32) IncenterCornerAnim {
	var centroid = tri.calcCentroid()
	var bisectorVec = findAngleBisectorVec(cornerP, p1, p2, centroid)
	var p = bisectorVec.Add(&cornerP)

	var baseEq = CalcLinearEquation(cornerP, p)
	var oppositeEq = CalcLinearEquation(p1, p2)
	var crossing = CalcCrossingPoint(baseEq, oppositeEq)
	var vec = crossing.Sub(&cornerP)
	var scaled = vec.Mul(1.2)
	var newC = cornerP.Add(&scaled)

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
	drawCartesianLine(this.cornerP, p, ctx, 2, glm.Vec4{0, 0.5, 0, 3})
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
	animDur     float32
}

func newIncenterAnim(tri *Tri, animDur float32) IncenterAnim {
	return IncenterAnim{tri: tri, currState: 0, animDur: animDur}
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
	this.cornerAnims[0] = newCornerAnim(this.tri.Points[0], this.tri.Points[1], this.tri.Points[2], this.tri, this.animDur)
	this.cornerAnims[1] = newCornerAnim(this.tri.Points[1], this.tri.Points[0], this.tri.Points[2], this.tri, this.animDur)
	this.cornerAnims[2] = newCornerAnim(this.tri.Points[2], this.tri.Points[0], this.tri.Points[1], this.tri, this.animDur)
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
