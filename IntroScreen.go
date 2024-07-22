package main

import (
	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type IntroScreen struct {
	duration closedGL.Animation
	ctx      *closedGL.ClosedGLContext
}

func newIntroScreen(duration float32, ctx *closedGL.ClosedGLContext) IntroScreen {
	return IntroScreen{
		closedGL.NewAnimation(0, 1, duration, false, false), ctx,
	}
}

func (this *IntroScreen) Process(delta float32) { this.duration.Process(delta) }
func (this *IntroScreen) Draw() {
	this.ctx.ClearBG(glm.Vec4{0, 0, 0, 0})
	this.ctx.Text.DrawText(100, 200, "Journey to", 3)
	this.ctx.Text.DrawText(50, 300, "the center of a triangle", 2)

}
func (this *IntroScreen) IsFinished() bool { return this.duration.IsFinished() }
