// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
	`bytes`
	`regexp`
)

const (
	isCapture   = 0x00F00000
	isPromo     = 0x0F000000
	isCastle    = 0x10000000
	isEnpassant = 0x20000000
)

// Bits 00:00:00:FF => Source square (0 .. 63).
// Bits 00:00:FF:00 => Destination square (0 .. 63).
// Bits 00:0F:00:00 => Piece making the move.
// Bits 00:F0:00:00 => Captured piece if any.
// Bits 0F:00:00:00 => Promoted piece if any.
// Bits F0:00:00:00 => Castle and en-passant flags.
type Move uint32

func NewMove(p *Position, from, to int) Move {
	piece, capture := p.pieces[from], p.pieces[to]

	if p.enpassant != 0 && to == p.enpassant {
		capture = pawn(piece.color() ^ 1)
	}

	return Move(from | (to << 8) | (int(piece) << 16) | (int(capture) << 20))
}

func NewPawnMove(p *Position, square, target int) Move {
	if Abs(square - target) == 16 {

		// Check if pawn jump causes en-passant. This is done by verifying
		// whether enemy pawns occupy squares ajacent to the target square.
		pawns := p.outposts[pawn(p.color ^ 1)]
		if pawns & maskIsolated[Col(target)] & maskRank[Row(target)] != 0 {
			return NewEnpassant(p, square, target)
		}
	}

	return NewMove(p, square, target)
}

func NewEnpassant(p *Position, from, to int) Move {
	return Move(from | (to << 8) | (int(p.pieces[from]) << 16) | isEnpassant)
}

func NewCastle(p *Position, from, to int) Move {
	return Move(from | (to << 8) | (int(p.pieces[from]) << 16) | isCastle)
}

func NewPromotion(p *Position, square, target int) (Move, Move, Move, Move) {
	return NewMove(p, square, target).promote(Queen),
	       NewMove(p, square, target).promote(Rook),
	       NewMove(p, square, target).promote(Bishop),
	       NewMove(p, square, target).promote(Knight)
}

// Decodes a string in long algebraic notation and returns a move.
func NewMoveFromString(p *Position, e2e4 string) (move Move) {
	re := regexp.MustCompile(`([KkQqRrBbNn]?)([a-h])([1-8])-?([a-h])([1-8])([QqRrBbNn]?)`)
	arr := re.FindStringSubmatch(e2e4)

	if len(arr) > 0 {
		name := arr[1]
		from := Square(int(arr[3][0]-'1'), int(arr[2][0]-'a'))
		to := Square(int(arr[5][0]-'1'), int(arr[4][0]-'a'))
		promo := arr[6]

		var piece Piece
		switch name {
		case `K`, `k`:
			piece = king(p.color)
		case `Q`, `q`:
			piece = queen(p.color)
		case `R`, `r`:
			piece = rook(p.color)
		case `B`, `b`:
			piece = bishop(p.color)
		case `N`, `n`:
			piece = knight(p.color)
		default:
			piece = p.pieces[from] // <-- Makes piece character optional.
		}
		if (p.pieces[from] != piece) || (p.targets(from)&bit[to] == 0) {
			move = 0 // Invalid move.
		} else {
			move = NewMove(p, from, to)
			if len(promo) > 0 {
				switch promo {
				case `Q`, `q`:
					move = move.promote(Queen)
				case `R`, `r`:
					move = move.promote(Rook)
				case `B`, `b`:
					move = move.promote(Bishop)
				case `N`, `n`:
					move = move.promote(Knight)
				default:
					move = 0
				}
			}
		}
	} else if e2e4 == `0-0` || e2e4 == `0-0-0` {
		from := p.king[p.color]
		to := G1
		if e2e4 == `0-0-0` {
			to = C1
		}
		if p.color == Black {
			to += 56
		}
		move = NewCastle(p, from, to)
		if !move.isCastle() {
			move = 0
		}
	}
	return
}

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

// Capture value based on most valueable victim/least valueable attacker.
func (m Move) value() int {
	return pieceValue[m.capture()] - m.piece().kind()
}

func (m Move) isCastle() bool {
	return m & isCastle != 0
}

func (m Move) isEnpassant() bool {
	return m & isEnpassant != 0
}

// Returns true if the move doesn't change material balance.
func (m Move) isQuiet() bool {
	return m & (isCapture | isPromo) == 0
}

// Returns string representation of the move in long coordinate notation as
// expected by UCI, ex. `g1f3`, `e4d5` or `h7h8q`.
func (m Move) notation() string {
	var buffer bytes.Buffer

	from, to, _, _ := m.split()
	buffer.WriteByte(byte(Col(from)) + 'a')
	buffer.WriteByte(byte(Row(from)) + '1')
	buffer.WriteByte(byte(Col(to)) + 'a')
	buffer.WriteByte(byte(Row(to)) + '1')
	if m & isPromo != 0 {
		buffer.WriteByte(m.promo().char() + 32)
	}

	return buffer.String()
}

// By default the move is represented in long algebraic notation, ex. `Ng1-f3`,
// `e4xd5` or `h7-h8Q`. This is used in tests, REPL, and when displaying
// principal variation.
func (m Move) String() (str string) {
	var buffer bytes.Buffer

	from, to, piece, capture := m.split()
	if m.isCastle() {
		if to > from {
			return `0-0`
		}
		return `0-0-0`
	}

	if !piece.isPawn() {
		if Settings.Fancy { // Figurine notation is more readable with extra space.
			buffer.WriteString(piece.String() + ` `)
		} else {
			buffer.WriteByte(piece.char())
		}
	}
	buffer.WriteByte(byte(Col(from)) + 'a')
	buffer.WriteByte(byte(Row(from)) + '1')
	if capture == 0 {
		buffer.WriteByte('-')
	} else {
		buffer.WriteByte('x')
	}
	buffer.WriteByte(byte(Col(to)) + 'a')
	buffer.WriteByte(byte(Row(to)) + '1')
	if m & isPromo != 0 {
		buffer.WriteByte(m.promo().char())
	}

	return buffer.String()
}
