package main

import (
	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type CircumcenterSideAnim struct {
	basePoint, vec glm.Vec2
	circumcenter   glm.Vec2
	p1, p2         glm.Vec2
	animR          closedGL.Animation
	targetR        float32
	animUnit       closedGL.Animation
}

func newCircumCenterSideAnim(p1, p2 glm.Vec2) CircumcenterSideAnim {
	var vec = p2.Sub(&p1)
	var perp = vec.Perp()
	perp = perp.Normalized()
	var mp = closedGL.MiddlePoint(p1, p2)

	return CircumcenterSideAnim{
		basePoint: mp,
		vec:       perp,
		p1:        p1,
		p2:        p2,
	}
}

func (this *CircumcenterSideAnim) setCircumcenter(circumcenter glm.Vec2) {
	var mp = closedGL.MiddlePoint(this.p1, this.p2)

	var diff = mp.Sub(&this.p1)
	this.circumcenter = circumcenter

	this.animUnit = closedGL.NewAnimation(0, 1.1, 1, false, false)
	this.targetR = diff.Len()
	this.animR = closedGL.NewAnimation(0, this.targetR*1, 3, false, false)

}

func (this *CircumcenterSideAnim) process(delta float32) {
	this.animR.Process(delta)
	if this.animR.GetValue() >= this.targetR {
		this.animUnit.Process(delta)
	}
}

func (this *CircumcenterSideAnim) draw(ctx *closedGL.ClosedGLContext) {
	var bc = glm.Vec4{1, 1, 0, 1}
	if !this.animR.IsFinished() {
		drawCartesianCircle(this.p1, ctx, glm.Vec4{0, 0, 0, 0}, bc, 3, this.animR.GetValue(), 2)
		drawCartesianCircle(this.p2, ctx, glm.Vec4{0, 0, 0, 0}, bc, 3, this.animR.GetValue(), 2)
	}
	var vec2 = this.circumcenter.Sub(&this.basePoint)
	var pForward = this.basePoint
	pForward.AddScaledVec(this.animUnit.GetValue(), &vec2)
	var pBack = this.basePoint
	pBack.AddScaledVec(-this.animUnit.GetValue(), &vec2)
	drawCartesianLine(this.basePoint, pForward, ctx, 2, glm.Vec4{1, 1, 0, 1})
	drawCartesianLine(this.basePoint, pBack, ctx, 2, glm.Vec4{1, 1, 0, 1})

}

type CircumcenterAnim struct {
	tri          *Tri
	anims        [3]CircumcenterSideAnim
	circumcenter glm.Vec2
	currState    int
	animCenter   closedGL.Animation
}

func newCircumCenterAnim(tri *Tri) CircumcenterAnim {
	return CircumcenterAnim{
		tri:   tri,
		anims: [3]CircumcenterSideAnim{},
	}

}

func (this *CircumcenterAnim) init() {
	this.anims[0] = newCircumCenterSideAnim(this.tri.Points[0], this.tri.Points[1])
	this.anims[1] = newCircumCenterSideAnim(this.tri.Points[0], this.tri.Points[2])
	this.anims[2] = newCircumCenterSideAnim(this.tri.Points[1], this.tri.Points[2])

}

func (this *CircumcenterAnim) setCircumcenter(circumCenter glm.Vec2) {
	this.anims[0].setCircumcenter(circumCenter)
	this.anims[1].setCircumcenter(circumCenter)
	this.anims[2].setCircumcenter(circumCenter)

	var vec = circumCenter.Sub(&this.tri.Points[0])
	this.animCenter = closedGL.NewAnimation(0, vec.Len(), 1, false, false)
}

func (this *CircumcenterAnim) draw() {
	for i := 0; i < 3; i++ {
		if this.currState >= i {
			this.anims[i].draw(this.tri.Ctx)
			if this.anims[i].animUnit.IsFinished() {
				this.currState = i + 1
			}
		}
	}
	if this.currState == 3 {
		this.tri.drawCircumCenter()
		drawCartesianCircle(this.tri.calcCircumcenter(), this.tri.Ctx, glm.Vec4{1, 0, 0, 0}, glm.Vec4{1, 1, 0, 1}, 3, this.animCenter.GetValue(), 3)
	}
}

func (this *CircumcenterAnim) process(delta float32) {
	for i := 0; i < 3; i++ {
		if this.currState >= i {
			this.anims[i].process(delta)
		}
	}
	if this.currState == 3 {
		this.animCenter.Process(delta)
	}
}

func (this *CircumcenterAnim) isFinished() bool {
	return this.animCenter.IsFinished()
}
