package lape

import (
        `fmt`
)

type Move struct {
        From     int
        To       int
        Piece    Piece
        Captured Piece
}

func (m *Move)Initialize(from, to int, moved, captured Piece) *Move {
        m.From = from
        m.To = to
        m.Piece = moved
        m.Captured = captured

        return m
}

func (m *Move)IsKingSideCastle() bool {
        return m.Piece.IsKing() && ((m.Piece.IsWhite() && m.From == 4 && m.To == 6) || (m.Piece.IsBlack() && m.From == 60 && m.To == 62))
}

func (m *Move)IsQueenSideCastle() bool {
        return m.Piece.IsKing() && ((m.Piece.IsWhite() && m.From == 4 && m.To == 2) || (m.Piece.IsBlack() && m.From == 60 && m.To == 58))
}

func (m *Move)IsCastle() bool {
        return m.IsKingSideCastle() || m.IsQueenSideCastle()
}

func (m *Move)String() string {
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