package main

import (
	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type TriState interface {
	process(delta float32)
	init()
	draw()
}
type Tri struct {
	Ctx       *closedGL.ClosedGLContext
	Points    [3]glm.Vec2 //cartesian
	Colour    glm.Vec4
	CurrState TriState
	centres   [4]glm.Vec2 //centroid,circumcenter,incenter,orthocenter
}

func newTri(ssPoints [3]glm.Vec2, ctx *closedGL.ClosedGLContext, colour glm.Vec4) Tri {
	var retTri = Tri{Ctx: ctx, Colour: colour}
	for i := 0; i < len(ssPoints); i++ {
		retTri.Points[i] = SSToCartesianPoint(ssPoints[i], ctx.Window.Wh)
	}
	retTri.centres[0] = retTri.calcCentroid()
	retTri.centres[1] = retTri.calcCircumcenter()
	retTri.centres[2] = retTri.calcIncenter()
	retTri.centres[3] = retTri.calcOrthocenter()
	return retTri
}

func (this *Tri) startCircumCenterAnim() {
	var state = newCircumCenterAnim(this)
	state.init()
	state.setCircumcenter(this.calcCircumcenter())
	this.CurrState = &state
}
func (this *Tri) startCentroidAnim() {
	var state = newCentroidAnim(this)
	state.init()
	this.CurrState = &state
}

func (this *Tri) startIncenterAnim() {
	var state = newIncenterAnim(this)
	state.init()
	this.CurrState = &state
}

func (this *Tri) startOrthocenterAnim() {
	var state = newOrthoCenterAnim(this)
	state.init()
	var ortho = this.calcOrthocenter()
	state.setOrthocenter(ortho)
	this.CurrState = &state
}

func (this *Tri) draw() {

	var ssPoints [3]glm.Vec2 = [3]glm.Vec2{}
	for i := 0; i < len(this.Points); i++ {
		ssPoints[i] = SSToCartesianPoint(this.Points[i], this.Ctx.Window.Wh)
	}
	this.Ctx.DrawTriangle(ssPoints, this.Colour, 0)

	if this.CurrState != nil {
		this.CurrState.draw()
	}
	drawCartesianCircle(this.calcCircumcenter(), this.Ctx, glm.Vec4{1, 1, 1, 1}, glm.Vec4{1, 1, 1, 1}, 2, 10, 3)
}

func (this *Tri) process(delta float32) {
	if this.CurrState != nil {
		this.CurrState.process(delta)
	}
}

func (this *Tri) calcCentroid() glm.Vec2 {
	var mp1 = CalcMiddlePoint(this.Points[0], this.Points[1])
	var mp2 = CalcMiddlePoint(this.Points[0], this.Points[2])
	var eq1 = CalcLinearEquation(mp1, this.Points[2])
	var eq2 = CalcLinearEquation(mp2, this.Points[1])
	var p = CalcCrossingPoint(eq1, eq2)
	return p
}

func (this *Tri) calcIncenter() glm.Vec2 {
	var anim = newIncenterAnim(this)
	anim.init()
	var eq1 = CalcLinearEquation(anim.cornerAnims[0].cornerP, anim.cornerAnims[0].dirP)
	var eq2 = CalcLinearEquation(anim.cornerAnims[1].cornerP, anim.cornerAnims[1].dirP)
	return CalcCrossingPoint(eq1, eq2)
}

func (this *Tri) calcCircumcenter() glm.Vec2 {
	var anim = newCircumCenterAnim(this)
	anim.init()

	var p2 = anim.anims[0].basePoint
	p2.AddScaledVec(50, &anim.anims[0].vec)

	var p3 = anim.anims[1].basePoint
	p3.AddScaledVec(50, &anim.anims[1].vec)

	var test = CalcLinearEquation(anim.anims[0].basePoint, p2)
	var test2 = CalcLinearEquation(anim.anims[1].basePoint, p3)

	return CalcCrossingPoint(test, test2)
}

func (this *Tri) calcOrthocenter() glm.Vec2 {
	var anim = newOrthoCenterAnim(this)
	anim.init()
	var other1 = anim.anims[0].corner
	other1.AddScaledVec(1, &anim.anims[0].vec)
	var eq1 = CalcLinearEquation(anim.anims[0].corner, other1)

	var other2 = anim.anims[1].corner
	other2.AddScaledVec(1, &anim.anims[1].vec)
	var eq2 = CalcLinearEquation(anim.anims[1].corner, other2)
	return CalcCrossingPoint(eq1, eq2)
}

func (this *Tri) drawCentroid() {
	this.Ctx.DrawCircle(CartesianToSSPoint(this.centres[0], this.Ctx.Window.Wh), glm.Vec4{1, 0, 0, 1}, glm.Vec4{1, 1, 0, 1}, 10, 3, 1)
}

func (this *Tri) drawCircumCenter() {
	this.Ctx.DrawCircle(CartesianToSSPoint(this.centres[1], this.Ctx.Window.Wh), glm.Vec4{1, 1, 1, 1}, glm.Vec4{1, 1, 0, 1}, 10, 3, 2)
}

func (this *Tri) drawIncenter() {
	this.Ctx.DrawCircle(CartesianToSSPoint(this.centres[2], this.Ctx.Window.Wh), glm.Vec4{0, 1, 0, 1}, glm.Vec4{1, 1, 0, 1}, 10, 3, 6)
}

func (this *Tri) drawOrthocenter() {
	this.Ctx.DrawCircle(CartesianToSSPoint(this.centres[3], this.Ctx.Window.Wh), glm.Vec4{0, 0.5, 0.5, 1}, glm.Vec4{1, 1, 0, 1}, 10, 3, 6)
}

func (this *Tri) drawCentroidLines() {

	var mp1 = CalcMiddlePoint(this.Points[0], this.Points[1])
	var mp2 = CalcMiddlePoint(this.Points[0], this.Points[2])
	var mp3 = CalcMiddlePoint(this.Points[1], this.Points[2])

	var offset = glm.Vec2{0.1, 0}
	var offsetP = mp3.Add(&offset)

	var eq1 = CalcLinearEquation(this.Points[0], offsetP)

	var eq2 = CalcLinearEquation(this.Points[2], mp1)
	var eq3 = CalcLinearEquation(this.Points[1], mp2)

	var p = glm.Vec2{this.Points[0][0], this.Points[0][0]*eq1[0] + eq1[1]}
	var p1 = glm.Vec2{offsetP[0], offsetP[0]*eq1[0] + eq1[1]}

	var p3 = glm.Vec2{this.Points[2][0], this.Points[2][0]*eq2[0] + eq2[1]}
	var p4 = glm.Vec2{mp1[0], mp1[0]*eq2[0] + eq2[1]}

	var p5 = glm.Vec2{this.Points[1][0], this.Points[1][0]*eq3[0] + eq3[1]}
	var p6 = glm.Vec2{mp2[0], mp2[0]*eq3[0] + eq3[1]}

	var cp = CalcCrossingPoint(eq1, eq2)
	var cp2 = CalcCrossingPoint(eq1, eq3)
	var cp3 = CalcCrossingPoint(eq2, eq3)
	_ = cp3

	this.Ctx.DrawCircle(CartesianToSSPoint(mp1, this.Ctx.Window.Wh), glm.Vec4{1, 0, 0, 1}, glm.Vec4{1, 0, 0, 1}, 20, 1, 2)
	this.Ctx.DrawCircle(CartesianToSSPoint(mp2, this.Ctx.Window.Wh), glm.Vec4{1, 0, 0, 1}, glm.Vec4{1, 0, 0, 1}, 20, 1, 2)
	this.Ctx.DrawCircle(CartesianToSSPoint(mp3, this.Ctx.Window.Wh), glm.Vec4{1, 0, 0, 1}, glm.Vec4{1, 0, 0, 1}, 20, 1, 2)

	this.Ctx.DrawLine(CartesianToSSPoint(p, this.Ctx.Window.Wh), CartesianToSSPoint(p1, this.Ctx.Window.Wh), glm.Vec4{1, 1, 0, 1}, glm.Vec4{1, 1, 0, 1}, 1)
	this.Ctx.DrawLine(CartesianToSSPoint(p3, this.Ctx.Window.Wh), CartesianToSSPoint(p4, this.Ctx.Window.Wh), glm.Vec4{1, 1, 0, 1}, glm.Vec4{1, 1, 0, 1}, 1)
	this.Ctx.DrawLine(CartesianToSSPoint(p5, this.Ctx.Window.Wh), CartesianToSSPoint(p6, this.Ctx.Window.Wh), glm.Vec4{1, 1, 0, 1}, glm.Vec4{1, 1, 0, 1}, 1)
	this.Ctx.DrawCircle(CartesianToSSPoint(cp, this.Ctx.Window.Wh), glm.Vec4{1, 0, 0, 1}, glm.Vec4{0, 1, 0, 1}, 10, 3, 2)
	this.Ctx.DrawCircle(CartesianToSSPoint(cp2, this.Ctx.Window.Wh), glm.Vec4{0, 1, 0, 1}, glm.Vec4{0, 1, 0, 1}, 10, 3, 2)
	this.Ctx.DrawCircle(CartesianToSSPoint(cp3, this.Ctx.Window.Wh), glm.Vec4{0, 0, 1, 1}, glm.Vec4{0, 1, 0, 1}, 10, 2, 2)
}
