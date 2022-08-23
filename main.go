package main

import (
	"chessgame/game"
	"chessgame/ui2d"
)

func main() {
	// Todo When we need multiple UI support - refactor event polling into it's own component
	// and run it only on the main thread
	game := game.NewGame(1)

	go func() {
		game.Run()
	}()

	ui := ui2d.NewUI(game.InputChan, game.LevelChans[0])
	ui.Run()

}
