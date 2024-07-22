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
	var typewriter = newTypewriterAnim(getExeBasePath()+"/assets/intro.txt", ctx, 0.05)
	_ = typewriter
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
	_ = vertices
	var texts = []string{"1. Centroid", "2. Circumcenter", "3. Incenter", "4. Orthocenter"}
	var start = newIntroScreen(2, ctx)
	var sm = newStateMachine()
	_, _ = texts, start
	sm.addState(&typewriter)
	sm.addState(&start)

	for i := 0; i < len(vertices); i++ {
		if i%3 == 0 {
			var scrn = newTextScreen(texts[i/3], 2, 3, ctx)
			sm.addState(&scrn)

		}
		var tri = createTri(vertices[i], i/3, ctx, 3)
		sm.addState(&tri)
	}
	var grand = newTextScreen("The Grand Tour", 2, 3, ctx)
	sm.addState(&grand)

	return sm
}

func updateConfig() {
	var path = getExeBasePath() + "/assets/config.ini"
	var c, err = os.ReadFile(path)
	if err != nil {
		println("err", err.Error())
	}
	var new = strings.Replace(string(c), "showed_intro=false", "showed_intro=true", 1)
	var f, _ = os.Create(path)
	f.WriteString(new)
}
func main() {
	runtime.LockOSThread()
	var openGL = closedGL.InitClosedGL(800, 600, "Centre of a triangle")
	var config = openGL.Config

	var sm = buildSm(&openGL)
	_ = sm
	var tri = newTri([3]glm.Vec2{{100, 100}, {100, 500}, {700, 400}}, &openGL, rgbToColour(67, 61, 139))
	tri.drawCenters = false
	tri.showUi = false
	tri.startCentroidAnim(2)

	openGL.Window.SetMouseButtonCB(tri.mouseCB)
	openGL.Window.SetCursorPosCallback(tri.cursorCB)

	if config["bgm"] != "" && strToBool(config["bgm"]) {
		openGL.PlayMusic("bgm", 0.5)
	}

	closedGL.SetWireFrameMode(true)

	var cbs = [](func(tri *Tri, animDur float32)){
		func(tri *Tri, animDur float32) { tri.startCentroidAnim(animDur) },
		func(tri *Tri, animDur float32) { tri.startCircumCenterAnim(animDur) },
		func(tri *Tri, animDur float32) { tri.startIncenterAnim(animDur) },
		func(tri *Tri, animDur float32) { tri.startOrthocenterAnim(animDur) },
	}
	var newTargetVertices = [][3]glm.Vec2{
		{{100, 100}, {100, 500}, {677, 101}},
		{{692, 365}, {100, 500}, {677, 101}},
		{{692, 365}, {100, 500}, {181, 277}},
		{{692, 365}, {100, 500}, {423, 66}},
		{{577, 366}, {175, 388}, {357, 32}},
	}

	_ = cbs
	var grandPt1Finished bool
	var grandPt2Finished bool
	var smFinished bool
	if config["showed_intro"] != "" {
		var val = strToBool(config["showed_intro"])
		grandPt1Finished = val
		grandPt2Finished = val
		smFinished = val
		if val == true {
			tri.setVertices(newTargetVertices[len(newTargetVertices)-1])
			tri.drawCenters = true
			tri.showUi = true
		}
	}

	var currTargetVertices = 0

	var unitAnim = closedGL.NewAnimation(0, 1, 3, false, false)
	var origVertices = tri.convertPointsToSS()

	for !openGL.Window.Window.ShouldClose() {

		var delta = float32(openGL.FPSCounter.Delta)
		if !sm.isFinished() {
			sm.process(delta)
		} else {
			if !smFinished {
				openGL.EndMusic("bgm")
				openGL.PlayMusic("bgm2", 0.5)
			}
			smFinished = true
			if !grandPt2Finished && grandPt1Finished {
				unitAnim.Process(delta)

				tri.setVertices([3]glm.Vec2{
					LerpVec2(origVertices[0], newTargetVertices[currTargetVertices][0], unitAnim.GetValue()),
					LerpVec2(origVertices[1], newTargetVertices[currTargetVertices][1], unitAnim.GetValue()),
					LerpVec2(origVertices[2], newTargetVertices[currTargetVertices][2], unitAnim.GetValue()),
				})
				if unitAnim.IsFinished() {
					unitAnim = closedGL.NewAnimation(0, 1, 4, false, false)
					origVertices = tri.convertPointsToSS()
					currTargetVertices++
					if currTargetVertices >= len(newTargetVertices) {
						grandPt2Finished = true
						tri.showUi = true
						updateConfig()
					}
				}
			}
			tri.Process(delta)
		}
		if !grandPt1Finished {
			for i := 0; i < len(tri.currAnims)-1; i++ {
				if tri.currAnims[i] == nil {
					break
				}
				if tri.currAnims[i].IsFinished() && tri.currAnims[i+1] == nil {
					cbs[i+1](&tri, 0.4)
				}
			}
			grandPt1Finished = tri.currAnims[3] != nil && tri.currAnims[3].IsFinished()
			if grandPt1Finished {
				tri.currAnims = [4]TriState{}
				tri.drawCenters = true
			}
		}
		openGL.BeginDrawing()
		openGL.ClearBG(rgbToColour(23, 21, 59))
		if !smFinished {
			sm.draw()
		} else {
			tri.Draw()
		}
		openGL.DrawFPS(500, 0, 1)

		openGL.EndDrawing()
		openGL.Process()
	}
}
