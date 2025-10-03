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

			if dbgf.allowPlayerMove {
				player.GetKnockedBack(enemies[i], dt)
			}
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

					// add  knockback to enemy on projectile hit
					// TODO: variadic knockback depending on projectile type, damage, level, etc?

					enemies[j].GetKnockedBackProjectile(projectiles[i], dt)

					// when we collide with enemy, IF attack is lightning, spawn two other projectiles to next two nearest enemies...

					// find two nearest enemies to collision point...

					// TODO: CHAIN LIGHTNING BLOCK
					// TODO: work out the idea below (spawn new chain lightning proc on chain lightning hit)
					// if projectiles[i].kind == LIGHTNING && len(enemies) > 5 {
					// 	closest := enemies[0]
					// 	secondClosest := enemies[1]
					// 	closestDistance := rl.Vector2Distance(projectiles[i].position, enemies[0].position)

					// 	for k := 0; k < len(enemies); k++ {
					// 		if enemies[k].hitpoints > 0 {
					// 			distanceTo := rl.Vector2Distance(projectiles[i].position, enemies[k].position)

					// 			if distanceTo < closestDistance {
					// 				secondClosest = closest
					// 				closest = enemies[k]
					// 				closestDistance = distanceTo
					// 			}
					// 		}

					// 	}

					// 	newProj := projectiles[i].Clone()
					// 	newProj.direction = rl.Vector2Subtract(closest.center, projectiles[i].position)
					// 	newProj2 := projectiles[i].Clone()
					// 	newProj2.direction = rl.Vector2Subtract(secondClosest.center, projectiles[i].position)

					// 	projectiles = append(projectiles, newProj)
					// 	projectiles = append(projectiles, newProj2)
					// }

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
