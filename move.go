// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
        `fmt`
        `regexp`
)

type Move struct {
        from     int
        to       int
        score    int
        piece    Piece
        captured Piece
        promoted Piece
}

func NewMove(p *Position, from, to int) *Move {
        move := new(Move)

        move.from = from
        move.to = to
        move.piece = p.pieces[from]
        move.captured = p.pieces[to]

        if p.enpassant != 0 && to == p.enpassant {
                move.captured = Pawn(p.color^1)
        }

        if move.captured == 0 {
                move.score = move.calculateScore(p)
        } else {
                move.score = move.calculateValue()
        }

        return move
}

func NewMoveFromString(e2e4 string, p *Position) (move *Move) {
	re := regexp.MustCompile(`([KkQqRrBbNn]?)([a-h])([1-8])-?([a-h])([1-8])([QqRrBbNn]?)`)
	arr := re.FindStringSubmatch(e2e4)

	if len(arr) > 0 {
		name  := arr[1]
		from  := Square(int(arr[3][0]-'1'), int(arr[2][0]-'a'))
		to    := Square(int(arr[5][0]-'1'), int(arr[4][0]-'a'))
		promo := arr[6]

		var piece Piece
		switch name {
		case `K`, `k`:
			piece = King(p.color)
		case `Q`, `q`:
			piece = Queen(p.color)
		case `R`, `r`:
			piece = Rook(p.color)
		case `B`, `b`:
			piece = Bishop(p.color)
		case `N`, `n`:
			piece = Knight(p.color)
		default:
			piece = Pawn(p.color)
		}
                if (p.pieces[from] != piece) || (p.targets[from] & Shift(to) == 0) {
                        move = nil
                } else {
                        move = NewMove(p, from, to)
                        if len(promo) > 0 {
                                switch promo {
                                case `Q`, `q`:
                                        move.promote(QUEEN)
                                case `R`, `r`:
                                        move.promote(ROOK)
                                case `B`, `b`:
                                        move.promote(BISHOP)
                                case `N`, `n`:
                                        move.promote(KNIGHT)
                                default:
                                        move = nil
                                }
                        }
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
                move = NewMove(p, from, to)
                if !move.isCastle() {
                        move = nil
                }
	}
	return
}

func (m *Move) promote(kind int) *Move {
        m.promoted = Piece(kind | m.piece.color())

        return m
}

func (m *Move) is(move *Move) bool {
        return m.from == move.from  &&
                 m.to == move.to    &&
              m.piece == move.piece &&
           m.captured == m.captured &&
           m.promoted == m.promoted
}

func (m *Move) calculateScore(position *Position) int {
	var midgame, endgame int
	square := flip[m.piece.color()][m.to]

	switch m.piece.kind() {
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

// PxQ, NxQ, BxQ, RxQ, QxQ, KxQ => where => QUEEN  = 5 << 1 // 10
// PxR, NxR, BxR, RxR, QxR, KxR             ROOK   = 4 << 1 // 8
// PxB, NxB, BxB, RxB, QxB, KxB             BISHOP = 3 << 1 // 6
// PxN, NxN, BxN, RxN, QxN, KxN             KNIGHT = 2 << 1 // 4
// PxP, NxP, BxP, RxP, QxP, KxP             PAWN   = 1 << 1 // 2
func (m *Move) calculateValue() int {
        if m.captured == 0 || m.captured.kind() == KING {
                return 0
        }

        victim := (QUEEN - m.captured.kind()) / PAWN
        attacker := m.piece.kind() / PAWN - 1

        return victimAttacker[victim][attacker]
}

func (m *Move) isKingSideCastle() bool {
        return m.piece.isKing() && ((m.piece.isWhite() && m.from == E1 && m.to == G1) || (m.piece.isBlack() && m.from == E8 && m.to == G8))
}

func (m *Move) isQueenSideCastle() bool {
        return m.piece.isKing() && ((m.piece.isWhite() && m.from == E1 && m.to == C1) || (m.piece.isBlack() && m.from == E8 && m.to == C8))
}

func (m *Move) isCastle() bool {
        return m.isKingSideCastle() || m.isQueenSideCastle()
}

func (m *Move) isEnpassant(opponentPawns Bitmask) bool {
        color := m.piece.color()

        if m.piece.isPawn() && Row(m.from) == [2]int{1,6}[color] && Row(m.to) == [2]int{3,4}[color] {
                switch col := Col(m.to); col {
                case 0:
                        return opponentPawns.IsSet(m.to + 1)
                case 7:
                        return opponentPawns.IsSet(m.to - 1)
                default:
                        return opponentPawns.IsSet(m.to + 1) || opponentPawns.IsSet(m.to - 1)
                }
        }
        return false
}

func (m *Move) isEnpassantCapture(enpassant int) bool {
        return m.piece.isPawn() && m.to == enpassant
}

func (m *Move) String() string {

        if !m.isCastle() {
                col := [2]int{ Col(m.from) + 'a', Col(m.to) + 'a' }
                row := [2]int{ Row(m.from) + 1, Row(m.to) + 1 }

                capture := '-'
                if m.captured != 0 {
                        capture = 'x'
                }
                piece, promoted := m.piece.String(), m.promoted.String()
                format := `%c%d%c%c%d%s`

                if m.piece.isPawn() { // Skip piece name if it's a pawn.
                        return fmt.Sprintf(format, col[0], row[0], capture, col[1], row[1], promoted)
                } else {
                        if Settings.Fancy { // Fancy notation is more readable with extra space.
                                format = `%s ` + format
                        } else {
                                format = `%s` + format
                        }
                        return fmt.Sprintf(format, piece, col[0], row[0], capture, col[1], row[1], promoted)
                }
        } else if m.isKingSideCastle() {
                return `0-0`
        } else {
                return `0-0-0`
        }
}