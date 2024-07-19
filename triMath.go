package main

import (
	"github.com/EngoEngine/glm"
	"github.com/EngoEngine/math"
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
	if p1[1] == p2[1] {
		p2 = glm.Vec2{p2[0], p2[1] + 0.1}
	}
	if p1[0] == p2[0] {
		p2 = glm.Vec2{p2[0] + 0.1, p2[1]}
	}

	var dy = p2[1] - p1[1]
	var dx = p2[0] - p1[0]
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
