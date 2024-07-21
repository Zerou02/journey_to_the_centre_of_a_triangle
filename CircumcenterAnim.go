package main

import (
	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type CircumcenterSideAnim struct {
	basePoint    glm.Vec2
	circumcenter glm.Vec2
	p1, p2       glm.Vec2
	animR        closedGL.Animation
	animDur      float32
	animLine     closedGL.Animation
	ctx          *closedGL.ClosedGLContext
}

func newCircumCenterSideAnim(p1, p2 glm.Vec2, tri *Tri, animDur float32) CircumcenterSideAnim {
	var mp = closedGL.MiddlePoint(p1, p2)

	var dirVec = mp.Sub(&p1)
	var targetR = dirVec.Len()
	_ = targetR
	var circumcenter = tri.calcCircumcenter()
	return CircumcenterSideAnim{
		basePoint:    mp,
		p1:           p1,
		p2:           p2,
		ctx:          tri.Ctx,
		circumcenter: circumcenter,
		animDur:      animDur,
		animR:        closedGL.NewAnimation(0, targetR, animDur, false, false),
		animLine:     closedGL.NewAnimation(0, 1.1, animDur, false, false),
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

	if !this.animR.IsFinished() {
		drawCartesianCircle(this.p1, this.ctx, glm.Vec4{1, 0, 0, 0}, borderC, 3, this.animR.GetValue(), 3)
		drawCartesianCircle(this.p2, this.ctx, glm.Vec4{1, 0, 0, 0}, borderC, 3, this.animR.GetValue(), 3)
	}
	drawCartesianLine(this.basePoint, LerpVec2(this.basePoint, this.circumcenter, this.animLine.GetValue()), this.ctx, 2, glm.Vec4{1, 1, 0, 1})
	drawCartesianLine(this.basePoint, LerpVec2(this.basePoint, this.circumcenter, -this.animLine.GetValue()), this.ctx, 2, glm.Vec4{1, 1, 0, 1})
}

func (this *CircumcenterSideAnim) IsFinished() bool {
	return this.animLine.IsFinished()
}

type CircumcenterAnim struct {
	tri        *Tri
	machine    StateMachine
	animCenter closedGL.Animation
	animDur    float32
	endTimer   closedGL.Timer
}

func newCircumCenterAnim(tri *Tri, animDur float32) CircumcenterAnim {
	var c = tri.calcCircumcenter()
	var vec = c.Sub(&tri.Points[0])
	var test = newCircumCenterSideAnim(tri.Points[0], tri.Points[1], tri, animDur)
	var test2 = newCircumCenterSideAnim(tri.Points[0], tri.Points[2], tri, animDur)
	var test3 = newCircumCenterSideAnim(tri.Points[1], tri.Points[2], tri, animDur)
	var machine = newStateMachine()
	machine.addState(&test).addState(&test2).addState(&test3)

	return CircumcenterAnim{
		tri:        tri,
		animDur:    animDur,
		animCenter: closedGL.NewAnimation(0, vec.Len(), animDur, false, false),
		machine:    machine,
		endTimer:   closedGL.NewTimer(0.5, false),
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
	return this.endTimer.IsTick()
}

func (this *CircumcenterAnim) init() {

}
