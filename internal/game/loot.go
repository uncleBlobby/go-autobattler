package game

import rl "github.com/gen2brain/raylib-go/raylib"

type Collectible interface {
	Pickup(player *Player)
}

type Loot struct {
	position  rl.Vector2
	radius    float32
	xpValue   int
	collected bool
	color     rl.Color
	// sprite    rl.Texture2D
	anim AnimatedSprite
}

type Magnet struct {
	position  rl.Vector2
	collected bool
	radius    float32
	color     rl.Color
}

type AnimatedSprite struct {
	sprite       rl.Texture2D
	spriteSize   rl.Vector2
	numFrames    int
	frameSpeed   int
	currentFrame int
	frameCounter float32
}

func CreateLootExpItem(spawnPos rl.Vector2, expValue int32) Loot {

	//sprite := rl.LoadTexture("assets/sprites/gem-green.png")

	anim := AnimatedSprite{
		sprite:       lootGemSpriteSheet,
		spriteSize:   rl.Vector2{16, 16},
		numFrames:    7,
		frameSpeed:   7,
		currentFrame: 0,
		frameCounter: 0,
	}

	gem := Loot{
		position:  spawnPos,
		radius:    8,
		xpValue:   int(expValue),
		color:     rl.Green,
		collected: false,
		anim:      anim,
	}
	return gem
}

func CreateLootMagnetItem(spawnPos rl.Vector2) Magnet {
	magnet := Magnet{
		position:  spawnPos,
		collected: false,
		color:     rl.Purple,
		radius:    8,
	}

	return magnet
}

func (m *Magnet) Pickup(player *Player) {

	m.collected = true
	player.buffs = append(player.buffs, Buff{
		duration:             2,
		pickupRadiusBonus:    1000,
		pickupRadiusOriginal: player.pickupRadius,
	})
}

func DrawAllLoot() {
	for i := 0; i < len(loot); i++ {
		if !loot[i].collected {
			// rl.DrawCircleV(loot[i].position, loot[i].radius, loot[i].color)

			src := rl.Rectangle{0, float32(loot[i].anim.currentFrame) * loot[i].anim.spriteSize.Y, loot[i].anim.spriteSize.X, loot[i].anim.spriteSize.Y}
			dst := rl.Rectangle{loot[i].position.X, loot[i].position.Y, 32, 32}
			og := rl.Vector2{dst.Width / 2, dst.Height / 2}

			// rl.DrawTextureEx(loot[i].sprite, loot[i].position, 0, 1, rl.White)
			rl.DrawTexturePro(loot[i].anim.sprite, src, dst, og, 0, rl.White)
		}
	}

	for i := 0; i < len(magnets); i++ {
		if !magnets[i].collected {
			rl.DrawCircleV(magnets[i].position, magnets[i].radius, magnets[i].color)
		}
	}
}

func UpdateAllLoot(dt float32) {
	for i := 0; i < len(loot); i++ {
		if !loot[i].collected {
			// move toward player if within wider pickup radius

			loot[i].anim.frameCounter += dt
			if loot[i].anim.frameCounter >= (1.0 / float32(loot[i].anim.frameSpeed)) {
				loot[i].anim.frameCounter = 0
				loot[i].anim.currentFrame++
				if loot[i].anim.currentFrame >= loot[i].anim.numFrames {
					loot[i].anim.currentFrame = 0
				}
			}

			if rl.CheckCollisionCircles(loot[i].position, loot[i].radius, player.center, player.pickupRadius) {
				dir := rl.Vector2Subtract(player.center, loot[i].position)
				dir = rl.Vector2Normalize(dir)

				loot[i].position.X += (250 * dir.X) * dt
				loot[i].position.Y += (250 * dir.Y) * dt
			}

			if rl.CheckCollisionCircleRec(loot[i].position, loot[i].radius, player.collider) {
				rl.PlaySound(player.expPickupSfx)
				loot[i].collected = true
				player.exp += int32(loot[i].xpValue)
				ft := CreateFloatingEXPText(player.position, loot[i].xpValue)
				fcts = append(fcts, &ft)
			}
		}
	}

	for i := 0; i < len(magnets); i++ {
		if !magnets[i].collected {

			if rl.CheckCollisionCircles(magnets[i].position, magnets[i].radius, player.center, player.pickupRadius) {
				dir := rl.Vector2Subtract(player.center, magnets[i].position)
				dir = rl.Vector2Normalize(dir)

				magnets[i].position.X += (250 * dir.X) * dt
				magnets[i].position.Y += (250 * dir.Y) * dt
			}

			if rl.CheckCollisionCircleRec(magnets[i].position, magnets[i].radius, player.collider) {
				magnets[i].Pickup(player)
				magnets[i].collected = true
			}
		}
	}
}
