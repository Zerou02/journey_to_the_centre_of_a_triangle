package journey

import "github.com/EngoEngine/glm"

func SSToCartesianVec(vec glm.Vec2, wh float32) glm.Vec2 {
	return glm.Vec2{vec[0], wh - vec[1]}
}
