// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

type Piece uint8

const (
        Pawn   =  2;  BlackPawn   = Pawn   | 1;
        Knight =  4;  BlackKnight = Knight | 1;
        Bishop =  6;  BlackBishop = Bishop | 1;
        Rook   =  8;  BlackRook   = Rook   | 1;
        Queen  = 10;  BlackQueen  = Queen  | 1;
        King   = 12;  BlackKing   = King   | 1;
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

// Returns intrinsic piece value for the middlegame and
// the endgame.
func (p Piece) value() (int, int) {
	switch p.kind() {
        case Pawn:
                return valuePawn.midgame, valuePawn.endgame
        case Knight:
                return valueKnight.midgame, valueKnight.endgame
        case Bishop:
                return valueBishop.midgame, valueBishop.endgame
        case Rook:
                return valueRook.midgame, valueRook.endgame
        case Queen:
                return valueQueen.midgame, valueQueen.endgame
        }
        return 0, 0
}

// Returns bonus points for a piece at the given square.
func (p Piece) bonus(square int) (int, int) {
	switch p.kind() {
        case Pawn:
                return bonusPawn[0][square], bonusPawn[1][square]
        case Knight:
                return bonusKnight[0][square], bonusKnight[1][square]
        case Bishop:
                return bonusBishop[0][square], bonusBishop[1][square]
        case Rook:
                return bonusRook[0][square], bonusRook[1][square]
        case Queen:
                return bonusQueen[0][square], bonusQueen[1][square]
        case King:
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
        switch(p.kind()) {
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
        switch(p.kind()) {
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
