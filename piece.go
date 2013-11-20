package lape

import ()

type Piece uint8

const (
        KING   = 1 << 7
        QUEEN  = 1 << 6
        ROOK   = 1 << 5
        BISHOP = 1 << 4
        KNIGHT = 1 << 3
        PAWN   = 1 << 2
)


func (p Piece)King(color int) Piece {
        return Piece(color | KING)
}

func (p Piece)Queen(color int) Piece {
        return Piece(color | QUEEN)
}

func (p Piece)Rook(color int) Piece {
        return Piece(color | ROOK)
}

func (p Piece)Bishop(color int) Piece {
        return Piece(color | BISHOP)
}

func (p Piece)Kinight(color int) Piece {
        return Piece(color | KNIGHT)
}

func (p Piece)Pawn(color int) Piece {
        return Piece(color | PAWN)
}

func (p Piece)Color() int {
        return int(p) & 0x01
}

func (p Piece)Type() int {
        return int(p) & 0xF7
}

func (p Piece)IsWhite() bool {
        return p & 0x01 == 0
}

func (p Piece)IsBlack() bool {
        return p & 0x01 == 1
}

func (p Piece)IsKing() bool {
        return p & 0xF7 == KING
}

func (p Piece)IsQueen() bool {
        return p & 0xF7 == QUEEN
}

func (p Piece)IsRook() bool {
        return p & 0xF7 == ROOK
}

func (p Piece)IsBishop() bool {
        return p & 0xF7 == BISHOP
}

func (p Piece)IsKinight() bool {
        return p & 0xF7 == KNIGHT
}

func (p Piece)IsPawn() bool {
        return p & 0xF7 == PAWN
}

func (p Piece)ToString() string {
        switch p {
        case KING:
                return "\u2654"
        case KING | 1:
                return "\u265A"
        case QUEEN:
                return "\u2655"
        case QUEEN | 1:
                return "\u265B"
        case ROOK:
                return "\u2656"
        case ROOK | 1:
                return "\u265C"
        case BISHOP:
                return "\u2657"
        case BISHOP | 1:
                return "\u265D"
        case KNIGHT:
                return "\u2658"
        case KNIGHT | 1:
                return "\u265E"
        case PAWN:
                return "\u2659"
        case PAWN | 1:
                return "\u265F"
        }
        return "\u22C5"
}
