// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
        `fmt`
        `regexp`
)

type Move struct {
        From     int
        To       int
        Piece    Piece
        Captured Piece
        Promoted Piece
}

func NewMove(from, to int, moved, captured Piece) *Move {
        move := new(Move)

        move.From = from
        move.To = to
        move.Piece = moved
        move.Captured = captured

        return move
}

func NewMoveFromString(e2e4 string, p *Position) (move *Move) {
	re := regexp.MustCompile(`([KQRBN]?)([a-h])([1-8])-?([a-h])([1-8])[QRBN]?`)
	arr := re.FindStringSubmatch(e2e4)

	if len(arr) > 0 {
		piece   := arr[1]
		from    := Square(int(arr[3][0]-'1'), int(arr[2][0]-'a'))
                to      := Square(int(arr[5][0]-'1'), int(arr[4][0]-'a'))
		capture := p.pieces[to]

		switch piece {
		case `K`:
			move = NewMove(from, to, King(p.color), capture)
		case `Q`:
			move = NewMove(from, to, Queen(p.color), capture)
		case `R`:
			move = NewMove(from, to, Rook(p.color), capture)
		case `B`:
			move = NewMove(from, to, Bishop(p.color), capture)
		case `N`:
			move = NewMove(from, to, Knight(p.color), capture)
		default:
			move = NewMove(from, to, Pawn(p.color), capture)
		}
                if (p.pieces[from] != move.Piece) || (p.targets[from] & Shift(to) == 0) {
                        move = nil
                }
	} else if e2e4 == `0-0` || e2e4 == `0-0-0` {
                from := p.outposts[King(p.color)].FirstSet()
                to := G1
                if e2e4 == `0-0-0` {
                        to = C1
                }
                if p.color == BLACK {
                        to += 56
                }
                move = NewMove(from, to, King(p.color), 0)
                if !move.IsCastle() {
                        move = nil
                }
	}
	return
}

func (m *Move) score(position *Position) int {
	var square, midgame, endgame int

	if m.Piece.IsBlack() {
		square = m.To
	} else {
		square = flip[m.To]
	}

	switch m.Piece.Kind() {
        case PAWN:
                midgame += bonusPawn[0][square]
                endgame += bonusPawn[1][square]
        case KNIGHT:
                midgame += bonusKnight[0][square]
                endgame += bonusKnight[1][square]
        case BISHOP:
                midgame += bonusBishop[0][square]
                endgame += bonusBishop[1][square]
        case ROOK:
                midgame += bonusRook[0][square]
                endgame += bonusRook[1][square]
        case QUEEN:
                midgame += bonusQueen[0][square]
                endgame += bonusQueen[1][square]
        case KING:
                midgame += bonusKing[0][square]
                endgame += bonusKing[1][square]
        }

	return (midgame * position.stage + endgame * (256 - position.stage)) / 256
}

func (m *Move) Promote(kind int) *Move {
        m.Promoted = Piece(kind | m.Piece.Color())

        return m
}

func (m *Move) IsKingSideCastle() bool {
        return m.Piece.IsKing() && ((m.Piece.IsWhite() && m.From == E1 && m.To == G1) || (m.Piece.IsBlack() && m.From == E8 && m.To == G8))
}

func (m *Move) IsQueenSideCastle() bool {
        return m.Piece.IsKing() && ((m.Piece.IsWhite() && m.From == E1 && m.To == C1) || (m.Piece.IsBlack() && m.From == E8 && m.To == C8))
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
                col := [2]int{ Col(m.From) + 'a', Col(m.To) + 'a' }
                row := [2]int{ Row(m.From) + 1, Row(m.To) + 1 }

                capture := '-'
                if m.Captured != 0 {
                        capture = 'x'
                }
                piece, promoted := m.Piece.String(), m.Promoted.String()
                format := `%c%d%c%c%d%s`

                if m.Piece.IsPawn() { // Skip piece name if it's a pawn.
                        return fmt.Sprintf(format, col[0], row[0], capture, col[1], row[1], promoted)
                } else {
                        if Settings.Fancy { // Fancy notation is more readable with extra space.
                                format = `%s ` + format
                        } else {
                                format = `%s` + format
                        }
                        return fmt.Sprintf(format, piece, col[0], row[0], capture, col[1], row[1], promoted)
                }
        } else if m.IsKingSideCastle() {
                return `0-0`
        } else {
                return `0-0-0`
        }
}