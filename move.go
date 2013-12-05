package donna

import (
        `fmt`
)

type Move struct {
        From     int
        To       int
        Piece    Piece
        Captured Piece
}

func NewMove(from, to int, moved, captured Piece) *Move {
        move := new(Move)

        move.From = from
        move.To = to
        move.Piece = moved
        move.Captured = captured

        return move
}

func (m *Move) IsKingSideCastle() bool {
        return m.Piece.IsKing() && ((m.Piece.IsWhite() && m.From == 4 && m.To == 6) || (m.Piece.IsBlack() && m.From == 60 && m.To == 62))
}

func (m *Move) IsQueenSideCastle() bool {
        return m.Piece.IsKing() && ((m.Piece.IsWhite() && m.From == 4 && m.To == 2) || (m.Piece.IsBlack() && m.From == 60 && m.To == 58))
}

func (m *Move) IsCastle() bool {
        return m.IsKingSideCastle() || m.IsQueenSideCastle()
}

func (m *Move) IsTwoSquarePawnAdvance() bool {
        rowFrom, rowTo := Row(m.From), Row(m.To)
        return m.Piece.IsPawn() && ((m.Piece.IsWhite() && rowFrom == 1 && rowTo == 3) || (m.Piece.IsBlack() && rowFrom == 6 && rowTo == 4))
}

func (m *Move) IsCrossing(enpassant Bitmask) bool {
        return m.Piece.IsPawn() && Bitmask(1 << uint(m.To)) == enpassant
}

func (m *Move) String() string {
        if !m.IsCastle() {
                col := [2]int{ Column(m.From) + 'a', Column(m.To) + 'a' }
                row := [2]int{ Row(m.From) + 1, Row(m.To) + 1 }

                capture := '-'
                if m.Captured != 0 {
                        capture = 'x'
                }
                piece := ``
                if !m.Piece.IsPawn() {
                        piece = m.Piece.String()
                }
                format := `%s %c%d%c%c%d` // More readable with extra space.
                if !Settings.Fancy {
                        format = `%s%c%d%c%c%d`
                }
                return fmt.Sprintf(format, piece, col[0], row[0], capture, col[1], row[1])
        } else if m.IsKingSideCastle() {
                return `0-0`
        } else {
                return `0-0-0`
        }
}