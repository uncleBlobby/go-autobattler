package game

import rl "github.com/gen2brain/raylib-go/raylib"

type Enemy struct {
	Actor
	player            *Player
	shootCooldown     float32
	timeSinceLastShot float32
	size              int32
	hitpoints         int32
}

func (e *Enemy) Update(dt float32) {
	e.center = rl.Vector2{X: e.position.X + float32(e.size/2), Y: e.position.Y + float32(e.size/2)}
	e.collider = rl.Rectangle{X: e.position.X, Y: e.position.Y, Width: float32(e.size), Height: float32(e.size)}

	//e.direction = rl.Vector2{X: 0, Y: 0}

	// IF PLAYER IS NIL FOR SOME REASON, STOP
	if e.player == nil {
		return
	}

	// get direction to player IF PLAYER IS NOT NIL
	e.direction = rl.Vector2Subtract(e.player.position, e.position)
	e.direction = rl.Vector2Normalize(e.direction)

	e.position.X += (e.speed * e.direction.X) * dt
	e.position.Y += (e.speed * e.direction.Y) * dt

	e.timeSinceLastShot += dt

	// if e.timeSinceLastShot >= e.shootCooldown {
	// 	//e.Shoot()
	// 	e.timeSinceLastShot = 0
	// }
}

func (e *Enemy) Draw() {

	//player body

	rl.DrawRectangle(int32(e.position.X), int32(e.position.Y), e.size, e.size, rl.Red)

}

func UpdateAllEnemies(dt float32) {
	for i := 0; i < len(enemies); i++ {
		enemies[i].Update(dt)
	}
}

func DrawAllEnemies() {
	for i := 0; i < len(enemies); i++ {
		enemies[i].Draw()
	}
}

func SpawnEnemy(player *Player) *Enemy {
	enemy := Enemy{
		Actor: Actor{
			position: PositionJustOutsideScreen(),
			speed:    75,
		},
		size:          25,
		player:        player,
		shootCooldown: 2,
		hitpoints:     10,
	}

	return &enemy
}

func (e *Enemy) GetKnockedBack(e2 *Enemy, dt float32) {
	knockDirection := rl.Vector2Subtract(e.position, e2.position)
	knockDirection = rl.Vector2Normalize(knockDirection)

	knockForce := 150

	e.position.X += (float32(knockForce) * knockDirection.X) * dt
	e.position.Y += (float32(knockForce) * knockDirection.Y) * dt
}
