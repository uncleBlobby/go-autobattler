package game

import rl "github.com/gen2brain/raylib-go/raylib"

type Actor struct {
	position  rl.Vector2
	center    rl.Vector2
	direction rl.Vector2
	collider  rl.Rectangle
	speed     float32
	isDead    bool
}
