package main

import (
	"runtime"

	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

func main() {
	runtime.LockOSThread()
	var openGL = closedGL.InitClosedGL(800, 600)
	openGL.LimitFPS(false)
	var tri = newTri([3]glm.Vec2{{400, 100}, {100, 500}, {700, 500}}, &openGL, glm.Vec4{0, 1, 1, 1})
	//var tri = newTri([3]glm.Vec2{{100, 000}, {500, 700}, {700, 500}}, &openGL, glm.Vec4{0, 1, 1, 1})

	//	tri.startIncenterAnim()
	tri.calcCircumcenter()
	tri.startCircumCenterAnim()
	closedGL.SetWireFrameMode(true)
	for !openGL.Window.Window.ShouldClose() {
		tri.process(float32(openGL.FPSCounter.Delta))

		openGL.BeginDrawing()
		openGL.ClearBG()
		tri.draw()
		tri.drawCentroid()
		openGL.DrawFPS(500, 0)
		openGL.EndDrawing()
		openGL.Process()
	}
}
