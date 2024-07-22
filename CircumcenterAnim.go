package main

import (
	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type CircumcenterSideAnim struct {
	basePoint     glm.Vec2
	p1, p2        glm.Vec2
	pointToLerpTo glm.Vec2
	animR         closedGL.Animation
	animDur       float32
	animLine      closedGL.Animation
	ctx           *closedGL.ClosedGLContext
}

func newCircumCenterSideAnim(p1, p2, p3 glm.Vec2, tri *Tri, animDur float32) CircumcenterSideAnim {
	var mp = closedGL.MiddlePoint(p1, p2)

	var dirVec = mp.Sub(&p1)
	var targetR = dirVec.Len()
	var circumcenter = tri.calcCircumcenter()
	var perpLine = CalcPerpLineVec(p1, p2)

	var eq = CalcLinearEquation(p1, p3)
	var eq2 = CalcLinearEquation(p2, p3)

	var cp = CalcCrossingPoint(perpLine, eq)
	var cp2 = CalcCrossingPoint(perpLine, eq2)

	var cpVec = mp.Sub(&cp)
	var cp2Vec = mp.Sub(&cp2)

	if cp2Vec.Len() < cpVec.Len() {
		cpVec = cp2Vec
		cp = cp2
	}

	var circumVec = mp.Sub(&circumcenter)
	var otherP = circumcenter
	if cpVec.Len() > circumVec.Len() {
		otherP = cp
	}
	return CircumcenterSideAnim{
		basePoint:     mp,
		p1:            p1,
		p2:            p2,
		ctx:           tri.Ctx,
		pointToLerpTo: otherP,
		animDur:       animDur,
		animR:         closedGL.NewAnimation(0, targetR, animDur, false, false),
		animLine:      closedGL.NewAnimation(0, 1.2, animDur, false, false),
	}
}

func (this *CircumcenterSideAnim) Process(delta float32) {
	this.animR.Process(delta)
	if this.animR.IsFinished() {
		this.animLine.Process(delta)
	}
}

func (this *CircumcenterSideAnim) Draw() {
	var borderC = glm.Vec4{1, 1, 0, 1}

	if !this.animR.IsFinished() && this.animR.GetValue() >= glm.Epsilon {
		drawCartesianCircle(this.p1, this.ctx, glm.Vec4{1, 0, 0, 0}, borderC, 3, this.animR.GetValue(), 3)
		drawCartesianCircle(this.p2, this.ctx, glm.Vec4{1, 0, 0, 0}, borderC, 3, this.animR.GetValue(), 3)

	}
	drawCartesianLine(this.basePoint, LerpVec2(this.basePoint, this.pointToLerpTo, this.animLine.GetValue()), this.ctx, 3, glm.Vec4{1, 1, 0, 1})
	drawCartesianLine(this.basePoint, LerpVec2(this.basePoint, this.pointToLerpTo, -this.animLine.GetValue()), this.ctx, 3, glm.Vec4{1, 1, 0, 1})

}

func (this *CircumcenterSideAnim) IsFinished() bool {
	return this.animLine.IsFinished()
}

type CircumcenterAnim struct {
	tri        *Tri
	machine    StateMachine
	animCenter closedGL.Animation
	animDur    float32
	endTimer   closedGL.Animation
}

func newCircumCenterAnim(tri *Tri, animDur float32) CircumcenterAnim {
	var c = tri.calcCircumcenter()
	var vec = c.Sub(&tri.Points[0])
	var test = newCircumCenterSideAnim(tri.Points[0], tri.Points[1], tri.Points[2], tri, animDur)
	var test2 = newCircumCenterSideAnim(tri.Points[0], tri.Points[2], tri.Points[1], tri, animDur)
	var test3 = newCircumCenterSideAnim(tri.Points[1], tri.Points[2], tri.Points[0], tri, animDur)
	var machine = newStateMachine()
	machine.addState(&test).addState(&test2).addState(&test3)

	return CircumcenterAnim{
		tri:        tri,
		animDur:    animDur,
		animCenter: closedGL.NewAnimation(0, vec.Len(), animDur, false, false),
		machine:    machine,
		endTimer:   closedGL.NewAnimation(1, 1, 1, false, false),
	}
}

func (this *CircumcenterAnim) Draw() {
	this.machine.drawAll()
	if this.machine.isFinished() {
		this.tri.drawCircumCenter()
		drawCartesianCircle(this.tri.calcCircumcenter(), this.tri.Ctx, glm.Vec4{1, 1, 0, 0}, glm.Vec4{1, 1, 0, 1}, 3, this.animCenter.GetValue(), 3)
	}
}

func (this *CircumcenterAnim) Process(delta float32) {
	this.machine.process(delta)
	if this.machine.isFinished() {
		this.animCenter.Process(delta)
	}
	if this.animCenter.IsFinished() {
		this.endTimer.Process(delta)
	}
}

func (this *CircumcenterAnim) IsFinished() bool {
	return this.endTimer.IsFinished()
}

func (this *CircumcenterAnim) init() {

}
