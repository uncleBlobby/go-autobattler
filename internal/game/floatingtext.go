package game

import (
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type FloatingText struct {
	basePosition  rl.Vector2
	text          string
	lifetime      float32
	totalLifetime float32
	startSize     float32
	endSize       float32
	velocity      rl.Vector2
	color         rl.Color
}

func CreateFloatingText(startPos rl.Vector2, damage int) FloatingText {
	ft := FloatingText{
		basePosition:  rl.Vector2{X: startPos.X, Y: startPos.Y - 25},
		text:          fmt.Sprintf("%d", damage),
		lifetime:      1,
		totalLifetime: 1,
		startSize:     40,
		endSize:       10,
		color:         rl.Black,
		velocity:      rl.Vector2{X: 0, Y: -40},
	}

	// randX := rand.Intn(int(startPos.X)-10, int(startPos.X)+10)
	// randX := rand.Int()
	randX := rand.Intn(25) - 25

	ft.basePosition.X += float32(randX)

	return ft
}

func CreateFloatingCRITText(startPos rl.Vector2) FloatingText {
	ft := FloatingText{
		basePosition:  rl.Vector2{X: startPos.X, Y: startPos.Y - 25},
		text:          ("CRIT"),
		lifetime:      1,
		totalLifetime: 1,
		startSize:     80,
		endSize:       10,
		color:         rl.Black,
		velocity:      rl.Vector2{X: 0, Y: -40},
	}

	ft.color.A = uint8(126)

	// randX := rand.Intn(int(startPos.X)-10, int(startPos.X)+10)
	// randX := rand.Int()
	randX := rand.Intn(25) - 25

	ft.basePosition.X += float32(randX)

	return ft
}

func CreateFloatingEXPText(startPos rl.Vector2, amount int) FloatingText {
	ft := FloatingText{
		basePosition:  rl.Vector2{X: startPos.X, Y: startPos.Y - 25},
		text:          fmt.Sprintf("+%d EXP", amount),
		lifetime:      1,
		totalLifetime: 1,
		startSize:     40,
		endSize:       10,
		color:         rl.Green,
		velocity:      rl.Vector2{X: 0, Y: -40},
	}

	// randX := rand.Intn(int(startPos.X)-10, int(startPos.X)+10)
	// randX := rand.Int()
	randX := rand.Intn(25) - 25

	ft.basePosition.X += float32(randX)

	return ft
}

func DrawFloatingTexts(dt float32) {
	for i := 0; i < len(fcts); i++ {
		fcts[i].lifetime -= dt
		if fcts[i].lifetime < 0 {
			fcts[i].lifetime = 0
		}
		t := 1 - (fcts[i].lifetime / fcts[i].totalLifetime)
		t = rl.Clamp(t, 0, 1)
		ease := EaseOutQuad(t)
		origin := fcts[i].basePosition
		offset := rl.Vector2{X: fcts[i].velocity.X * ease, Y: fcts[i].velocity.Y * ease}
		currentPos := rl.Vector2Add(origin, offset)
		fontSize := LERP(fcts[i].startSize, fcts[i].endSize, ease)
		c := fcts[i].color
		c.A = uint8(255 * (fcts[i].lifetime / fcts[i].totalLifetime))
		rl.DrawText(fcts[i].text, int32(currentPos.X), int32(currentPos.Y), int32(fontSize), c)
	}
}
