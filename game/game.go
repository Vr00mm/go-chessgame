package game

import (
	"encoding/csv"
	"fmt"
	"log"
	"strings"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Game struct {
	LevelChans []chan *Board
	InputChan  chan *Input
	Board      *Board
}

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

const (
	MovePiece GameEvent = iota
)

const (
	None InputType = iota
	QuitGame
	CloseWindow
	Search //temporary
)

const (
	BKing   string = "b_king"
	BQueen         = "b_queen"
	BRook          = "b_rook"
	BBishop        = "b_bishop"
	BKnight        = "b_knight"
	BPawn          = "b_pawn"
	WKing          = "w_king"
	WQueen         = "w_queen"
	WRook          = "w_rook"
	WBishop        = "w_bishop"
	WKnight        = "w_knight"
	WPawn          = "w_pawn"
)

const (
	BoardInit string = `b_rook;b_knight;b_bishop;b_queen;b_king;b_bishop;b_knight;b_rook
b_pawn;b_pawn;b_pawn;b_pawn;b_pawn;b_pawn;b_pawn;b_pawn
;;;;;;;
;;;;;;;
;;;;;;;
;;;;;;;
w_pawn;w_pawn;w_pawn;w_pawn;w_pawn;w_pawn;w_pawn;w_pawn
w_rook;w_knight;w_bishop;w_king;w_queen;w_bishop;w_knight;w_rook`
)

type InputType int

type Input struct {
	Typ          InputType
	Item         *Piece
	BoardChannel chan *Board
}

type Tile struct {
	Rune        string
	OverlayRune string
}

type Pos struct {
	X, Y int
}

type Piece struct {
	Pos
	Name       string
	Rune       string
	Team       string
	MooveCount int
}

type GameEvent int

type Board struct {
	Map              [][]Tile
	BPlayer, WPlayer Player
	Pieces           map[Pos][]*Piece
	Events           []string
	EventPos         int
	Debug            map[Pos]bool
	LastEvent        GameEvent
}

type Player struct {
	Name         string
	Team         float64
	Actions      []string
	ActionPoints float64
	Pieces       []*Piece
}

func NewGame(numWindows int) *Game {
	levelChans := make([]chan *Board, numWindows)
	for i := range levelChans {
		levelChans[i] = make(chan *Board)
	}
	inputChan := make(chan *Input)
	board := loadBoard()

	game := &Game{levelChans, inputChan, board}
	return game
}

func loadBoard() *Board {

	BPLayer := Player{}
	BPLayer.Name = "Player 1"

	WPLayer := Player{}
	WPLayer.Name = "Player 2"

	NbBlocHeight := 8
	NbBlocLenght := 8

	level := &Board{}
	level.Debug = make(map[Pos]bool)
	level.Events = make([]string, 10)
	level.BPlayer = BPLayer
	level.WPlayer = WPLayer

	level.Map = make([][]Tile, NbBlocHeight)
	level.Pieces = make(map[Pos][]*Piece)

	boardInit := strings.NewReader(BoardInit)
	csvReader := csv.NewReader(boardInit)
	csvReader.Comma = ';'
	csvReader.Comment = '#'
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	k := 0
	for x := 0; x < len(level.Map); x++ {
		level.Map[x] = make([]Tile, NbBlocLenght)
		for y := 0; y < len(level.Map[x]); y++ {
			ground := ""
			if k%2 == 1 {
				ground = "bg_white"
			} else {
				ground = "bg_black"
			}

			var t Tile
			pos := Pos{y, x}

			c := data[x][y]
			switch c {
			case "":
				t.Rune = ground
			default:
				level.Pieces[pos] = append(level.Pieces[pos], NewItem(c, pos))
				t.Rune = ground
			}

			level.Map[x][y] = t
			k++
		}
		k++
	}
	fmt.Printf("%+v\n", level.Map)

	return level
}

func (game *Game) Run() {
	fmt.Println("Starting...")

	count := 0
	for _, lchan := range game.LevelChans {
		lchan <- game.Board
	}

	//for input := range game.InputChan {
	for range game.InputChan {

		//game.handleInput(input)

		//game.Level.AddEvent("move:" + strconv.Itoa(count))
		count++

		if len(game.LevelChans) == 0 {
			return
		}

		for _, lchan := range game.LevelChans {
			lchan <- game.Board
		}
	}

}
