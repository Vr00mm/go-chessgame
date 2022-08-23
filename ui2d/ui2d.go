package ui2d

import (
	"chessgame/game"
	"math/rand"
	"strconv"
	"os"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func init() {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	err = ttf.Init()
	if err != nil {
		panic(err)
	}

	err = mix.Init(mix.INIT_OGG)
	//SDL Bug here, ignoring error
	/*if err != nil {
		panic(err)
	}*/
}

type mouseState struct {
	leftButton  bool
	rightButton bool
	pos         game.Pos
}

type uiState int

type sounds struct {
	moovingPiece []*mix.Chunk
}

const BlocSize int = 128

const (
	UIMain uiState = iota
	UIInventory
)

type ui struct {
	state             uiState
	draggedItem       *game.Piece
	sounds            sounds
	winWidth          int
	winHeight         int
	renderer          *sdl.Renderer
	window            *sdl.Window
	texturesIndex     map[string]*sdl.Texture
	prevKeyboardState []uint8
	keyboardState     []uint8
	centerX           int
	centerY           int
	r                 *rand.Rand
	levelChan         chan *game.Board
	inputChan         chan *game.Input
	fontSmall         *ttf.Font
	fontMedium        *ttf.Font
	fontLarge         *ttf.Font

	eventBackground           *sdl.Texture
	groundInventoryBackground *sdl.Texture
	slotBackground            *sdl.Texture

	str2TexSmall  map[string]*sdl.Texture
	str2TexMedium map[string]*sdl.Texture
	str2TexLarge  map[string]*sdl.Texture

	currentMouseState *mouseState
	prevMouseState    *mouseState
}

func NewUI(inputChan chan *game.Input, levelChan chan *game.Board) *ui {

	ui := &ui{}
	ui.state = UIMain
	ui.str2TexSmall = make(map[string]*sdl.Texture)
	ui.str2TexMedium = make(map[string]*sdl.Texture)
	ui.str2TexLarge = make(map[string]*sdl.Texture)
	ui.inputChan = inputChan
	ui.levelChan = levelChan
	ui.r = rand.New(rand.NewSource(1))
	ui.winHeight = 8 * BlocSize
	ui.winWidth = 8 * BlocSize
	window, err := sdl.CreateWindow("ChessGame", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		int32(ui.winWidth), int32(ui.winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	ui.window = window

	ui.renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	//sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	ui.loadTextures()

	ui.keyboardState = sdl.GetKeyboardState()
	ui.prevKeyboardState = make([]uint8, len(ui.keyboardState))
	for i, v := range ui.keyboardState {
		ui.prevKeyboardState[i] = v
	}

	ui.centerX = -1
	ui.centerY = -1

	ui.fontSmall, err = ttf.OpenFont("test.ttf", int(float64(ui.winWidth)*.015))
	if err != nil {
		panic(err)
	}

	ui.fontMedium, err = ttf.OpenFont("test.ttf", 32)
	if err != nil {
		panic(err)
	}

	ui.fontLarge, err = ttf.OpenFont("test.ttf", 64)
	if err != nil {
		panic(err)
	}

	ui.eventBackground = ui.GetSinglePixelTex(sdl.Color{R: 0, G: 0, B: 0, A: 128})
	ui.eventBackground.SetBlendMode(sdl.BLENDMODE_BLEND)

	ui.groundInventoryBackground = ui.GetSinglePixelTex(sdl.Color{R: 149, G: 84, B: 19, A: 200})
	ui.groundInventoryBackground.SetBlendMode(sdl.BLENDMODE_BLEND)

	ui.slotBackground = ui.GetSinglePixelTex(sdl.Color{R: 0, G: 0, B: 0, A: 0})

	//if( Mix_OpenAudio( 22050, MIX_DEFAULT_FORMAT, 2, 4096 ) == -1 )
	err = mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096)
	if err != nil {
		panic(err)
	}
	/*mus, err := mix.LoadMUS("ui2d/assets/ambient.ogg")
	if err != nil {
		panic(err)
	}
	mus.Play(-1)
	*/

	moovingPieceBase := "sounds"+ string(os.PathSeparator)  +"footstep0"
	for i := 0; i < 10; i++ {
		moovePieceFile := moovingPieceBase + strconv.Itoa(i) + ".ogg"
		moovePieceSound, err := mix.LoadWAV(moovePieceFile)
		if err != nil {
			panic(err)
		}
		ui.sounds.moovingPiece = append(ui.sounds.moovingPiece, moovePieceSound)
	}

	return ui
}

func getMouseState() *mouseState {
	mouseX, mouseY, mouseButtonState := sdl.GetMouseState()
	leftButton := mouseButtonState & sdl.ButtonLMask()
	rightButton := mouseButtonState & sdl.ButtonRMask()
	var result mouseState
	result.pos = game.Pos{X: int(mouseX), Y: int(mouseY)}
	result.leftButton = !(leftButton == 0)
	result.rightButton = !(rightButton == 0)

	return &result
}

func (ui *ui) Run() {
	var newLevel *game.Board
	ui.prevMouseState = getMouseState()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_CLOSE {
					ui.inputChan <- &game.Input{Typ: game.CloseWindow, BoardChannel: ui.levelChan}
				}
			}
		}
		ui.currentMouseState = getMouseState()

		// Suspect quick keypresses sometimes cause channel gridlock
		var ok bool
		select {
		case newLevel, ok = <-ui.levelChan:
			if ok {
				switch newLevel.LastEvent {
				case game.MovePiece:
					playRandomSound(ui.sounds.moovingPiece, 10)
				default:
					//add more sounds
				}

			}
		default:
		}

		ui.DrawTest(newLevel)
		ui.renderer.Present()

		/*
			var input game.Input

			item := ui.CheckGroundItems(newLevel)
			if item != nil {
				input.Typ = game.TakeItem
				input.Item = item
			}
			if sdl.GetKeyboardFocus() == ui.window || sdl.GetMouseFocus() == ui.window {

				if ui.keyDownOnce(sdl.SCANCODE_UP) {
					input.Typ = game.Up
				} else if ui.keyDownOnce(sdl.SCANCODE_DOWN) {
					input.Typ = game.Down
				} else if ui.keyDownOnce(sdl.SCANCODE_LEFT) {
					input.Typ = game.Left
				} else if ui.keyDownOnce(sdl.SCANCODE_RIGHT) {
					input.Typ = game.Right
				} else if ui.keyDownOnce(sdl.SCANCODE_T) {
					input.Typ = game.TakeAll
				} else if ui.keyDownOnce(sdl.SCANCODE_I) {
					fmt.Println("I")
					if ui.state == UIMain {
						ui.state = UIInventory
					} else {
						ui.state = UIMain
					}
				}

				for i, v := range ui.keyboardState {
					ui.prevKeyboardState[i] = v
				}

				if input.Typ != game.None {
					ui.inputChan <- &input
				}
			}
		*/
		ui.prevMouseState = ui.currentMouseState
		sdl.Delay(10)

	}

}
