package lape

import ()

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
