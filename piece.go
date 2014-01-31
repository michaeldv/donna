// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

type Piece uint8

const (
        NONE   = iota
        PAWN   = 1 << 1 // 2
        KNIGHT = 2 << 1 // 4
        BISHOP = 3 << 1 // 6
        ROOK   = 4 << 1 // 8
        QUEEN  = 5 << 1 // 10
        KING   = 6 << 1 // 12
)

func King(color int) Piece {
        return Piece(color | KING)
}

func Queen(color int) Piece {
        return Piece(color | QUEEN)
}

func Rook(color int) Piece {
        return Piece(color | ROOK)
}

func Bishop(color int) Piece {
        return Piece(color | BISHOP)
}

func Knight(color int) Piece {
        return Piece(color | KNIGHT)
}

func Pawn(color int) Piece {
        return Piece(color | PAWN)
}

// Returns intrinsic piece value for the middlegame and
// the endgame.
func (p Piece) value() (int, int) {
	switch p.kind() {
        case PAWN:
                return valuePawn.midgame, valuePawn.endgame
        case KNIGHT:
                return valueKnight.midgame, valueKnight.endgame
        case BISHOP:
                return valueBishop.midgame, valueBishop.endgame
        case ROOK:
                return valueRook.midgame, valueRook.endgame
        case QUEEN:
                return valueQueen.midgame, valueQueen.endgame
        }
        return 0, 0
}

// Returns bonus points for a piece at the given square.
func (p Piece) bonus(square int) (int, int) {
	switch p.kind() {
        case PAWN:
                return bonusPawn[0][square], bonusPawn[1][square]
        case KNIGHT:
                return bonusKnight[0][square], bonusKnight[1][square]
        case BISHOP:
                return bonusBishop[0][square], bonusBishop[1][square]
        case ROOK:
                return bonusRook[0][square], bonusRook[1][square]
        case QUEEN:
                return bonusQueen[0][square], bonusQueen[1][square]
        case KING:
                return bonusKing[0][square], bonusKing[1][square]
        }
        return 0, 0
}

// return Piece - 1 when color is White(0)
//        Piece - 3 when color is Black(1)
func (p Piece) polyglot() int {
        return int(p) - 1 - 2 * p.color()
}

func (p Piece) color() int {
        return int(p) & 0x01
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
        return p & 0xFE == KING
}

func (p Piece) isQueen() bool {
        return p & 0xFE == QUEEN
}

func (p Piece) isRook() bool {
        return p & 0xFE == ROOK
}

func (p Piece) isBishop() bool {
        return p & 0xFE == BISHOP
}

func (p Piece) isKnight() bool {
        return p & 0xFE == KNIGHT
}

func (p Piece) isPawn() bool {
        return p & 0xFE == PAWN
}

func (p Piece) String() string {
        color := p.color()
        switch(p.kind()) {
        case KING:
                if Settings.Fancy {
                        return []string{"\u2654", "\u265A"}[color]
                } else {
                        return []string{`K`, `k`}[color]
                }
        case QUEEN:
                if Settings.Fancy {
                        return []string{"\u2655", "\u265B"}[color]
                } else {
                        return []string{`Q`, `q`}[color]
                }
        case ROOK:
                if Settings.Fancy {
                        return []string{"\u2656", "\u265C"}[color]
                } else {
                        return []string{`R`, `r`}[color]
                }
        case BISHOP:
                if Settings.Fancy {
                        return []string{"\u2657", "\u265D"}[color]
                } else {
                        return []string{`B`, `b`}[color]
                }
        case KNIGHT:
                if Settings.Fancy {
                        return []string{"\u2658", "\u265E"}[color]
                } else {
                        return []string{`N`, `n`}[color]
                }
        case PAWN:
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
        switch(p.kind()) {
        case KING:
                return `K`
        case QUEEN:
                return `Q`
        case ROOK:
                return `R`
        case BISHOP:
                return `B`
        case KNIGHT:
                return `N`
        }
        return ``
}
