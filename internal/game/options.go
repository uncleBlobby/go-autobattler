package game

import (
	"fmt"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawDebugOptionsScreen() {
	screenColor := rl.LightGray
	screenColor.A = uint8(128)
	rl.DrawRectangle(50, 50, 1920-100, 1080-100, screenColor)

	rl.DrawText("Debug Options", 800, 500, 24, rl.Black)
	rl.DrawText("pick some new skills or something...", 820, 600, 20, rl.Black)

	// rl.DrawRectangle(750, 900, 100, 100, rl.Orange)

	if rl.IsKeyReleased(rl.KeyBackspace) || rl.IsKeyReleased(rl.KeyF3) {
		showDebugOptionsScreen = false
	}

	weaponsOn := rg.CheckBox(rl.Rectangle{X: 200, Y: 200, Width: 100, Height: 50}, "weapons on/off", dbgf.enableWeapons)
	rl.DrawText(fmt.Sprintf("%v", dbgf.enableWeapons), 250, 200, 12, rl.Black)

	dbgf.enableWeapons = weaponsOn

	spellsOn := rg.CheckBox(rl.Rectangle{X: 200, Y: 300, Width: 100, Height: 50}, "spells on/off", dbgf.enableSpells)
	rl.DrawText(fmt.Sprintf("%v", dbgf.enableSpells), 250, 300, 12, rl.Black)

	dbgf.enableSpells = spellsOn

	playerDamageOn := rg.CheckBox(rl.Rectangle{X: 200, Y: 400, Width: 100, Height: 50}, "damage on/off", dbgf.allowPlayerDamage)
	rl.DrawText(fmt.Sprintf("%v", dbgf.allowPlayerDamage), 250, 400, 12, rl.Black)

	dbgf.allowPlayerDamage = playerDamageOn

	playerMoveOn := rg.CheckBox(rl.Rectangle{X: 200, Y: 500, Width: 100, Height: 50}, "player move on/off", dbgf.allowPlayerMove)
	rl.DrawText(fmt.Sprintf("%v", dbgf.allowPlayerMove), 250, 500, 12, rl.Black)

	dbgf.allowPlayerMove = playerMoveOn
}

type DebugOptionToggle struct {
	optionName     string
	optionValue    bool
	buttonPosition rl.Vector2
	buttonSize     rl.Vector2
	drawCheckbox   func(rl.Rectangle, string, bool) bool
}

func InitWeaponDebugToggle(initialOptions DebugFlags) DebugOptionToggle {
	var weaponDebugToggle = DebugOptionToggle{
		optionName:     "Weapons On/Off",
		optionValue:    initialOptions.enableWeapons,
		buttonPosition: rl.Vector2{X: 200, Y: 200},
		buttonSize:     rl.Vector2{X: 100, Y: 50},
		drawCheckbox:   rg.CheckBox,
	}

	return weaponDebugToggle
}
