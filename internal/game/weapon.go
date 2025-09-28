package game

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Cooldown struct {
	timeSinceShot float32
	duration      float32
}

type Weapon interface {
	Shoot(enemy *Enemy)
	Update(dt float32, position rl.Vector2)
	ReduceCooldownDuration(amount float32)
	AdjustBaseDamage(amount float32)
}

type BaseWeapon struct {
	position    rl.Vector2
	critChance  float32
	cooldown    Cooldown
	soundEffect rl.Sound
	baseDamage  int
}

func (b *BaseWeapon) Update(dt float32, position rl.Vector2) {
	b.position = position
	b.cooldown.timeSinceShot += dt
}

func (b *BaseWeapon) ReduceCooldownDuration(amount float32) {
	factor := 1 - amount
	b.cooldown.duration *= factor

	if b.cooldown.duration <= 0 {
		b.cooldown.duration = 0.1
	}
}

func (b *BaseWeapon) AdjustBaseDamage(amount float32) {
	b.baseDamage = b.baseDamage + int(float32(b.baseDamage)*amount)
}

func (s *Shotgun) AdjustBaseDamage(amount float32) {
	s.baseDamage = s.baseDamage + int(float32(s.baseDamage)*amount)
}

func (s *SMG) AdjustBaseDamage(amount float32) {
	s.baseDamage = s.baseDamage + int(float32(s.baseDamage)*amount)
}

func (b *BaseWeapon) Shoot(enemy *Enemy) {
	if enemy == nil {
		return
	}

	if b.cooldown.timeSinceShot >= b.cooldown.duration {
		dirToTarget := rl.Vector2Subtract(enemy.center, b.position)

		proj := Projectile{
			position:   b.position,
			direction:  dirToTarget,
			owner:      PLAYER,
			radius:     3,
			speed:      500,
			critChance: b.critChance,
			damage:     b.baseDamage,
		}

		projectiles = append(projectiles, &proj)

		rl.PlaySound(b.soundEffect)

		b.cooldown.timeSinceShot = 0
	}
}

type Shotgun struct {
	position       rl.Vector2
	critChance     float32
	cooldown       Cooldown
	numProjectiles int
	soundEffect    rl.Sound
	baseDamage     int
}

type SMG struct {
	position            rl.Vector2
	critChance          float32
	cooldown            Cooldown
	magazineSize        int
	shotsFiredThisBurst int
	rateOfFire          float32
	timeSinceLastRound  float32
	soundEffect         rl.Sound
	baseDamage          int
}

func (s *Shotgun) Update(dt float32, position rl.Vector2) {
	s.position = position
	s.cooldown.timeSinceShot += dt
}

func (s *Shotgun) ReduceCooldownDuration(amount float32) {
	factor := 1 - amount
	s.cooldown.duration *= factor

	if s.cooldown.duration <= 0 {
		s.cooldown.duration = 0.1
	}
}

func (s *Shotgun) Shoot(enemy *Enemy) {
	if enemy == nil {
		return
	}

	if s.cooldown.timeSinceShot >= s.cooldown.duration {

		spreadAngle := rl.Deg2rad * 25

		dirToTarget := rl.Vector2Subtract(enemy.center, s.position)

		for i := 0; i < s.numProjectiles; i++ {
			t := (float32(i) / (float32(s.numProjectiles) - 1))
			jitter := (rand.Float32() - 0.5) * (5 * rl.Deg2rad)

			baseAngle := (t - 0.5) * float32(spreadAngle)

			angle := baseAngle + jitter

			projDir := rl.Vector2Rotate(dirToTarget, angle)

			proj := Projectile{
				position:   s.position,
				direction:  projDir,
				owner:      PLAYER,
				radius:     3,
				speed:      500,
				critChance: s.critChance,
				damage:     s.baseDamage,
			}

			projectiles = append(projectiles, &proj)
		}

		rl.PlaySound(s.soundEffect)
		s.cooldown.timeSinceShot = 0
	}
}

func (s *SMG) Update(dt float32, position rl.Vector2) {
	s.position = position
	s.cooldown.timeSinceShot += dt
	s.timeSinceLastRound += dt
}

func (s *SMG) ReduceCooldownDuration(amount float32) {
	factor := 1 - amount
	s.cooldown.duration *= factor

	if s.cooldown.duration <= 0 {
		s.cooldown.duration = 0.1
	}
}

func (s *SMG) Shoot(enemy *Enemy) {
	if enemy == nil {
		return
	}

	if s.cooldown.timeSinceShot >= s.cooldown.duration {

		dirToTarget := rl.Vector2Subtract(enemy.center, s.position)

		if s.shotsFiredThisBurst <= s.magazineSize && s.timeSinceLastRound > 0.05 {

			proj := Projectile{
				position:   s.position,
				direction:  dirToTarget,
				owner:      PLAYER,
				radius:     3,
				speed:      500,
				critChance: s.critChance,
				damage:     s.baseDamage,
			}

			projectiles = append(projectiles, &proj)

			if !rl.IsSoundPlaying(s.soundEffect) {
				rl.PlaySound(s.soundEffect)
			}

			s.shotsFiredThisBurst += 1
			s.timeSinceLastRound = 0
		}

		if s.shotsFiredThisBurst > s.magazineSize {
			s.shotsFiredThisBurst = 0
			s.cooldown.timeSinceShot = 0
		}
	}
}
