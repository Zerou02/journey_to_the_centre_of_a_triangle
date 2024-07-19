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
	retVal.findR(p1, p2)
	return retVal
}

func (this *OrthoCenterSideAnim) setOrthocenter(orthocenter glm.Vec2) {
	/* var a = orthocenter.Sub(&this.corner)
	if a.Len() > this.vec.Len() {
		this.vec = a
	} */
}

func (this *OrthoCenterSideAnim) findR(p1, p2 glm.Vec2) {
	var r float32 = 0
	var step float32 = 1
	var targetPoint = glm.Vec2{0, 0}

	if p1[0] == p2[0] {
		targetPoint = glm.Vec2{p1[0], this.oppositeSide[0]*p1[0] + this.oppositeSide[1]}
	}
	var currOffsets = LineCircleIntersection(r, this.oppositeSide, this.corner)
	for distToLine(p1, p2, targetPoint) > glm.Epsilon {
		currOffsets = LineCircleIntersection(r, this.oppositeSide, this.corner)
		var len = len(currOffsets)
		if len == 2 {
			if step > glm.Epsilon {
				r -= step
			} else {
				r -= glm.Epsilon
			}
			step *= 0.1
			targetPoint = currOffsets[0]
			if step < glm.Epsilon {
				break
			}
		}
		if len == 0 {
			if step < glm.Epsilon {
				r += glm.Epsilon
			} else {
				var rNew = step + r
				if rNew == r {
					var i float32 = 0
					for rNew == r {
						i++
						rNew = i*step + r
					}
					r = rNew
					break
				}
				r = rNew
			}
		}
		if len == 1 {
			targetPoint = currOffsets[0]
			break
		}
	}
	if len(currOffsets) == 2 {
		targetPoint = closedGL.MiddlePoint(currOffsets[0], currOffsets[1])
	}
	var a = targetPoint.Sub(&this.corner)
	this.vec = a
	this.rAnim = closedGL.NewAnimation(0, r, 1, false, false)
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
		drawCartesianCircle(this.corner, ctx, glm.Vec4{1, 0, 0, 0}, glm.Vec4{1, 0, 0, 1}, 3, this.rAnim.GetValue(), float32(depth))
	}
	var p = this.corner
	p.AddScaledVec(this.unitAnim.GetValue(), &this.vec)
	drawCartesianLine(this.corner, p, ctx, int(depth))
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

func (this *OrthocenterAnim) draw() {

	for i := 0; i < 3; i++ {
		if this.currState >= i {
			this.anims[i].draw(this.tri.Ctx, i+3)
		}
	}
	if this.currState == 3 {
		this.tri.drawOrthocenter()
	}

}

func (this *OrthocenterAnim) process(delta float32) {
	for i := 0; i < 3; i++ {
		if this.currState >= i {
			this.anims[i].process(delta)
			if this.anims[i].unitAnim.IsFinished() {
				this.currState = i + 1
			}
		}
	}
	if this.currState == 3 {
	}
}
