package game

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Buff struct {
	duration             float32
	timeSinceStart       float32
	pickupRadiusBonus    float32
	pickupRadiusOriginal float32
}

type Player struct {
	Actor
	shootCooldown float32
	// timeSinceLastShot float32
	size           int32
	hitpoints      int32
	exp            int32
	expToLevel     int32
	expToPrevLevel int32
	level          int32
	pickupRadius   float32
	weapons        []Weapon
	spells         []Spell
	buffs          []Buff
	levelUpSfx     rl.Sound
	expPickupSfx   rl.Sound
}

func (p *Player) InitBaseWeapon() {

	sfx := rl.LoadSound("assets/sounds/fixed.wav")

	rl.SetSoundVolume(sfx, 0.7)

	b := BaseWeapon{position: p.center,
		critChance: 0.5,
		baseDamage: 10,
		cooldown: Cooldown{
			timeSinceShot: 0,
			duration:      2,
		},
		soundEffect: sfx}
	p.weapons = append(p.weapons, &b)
}

func (p *Player) InitShotgunWeapon() {

	sfx := rl.LoadSound("assets/sounds/shotgun-firing.mp3")
	rl.SetSoundVolume(sfx, 0.5)

	s := Shotgun{position: p.center, critChance: 0.25, baseDamage: 5, numProjectiles: 3, cooldown: Cooldown{
		timeSinceShot: 0,
		duration:      5,
	},
		soundEffect: sfx}
	p.weapons = append(p.weapons, &s)
}

func (p *Player) InitSMGWeapon() {

	sfx := rl.LoadSound("assets/sounds/smg-firing.mp3")
	rl.SetSoundVolume(sfx, 0.5)

	s := SMG{position: p.center, critChance: 0.1, baseDamage: 2, cooldown: Cooldown{
		timeSinceShot: 0,
		duration:      3,
	}, magazineSize: 10, shotsFiredThisBurst: 0, rateOfFire: 0.01, timeSinceLastRound: 0, soundEffect: sfx}
	p.weapons = append(p.weapons, &s)
}

func (p *Player) Shoot(enemy *Enemy) {

	if dbgf.enableWeapons {
		for i := 0; i < len(p.weapons); i++ {
			p.weapons[i].Shoot(enemy)
		}

	}

	if dbgf.enableSpells {
		for i := 0; i < len(p.spells); i++ {
			p.spells[i].Shoot(p.center, enemy)
		}
	}

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
		p.speed = 500
	} else {
		p.speed = 300
	}

	p.position.X += (p.speed * p.direction.X) * dt
	p.position.Y += (p.speed * p.direction.Y) * dt

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
		//fmt.Println(closest)
		p.Shoot(closest)
		// UNUSED(closest)
	}

	// p.timeSinceLastShot = 0

	// check level up

	p.LevelUpRoutine()

	// update all weapons

	for i := 0; i < len(p.weapons); i++ {
		p.weapons[i].Update(dt, p.position)
	}

	// update all spells

	for i := 0; i < len(p.spells); i++ {
		p.spells[i].Update(dt)
	}

	// update buffs (TODO: fix the hack below)

	for i := 0; i < len(p.buffs); i++ {
		p.buffs[i].timeSinceStart += dt
		p.pickupRadius = p.buffs[i].pickupRadiusBonus

		if p.buffs[i].timeSinceStart > p.buffs[i].duration {
			p.pickupRadius = p.buffs[i].pickupRadiusOriginal
			p.buffs = append(p.buffs[:i], p.buffs[i+1:]...)
		}
	}

	if p.hitpoints <= 0 {
		gameOver = true
	}
}

func (p *Player) LevelUpRoutine() {
	if p.exp > p.expToLevel {
		p.level += 1
		p.expToPrevLevel = p.expToLevel
		p.expToLevel = p.level * p.expToLevel
		p.shootCooldown = 1 - (float32(p.level-1) * 0.15)
		//p.shootCooldown = float32(math.Min(float64(p.shootCooldown)))
		p.shootCooldown = float32(math.Max(float64(p.shootCooldown), 0.1))
		p.pickupRadius *= 1.1

		for i := 0; i < len(p.weapons); i++ {
			p.weapons[i].ReduceCooldownDuration(0.05)
			p.weapons[i].AdjustBaseDamage(0.5)
		}

		showLevelUpScreen = true
		rl.PlaySound(p.levelUpSfx)
	}
}

func (p *Player) Draw() {
	rl.DrawRectangle(int32(p.position.X), int32(p.position.Y), p.size, p.size, rl.Blue)

	// health bar

	rl.DrawRectangle(10, 900, 20, 100, rl.Red)

	healthFillPerc := float32((float32(p.hitpoints) / float32(100)))

	redHeight := healthFillPerc * 100
	whiteHeight := 100 - redHeight

	UNUSED(healthFillPerc)

	rl.DrawRectangle(12, 902, 16, int32(whiteHeight), rl.RayWhite)

	// exp bar

	rl.DrawRectangle(10, 1000, 1900, 15, rl.Blue)
	rl.DrawRectangle(15, 1002, 1890, 12, rl.RayWhite)

	// exp bar fill

	// fill percentage == (player.exp / player.expToLevel) * 100

	fillPerc := ((float32(p.exp) - float32(p.expToPrevLevel)) / (float32(p.expToLevel) - float32(p.expToPrevLevel)))

	rl.DrawRectangle(15, 1002, int32(1890*fillPerc), 12, rl.Blue)
	rl.DrawText(fmt.Sprintf("Level: %d", p.level), (1920/2)-50, 1004, 12, rl.Black)
}

func (p *Player) GetKnockedBack(e *Enemy, dt float32) {

	knockDirection := rl.Vector2Subtract(p.position, e.position)
	knockDirection = rl.Vector2Normalize(knockDirection)

	knockForce := 150

	p.position.X += (float32(knockForce) * knockDirection.X) * dt
	p.position.Y += (float32(knockForce) * knockDirection.Y) * dt
}
