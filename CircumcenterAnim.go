package main

import (
	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type CircumcenterSideAnim struct {
	basePoint, vec    glm.Vec2
	circumcenter      glm.Vec2
	p1, p2            glm.Vec2
	animX, animY      closedGL.Animation
	animR             closedGL.Animation
	targetR, targetPD float32
	animCenter        closedGL.Animation
	vec2              glm.Vec2
	animUnit          closedGL.Animation
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

	this.vec2 = this.circumcenter.Sub(&this.basePoint)
	var centerDiff = this.p1.Sub(&this.circumcenter)
	this.animX = closedGL.NewAnimation(this.basePoint[0], this.circumcenter[0], 1, false, false)
	this.animY = closedGL.NewAnimation(this.basePoint[1], this.circumcenter[1], 1, false, false)
	this.animUnit = closedGL.NewAnimation(0, 1.1, 1, false, false)
	this.targetR = diff.Len()
	this.animR = closedGL.NewAnimation(0, this.targetR*1.1, 3, false, false)
	this.animCenter = closedGL.NewAnimation(0, centerDiff.Len(), 3, false, false)

}

func (this *CircumcenterSideAnim) process(delta float32) {
	this.animR.Process(delta)
	if this.animR.GetValue() >= this.targetR {
		this.animX.Process(delta)
		this.animY.Process(delta)
		this.animUnit.Process(delta)
	}
	if this.animUnit.GetValue() > 1 {
		this.animCenter.Process(delta)
	}
}

func (this *CircumcenterSideAnim) draw(ctx *closedGL.ClosedGLContext, depth int) {
	var p = glm.Vec2{this.animX.GetValue(), this.animY.GetValue()}
	var diff = p.Sub(&this.basePoint)
	var p2 = this.basePoint.Sub(&diff)
	_ = p2
	var bc = glm.Vec4{1, 0, 0, 1}
	if !this.animR.IsFinished() {
		drawCartesianCircle(this.p1, ctx, glm.Vec4{0, 0, 0, 0}, bc, 3, this.animR.GetValue(), float32(depth))
		drawCartesianCircle(this.p2, ctx, glm.Vec4{0, 0, 0, 0}, bc, 3, this.animR.GetValue(), float32(depth))
	}
	if this.animUnit.IsFinished() {
		drawCartesianCircle(this.circumcenter, ctx, glm.Vec4{0, 0, 0, 0}, bc, 3, this.animCenter.GetValue(), float32(depth))
	}

	var p3 = this.basePoint
	p3.AddScaledVec(this.animUnit.GetValue(), &this.vec2)
	var p4 = this.basePoint
	p4.AddScaledVec(-this.animUnit.GetValue(), &this.vec2)
	drawCartesianLine(this.basePoint, p3, ctx, depth)
	drawCartesianLine(this.basePoint, p4, ctx, depth)

}

type CircumcenterAnim struct {
	tri          *Tri
	anims        [3]CircumcenterSideAnim
	circumcenter glm.Vec2
}

func newCircumCenterAnim(tri *Tri) CircumcenterAnim {
	return CircumcenterAnim{tri: tri, anims: [3]CircumcenterSideAnim{}}
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
}

func (this *CircumcenterAnim) draw() {
	for i := 0; i < 3; i++ {
		this.anims[i].draw(this.tri.Ctx, i+3)
	}
}

func (this *CircumcenterAnim) process(delta float32) {
	for i := 0; i < 3; i++ {
		this.anims[i].process(delta)
	}
}
