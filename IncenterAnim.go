package main

import (
	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type IncenterCornerAnim struct {
	corner, p1, p2, cp glm.Vec2
	unitAnim           closedGL.Animation
	ctx                *closedGL.ClosedGLContext
}

func newCornerAnim(cornerP, p1, p2 glm.Vec2, tri *Tri, animDur float32) IncenterCornerAnim {
	var centroid = tri.calcCentroid()
	var bisectorVec = findAngleBisectorVec(cornerP, p1, p2, centroid)
	var p = bisectorVec.Add(&cornerP)

	var baseEq = CalcLinearEquation(cornerP, p)
	var oppositeEq = CalcLinearEquation(p1, p2)
	var crossing = CalcCrossingPoint(baseEq, oppositeEq)

	return IncenterCornerAnim{
		corner:   cornerP,
		p1:       p1,
		p2:       p2,
		cp:       crossing,
		unitAnim: closedGL.NewAnimation(0, 1.1, animDur, false, false),
		ctx:      tri.Ctx,
	}
}

func (this *IncenterCornerAnim) Process(delta float32) {
	this.unitAnim.Process(delta)
}

func (this *IncenterCornerAnim) IsFinished() bool {
	return this.unitAnim.IsFinished()
}

func (this *IncenterCornerAnim) Draw() {
	if this.unitAnim.GetValue() < glm.Epsilon {
		return
	}

	var p = LerpVec2(this.corner, this.cp, this.unitAnim.GetValue())
	var dist = distToLine(this.corner, this.p2, p)
	drawCartesianLine(this.corner, p, this.ctx, 2, glm.Vec4{0, 0.5, 0, 1})
	drawCartesianCircle(p, this.ctx, glm.Vec4{0, 0.5, 0, 1}, glm.Vec4{0, 0.5, 0, 0}, 3, 10, 3)
	if !this.unitAnim.IsFinished() {
		drawCartesianCircle(p, this.ctx, glm.Vec4{0, 0.5, 0, 0}, glm.Vec4{0, 0.5, 0, 0}, 3, dist, 3)
	}
}

type IncenterAnim struct {
	tri      *Tri
	centroid glm.Vec2
	machine  StateMachine
	unitAnim closedGL.Animation
	animDur  float32
}

func newIncenterAnim(tri *Tri, animDur float32) IncenterAnim {
	var machine = newStateMachine()
	var cornerAnim0 = newCornerAnim(tri.Points[0], tri.Points[1], tri.Points[2], tri, animDur)
	var cornerAnim1 = newCornerAnim(tri.Points[1], tri.Points[0], tri.Points[2], tri, animDur)
	var cornerAnim2 = newCornerAnim(tri.Points[2], tri.Points[0], tri.Points[1], tri, animDur)
	machine.addState(&cornerAnim0).addState(&cornerAnim1).addState(&cornerAnim2)
	return IncenterAnim{
		tri:      tri,
		animDur:  animDur,
		machine:  machine,
		unitAnim: closedGL.NewAnimation(0, 1, 1, false, false),
		centroid: tri.calcCentroid(),
	}
}

func (this *IncenterAnim) Process(delta float32) {
	this.machine.process(delta)
	if this.machine.isFinished() {
		this.unitAnim.Process(delta)
	}
}
func (this *IncenterAnim) init() {
}

func (this *IncenterAnim) Draw() {

	this.machine.drawAll()
	if this.machine.isFinished() {
		var incenter = this.tri.calcIncenter()
		var p = findLineCircleIntersectionPoint(this.tri.calcIncenter(), this.tri.Points[0], this.tri.Points[1])
		var len = incenter.Sub(&p)
		drawCartesianCircle(incenter, this.tri.Ctx, glm.Vec4{0, 1, 0, 0}, glm.Vec4{0, 0.5, 0, 1}, 3, this.unitAnim.GetValue()*len.Len(), 3)

	}
}

func (this *IncenterAnim) IsFinished() bool {
	return this.unitAnim.IsFinished()
}
