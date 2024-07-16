package main

import (
	"runtime"

	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

func main() {
	runtime.LockOSThread()
	var openGL = closedGL.InitClosedGL(800, 600)
	closedGL.LimitFPS(false)
	closedGL.PrintlnVec2(SSToCartesianVec)
	for !openGL.Window.Window.ShouldClose() {

		openGL.BeginDrawing()
		closedGL.ClearBG()
		openGL.DrawTriangle([3]glm.Vec2{{100, 100}, {50, 250}, {150, 250}}, glm.Vec4{1, 1, 1, 1})
		openGL.DrawFPS(500, 0)
		openGL.EndDrawing()
		openGL.Process()
	}
	println("test")
}
