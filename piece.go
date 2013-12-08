package donna

import ()

type Piece uint8

const (
        NONE   = iota
        PAWN   = 2 << 1
        KNIGHT = 3 << 1
        BISHOP = 4 << 1
        ROOK   = 5 << 1
        QUEEN  = 6 << 1
        KING   = 7 << 1
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

func (p Piece)Color() int {
        return int(p) & 0x01
}

func (p Piece)Kind() int {
        return int(p) & 0xFE
}

func (p Piece)IsWhite() bool {
        return p & 0x01 == 0
}

func (p Piece)IsBlack() bool {
        return p & 0x01 == 1
}

func (p Piece)IsKing() bool {
        return p & 0xFE == KING
}

func (p Piece)IsQueen() bool {
        return p & 0xFE == QUEEN
}

func (p Piece)IsRook() bool {
        return p & 0xFE == ROOK
}

func (p Piece)IsBishop() bool {
        return p & 0xFE == BISHOP
}

func (p Piece)IsKnight() bool {
        return p & 0xFE == KNIGHT
}

func (p Piece)IsPawn() bool {
        return p & 0xFE == PAWN
}

func (p Piece)String() string {
        color := p.Color()
        switch(p.Kind()) {
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
        // case PAWN:
        //         if Settings.Fancy {
        //                 return []string{"\u2659", "\u265F"}[color]
        //         } else {
        //                 return []string{`P`, `p`}[color]
        //         }
        }
        return ``
}
