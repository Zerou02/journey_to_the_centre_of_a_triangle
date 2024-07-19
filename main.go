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
	//var tri = newTri([3]glm.Vec2{{400, 50}, {100, 500}, {700, 500}}, &openGL, glm.Vec4{0, 1, 1, 1})
	var tri = newTri([3]glm.Vec2{{100, 500}, {500, 400}, {700, 600}}, &openGL, glm.Vec4{0, 1, 1, 1})

	tri.startIncenterAnim()
	tri.startOrthocenterAnim()
	//tri.startCentroidAnim()
	//tri.startCircumCenterAnim()
	closedGL.SetWireFrameMode(true)
	/* for _, x := range LineCircleIntersection(1, glm.Vec2{3, 17}, glm.Vec2{-5, -1}) {
		closedGL.PrintlnVec2(x)
	}
	*/
	for !openGL.Window.Window.ShouldClose() {
		tri.process(float32(openGL.FPSCounter.Delta))

		openGL.BeginDrawing()
		openGL.ClearBG()
		tri.draw()
		tri.drawCentroid()
		tri.drawCircumCenter()
		tri.drawIncenter()
		tri.drawOrthocenter()
		openGL.DrawFPS(500, 0)
		openGL.EndDrawing()
		openGL.Process()
	}
}
