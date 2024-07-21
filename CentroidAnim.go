package main

import (
	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type AnimSideCentroid struct {
	corner, p1, p2, mp      glm.Vec2
	movToMpAnim, radiusAnim closedGL.Animation
	ctx                     *closedGL.ClosedGLContext
}

func newAnimSideCentroid(cornerP, p1, p2 glm.Vec2, ctx *closedGL.ClosedGLContext, animDur float32) AnimSideCentroid {
	var mp = CalcMiddlePoint(p1, p2)
	var distVec = mp.Sub(&p1)
	var factor float32 = 1.1
	return AnimSideCentroid{
		corner:      cornerP,
		p1:          p1,
		p2:          p2,
		mp:          mp,
		radiusAnim:  closedGL.NewAnimation(0, distVec.Len()*factor, animDur, false, false),
		movToMpAnim: closedGL.NewAnimation(0, 1, (animDur / (distVec.Len() * factor) * distVec.Len()), false, false),
		ctx:         ctx,
	}
}

func (this *AnimSideCentroid) Process(delta float32) {
	this.movToMpAnim.Process(delta)
	this.radiusAnim.Process(delta)

}

func (this *AnimSideCentroid) Draw() {
	var targetP1 = LerpVec2(this.p1, this.mp, this.movToMpAnim.GetValue())
	var targetP2 = LerpVec2(this.p2, this.mp, this.movToMpAnim.GetValue())

	drawCartesianCircle(targetP1, this.ctx, glm.Vec4{1, 0, 0, 1}, glm.Vec4{1, 0, 0, 1}, 2, 10, 3)
	drawCartesianCircle(targetP2, this.ctx, glm.Vec4{1, 0, 0, 1}, glm.Vec4{1, 0, 0, 1}, 2, 10, 3)
	if !this.radiusAnim.IsFinished() {
		drawCartesianCircle(this.p1, this.ctx, glm.Vec4{0, 0, 0, 0}, glm.Vec4{1, 0, 0, 1}, 2, this.radiusAnim.GetValue(), 3)
		drawCartesianCircle(this.p2, this.ctx, glm.Vec4{0, 0, 0, 0}, glm.Vec4{1, 0, 0, 1}, 2, this.radiusAnim.GetValue(), 3)
	}
	drawCartesianLine(this.corner, targetP1, this.ctx, 2, glm.Vec4{1, 0, 0, 1})
	drawCartesianLine(this.corner, targetP2, this.ctx, 2, glm.Vec4{1, 0, 0, 1})
}

func (this *AnimSideCentroid) IsFinished() bool { return this.radiusAnim.IsFinished() }

type CentroidAnim struct {
	tri      *Tri
	machine  StateMachine
	endTimer closedGL.Timer
}

func newCentroidAnim(tri *Tri) CentroidAnim {
	var machine = newStateMachine()
	var side0 = newAnimSideCentroid(tri.Points[0], tri.Points[1], tri.Points[2], tri.Ctx, 1)
	var side1 = newAnimSideCentroid(tri.Points[1], tri.Points[0], tri.Points[2], tri.Ctx, 1)
	var side2 = newAnimSideCentroid(tri.Points[2], tri.Points[0], tri.Points[1], tri.Ctx, 1)
	machine.addState(&side0).addState(&side1).addState(&side2)
	var retVal = CentroidAnim{tri: tri, machine: machine, endTimer: closedGL.NewTimer(0.5, false)}
	return retVal
}

func (this *CentroidAnim) init() {

}

func (this *CentroidAnim) Draw() {
	this.machine.drawAll()
	if this.machine.isFinished() {
		this.tri.drawCentroid()
	}
}

func (this *CentroidAnim) Process(delta float32) {
	this.machine.process(delta)
	if this.machine.isFinished() {
		this.endTimer.Process(delta)
	}
}

func (this *CentroidAnim) IsFinished() bool {
	return this.endTimer.IsTick()
}
