package main

import (
	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type TriState interface {
	Process(delta float32)
	init()
	Draw()
	IsFinished() bool
}
type Tri struct {
	Ctx       *closedGL.ClosedGLContext
	Points    [3]glm.Vec2 //cartesian
	Colour    glm.Vec4
	currAnims [4]TriState
	centres   [4]glm.Vec2 //centroid,circumcenter,incenter,orthocenter

	buttonVecs        [5]glm.Vec4
	mouseDown         bool
	mousePos          glm.Vec2
	currentDraggedIdx int
	drawCenters       bool
	showUi            bool
}

func newTri(ssPoints [3]glm.Vec2, ctx *closedGL.ClosedGLContext, colour glm.Vec4) Tri {
	var retTri = Tri{Ctx: ctx, Colour: colour, currentDraggedIdx: -1, buttonVecs: [5]glm.Vec4{}, drawCenters: true, showUi: false}
	var baseX float32 = 450
	var baseY float32 = 550
	var gap float32 = 50
	var dim float32 = 25
	for i := 0; i < len(retTri.buttonVecs); i++ {
		retTri.buttonVecs[i] = glm.Vec4{baseX + float32(i)*(gap+dim), baseY, dim, dim}
	}

	retTri.setVertices(ssPoints)
	return retTri
}

func (this *Tri) setVertices(ssPoints [3]glm.Vec2) {
	for i := 0; i < len(ssPoints); i++ {
		this.Points[i] = SSToCartesianPoint(ssPoints[i], this.Ctx.Window.Wh)
	}
	this.currAnims = [4]TriState{}
	this.calcCentres()
}
func (this *Tri) calcCentres() {
	this.centres[0] = this.calcCentroid()
	this.centres[1] = this.calcCircumcenter()
	this.centres[2] = this.calcIncenter()
	this.centres[3] = this.calcOrthocenter()
}

func (this *Tri) IsFinished() bool {
	var oneFinished = true
	for _, x := range this.currAnims {
		if x != nil && !x.IsFinished() {
			oneFinished = false
			break
		}
	}
	return !this.showUi && oneFinished
}

func (this *Tri) startCircumCenterAnim(animDur float32) {
	var state = newCircumCenterAnim(this, animDur)
	this.currAnims[1] = &state
}
func (this *Tri) startCentroidAnim(animDur float32) {
	var state = newCentroidAnim(this, animDur)
	state.init()
	this.currAnims[0] = &state
}

func (this *Tri) startIncenterAnim(animDur float32) {
	var state = newIncenterAnim(this, animDur)
	state.init()
	this.currAnims[2] = &state
}

func (this *Tri) startOrthocenterAnim(animDur float32) {
	var state = newOrthoCenterAnim(this, animDur)
	state.init()
	var ortho = this.calcOrthocenter()
	state.setOrthocenter(ortho)
	this.currAnims[3] = &state

}

func (this *Tri) Draw() {

	println("--")
	var ssPoints [3]glm.Vec2 = [3]glm.Vec2{}
	for i, x := range this.Points {
		ssPoints[i] = SSToCartesianPoint(x, this.Ctx.Window.Wh)
		var cartMouse = SSToCartesianPoint(this.mousePos, this.Ctx.Window.Wh)
		var dist = cartMouse.Sub(&this.Points[i])
		closedGL.PrintlnVec2(ssPoints[i])
		if dist.Len() < 10 && this.showUi {
			drawCartesianCircle(x, this.Ctx, glm.Vec4{0, 0.25, 0, 1}, glm.Vec4{1, 1, 1, 1}, 6, 10, 3)
		} else if this.showUi {
			drawCartesianCircle(x, this.Ctx, glm.Vec4{1, 1, 1, 1}, glm.Vec4{1, 1, 1, 1}, 6, 10, 3)
		}
	}
	this.Ctx.DrawTriangle(ssPoints, this.Colour, 0)

	drawCartesianLine(this.Points[0], this.Points[1], this.Ctx, 1, glm.Vec4{1, 1, 1, 1})
	drawCartesianLine(this.Points[0], this.Points[2], this.Ctx, 1, glm.Vec4{1, 1, 1, 1})
	drawCartesianLine(this.Points[1], this.Points[2], this.Ctx, 1, glm.Vec4{1, 1, 1, 1})

	for i := 0; i < len(this.currAnims); i++ {
		if this.currAnims[i] != nil {
			this.currAnims[i].Draw()
		}
	}

	var colours = []glm.Vec4{
		{1, 0, 0, 0.5}, {1, 1, 0, 0.5}, {0, 0.5, 0, 0.5}, rgbaToColour(54, 194, 206, 0.5),
	}
	if this.drawCenters {
		colours = append(colours, glm.Vec4{1, 1, 1, 1})
	} else {
		colours = append(colours, glm.Vec4{0.35, 0.35, 0.35, 1})

	}
	if this.showUi {
		for i, x := range this.buttonVecs {
			this.Ctx.DrawRect(x, colours[i], 10)
		}
	}
	if this.drawCenters {
		this.drawCentroid()
		this.drawCircumCenter()
		this.drawIncenter()
		this.drawOrthocenter()
	}
}

func (this *Tri) Process(delta float32) {

	for i := 0; i < len(this.currAnims); i++ {
		if this.currAnims[i] != nil {
			this.currAnims[i].Process(delta)
		}
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
	var v1 = findAngleBisectorVec(this.Points[0], this.Points[1], this.Points[2], this.calcCentroid())
	var v2 = findAngleBisectorVec(this.Points[1], this.Points[0], this.Points[2], this.calcCentroid())
	var oppositeP1 = this.Points[0].Add(&v1)
	var oppositeP2 = this.Points[1].Add(&v2)

	var eq1 = CalcLinearEquation(this.Points[0], oppositeP1)
	var eq2 = CalcLinearEquation(this.Points[1], oppositeP2)

	return CalcCrossingPoint(eq1, eq2)
}

func (this *Tri) calcCircumcenter() glm.Vec2 {
	var line1 = CalcPerpLineVec(this.Points[0], this.Points[1])
	var line2 = CalcPerpLineVec(this.Points[1], this.Points[2])
	return CalcCrossingPoint(line1, line2)
}

func (this *Tri) calcOrthocenter() glm.Vec2 {
	var anim = newOrthoCenterAnim(this, 1)
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
	this.Ctx.DrawCircle(CartesianToSSPoint(this.centres[0], this.Ctx.Window.Wh), glm.Vec4{1, 0, 0, 0.5}, glm.Vec4{1, 1, 0, 0.5}, 10, 3, 6)
}

func (this *Tri) drawCircumCenter() {
	this.Ctx.DrawCircle(CartesianToSSPoint(this.centres[1], this.Ctx.Window.Wh), glm.Vec4{1, 1, 0, 0.5}, glm.Vec4{1, 1, 1, 0.5}, 10, 3, 6)
}

func (this *Tri) drawIncenter() {
	this.Ctx.DrawCircle(CartesianToSSPoint(this.centres[2], this.Ctx.Window.Wh), glm.Vec4{0, 0.5, 0, 0.5}, glm.Vec4{1, 1, 1, 0.5}, 10, 3, 6)
}

func (this *Tri) drawOrthocenter() {
	this.Ctx.DrawCircle(CartesianToSSPoint(this.centres[3], this.Ctx.Window.Wh), rgbaToColour(54, 194, 206, 0.5), glm.Vec4{1, 1, 1, 0.5}, 10, 3, 6)
}

func (this *Tri) mouseCB(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if button != glfw.MouseButton1 {
		return
	}
	if action == glfw.Press && this.showUi {
		this.mouseDown = true

		var threshold float32 = 10
		var cart = SSToCartesianPoint(this.mousePos, this.Ctx.Window.Wh)
		for i, x := range this.Points {
			if distBetweenPoints(cart, x) < float32(threshold) {
				this.currentDraggedIdx = i
				break
			}
		}
		var cbs = [](func()){
			func() { this.startCentroidAnim(1) },
			func() { this.startCircumCenterAnim(1) },
			func() { this.startIncenterAnim(1) },
			func() { this.startOrthocenterAnim(1) },
			func() { this.drawCenters = !this.drawCenters },
		}
		for i, x := range this.buttonVecs {
			if pointInRect(this.mousePos, x) {
				cbs[i]()
			}
		}
	}
	if action == glfw.Release {
		this.currentDraggedIdx = -1
		this.mouseDown = false
	}
}

func (this *Tri) cursorCB(w *glfw.Window, xpos float64, ypos float64) {
	this.mousePos = glm.Vec2{float32(xpos), float32(ypos)}
	var newPoints = [3]glm.Vec2{}
	if this.currentDraggedIdx != -1 {
		for i := 0; i < len(this.Points); i++ {
			if i == this.currentDraggedIdx {
				newPoints[i] = glm.Vec2{
					closedGL.Clamp(0, this.Ctx.Window.Ww, this.mousePos[0]),
					closedGL.Clamp(0, this.Ctx.Window.Wh, this.mousePos[1]),
				}
				continue
			}
			newPoints[i] = CartesianToSSPoint(this.Points[i], this.Ctx.Window.Wh)
		}
		this.setVertices(newPoints)
	}
}
