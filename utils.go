package main

import (
	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

func drawCartesianLine(cp1, cp2 glm.Vec2, ctx *closedGL.ClosedGLContext, depth int) {
	ctx.DrawLine(CartesianToSSPoint(cp1, ctx.Window.Wh), CartesianToSSPoint(cp2, ctx.Window.Wh), glm.Vec4{1, 1, 1, 1}, glm.Vec4{1, 1, 1, 1}, depth)
}

func drawCartesianCircle(p glm.Vec2, ctx *closedGL.ClosedGLContext, colour, borderColour glm.Vec4, depth int, radius, borderThickness float32) {
	ctx.DrawCircle(CartesianToSSPoint(p, ctx.Window.Wh), colour, borderColour, radius, borderThickness, depth)
}
func lerpAnim(startVal, oldTime, newTime float32) float32 {
	return startVal / oldTime * newTime
}
