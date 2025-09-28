package game

import (
	"fmt"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var projectiles = []*Projectile{}
var enemies = []*Enemy{}
var fcts = []*FloatingText{}
var loot = []*Loot{}
var magnets = []*Magnet{}
var enemySpawnTimer float32 = 0
var enemySpawnCooldown float32 = 2
var gameTime float32 = 0
var player *Player

var showLevelUpScreen bool = false
var gameOver bool = false

type Game struct{}

func InitGame() *Game {
	levelUpSfx := rl.LoadSound("assets/sounds/level-up-sequence.mp3")
	xpPickupFx := rl.LoadSound("assets/sounds/retro-coin.mp3")

	rl.SetSoundVolume(levelUpSfx, 0.9)
	rl.SetSoundVolume(xpPickupFx, 0.4)

	player = &Player{
		Actor: Actor{
			position: rl.Vector2{X: 200, Y: 200},
			speed:    300,
		},
		size:           25,
		shootCooldown:  1,
		exp:            0,
		level:          1,
		expToLevel:     100,
		expToPrevLevel: 0,
		hitpoints:      100,
		pickupRadius:   75,
		levelUpSfx:     levelUpSfx,
		expPickupSfx:   xpPickupFx,
	}

	player.InitBaseWeapon()
	player.InitShotgunWeapon()
	player.InitSMGWeapon()

	InitMagicMissile(player)

	return &Game{}
}

func (g *Game) Draw(dt float32) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	rl.DrawText(fmt.Sprintf("FPS:\t\t\t%d", rl.GetFPS()), 5, 5, 16, rl.Black)
	rl.DrawText(fmt.Sprintf("PLAYER_HP:\t%d", player.hitpoints), 5, 25, 16, rl.Black)

	DrawAllLoot()

	player.Draw()

	DrawAllProjectiles()

	DrawAllEnemies()

	DrawFloatingTexts(dt)

	if showLevelUpScreen {
		DrawLevelUpScreen()
	}

	if gameOver {
		DrawGameOverScreen()
	}

	rl.EndDrawing()
}

func (g *Game) Update(dt float32) {

	gameTime += dt

	wave := int(gameTime / 30)
	enemySpawnCooldown = 2.0 - (float32(wave) * 0.2)
	enemySpawnCooldown = float32(math.Max(float64(enemySpawnCooldown), 0.2))

	if !showLevelUpScreen && !gameOver {
		player.Update(dt)

		UpdateAllLoot(dt)

		EnemyCollisionLoop(dt)

		ProjectileCollisionLoop(dt)

		UpdateAllEnemies(dt)
	}

	// clean up

	CleanUpDeadEntities()

	// enemy spawn routine

	enemySpawnTimer += dt
	if enemySpawnTimer > enemySpawnCooldown {
		// spawn enemy
		// log.Printf("SHOULD SPAWN NEW ENEMY...")

		for range 5 {
			enemies = append(enemies, SpawnEnemy(player))
		}

		// reset timer
		enemySpawnTimer = 0
	}
}

func CleanUpDeadEntities() {
	for i := 0; i < len(projectiles); i++ {
		if projectiles[i].isDead {
			projectiles = append(projectiles[:i], projectiles[i+1:]...)
		}
	}

	// clean up dead enemies

	for i := 0; i < len(enemies); i++ {
		if enemies[i].isDead {
			enemies = append(enemies[:i], enemies[i+1:]...)
		}
	}

	// clean up collected loot

	for i := 0; i < len(loot); i++ {
		if loot[i].collected {
			loot = append(loot[:i], loot[i+1:]...)
		}
	}

	for i := 0; i < len(magnets); i++ {
		if magnets[i].collected {
			magnets = append(magnets[:i], magnets[i+1:]...)
		}
	}
}

func EaseOutQuad(t float32) float32 {
	return 1 - (1-t)*(1-t)
}

func LERP(f1, f2, amount float32) float32 {
	return f1 + amount*(f2-f1)
}

func PositionJustOutsideScreen() rl.Vector2 {
	const width = 1920
	const height = 1080
	const perim = 2 * (width + height)

	r := rand.Intn(perim)

	if r > 0 && r < width {
		return rl.Vector2{X: float32(width - r), Y: height}
	}

	if r > width && r < (width+height) {
		return rl.Vector2{X: 0, Y: float32(r - width)}
	}

	if r > (width+height) && r < (2*width+height) {
		return rl.Vector2{X: float32(r - (width + height)), Y: 0}
	}

	if r > (2*width+height) && r < (2*width+2*height) {
		return rl.Vector2{X: float32(width), Y: float32(height - (r - (2*width + height)))}
	}

	return rl.Vector2{X: -1, Y: -1}
}

func DrawLevelUpScreen() {

	screenColor := rl.LightGray
	screenColor.A = uint8(128)
	rl.DrawRectangle(50, 50, 1920-100, 1080-100, screenColor)

	rl.DrawText("LEVEL UP!!!!", 800, 500, 60, rl.Black)
	rl.DrawText("pick some new skills or something...", 820, 600, 20, rl.Black)

	// rl.DrawRectangle(750, 900, 100, 100, rl.Orange)

	if rl.IsKeyReleased(rl.KeyBackspace) {
		showLevelUpScreen = false
	}
}

func DrawGameOverScreen() {
	screenColor := rl.LightGray
	screenColor.A = uint8(128)
	rl.DrawRectangle(50, 50, 1920-100, 1080-100, screenColor)

	rl.DrawText("GAME OVER!!!!", 800, 500, 60, rl.Black)
}

func UNUSED(a any) {
	return
}
