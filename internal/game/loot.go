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
}

type Magnet struct {
	position  rl.Vector2
	collected bool
	radius    float32
	color     rl.Color
}

func CreateLootExpItem(spawnPos rl.Vector2, expValue int32) Loot {
	gem := Loot{
		position:  spawnPos,
		radius:    8,
		xpValue:   int(expValue),
		color:     rl.Green,
		collected: false,
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
			rl.DrawCircleV(loot[i].position, loot[i].radius, loot[i].color)
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
