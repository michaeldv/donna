package lape

import ()

type Piece uint8

const (
        PAWN   = 1 << 2
        KNIGHT = 1 << 3
        BISHOP = 1 << 4
        ROOK   = 1 << 5
        QUEEN  = 1 << 6
        KING   = 1 << 7
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
