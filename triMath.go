package main

import (
	"github.com/EngoEngine/glm"
	"github.com/EngoEngine/math"
	"github.com/Zerou02/closedGL/closedGL"
)

func SSToCartesianPoint(vec glm.Vec2, wh float32) glm.Vec2 {
	return glm.Vec2{vec[0], wh - vec[1]}
}

func CartesianToSSPoint(vec glm.Vec2, wh float32) glm.Vec2 {
	return glm.Vec2{vec[0], wh - vec[1]}
}

// Result in rad
func AngleTo(vec, to glm.Vec2) float32 {
	var dot = vec.Dot(&to)
	var lenVec = vec.Len()
	var lenTo = to.Len()
	return math.Acos(dot / (lenVec * lenTo))
}

func DegToRad(deg float32) float32 {
	return deg * (math.Pi / 180.0)
}

func RadToDeg(rad float32) float32 {
	return rad * (180.0 / math.Pi)
}

// ccw
func Rotate(rad float32, vec glm.Vec2) glm.Vec2 {
	var rotMat = glm.Rotate2D(rad)
	return rotMat.Mul2x1(&vec)
}

func CalcMiddlePoint(p1, p2 glm.Vec2) glm.Vec2 {
	return glm.Vec2{(p1[0] + p2[0]) / 2, (p1[1] + p2[1]) / 2}
}

// m,n
func CalcLinearEquation(p1, p2 glm.Vec2) glm.Vec2 {
	var offsetP2 = p2
	if p1[1] == p2[1] {
		offsetP2 = glm.Vec2{p2[0], p2[1] + 0.1}
	}
	if glm.FloatEqualThreshold(p1[0], p2[0], 0.01) {
		offsetP2 = glm.Vec2{p2[0] + 0.1, p2[1]}
	}

	var dy = offsetP2[1] - p1[1]
	var dx = offsetP2[0] - p1[0]
	var m = dy / dx
	var n = p1[1] - m*p1[0]
	return glm.Vec2{m, n}
}

// x,y
func CalcCrossingPoint(linEq1, linEq2 glm.Vec2) glm.Vec2 {
	var x = (linEq1[1] - linEq2[1]) / (linEq2[0] - linEq1[0])
	return glm.Vec2{x, linEq1[0]*x + linEq1[1]}
}

func distBetweenPoints(p1, p2 glm.Vec2) float32 {
	var diff = p1.Sub(&p2)
	return math.Abs(diff.Len())
}

func distToLine(lineP1, lineP2, p glm.Vec2) float32 {
	var num = math.Abs((lineP2[1]-lineP1[1])*p[0] - (lineP2[0]-lineP1[0])*p[1] + lineP2[0]*lineP1[1] - lineP2[1]*lineP1[0])
	var denom = distBetweenPoints(lineP1, lineP2)
	return num / denom
}

// ax+by+c = 0
func ConvertToStandardForm(eq glm.Vec2) glm.Vec3 {
	return glm.Vec3{-eq[0], 1, -eq[1]}
}

// https://cp-algorithms.com/geometry/circle-line-intersection.html
func LineCircleIntersection(r float32, eq glm.Vec2, circlePos glm.Vec2) []glm.Vec2 {
	var retArr []glm.Vec2 = []glm.Vec2{}

	//translation
	var normEq = glm.Vec2{eq[0], eq[1] + circlePos[0]*eq[0] - circlePos[1]}
	var stdForm = ConvertToStandardForm(normEq)
	var a = stdForm[0]
	var b = stdForm[1]
	var c = stdForm[2]
	var x0 = -a * c / (a*a + b*b)
	var y0 = -b * c / (a*a + b*b)
	if c*c > r*r*(a*a+b*b)+glm.Epsilon {
		//do nothing
	} else if math.Abs(c*c-r*r*(a*a+b*b)) < glm.Epsilon {
		retArr = append(retArr, glm.Vec2{x0, y0})
	} else {
		var d = r*r - c*c/(a*a+b*b)
		var mult = math.Sqrt(d / (a*a + b*b))
		var ax = x0 + b*mult
		var bx = x0 - b*mult
		var ay = y0 - a*mult
		var by = y0 + a*mult
		retArr = append(retArr, glm.Vec2{ax, ay})
		retArr = append(retArr, glm.Vec2{bx, by})
	}

	for i := 0; i < len(retArr); i++ {
		retArr[i] = retArr[i].Add(&circlePos)
	}
	return retArr
}

func pointInRect(p glm.Vec2, rect glm.Vec4) bool {
	return p[0] >= rect[0] && p[1] >= rect[1] && p[0] < rect[0]+rect[2] && p[1] < rect[1]+rect[3]
}

// works if and only if there is exactly one
func findLineCircleIntersectionPoint(cp, p1, p2 glm.Vec2) glm.Vec2 {
	var r float32 = 0
	var step float32 = 1
	var targetPoint = glm.Vec2{0, 0}

	var oppositeSide = CalcLinearEquation(p1, p2)
	var currOffsets = LineCircleIntersection(r, oppositeSide, cp)
	for distToLine(p1, p2, targetPoint) > glm.Epsilon {
		currOffsets = LineCircleIntersection(r, oppositeSide, cp)
		var len = len(currOffsets)
		if len == 2 {
			if step > glm.Epsilon {
				r -= step
			} else {
				r -= glm.Epsilon
			}
			step *= 0.1
			targetPoint = currOffsets[0]
			if step < glm.Epsilon {
				break
			}
		}
		if len == 0 {
			if step < glm.Epsilon {
				r += glm.Epsilon
			} else {
				var rNew = step + r
				if rNew == r {
					break
				}
				r = rNew
			}
		}
		if len == 1 {
			targetPoint = currOffsets[0]
			break
		}
	}
	if len(currOffsets) == 2 {
		targetPoint = closedGL.MiddlePoint(currOffsets[0], currOffsets[1])
	}
	return targetPoint
}
