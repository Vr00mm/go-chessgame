package game

import "strings"

type ItemType int

const (
	Pawn ItemType = iota
	Rook
	Bishop
	Queen
	King
	Knight
)

func NewItem(item string, p Pos) *Piece {

	team := strings.Split(item, "_")
	switch item {
	case "b_king":
		return newKing(team[0], p)
	case "w_king":
		return newKing(team[0], p)
	case "b_queen":
		return newQueen(team[0], p)
	case "w_queen":
		return newQueen(team[0], p)
	case "b_bishop":
		return newBishop(team[0], p)
	case "w_bishop":
		return newBishop(team[0], p)
	case "b_knight":
		return newKnight(team[0], p)
	case "w_knight":
		return newKnight(team[0], p)
	case "b_rook":
		return newRook(team[0], p)
	case "w_rook":
		return newRook(team[0], p)
	case "b_pawn":
		return newPawn(team[0], p)
	case "w_pawn":
		return newPawn(team[0], p)
	default:
		panic("Cannot Load item not known:" + item)
	}

}

func newPawn(t string, p Pos) *Piece {
	return &Piece{p, "Pawn", "pawn", t, 0}
}

func newRook(t string, p Pos) *Piece {
	return &Piece{p, "Rook", "rook", t, 0}
}

func newBishop(t string, p Pos) *Piece {
	return &Piece{p, "Bishop", "bishop", t, 0}
}

func newQueen(t string, p Pos) *Piece {
	return &Piece{p, "Queen", "queen", t, 0}
}

func newKing(t string, p Pos) *Piece {
	return &Piece{p, "King", "king", t, 0}
}

func newKnight(t string, p Pos) *Piece {
	return &Piece{p, "Knight", "knight", t, 0}
}
