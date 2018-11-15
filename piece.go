// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

type Piece int

const (
	White = iota	// 0
	Black		// 1
	Pawn		// 2
	BlackPawn	// 3
	Knight		// 4
	BlackKnight	// 5
	Bishop		// 6
	BlackBishop	// 7
	Rook		// 8
	BlackRook	// 9
	Queen		// 10
	BlackQueen	// 11
	King		// 12
	BlackKing	// 13
)

func king(color int) Piece {
	return Piece(color | King)
}

func queen(color int) Piece {
	return Piece(color | Queen)
}

func rook(color int) Piece {
	return Piece(color | Rook)
}

func bishop(color int) Piece {
	return Piece(color | Bishop)
}

func knight(color int) Piece {
	return Piece(color | Knight)
}

func pawn(color int) Piece {
	return Piece(color | Pawn)
}

func (p Piece) noneʔ() bool {
	return p == Piece(0)
}

func (p Piece) someʔ() bool {
	return p != Piece(0)
}

func (p Piece) polyglot(sq Square) uint64 {
	return polyglotRandom[p-2][sq]
}

func (p Piece) color() int {
	return int(p) & 1
}

func (p Piece) kind() int {
	return int(p) & 0xFE
}

func (p Piece) value() int {
	return pieceValue[p]
}

func (p Piece) whiteʔ() bool {
	return p & 1 == White
}

func (p Piece) blackʔ() bool {
	return p & 1 == Black
}

func (p Piece) kingʔ() bool {
	return p & 0xFE == King
}

func (p Piece) queenʔ() bool {
	return p & 0xFE == Queen
}

func (p Piece) rookʔ() bool {
	return p & 0xFE == Rook
}

func (p Piece) bishopʔ() bool {
	return p & 0xFE == Bishop
}

func (p Piece) knightʔ() bool {
	return p & 0xFE == Knight
}

func (p Piece) pawnʔ() bool {
	return p & 0xFE == Pawn
}

// Returns colorless ASCII code for the piece.
func (p Piece) char() byte {
	return []byte{ 0, 0, 0, 0, 'N', 'N', 'B', 'B', 'R', 'R', 'Q', 'Q', 'K', 'K' }[p]
}

func (p Piece) String() string {
	plain := []string{ ` `, ` `, `P`, `p`, `N`, `n`, `B`, `b`, `R`, `r`, `Q`, `q`, `K`, `k` }
	fancy := []string{ ` `, ` `, "\u2659", "\u265F", "\u2658", "\u265E", "\u2657", "\u265D", "\u2656", "\u265C", "\u2655", "\u265B", "\u2654", "\u265A" }

	if engine.fancyʔ {
		return fancy[p]
	}
	return plain[p]
}
