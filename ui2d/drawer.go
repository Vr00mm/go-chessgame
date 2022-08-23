package ui2d

import (
	"chessgame/game"
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

func (ui *ui) DrawTest(level *game.Board) {

	ui.renderer.Clear()
	ui.r.Seed(1)
	// Render Map Tiles
	for x, row := range level.Map {
		for y, tile := range row {

			texture, _ := ui.renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888,
				sdl.TEXTUREACCESS_STATIC, int32(BlocSize), int32(BlocSize))
			dstRect := sdl.Rect{X: int32(x * BlocSize), Y: int32(y * BlocSize), W: int32(BlocSize), H: int32(BlocSize)}

			if tile.Rune == "bg_black" {
				ui.renderer.SetDrawColor(118, 150, 86, 255)
			} else {
				ui.renderer.SetDrawColor(238, 238, 210, 255)
			}

			ui.renderer.SetRenderTarget(texture)
			ui.renderer.FillRect(&dstRect)
			ui.renderer.SetRenderTarget(nil)
		}
	}
	// Render Pieces
	for pos, piece := range level.Pieces {
		imageName := piece[0].Team + "_" + piece[0].Rune
		texture := ui.texturesIndex[imageName]
		rect := sdl.Rect{X: int32(pos.X * BlocSize), Y: int32(pos.Y) * int32(BlocSize), W: int32(BlocSize), H: int32(BlocSize)}
		err := ui.renderer.Copy(texture, nil, &rect)
		if err != nil {
			panic(err)
		}
	}

}

func (ui *ui) Draw(level *game.Board) {

	ui.renderer.Clear()
	ui.r.Seed(1)
	// Render Map Tiles
	for x, row := range level.Map {
		for y, tile := range row {

			dstRect := sdl.Rect{X: int32(x * 32), Y: int32(y * 32), W: 128, H: 128}

			fmt.Println("rendering: ", tile.Rune)
			texture := ui.texturesIndex[tile.Rune]
			err := ui.renderer.Copy(texture, nil, &dstRect)
			if err != nil {
				panic(err)
			}

			if tile.OverlayRune != "" {
				// Todo what if there are multiple variants for overlay images?
				texture := ui.texturesIndex[tile.OverlayRune]
				err := ui.renderer.Copy(texture, nil, &dstRect)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	/* Event UI Begin
	textStart := int32(float64(ui.winHeight) * .68)
	textWidth := int32(float64(ui.winWidth) * .25)

	ui.renderer.Copy(ui.eventBackground, nil, &sdl.Rect{X: 0, Y: textStart, W: textWidth, H: int32(ui.winHeight) - textStart})
	i := level.EventPos
	count := 0
	_, fontSizeY, _ := ui.fontSmall.SizeUTF8("A")
	for {
		event := level.Events[i]
		if event != "" {
			tex := ui.stringToTexture(event, sdl.Color{R: 255, G: 0, B: 0, A: 0}, FontSmall)
			_, _, w, h, _ := tex.Query()
			ui.renderer.Copy(tex, nil, &sdl.Rect{X: 5, Y: int32(count*fontSizeY) + textStart, W: w, H: h})
		}
		i = (i + 1) % (len(level.Events))
		count++
		if i == level.EventPos {
			break
		}
	}
	// Event UI End
	*/

}
