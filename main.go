package main

import (
	"os"
	"runtime"
	"strings"

	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

func rgbaToColour(r, g, b, a float32) glm.Vec4 {
	return glm.Vec4{r / 255, g / 255, b / 255, a}
}

func rgbToColour(r, g, b float32) glm.Vec4 {
	return glm.Vec4{r / 255, g / 255, b / 255, 1}
}

func getExeBasePath() string {
	var exePath, _ = (os.Executable())
	var b = strings.Split(exePath, "/")
	b[len(b)-1] = ""
	return strings.Join(b, "/")
}

func strToBool(str string) bool {
	if str == "true" {
		return true
	} else {
		return false
	}
}

func createTri(vertices [3]glm.Vec2, anim int, ctx *closedGL.ClosedGLContext, animDur float32) Tri {
	var tri = newTri(vertices, ctx, rgbToColour(67, 61, 139))
	tri.showUi = false
	tri.drawCenters = false
	var cbs = [](func(tri *Tri)){
		func(tri *Tri) { tri.startCentroidAnim(animDur) },
		func(tri *Tri) { tri.startCircumCenterAnim(animDur) },
		func(tri *Tri) { tri.startIncenterAnim(animDur) },
		func(tri *Tri) { tri.startOrthocenterAnim(animDur) },
	}
	cbs[anim](&tri)
	return tri
}

func buildSm(ctx *closedGL.ClosedGLContext) StateMachine {
	var vertices = [][3]glm.Vec2{
		{{320, 76}, {60, 327}, {647, 305}},
		{{513, 503}, {384, 280}, {572, 308}},
		{{647, 305}, {72, 71}, {116, 560}},

		{{579, 440}, {416, 82}, {174, 544}},
		{{623, 282}, {227, 186}, {53, 427}},
		{{639, 557}, {67, 153}, {67, 557}},

		{{535, 192}, {100, 500}, {700, 500}},
		{{423, 88}, {100, 374}, {460, 476}},
		{{268, 423}, {249, 511}, {745, 326}},

		{{547, 525}, {421, 23}, {335, 527}},
		{{527, 453}, {507, 159}, {147, 342}},
		{{382, 398}, {152, 401}, {453, 542}},
	}
	var texts = []string{"1. Centroid", "2. Circumcenter", "3. Incenter", "4. Orthocenter"}
	var sm = newStateMachine()
	for i := 0; i < len(vertices); i++ {
		if i%3 == 0 {
			var scrn = newTextScreen(texts[i/3], 2, 2, ctx)
			sm.addState(&scrn)

		}
		var tri = createTri(vertices[i], i/3, ctx, 3)
		sm.addState(&tri)
	}
	return sm
}
func main() {
	runtime.LockOSThread()
	var basePath = getExeBasePath()
	var typewriter = newTypewriterAnim(basePath + "/assets/intro.txt")
	_ = typewriter
	var openGL = closedGL.InitClosedGL(800, 600, "Centre of a triangle")
	var config = openGL.Config

	var sm = buildSm(&openGL)
	var tri = newTri([3]glm.Vec2{{400, 100}, {100, 100}, {200, 200}}, &openGL, rgbToColour(67, 61, 139))

	openGL.Window.SetMouseButtonCB(tri.mouseCB)
	openGL.Window.SetCursorPosCallback(tri.cursorCB)

	if config["bgm"] != "" && strToBool(config["bgm"]) {
		openGL.PlayMusic("bgm", 0.5)
	}

	closedGL.SetWireFrameMode(true)

	for !openGL.Window.Window.ShouldClose() {

		var delta = float32(openGL.FPSCounter.Delta)
		sm.process(delta)
		openGL.BeginDrawing()
		openGL.ClearBG(rgbToColour(23, 21, 59))

		/* 	if !typewriter.anim.IsFinished() {
			openGL.ClearBG(rgbToColour(0, 0, 0))
			typewriter.process(delta)
			typewriter.draw(&openGL)
		} else { */
		/* 	} */
		sm.draw()
		openGL.DrawFPS(500, 0, 1)

		openGL.EndDrawing()
		openGL.Process()
	}
}
