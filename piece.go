package donna

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
        if Settings.Fancy {
                switch p {
                case King(WHITE):
                        return "\u2654"
                case King(BLACK):
                        return "\u265A"
                case Queen(WHITE):
                        return "\u2655"
                case Queen(BLACK):
                        return "\u265B"
                case Rook(WHITE):
                        return "\u2656"
                case Rook(BLACK):
                        return "\u265C"
                case Bishop(WHITE):
                        return "\u2657"
                case Bishop(BLACK):
                        return "\u265D"
                case Knight(WHITE):
                        return "\u2658"
                case Knight(BLACK):
                        return "\u265E"
                case Pawn(WHITE):
                        return "\u2659"
                case Pawn(BLACK):
                        return "\u265F"
                }
        } else {
                switch(p.Kind()) {
                case KING:
                        return "K"
                case QUEEN:
                        return "Q"
                case ROOK:
                        return "R"
                case BISHOP:
                        return "B"
                case KNIGHT:
                        return "N"
                case PAWN:
                        return ""
                }
        }
        return "?"
}
