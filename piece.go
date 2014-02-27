// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

type Piece uint8

const (
        WhitePawn   =  2;  BlackPawn   =  3
        WhiteKnight =  4;  BlackKnight =  5
        WhiteBishop =  6;  BlackBishop =  7
        WhiteRook   =  8;  BlackRook   =  9
        WhiteQueen  = 10;  BlackQueen  = 11
        WhiteKing   = 12;  BlackKing   = 13
)

func King(color int) Piece {
        return Piece(color | WhiteKing)
}

func Queen(color int) Piece {
        return Piece(color | WhiteQueen)
}

func Rook(color int) Piece {
        return Piece(color | WhiteRook)
}

func Bishop(color int) Piece {
        return Piece(color | WhiteBishop)
}

func Knight(color int) Piece {
        return Piece(color | WhiteKnight)
}

func Pawn(color int) Piece {
        return Piece(color | WhitePawn)
}

// Returns intrinsic piece value for the middlegame and
// the endgame.
func (p Piece) value() (int, int) {
	switch p.kind() {
        case WhitePawn:
                return valuePawn.midgame, valuePawn.endgame
        case WhiteKnight:
                return valueKnight.midgame, valueKnight.endgame
        case WhiteBishop:
                return valueBishop.midgame, valueBishop.endgame
        case WhiteRook:
                return valueRook.midgame, valueRook.endgame
        case WhiteQueen:
                return valueQueen.midgame, valueQueen.endgame
        }
        return 0, 0
}

// Returns bonus points for a piece at the given square.
func (p Piece) bonus(square int) (int, int) {
	switch p.kind() {
        case WhitePawn:
                return bonusPawn[0][square], bonusPawn[1][square]
        case WhiteKnight:
                return bonusKnight[0][square], bonusKnight[1][square]
        case WhiteBishop:
                return bonusBishop[0][square], bonusBishop[1][square]
        case WhiteRook:
                return bonusRook[0][square], bonusRook[1][square]
        case WhiteQueen:
                return bonusQueen[0][square], bonusQueen[1][square]
        case WhiteKing:
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
        return p & 0xFE == WhiteKing
}

func (p Piece) isQueen() bool {
        return p & 0xFE == WhiteQueen
}

func (p Piece) isRook() bool {
        return p & 0xFE == WhiteRook
}

func (p Piece) isBishop() bool {
        return p & 0xFE == WhiteBishop
}

func (p Piece) isKnight() bool {
        return p & 0xFE == WhiteKnight
}

func (p Piece) isPawn() bool {
        return p & 0xFE == WhitePawn
}

func (p Piece) String() string {
        color := p.color()
        switch(p.kind()) {
        case WhiteKing:
                if Settings.Fancy {
                        return []string{"\u2654", "\u265A"}[color]
                } else {
                        return []string{`K`, `k`}[color]
                }
        case WhiteQueen:
                if Settings.Fancy {
                        return []string{"\u2655", "\u265B"}[color]
                } else {
                        return []string{`Q`, `q`}[color]
                }
        case WhiteRook:
                if Settings.Fancy {
                        return []string{"\u2656", "\u265C"}[color]
                } else {
                        return []string{`R`, `r`}[color]
                }
        case WhiteBishop:
                if Settings.Fancy {
                        return []string{"\u2657", "\u265D"}[color]
                } else {
                        return []string{`B`, `b`}[color]
                }
        case WhiteKnight:
                if Settings.Fancy {
                        return []string{"\u2658", "\u265E"}[color]
                } else {
                        return []string{`N`, `n`}[color]
                }
        case WhitePawn:
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
        case WhiteKing:
                return `K`
        case WhiteQueen:
                return `Q`
        case WhiteRook:
                return `R`
        case WhiteBishop:
                return `B`
        case WhiteKnight:
                return `N`
        }
        return ``
}
