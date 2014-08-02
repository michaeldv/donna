// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

type Piece uint8

const (
	Pawn        = 2
	Knight      = 4
	Bishop      = 6
	Rook        = 8
	Queen       = 10
	King        = 12
	BlackPawn   = Pawn | 1
	BlackKnight = Knight | 1
	BlackBishop = Bishop | 1
	BlackRook   = Rook | 1
	BlackQueen  = Queen | 1
	BlackKing   = King | 1
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

// Returns score points for a piece at given square.
func (p Piece) score(square int) Score {
	return pst[p][square]
}

func (p Piece) polyglot(square int) uint64 {
	return polyglotRandom[polyglotBase[p] + square]
}

func (p Piece) color() int {
	return int(p) & 1
}

func (p Piece) kind() int {
	return int(p) & 0xFE
}

func (p Piece) isWhite() bool {
	return p & 0x01 == 0
}

func (p Piece) isBlack() bool {
	return p & 0x01 == 1
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

func (p Piece) String() string {
	plain := []string{ ` `, ` `, `P`, `p`, `N`, `n`, `B`, `b`, `R`, `r`, `Q`, `q`, `K`, `k` }
	fancy := []string{ ` `, ` `, "\u2659", "\u265F", "\u2658", "\u265E", "\u2657", "\u265D", "\u2656", "\u265C", "\u2655", "\u265B", "\u2654", "\u265A" }

	if Settings.Fancy {
		return fancy[p]
	}
	return plain[p]
}

// Colorless ASCII representation (perfect for tests).
func (p Piece) s() string {
	switch p.kind() {
	case King:
		return `K`
	case Queen:
		return `Q`
	case Rook:
		return `R`
	case Bishop:
		return `B`
	case Knight:
		return `N`
	}
	return ``
}
