package main

import (
	"fmt"
	"log"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// type Game struct {
// 	player      *Player
// 	enemies     []*Enemy
// 	projectiles []*Projectile
// }

var projectiles = []*Projectile{}
var enemies = []*Enemy{}
var fcts = []*FloatingText{}
var loot = []*Loot{}
var enemySpawnTimer float32 = 0
var enemySpawnCooldown float32 = 2

type ProjectileOwnership int

const (
	PLAYER ProjectileOwnership = iota
	ENEMY
)

type Projectile struct {
	position       rl.Vector2
	direction      rl.Vector2
	targetPosition rl.Vector2
	speed          float32
	radius         float32
	isDead         bool
	owner          ProjectileOwnership
}

func (p *Projectile) Update(dt float32) {
	p.direction = rl.Vector2Normalize(p.direction)
	p.position.X += (p.speed * p.direction.X) * dt
	p.position.Y += (p.speed * p.direction.Y) * dt

	if p.position.Y < 0 || p.position.X < 0 || p.position.Y > 1080 || p.position.X > 1920 {
		p.isDead = true
	}
}

func (p *Projectile) Draw() {
	rl.DrawCircle(int32(p.position.X), int32(p.position.Y), 3, rl.Black)
}

type Player struct {
	Actor
	shootCooldown     float32
	timeSinceLastShot float32
	size              int32
	hitpoints         int32
	exp               int32
	expToLevel        int32
	expToPrevLevel    int32
	level             int32
}

type Enemy struct {
	Actor
	player            *Player
	shootCooldown     float32
	timeSinceLastShot float32
	size              int32
	hitpoints         int32
	// game              *Game
}

func (p *Player) Shoot(enemy *Enemy) {

	if enemy == nil {
		return
	}

	dirToTarget := rl.Vector2Subtract(enemy.position, p.position)

	proj := Projectile{
		position:  p.center,
		direction: dirToTarget,
		owner:     PLAYER,
		radius:    3,
		speed:     250,
	}

	projectiles = append(projectiles, &proj)
}

func (e *Enemy) Shoot() {
	//log.Printf("enemy shooting...")
	if e.player == nil {
		log.Printf("player is nil, aborting shoot")
		return
	}

	dirToPlayer := rl.Vector2Subtract(e.player.position, e.position)

	p := Projectile{
		position:  e.center,
		direction: dirToPlayer,
		owner:     ENEMY,
		radius:    3,
		speed:     250,
	}

	projectiles = append(projectiles, &p)
}

type Actor struct {
	position  rl.Vector2
	center    rl.Vector2
	direction rl.Vector2
	collider  rl.Rectangle
	speed     float32
	isDead    bool
}

func (p *Player) Draw() {
	rl.DrawRectangle(int32(p.position.X), int32(p.position.Y), p.size, p.size, rl.Blue)

	// exp bar

	rl.DrawRectangle(10, 1000, 1900, 15, rl.Blue)
	rl.DrawRectangle(15, 1002, 1890, 12, rl.RayWhite)

	// exp bar fill

	// fill percentage == (player.exp / player.expToLevel) * 100

	fillPerc := ((float32(p.exp) - float32(p.expToPrevLevel)) / (float32(p.expToLevel) - float32(p.expToPrevLevel)))

	rl.DrawRectangle(15, 1002, int32(1890*fillPerc), 12, rl.Blue)
}

func (e *Enemy) Draw() {

	//player body

	rl.DrawRectangle(int32(e.position.X), int32(e.position.Y), e.size, e.size, rl.Red)

}

func (p *Player) Update(dt float32) {
	p.center = rl.Vector2{X: p.position.X + float32(p.size/2), Y: p.position.Y + float32(p.size/2)}
	p.collider = rl.Rectangle{X: p.position.X, Y: p.position.Y, Width: float32(p.size), Height: float32(p.size)}

	p.direction = rl.Vector2{X: 0, Y: 0}
	if rl.IsKeyDown(rl.KeyA) {
		p.direction.X = -1
	}

	if rl.IsKeyDown(rl.KeyD) {
		p.direction.X = 1
	}

	if rl.IsKeyDown(rl.KeyW) {
		p.direction.Y = -1
	}

	if rl.IsKeyDown(rl.KeyS) {
		p.direction.Y = 1
	}

	if rl.IsKeyDown(rl.KeyLeftShift) {
		p.speed = 300
	} else {
		p.speed = 100
	}

	p.position.X += (p.speed * p.direction.X) * dt
	p.position.Y += (p.speed * p.direction.Y) * dt

	p.timeSinceLastShot += dt
	if p.timeSinceLastShot >= p.shootCooldown {

		if len(enemies) > 0 {

			// find closest enemy

			closest := enemies[0]
			closestDistance := rl.Vector2Distance(p.position, enemies[0].position)

			for i := 0; i < len(enemies); i++ {
				distanceTo := rl.Vector2Distance(p.position, enemies[i].position)

				if distanceTo < closestDistance {
					closest = enemies[i]
					closestDistance = distanceTo
				}
			}

			p.Shoot(closest)
		}

		p.timeSinceLastShot = 0
	}

	// check level up

	if p.exp > p.expToLevel {
		p.level += 1
		p.expToPrevLevel = p.expToLevel
		p.expToLevel = p.level * p.expToLevel
	}
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

	if e.timeSinceLastShot >= e.shootCooldown {
		e.Shoot()
		e.timeSinceLastShot = 0
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

func main() {
	rl.InitWindow(1920, 1080, "go atb")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	player := Player{
		Actor: Actor{
			position: rl.Vector2{X: 200, Y: 200},
			speed:    100,
		},
		size:           25,
		shootCooldown:  1,
		exp:            0,
		level:          1,
		expToLevel:     100,
		expToPrevLevel: 0,
	}

	enemy := Enemy{
		Actor: Actor{
			position: rl.Vector2{X: 500, Y: 500},
			speed:    75,
			center:   rl.Vector2{X: 500, Y: 500},
		},
		size:          25,
		player:        &player,
		shootCooldown: 0.01,
		hitpoints:     10,
	}

	enemies = append(enemies, &enemy)

	for !rl.WindowShouldClose() {

		dt := rl.GetFrameTime()

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawText(fmt.Sprintf("FPS:\t\t\t%d", rl.GetFPS()), 5, 5, 16, rl.Black)
		// rl.DrawText(fmt.Sprintf("ENTITIES:\t%d", ), 5, 5, 16, rl.Black)

		player.Update(dt)
		//enemy.Update(dt)

		// update all loot

		for i := 0; i < len(loot); i++ {
			if !loot[i].collected {
				if rl.CheckCollisionCircleRec(loot[i].position, loot[i].radius, player.collider) {
					loot[i].collected = true
					player.exp += int32(loot[i].xpValue)
					ft := CreateFloatingEXPText(player.position, loot[i].xpValue)
					fcts = append(fcts, &ft)
				}
			}
		}

		// update all projectiles

		for i := 0; i < len(projectiles); i++ {
			projectiles[i].Update(dt)

			if projectiles[i].owner == ENEMY {
				//collision against player, disappear bullet
				if rl.CheckCollisionCircleRec(projectiles[i].position, projectiles[i].radius, player.collider) {
					projectiles[i].isDead = true

					newFt := CreateFloatingText(player.position, 10)
					fcts = append(fcts, &newFt)
				}
			}

			for j := 0; j < len(enemies); j++ {
				if projectiles[i].owner == PLAYER {
					if rl.CheckCollisionCircleRec(projectiles[i].position, projectiles[i].radius, enemies[j].collider) {
						projectiles[i].isDead = true

						enemies[j].hitpoints -= 10

						if enemies[j].hitpoints <= 0 {
							enemies[j].isDead = true
							lt := CreateLootExpItem(enemies[j].center)
							loot = append(loot, &lt)
						}

						newFt := CreateFloatingText(enemies[j].position, 10)
						fcts = append(fcts, &newFt)
					}
				}
			}
		}

		for i := 0; i < len(enemies); i++ {
			enemies[i].Update(dt)
		}

		player.Draw()
		//enemy.Draw()

		for i := 0; i < len(projectiles); i++ {
			if !projectiles[i].isDead {
				projectiles[i].Draw()
			}

		}

		for i := 0; i < len(enemies); i++ {
			enemies[i].Draw()
		}

		// draw all the loots

		for i := 0; i < len(loot); i++ {
			if !loot[i].collected {
				rl.DrawCircleV(loot[i].position, loot[i].radius, loot[i].color)
			}
		}

		// draw floating texts

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

		rl.EndDrawing()

		// clean up dead projectiles

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

		// enemy spawn routine

		enemySpawnTimer += dt
		if enemySpawnTimer > enemySpawnCooldown {
			// spawn enemy
			log.Printf("SHOULD SPAWN NEW ENEMY...")
			enemies = append(enemies, SpawnEnemy(&player))

			// reset timer
			enemySpawnTimer = 0
		}
	}
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

func EaseOutQuad(t float32) float32 {
	return 1 - (1-t)*(1-t)
}

func LERP(f1, f2, amount float32) float32 {
	return f1 + amount*(f2-f1)
}

type Loot struct {
	position  rl.Vector2
	radius    float32
	xpValue   int
	collected bool
	color     rl.Color
}

func CreateLootExpItem(spawnPos rl.Vector2) Loot {
	gem := Loot{
		position:  spawnPos,
		radius:    8,
		xpValue:   10,
		color:     rl.Green,
		collected: false,
	}
	return gem
}
