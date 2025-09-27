package game

import rl "github.com/gen2brain/raylib-go/raylib"

type Loot struct {
	position  rl.Vector2
	radius    float32
	xpValue   int
	collected bool
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

func DrawAllLoot() {
	for i := 0; i < len(loot); i++ {
		if !loot[i].collected {
			rl.DrawCircleV(loot[i].position, loot[i].radius, loot[i].color)
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
				loot[i].collected = true
				player.exp += int32(loot[i].xpValue)
				ft := CreateFloatingEXPText(player.position, loot[i].xpValue)
				fcts = append(fcts, &ft)
			}
		}
	}
}
