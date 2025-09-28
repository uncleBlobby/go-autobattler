package game

import rl "github.com/gen2brain/raylib-go/raylib"

type Spell interface {
	Shoot(startPos rl.Vector2, enemy *Enemy)
	Update(dt float32)
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
}

func InitMagicMissile(p *Player) {

	sfx := rl.LoadSound("assets/sounds/magic-1.wav")
	rl.SetSoundVolume(sfx, 0.7)

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
	}

	p.spells = append(p.spells, &mm)
}

func (m *MagicMissile) Update(dt float32) {
	m.cooldown.timeSinceShot += dt
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
		}

		projectiles = append(projectiles, &proj)
		rl.PlaySound(m.soundEffect)
		m.cooldown.timeSinceShot = 0
	}
}
