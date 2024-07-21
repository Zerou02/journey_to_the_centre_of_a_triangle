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
func main() {
	runtime.LockOSThread()
	var basePath = getExeBasePath()
	var typewriter = newTypewriterAnim(basePath + "/assets/intro.txt")
	_ = typewriter
	var openGL = closedGL.InitClosedGL(800, 600, "Centre of a triangle")
	var config = openGL.Config

	var sm = newStateMachine()
	_ = sm
	var vertices = [][3]glm.Vec2{
		{{400, 50}, {100, 500}, {700, 500}},
		{{350, 50}, {156, 500}, {123, 300}},
		{{370, 350}, {576, 500}, {345, 300}},
	}
	var tri = newTri(vertices[0], &openGL, rgbToColour(67, 61, 139))
	var cbs = [](func()){
		func() { tri.startCentroidAnim() },
		func() { tri.startCircumCenterAnim() },
		func() { tri.startCentroidAnim() },
	}
	var currIdx = 0
	tri.setVertices(vertices[currIdx])
	cbs[currIdx]()
	tri.drawCenters = false
	openGL.Window.SetMouseButtonCB(tri.mouseCB)
	openGL.Window.SetCursorPosCallback(tri.cursorCB)
	if config["bgm"] != "" && strToBool(config["bgm"]) {
		openGL.PlayMusic("bgm", 0.5)
	}

	closedGL.SetWireFrameMode(true)

	for !openGL.Window.Window.ShouldClose() {

		var delta = float32(openGL.FPSCounter.Delta)
		_ = delta
		tri.Process(delta)
		tri.showUi = true
		/* if tri.IsFinished() {
			currIdx++
			if currIdx >= len(vertices) {
			} else {

				tri.setVertices(vertices[currIdx])
				cbs[currIdx]()
			}
		} */
		openGL.BeginDrawing()
		openGL.ClearBG(rgbToColour(0, 0, 0))

		/* 	if !typewriter.anim.IsFinished() {
			openGL.ClearBG(rgbToColour(0, 0, 0))
			typewriter.process(delta)
			typewriter.draw(&openGL)
		} else { */
		openGL.ClearBG(rgbToColour(23, 21, 59))
		tri.Draw()
		/* 	} */
		openGL.DrawFPS(500, 0, 1)

		openGL.EndDrawing()
		openGL.Process()
	}
}
