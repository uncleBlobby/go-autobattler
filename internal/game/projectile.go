package game

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ProjectileOwnership int

const (
	PLAYER ProjectileOwnership = iota
	ENEMY
)

type ProjectileKind int

const (
	DEFAULT ProjectileKind = iota
	BULLET
	MAGIC_MISSILE
	FIREBALL
	LIGHTNING
)

type Projectile struct {
	position  rl.Vector2
	direction rl.Vector2
	// targetPosition rl.Vector2
	damage     int
	speed      float32
	radius     float32
	isDead     bool
	owner      ProjectileOwnership
	critChance float32
	kind       ProjectileKind
	sprite     rl.Texture2D
	isAnim     bool
	anim       AnimatedSprite
}

func NewProjectile(k ProjectileKind) *Projectile {
	return &Projectile{
		kind: k,
	}
}

func (p *Projectile) Update(dt float32) {

	if p.isAnim {
		p.anim.frameCounter += dt
		if p.anim.frameCounter >= (1.0 / float32(p.anim.frameSpeed)) {
			p.anim.frameCounter = 0
			p.anim.currentFrame++
			if p.anim.currentFrame >= p.anim.numFrames {
				p.anim.currentFrame = 0
			}
		}
	}

	p.direction = rl.Vector2Normalize(p.direction)
	p.position.X += (p.speed * p.direction.X) * dt
	p.position.Y += (p.speed * p.direction.Y) * dt

	if p.IsOffScreen() {
		p.isDead = true
	}
}

func (p *Projectile) Clone() *Projectile {
	proj := Projectile{
		position:  p.position,
		direction: p.direction,
		// targetPosition rl.Vector2
		damage:     p.damage,
		speed:      p.speed,
		radius:     p.radius,
		isDead:     p.isDead,
		owner:      p.owner,
		critChance: p.critChance,
		kind:       p.kind,
		sprite:     p.sprite,
		isAnim:     p.isAnim,
		anim:       p.anim,
	}

	return &proj
}

func (p *Projectile) Draw() {

	switch p.kind {
	case MAGIC_MISSILE:

		src := rl.Rectangle{X: 0, Y: 0, Width: 16, Height: 16}
		dst := rl.Rectangle{X: p.position.X, Y: p.position.Y, Width: 32, Height: 32}
		og := rl.Vector2{X: 32 / 2, Y: 32 / 2}

		angle := math.Atan2(float64(p.direction.Y), float64(p.direction.X)) * rl.Rad2deg
		// rl.DrawTextureEx(p.sprite, p.position, float32(angle-90), 2, rl.White)
		rl.DrawTexturePro(p.sprite, src, dst, og, float32(angle-90), rl.White)
		rl.DrawCircle(int32(p.position.X), int32(p.position.Y), p.radius, rl.SkyBlue)

	case BULLET:
		rl.DrawCircle(int32(p.position.X), int32(p.position.Y), p.radius, rl.Black)

	case LIGHTNING:

		src := rl.Rectangle{X: 0, Y: float32(p.anim.currentFrame) * p.anim.spriteSize.Y, Width: p.anim.spriteSize.X, Height: p.anim.spriteSize.Y}
		dst := rl.Rectangle{X: p.position.X, Y: p.position.Y, Width: 32, Height: 32}
		og := rl.Vector2{X: dst.Width / 2, Y: dst.Height / 2}

		angle := math.Atan2(float64(p.direction.Y), float64(p.direction.X)) * rl.Rad2deg
		// rl.DrawTextureEx(p.sprite, p.position, float32(angle-90), 2, rl.White)
		rl.DrawTexturePro(p.sprite, src, dst, og, float32(angle-90), rl.White)

		rl.DrawCircle(int32(p.position.X), int32(p.position.Y), p.radius, rl.Black)

	case FIREBALL:
		src := rl.Rectangle{X: 0, Y: 0, Width: 16, Height: 16}
		dst := rl.Rectangle{X: p.position.X, Y: p.position.Y, Width: 32, Height: 32}
		og := rl.Vector2{X: 32 / 2, Y: 32 / 2}

		angle := math.Atan2(float64(p.direction.Y), float64(p.direction.X)) * rl.Rad2deg
		// rl.DrawTextureEx(p.sprite, p.position, float32(angle-90), 2, rl.White)
		rl.DrawTexturePro(p.sprite, src, dst, og, float32(angle-90), rl.White)
		rl.DrawCircle(int32(p.position.X), int32(p.position.Y), p.radius, rl.Red)

	default:
		rl.DrawCircle(int32(p.position.X), int32(p.position.Y), p.radius, rl.Black)
	}
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
