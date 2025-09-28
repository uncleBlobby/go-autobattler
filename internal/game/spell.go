package game

import rl "github.com/gen2brain/raylib-go/raylib"

type Spell interface {
	Shoot(startPos rl.Vector2, enemy *Enemy)
	Update(dt float32)
}

type ChainLightning struct {
	spellLevel       int
	critChance       float32
	cooldown         Cooldown
	soundEffect      rl.Sound
	textName         string
	baseDamage       int
	owner            ProjectileOwnership
	projectileRadius float32
	projectileSpeed  float32
	sprite           rl.Texture2D
}

type MagicMissile struct {
	// position    rl.Vector2
	spellLevel       int
	critChance       float32
	cooldown         Cooldown
	soundEffect      rl.Sound
	textName         string
	baseDamage       int
	owner            ProjectileOwnership
	projectileRadius float32
	projectileSpeed  float32
	sprite           rl.Texture2D
}

type Fireball struct {
	spellLevel       int
	critChance       float32
	cooldown         Cooldown
	soundEffect      rl.Sound
	textName         string
	baseDamage       int
	owner            ProjectileOwnership
	projectileRadius float32
	projectileSpeed  float32
	sprite           rl.Texture2D
}

func InitMagicMissile(p *Player) {

	sfx := rl.LoadSound("assets/sounds/magic-1.wav")
	rl.SetSoundVolume(sfx, 0.7)

	mmSprite := rl.LoadTexture("assets/sprites/magic-missile.png")

	mm := MagicMissile{
		spellLevel: 1,
		critChance: 0.1,
		cooldown: Cooldown{
			timeSinceShot: 0,
			duration:      2,
		},
		soundEffect:      sfx,
		textName:         "Magic Missile",
		baseDamage:       10,
		owner:            PLAYER,
		projectileRadius: 3,
		projectileSpeed:  500,
		sprite:           mmSprite,
	}

	p.spells = append(p.spells, &mm)
}

func InitChainLightning(p *Player) {
	sfx := rl.LoadSound("assets/sounds/magic-1.wav")
	rl.SetSoundVolume(sfx, 0.7)

	//mmSprite := rl.LoadTexture("assets/sprites/magic-missile.png")

	mm := ChainLightning{
		spellLevel: 1,
		critChance: 0.1,
		cooldown: Cooldown{
			timeSinceShot: 0,
			duration:      2,
		},
		soundEffect:      sfx,
		textName:         "Chain Lightning",
		baseDamage:       10,
		owner:            PLAYER,
		projectileRadius: 3,
		projectileSpeed:  500,
		sprite:           lightningSpriteSheet,
	}

	p.spells = append(p.spells, &mm)
}

func InitFireball(p *Player) {
	sfx := rl.LoadSound("assets/sounds/magic-1.wav")
	rl.SetSoundVolume(sfx, 0.7)

	mmSprite := rl.LoadTexture("assets/sprites/Fireball.png")

	mm := Fireball{
		spellLevel: 1,
		critChance: 0.1,
		cooldown: Cooldown{
			timeSinceShot: 0,
			duration:      2.5,
		},
		soundEffect:      sfx,
		textName:         "Fireball",
		baseDamage:       10,
		owner:            PLAYER,
		projectileRadius: 3,
		projectileSpeed:  500,
		sprite:           mmSprite,
	}

	p.spells = append(p.spells, &mm)
}

func (m *MagicMissile) Update(dt float32) {
	m.cooldown.timeSinceShot += dt
}

func (f *Fireball) Update(dt float32) {
	f.cooldown.timeSinceShot += dt
}

func (c *ChainLightning) Update(dt float32) {
	c.cooldown.timeSinceShot += dt
}

func (c *ChainLightning) Shoot(startPos rl.Vector2, enemy *Enemy) {
	if enemy == nil {
		return
	}

	if c.cooldown.timeSinceShot >= c.cooldown.duration {
		// pass direction to shoot as param instead of start / end pos?
		dirToTarget := rl.Vector2Subtract(enemy.center, startPos)

		anim := AnimatedSprite{
			sprite:       c.sprite,
			spriteSize:   rl.Vector2{X: 16, Y: 16},
			numFrames:    8,
			frameSpeed:   8,
			currentFrame: 0,
			frameCounter: 0,
		}

		proj := Projectile{
			position:   startPos,
			direction:  dirToTarget,
			owner:      c.owner,
			radius:     c.projectileRadius,
			speed:      c.projectileSpeed,
			critChance: c.critChance,
			kind:       LIGHTNING,
			sprite:     c.sprite,
			damage:     c.baseDamage,
			isAnim:     true,
			anim:       anim,
		}

		projectiles = append(projectiles, &proj)
		rl.PlaySound(c.soundEffect)
		c.cooldown.timeSinceShot = 0
	}
}

func (f *Fireball) Shoot(startPos rl.Vector2, enemy *Enemy) {
	if enemy == nil {
		return
	}

	if f.cooldown.timeSinceShot >= f.cooldown.duration {
		// pass direction to shoot as param instead of start / end pos?
		dirToTarget := rl.Vector2Subtract(enemy.center, startPos)

		proj := Projectile{
			position:   startPos,
			direction:  dirToTarget,
			owner:      f.owner,
			radius:     f.projectileRadius,
			speed:      f.projectileSpeed,
			critChance: f.critChance,
			kind:       FIREBALL,
			sprite:     f.sprite,
			damage:     f.baseDamage,
		}

		projectiles = append(projectiles, &proj)
		rl.PlaySound(f.soundEffect)
		f.cooldown.timeSinceShot = 0
	}
}

func (m *MagicMissile) Shoot(startPos rl.Vector2, enemy *Enemy) {
	if enemy == nil {
		return
	}

	if m.cooldown.timeSinceShot >= m.cooldown.duration {
		// pass direction to shoot as param instead of start / end pos?
		dirToTarget := rl.Vector2Subtract(enemy.center, startPos)

		proj := Projectile{
			position:   startPos,
			direction:  dirToTarget,
			owner:      m.owner,
			radius:     m.projectileRadius,
			speed:      m.projectileSpeed,
			critChance: m.critChance,
			kind:       MAGIC_MISSILE,
			sprite:     m.sprite,
			damage:     m.baseDamage,
		}

		projectiles = append(projectiles, &proj)
		rl.PlaySound(m.soundEffect)
		m.cooldown.timeSinceShot = 0
	}
}
