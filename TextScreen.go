package main

import (
	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type TextScreen struct {
	text     string
	scale    int
	duration closedGL.Animation
	ctx      *closedGL.ClosedGLContext
}

func newTextScreen(text string, scale int, duration float32, ctx *closedGL.ClosedGLContext) TextScreen {
	return TextScreen{
		text, scale, closedGL.NewAnimation(0, 1, duration, false, false), ctx,
	}
}

func (this *TextScreen) Process(delta float32) { this.duration.Process(delta) }
func (this *TextScreen) Draw() {
	this.ctx.ClearBG(glm.Vec4{0, 0, 0, 0})
	this.ctx.Text.DrawText(200, 200, this.text, float32(this.scale))
}
func (this *TextScreen) IsFinished() bool { return this.duration.IsFinished() }
