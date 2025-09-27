package game

import rl "github.com/gen2brain/raylib-go/raylib"

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
	critChance     float32
}

func (p *Projectile) Update(dt float32) {
	p.direction = rl.Vector2Normalize(p.direction)
	p.position.X += (p.speed * p.direction.X) * dt
	p.position.Y += (p.speed * p.direction.Y) * dt

	if p.IsOffScreen() {
		p.isDead = true
	}
}

func (p *Projectile) Draw() {
	rl.DrawCircle(int32(p.position.X), int32(p.position.Y), p.radius, rl.Black)
}

func DrawAllProjectiles() {
	for i := 0; i < len(projectiles); i++ {
		if !projectiles[i].isDead {
			projectiles[i].Draw()
		}

	}
}

func (p *Projectile) IsOffScreen() bool {
	if p.position.Y < 0 || p.position.X < 0 || p.position.Y > 1080 || p.position.X > 1920 {
		return true
	}

	return false
}
