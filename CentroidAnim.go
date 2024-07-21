package main

import (
	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type AnimSideCentroid struct {
	ps             [2]AnimPointCentroid
	freePoint      glm.Vec2
	cAnim          closedGL.Animation
	origP1, origP2 glm.Vec2
	allowedState   int
	currState      *int
}

func newAnimSideCentroid(p1, p2 glm.Vec2, freePoint glm.Vec2, allowedState int, currState *int) AnimSideCentroid {
	var mp = closedGL.MiddlePoint(p1, p2)
	var animDur float32 = 3
	var animP = AnimPointCentroid{
		animX: closedGL.NewAnimation(p1[0], mp[0], animDur, false, false),
		animY: closedGL.NewAnimation(p1[1], mp[1], animDur, false, false),
	}
	var animP2 = AnimPointCentroid{
		animX: closedGL.NewAnimation(p2[0], mp[0], animDur, false, false),
		animY: closedGL.NewAnimation(p2[1], mp[1], animDur, false, false),
	}
	var dist = distBetweenPoints(p1, mp)
	var time = animDur * 1.1
	var cAnim = closedGL.NewAnimation(0, lerpAnim(dist, animDur, time), time, false, false)

	return AnimSideCentroid{
		ps:           [2]AnimPointCentroid{animP, animP2},
		freePoint:    freePoint,
		cAnim:        cAnim,
		origP1:       p1,
		origP2:       p2,
		allowedState: allowedState,
		currState:    currState,
	}
}

func (this *AnimSideCentroid) Process(delta float32) {
	if *this.currState < this.allowedState {
		return
	}
	if *this.currState == 5 {
		this.ps[0].process(delta)
		this.ps[1].process(delta)
		this.cAnim.Process(delta)
		if this.cAnim.IsFinished() {
			*this.currState = 6
		}
	}
	if *this.currState == this.allowedState {
		this.ps[0].process(delta)
		this.ps[1].process(delta)
		this.cAnim.Process(delta)
		if this.cAnim.IsFinished() {
			*this.currState++
		}
	}

}

func (this *AnimSideCentroid) Draw(ctx *closedGL.ClosedGLContext) {
	var wh = ctx.Window.Wh
	var cc = glm.Vec4{1, 0, 0, 1}
	var c = glm.Vec4{1, 0, 0, 1}
	if *this.currState == 5 {
		this.ps[0].draw(ctx)
		this.ps[1].draw(ctx)

		ctx.DrawLine(CartesianToSSPoint(this.ps[0].p, wh), CartesianToSSPoint(this.freePoint, wh), c, c, 1)
		ctx.DrawLine(CartesianToSSPoint(this.ps[1].p, wh), CartesianToSSPoint(this.freePoint, wh), c, c, 1)

		ctx.DrawCircle(CartesianToSSPoint(this.origP1, wh), glm.Vec4{0, 0, 0, 0}, cc, this.cAnim.GetValue(), 5, 3)
		ctx.DrawCircle(CartesianToSSPoint(this.origP2, wh), glm.Vec4{0, 0, 0, 0}, cc, this.cAnim.GetValue(), 5, 3)
		return
	}
	if *this.currState == this.allowedState {
		this.ps[0].draw(ctx)
		this.ps[1].draw(ctx)

		ctx.DrawLine(CartesianToSSPoint(this.ps[0].p, wh), CartesianToSSPoint(this.freePoint, wh), c, c, 1)
		ctx.DrawLine(CartesianToSSPoint(this.ps[1].p, wh), CartesianToSSPoint(this.freePoint, wh), c, c, 1)

		ctx.DrawCircle(CartesianToSSPoint(this.origP1, wh), glm.Vec4{0, 0, 0, 0}, cc, this.cAnim.GetValue(), 5, 2)
		ctx.DrawCircle(CartesianToSSPoint(this.origP2, wh), glm.Vec4{0, 0, 0, 0}, cc, this.cAnim.GetValue(), 5, 2)
	}
	if *this.currState >= this.allowedState {
		ctx.DrawLine(CartesianToSSPoint(this.ps[0].p, wh), CartesianToSSPoint(this.freePoint, wh), c, c, 3)
		ctx.DrawLine(CartesianToSSPoint(this.ps[1].p, wh), CartesianToSSPoint(this.freePoint, wh), c, c, 3)
	}

}

type AnimPointCentroid struct {
	p            glm.Vec2
	animX, animY closedGL.Animation
}

func (this *AnimPointCentroid) process(delta float32) {
	this.animX.Process(delta)
	this.animY.Process(delta)
	this.p = glm.Vec2{this.animX.GetValue(), this.animY.GetValue()}
}

func (this *AnimPointCentroid) draw(ctx *closedGL.ClosedGLContext) {
	ctx.DrawCircle(CartesianToSSPoint(this.p, ctx.Window.Wh), glm.Vec4{1, 0, 0, 1}, glm.Vec4{1, 0, 0, 1}, 10, 1, 1)
}

type CentroidAnim struct {
	tri       *Tri
	sides     [3]AnimSideCentroid
	currState int
}

func newCentroidAnim(tri *Tri) CentroidAnim {
	var retVal = CentroidAnim{tri: tri, currState: 0}
	return retVal
}

func (this *CentroidAnim) init() {
	this.sides[0] = newAnimSideCentroid(this.tri.Points[0], this.tri.Points[1], this.tri.Points[2], 0, &this.currState)
	this.sides[1] = newAnimSideCentroid(this.tri.Points[0], this.tri.Points[2], this.tri.Points[1], 1, &this.currState)
	this.sides[2] = newAnimSideCentroid(this.tri.Points[1], this.tri.Points[2], this.tri.Points[0], 2, &this.currState)
}

func (this *CentroidAnim) Draw() {
	this.sides[0].Draw(this.tri.Ctx)
	this.sides[1].Draw(this.tri.Ctx)
	this.sides[2].Draw(this.tri.Ctx)
	if this.currState == 3 || this.currState == 6 {
		var lastMp = closedGL.MiddlePoint(this.tri.Points[1], this.tri.Points[2])
		this.tri.drawCentroid()
		this.tri.Ctx.DrawCircle(CartesianToSSPoint(lastMp, this.tri.Ctx.Window.Wh), glm.Vec4{1, 0, 0, 1}, glm.Vec4{1, 0, 0, 1}, 10, 3, 3)
	}
}

func (this *CentroidAnim) Process(delta float32) {
	this.sides[0].Process(delta)
	this.sides[1].Process(delta)
	this.sides[2].Process(delta)
}

func (this *CentroidAnim) IsFinished() bool {
	return this.currState == 3
}
