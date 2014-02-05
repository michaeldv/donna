// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`fmt`)

const (
        isCastle    = 0x10000000
        isEnpassant = 0x20000000
        isPawnJump  = 0x40000000
)

// Bits 00:00:00:FF => Dource square (0 .. 63).
// Bits 00:00:FF:00 => Destination square (0 .. 63).
// Bits 00:0F:00:00 => Piece making the move.
// Bits 00:F0:00:00 => Captured piece if any.
// Bits 0F:00:00:00 => Promoted piece if any.
// Bits F0:00:00:00 => Castle, en-passant, or pawn jump flags.
type Move uint32

func (m Move) from() int {
        return int(m & 0xFF)
}

func (m Move) to() int {
        return int((m >> 8) & 0xFF)
}

func (m Move) piece() Piece {
        return Piece((m >> 16) & 0x0F)
}

func (m Move) color() int {
        return int((m >> 16) & 1)
}

func (m Move) capture() Piece {
        return Piece((m >> 20) & 0x0F)
}

func (m Move) split() (from, to int, piece, capture Piece) {
        return int(m & 0xFF), int((m >> 8) & 0xFF), Piece((m >> 16) & 0x0F), Piece((m >> 20) & 0x0F)
}

func (m Move) promo() Piece {
        return Piece((m >> 24) & 0x0F)
}

func (m Move) promote(kind int) Move {
        piece := Piece(kind | m.color())
        return m | Move(int(piece) << 24)
}

func (m Move) izCastle() bool {
        return m & isCastle != 0
}

func (m Move) castle() Move {
        return m | isCastle
}

func (m Move) izEnpassant() bool {
        return m & isEnpassant != 0
}

func (m Move) enpassant() Move {
        return m | isEnpassant
}

func (m Move) izPawnJump() bool {
        return m & isPawnJump != 0
}

func (m Move) pawnJump() Move {
        return m | isPawnJump
}

// PxQ, NxQ, BxQ, RxQ, QxQ, KxQ => where => QUEEN  = 5 << 1 // 10
// PxR, NxR, BxR, RxR, QxR, KxR             ROOK   = 4 << 1 // 8
// PxB, NxB, BxB, RxB, QxB, KxB             BISHOP = 3 << 1 // 6
// PxN, NxN, BxN, RxN, QxN, KxN             KNIGHT = 2 << 1 // 4
// PxP, NxP, BxP, RxP, QxP, KxP             PAWN   = 1 << 1 // 2
func (m Move) value() int {
        capture := m.capture()
        if capture == 0 || capture.isKing() {
                return 0
        }

        victim := (QUEEN - capture.kind()) / PAWN
        attacker := m.piece().kind() / PAWN - 1

        return victimAttacker[victim][attacker]
}

func (m Move) is00() bool {
        from, to, piece, _ := m.split()
        return (piece == King(White) && from == E1 && to == G1) || (piece == King(Black) && from == E8 && to == G8)
}

func (m Move) is000() bool {
        from, to, piece, _ := m.split()
        return (piece == King(White) && from == E1 && to == C1) || (piece == King(Black) && from == E8 && to == C8)
}

func (m Move) isCastle() bool {
        return m.is00() || m.is000()
}

func (m Move) String() string {
        from, to, piece, capture := m.split()
        promo := m.promo().s()

        if !m.isCastle() {
                col := [2]int{ Col(from) + 'a', Col(to) + 'a' }
                row := [2]int{ Row(from) + 1, Row(to) + 1 }

                sign := '-'
                if capture != 0 || (piece.isPawn() && Col(from) != Col(to)) {
                        sign = 'x'
                }

                format := `%c%d%c%c%d%s`
                if piece.isPawn() { // Skip piece name if it's a pawn.
                        return fmt.Sprintf(format, col[0], row[0], sign, col[1], row[1], promo)
                } else {
                        if Settings.Fancy {
                                // Fancy notation is more readable with extra space.
                                return fmt.Sprintf(`%s ` + format, piece, col[0], row[0], sign, col[1], row[1], promo)
                        } else {
                                // Use uppercase letter to representa a piece regardless of its color.
                                return fmt.Sprintf(`%s` + format, piece.s(), col[0], row[0], sign, col[1], row[1], promo)
                        }
                }
        } else if m.is00() {
                return `0-0`
        } else {
                return `0-0-0`
        }
}