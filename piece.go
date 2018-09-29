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

func (p Piece) none() bool {
	return p == Piece(0)
}

func (p Piece) some() bool {
	return p != Piece(0)
}

func (p Piece) polyglot(square int) uint64 {
	return polyglotRandom[polyglotBase[p] + square]
}

func (p Piece) color() int {
	return int(p) & 1
}

func (p Piece) id() int {
	return int(p) >> 1
}

func (p Piece) kind() int {
	return int(p) & 0xFE
}

func (p Piece) value() int {
	return pieceValue[p.id()]
}

func (p Piece) isWhite() bool {
	return p & 1 == 0
}

func (p Piece) isBlack() bool {
	return p & 1 == 1
}

func (p Piece) isKing() bool {
	return p & 0xFE == King
}

func (p Piece) isQueen() bool {
	return p & 0xFE == Queen
}

func (p Piece) isRook() bool {
	return p & 0xFE == Rook
}

func (p Piece) isBishop() bool {
	return p & 0xFE == Bishop
}

func (p Piece) isKnight() bool {
	return p & 0xFE == Knight
}

func (p Piece) isPawn() bool {
	return p & 0xFE == Pawn
}

// Returns colorless ASCII code for the piece.
func (p Piece) char() byte {
	return []byte{ 0, 0, 0, 0, 'N', 'N', 'B', 'B', 'R', 'R', 'Q', 'Q', 'K', 'K' }[p]
}

func (p Piece) String() string {
	plain := []string{ ` `, ` `, `P`, `p`, `N`, `n`, `B`, `b`, `R`, `r`, `Q`, `q`, `K`, `k` }
	fancy := []string{ ` `, ` `, "\u2659", "\u265F", "\u2658", "\u265E", "\u2657", "\u265D", "\u2656", "\u265C", "\u2655", "\u265B", "\u2654", "\u265A" }

	if engine.fancy {
		return fancy[p]
	}
	return plain[p]
}
