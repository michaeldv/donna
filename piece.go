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

// Returns intrinsic piece value for the middlegame and the endgame.
func (p Piece) value() Score {
	switch p.kind() {
	case Pawn:
		return valuePawn
	case Knight:
		return valueKnight
	case Bishop:
		return valueBishop
	case Rook:
		return valueRook
	case Queen:
		return valueQueen
	}
	return Score{0, 0}
}

// Returns score points for a piece at given square.
func (p Piece) score(square int) Score {
	return pst[p][square]
}

// Converts a piece to "official" polyglot representation, i.e. returns (Piece - 1)
// when the color is 0 and (Piece - 3) when color is 1.
func (p Piece) polyglot() int {
	return int(p) - 1 - (int(p) & 1) << 1 // Fast int(p) - 1 - 2 * p.color()
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
	color := p.color()
	switch p.kind() {
	case King:
		if Settings.Fancy {
			return []string{"\u2654", "\u265A"}[color]
		} else {
			return []string{`K`, `k`}[color]
		}
	case Queen:
		if Settings.Fancy {
			return []string{"\u2655", "\u265B"}[color]
		} else {
			return []string{`Q`, `q`}[color]
		}
	case Rook:
		if Settings.Fancy {
			return []string{"\u2656", "\u265C"}[color]
		} else {
			return []string{`R`, `r`}[color]
		}
	case Bishop:
		if Settings.Fancy {
			return []string{"\u2657", "\u265D"}[color]
		} else {
			return []string{`B`, `b`}[color]
		}
	case Knight:
		if Settings.Fancy {
			return []string{"\u2658", "\u265E"}[color]
		} else {
			return []string{`N`, `n`}[color]
		}
	case Pawn:
		if Settings.Fancy {
			return []string{"\u2659", "\u265F"}[color]
		} else {
			return []string{`P`, `p`}[color]
		}
	}
	return ``
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
