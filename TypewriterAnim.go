package main

import (
	"os"
	"strings"

	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type TypewriterAnim struct {
	anim        closedGL.Animation
	lines       []string
	linesToDraw []string
	ctx         *closedGL.ClosedGLContext
	endAnim     closedGL.Animation
}

func newTypewriterAnim(path string, ctx *closedGL.ClosedGLContext, secondsPerChars float32) TypewriterAnim {
	var bytees, _ = os.ReadFile(path)
	var contents = string(bytees)
	var lines []string = []string{}
	var amountChars = 0
	for _, x := range strings.Split(contents, "\n") {
		amountChars += len(x)
		lines = append(lines, x)
	}
	var time = float32(amountChars) * secondsPerChars
	var anim = closedGL.NewAnimation(0, float32(amountChars), time, false, false)
	return TypewriterAnim{
		anim:    anim,
		lines:   lines,
		ctx:     ctx,
		endAnim: closedGL.NewAnimation(0, 1, 1, false, false),
	}
}

func (this *TypewriterAnim) Process(delta float32) {
	this.anim.Process(delta)
	var currTextLen = this.anim.GetValue()
	this.linesToDraw = []string{}

	var alreadyDrawn = 0
	for i := 0; i < len(this.lines); i++ {
		if alreadyDrawn >= int(currTextLen) {
			break
		}
		var line = this.lines[i]
		var lineToDraw = ""
		if len(line)+alreadyDrawn < int(currTextLen) {
			alreadyDrawn += len(line)
			lineToDraw = line
		} else {
			var copy = alreadyDrawn
			for j := 0; j < int(currTextLen)-copy; j++ {
				alreadyDrawn++
				lineToDraw += string(line[j])
			}
		}
		this.linesToDraw = append(this.linesToDraw, lineToDraw)
	}
	if this.anim.IsFinished() {
		this.endAnim.Process(delta)
	}
}

func (this *TypewriterAnim) Draw() {
	this.ctx.ClearBG(glm.Vec4{0, 0, 0, 0})
	for i, x := range this.linesToDraw {
		this.ctx.Text.DrawText(0, 100+i*50, x, 1)
	}
}

func (this *TypewriterAnim) IsFinished() bool {
	return this.endAnim.IsFinished()
}
