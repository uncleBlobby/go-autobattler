package game

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func EnemyCollisionLoop(dt float32) {
	for i := 0; i < len(enemies); i++ {
		if rl.CheckCollisionRecs(enemies[i].collider, player.collider) {
			// player take damage...

			//fmt.Printf("enemies[i].timeSinceLastShot: %v", enemies[i].timeSinceLastShot)

			if dbgf.allowPlayerDamage {
				if enemies[i].timeSinceLastShot >= enemies[i].shootCooldown {
					//fmt.Println("ENEMY DAMAGE PLAYER")
					player.hitpoints -= 10
					ft := CreateFloatingText(player.position, 10)
					fcts = append(fcts, &ft)
					enemies[i].timeSinceLastShot = 0
				}
			}

			// player get bumped back

			player.GetKnockedBack(enemies[i], dt)
		}

		for j := 0; j < len(enemies); j++ {
			if rl.CheckCollisionRecs(enemies[i].collider, enemies[j].collider) {

				enemies[j].GetKnockedBack(enemies[i], dt)
			}
		}
	}
}

func ProjectileCollisionLoop(dt float32) {
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

					// random chance for a CRIT (bullet doesn't die, but instead penetrates)..

					if rand.Float32() > projectiles[i].critChance {
						projectiles[i].isDead = true
					} else {
						critText := CreateFloatingCRITText(enemies[j].position)
						fcts = append(fcts, &critText)
					}

					enemies[j].hitpoints -= int32(projectiles[i].damage)

					if enemies[j].hitpoints <= 0 {
						enemies[j].isDead = true
						if rand.Float32() > 0.95 {
							mg := CreateLootMagnetItem(enemies[j].center)
							magnets = append(magnets, &mg)
							continue
						}
						lt := CreateLootExpItem(enemies[j].center, 10*player.level)
						loot = append(loot, &lt)

					}

					newFt := CreateFloatingText(enemies[j].position, projectiles[i].damage)
					fcts = append(fcts, &newFt)
				}
			}
		}
	}
}
