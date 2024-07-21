package main

import (
	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type OrthoCenterSideAnim struct {
	corner       glm.Vec2
	oppositeSide glm.Vec2
	vec          glm.Vec2
	unitAnim     closedGL.Animation
	rAnim        closedGL.Animation
}

func newOrthoSideCenterAnim(corner, p1, p2 glm.Vec2) OrthoCenterSideAnim {
	var retVal = OrthoCenterSideAnim{
		corner:       corner,
		oppositeSide: CalcLinearEquation(p1, p2),
		unitAnim:     closedGL.NewAnimation(0, 1, 1, false, false),
	}
	var p = findLineCircleIntersectionPoint(corner, p1, p2)
	var a = p.Sub(&corner)
	retVal.vec = a
	retVal.rAnim = closedGL.NewAnimation(0, retVal.vec.Len(), 1, false, false)
	return retVal
}

func (this *OrthoCenterSideAnim) setOrthocenter(orthocenter glm.Vec2) {
	var newVec = orthocenter.Sub(&this.corner)
	if newVec.Len() > this.vec.Len() {

		this.vec = orthocenter.Sub(&this.corner)
	}
}

func (this *OrthoCenterSideAnim) process(delta float32) {
	if !this.rAnim.IsFinished() {
		this.rAnim.Process(delta)
	} else {
		this.unitAnim.Process(delta)
	}
}

func (this *OrthoCenterSideAnim) draw(ctx *closedGL.ClosedGLContext, depth int) {
	if !this.rAnim.IsFinished() {
		drawCartesianCircle(this.corner, ctx, glm.Vec4{0, 0, 0, 0}, rgbToColour(54, 194, 206), 3, this.rAnim.GetValue(), float32(depth))
	}
	var p = this.corner
	var p2 = this.corner

	p.AddScaledVec(this.unitAnim.GetValue(), &this.vec)
	p2.AddScaledVec(-this.unitAnim.GetValue(), &this.vec)

	drawCartesianLine(this.corner, p, ctx, int(depth), rgbToColour(54, 194, 206))
	drawCartesianLine(this.corner, p2, ctx, int(depth), rgbToColour(54, 194, 206))

}

type OrthocenterAnim struct {
	tri       *Tri
	anims     [3]OrthoCenterSideAnim
	currState int
}

func newOrthoCenterAnim(tri *Tri) OrthocenterAnim {
	return OrthocenterAnim{
		tri:       tri,
		anims:     [3]OrthoCenterSideAnim{},
		currState: 0,
	}
}

func (this *OrthocenterAnim) setOrthocenter(orthocenter glm.Vec2) {
	this.anims[0].setOrthocenter(orthocenter)
	this.anims[1].setOrthocenter(orthocenter)
	this.anims[2].setOrthocenter(orthocenter)
}

func (this *OrthocenterAnim) init() {
	this.anims[0] = newOrthoSideCenterAnim(this.tri.Points[0], this.tri.Points[1], this.tri.Points[2])
	this.anims[1] = newOrthoSideCenterAnim(this.tri.Points[1], this.tri.Points[0], this.tri.Points[2])
	this.anims[2] = newOrthoSideCenterAnim(this.tri.Points[2], this.tri.Points[0], this.tri.Points[1])
}

func (this *OrthocenterAnim) Draw() {

	for i := 0; i < 3; i++ {
		if this.currState >= i {
			this.anims[i].draw(this.tri.Ctx, i+3)
		}
	}
	if this.currState == 3 {
		this.tri.drawOrthocenter()
	}

}

func (this *OrthocenterAnim) Process(delta float32) {
	for i := 0; i < 3; i++ {
		if this.currState >= i {
			this.anims[i].process(delta)
			if this.anims[i].unitAnim.IsFinished() {
				this.currState = i + 1
			}
		}
	}
}

func (this *OrthocenterAnim) IsFinished() bool {
	return this.currState == 3
}
