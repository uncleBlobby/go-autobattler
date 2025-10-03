package main

import (
	_ "net/http/pprof"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/uncleblobby/raylib-go-test/internal/game"
)

func main() {

	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	rl.InitAudioDevice()
	rl.SetMasterVolume(1.0)

	rl.InitWindow(1920, 1080, "go atb")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	game := game.InitGame()

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()

		game.Update(dt)
		game.Draw(dt)
	}
}
