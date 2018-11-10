// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

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

// Moving pianos is dangerous. Moving pianos are dangerous.
func NewMove(p *Position, from, to int) Move {
	piece, capture := p.pieces[from], p.pieces[to]

	if p.enpassant != 0 && to == p.enpassant && piece.pawnʔ() {
		capture = pawn(piece.color()^1)
	}

	return Move(from | (to << 8) | (int(piece) << 16) | (int(capture) << 20))
}

func NewPawnMove(p *Position, square, target int) Move {
	if abs(square - target) == 16 {

		// Check if pawn jump causes en-passant. This is done by verifying
		// whether enemy pawns occupy squares ajacent to the target square.
		pawns := p.outposts[pawn(p.color^1)]
		if pawns & maskIsolated[col(target)] & maskRank[row(target)] != 0 {
			return NewEnpassant(p, square, target)
		}
	}

	return NewMove(p, square, target)
}

func NewEnpassant(p *Position, from, to int) Move {
	return Move(from | (to << 8) | (int(p.pieces[from] << 16)) | isEnpassant)
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

// Decodes a string in coordinate notation and returns a move. The string is
// expected to be either 4 or 5 characters long (with promotion).
func NewMoveFromNotation(p *Position, e2e4 string) Move {
	from := square(int(e2e4[1] - '1'), int(e2e4[0] - 'a'))
	to := square(int(e2e4[3] - '1'), int(e2e4[2] - 'a'))

	// Check if this is a castle.
	if p.pieces[from].kingʔ() && abs(from - to) == 2 {
		return NewCastle(p, from, to)
	}

	// Special handling for pawn pushes because they might cause en-passant
	// and result in promotion.
	if p.pieces[from].pawnʔ() {
		move := NewPawnMove(p, from, to)
		if len(e2e4) > 4 {
			switch e2e4[4] {
			case 'q', 'Q':
				move = move.promote(Queen)
			case 'r', 'R':
				move = move.promote(Rook)
			case 'b', 'B':
				move = move.promote(Bishop)
			case 'n', 'N':
				move = move.promote(Knight)
			}
		}
		return move
	}

	return NewMove(p, from, to)
}

// Decodes a string in long algebraic notation and returns a move. All invalid
// moves are discarded and returned as Move(0).
func NewMoveFromString(p *Position, e2e4 string) (move Move, validMoves []Move) {
	re := regexp.MustCompile(`([KkQqRrBbNn]?)([a-h])([1-8])[-x]?([a-h])([1-8])([QqRrBbNn]?)\+?[!\?]{0,2}`)
	matches := re.FindStringSubmatch(e2e4)

	// Before returning the move make sure it is valid in current position.
	defer func() {
		gen := NewMoveGen(p).generateAllMoves().validOnly()
		validMoves = gen.allMoves()
		if move.someʔ() && !gen.amongValidʔ(move) {
			move = Move(0)
		}
	}()

	if len(matches) == 7 { // Full regex match.
		if letter := matches[1]; letter != `` {
			var piece Piece

			// Validate optional piece character to make sure the actual piece it
			// represents is there.
			switch letter {
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
			}
			square := square(int(matches[3][0] - '1'), int(matches[2][0] - 'a'))
			if p.pieces[square] != piece {
				move = Move(0)
				return
			}
		}
		move = NewMoveFromNotation(p, matches[2] + matches[3] + matches[4] + matches[5] + matches[6])
		return
	}

	// Special castle move notation.
	if e2e4 == `0-0` || e2e4 == `0-0-0` {
		kingside, queenside := p.canCastleʔ(p.color)
		if e2e4 == `0-0` && kingside {
			from, to := p.king[p.color&1], G1 + p.color * A8
			move = NewCastle(p, from, to)
			return
		}
		if e2e4 == `0-0-0` && queenside {
			from, to := p.king[p.color&1], C1 + p.color * A8
			move = NewCastle(p, from, to)
			return
		}
	}
	return
}

func (m Move) nullʔ() bool {
	return m == Move(0)
}

func (m Move) someʔ() bool {
	return m != Move(0)
}

func (m Move) from() int {
	return int(m & 0x3F)
}

func (m Move) to() int {
	return int((m >> 8) & 0x3F)
}

func (m Move) piece() Piece {
	return Piece((m >> 16) & 0x0F)
}

func (m Move) color() int {
	return int(m >> 16) & 1
}

func (m Move) capture() Piece {
	return Piece((m >> 20) & 0x0F)
}

func (m Move) split() (from, to int, piece, capture Piece) {
	return int(m & 0x3F), int((m >> 8) & 0x3F), Piece((m >> 16) & 0x0F), Piece((m >> 20) & 0x0F)
}

func (m Move) promo() Piece {
	return Piece((m >> 24) & 0x0F)
}

func (m Move) promote(kind int) Move {
	piece := Piece(kind | m.color())
	return m | Move(piece << 24)
}

// Capture value based on most valueable victim/least valueable attacker.
func (m Move) value() (value int) {
	value = pieceValue[m.capture().id()] - int(m.piece())
	if m.enpassantʔ() {
		value += valuePawn.midgame
	} else if m.promoʔ() {
		value += pieceValue[m.promo().id()] - valuePawn.midgame
	}
	return
}

func (m Move) castleʔ() bool {
	return m & isCastle != 0
}

func (m Move) captureʔ() bool {
	return m & isCapture != 0
}

func (m Move) enpassantʔ() bool {
	return m & isEnpassant != 0
}

func (m Move) promoʔ() bool {
	return m & isPromo != 0
}

// Returns true if the move doesn't change material balance.
func (m Move) quietʔ() bool {
	return m & (isCapture | isPromo) == 0 // | isEnpassant) == 0
}

// Returns true for pawn pushes beyond home half of the board.
func (m Move) pawnAdvanceʔ() bool {
	return m.piece().pawnʔ() && rank(m.color(), m.to()) > A4H4
}

// Returns true is the move is one of the killer moves at given ply.
func (m Move) killerʔ(ply int) bool {
	return m.someʔ() && (m == game.killers[ply][0] || m == game.killers[ply][1])
}

// Returns true if *non-evasion* move is valid, i.e. it is possible to make
// the move in current position without violating chess rules.
//
// If the king is in check move generator is expected to generate valid evasions
// where extra validation is not needed.
func (m Move) validʔ(p *Position, pins Bitmask) bool {
	our := m.color(); their := our^1
	from, to, piece, capture := m.split()

	// For rare en-passant pawn captures we validate the move by actually
	// making it, and then taking it back.
	if p.enpassant != 0 && to == p.enpassant && capture.pawnʔ() {
		position := p.makeMove(m)
		defer position.undoLastMove()
		return !position.inCheckʔ(our)
	}

	// King's move is valid when a) the move is a castle or b) the destination
	// square is not being attacked by the opponent.
	if piece.kingʔ() {
		return m.castleʔ() || !p.attackedʔ(their, to)
	}

	// For all other pieces the move is valid when it doesn't cause a
	// check. For pinned sliders this includes moves along the pinning
	// file, rank, or diagonal.
	return pins.noneʔ() || pins.offʔ(from) || maskLine[from][to].onʔ(p.king[our&1])
}

// Returns true if the move could have be generated on the given board. We
// use this to test whether cached or killer moves could be returned by the
// incremental move generator.
func (m Move) legitʔ(p *Position, pins Bitmask) bool {
	// from, to, piece, capture := m.split()
	our := m.color(); their := our^1
	from, to, piece, capture := m.split()

	// `from` must have a piece and `to` can't have a piece of the same color.
	if p.outposts[piece].offʔ(from) || p.outposts[our&1].onʔ(to) {
		return false
	}

	// First check pawn captures and pushes.
	if piece.pawnʔ() {
		if p.enpassant != 0 {
			return m.enpassantʔ() && p.enpassant == to
		}
		if capture.someʔ() {
			return pawnAttacks[our&1][from].onʔ(to) && p.outposts[capture].onʔ(to)
		} else {
			// If no capture then the target square should not be occupied.
			if p.board.onʔ(to) {
				return false
			}
			if row := rank(our, from); row == A7H7 && !m.promoʔ() {
				return false
			} else if push := from + up[our&1]; to == push {
				return true
			} else { // Must be pawn jump.
				return row == A1H1 && to == from + 2 * up[our&1] && p.board.offʔ(push)
			}
		}
	}

	// Anything pawn-related is now non-legit.
	if m.enpassantʔ() || m.promoʔ() {
		return false
	}

	// Captures should capture, non-captures should be free to move.
	if (capture.someʔ() && p.outposts[their&1].offʔ(to)) || (capture.noneʔ() && p.board.onʔ(to)) {
		return false
	}

	// Now check king moves including castles.
	if piece.kingʔ() {
		if m.castleʔ() {
			if from != homeKing[our&1] {
				return false
			}
			switch to {
			case G1, G8:
				if p.outposts[rook(our)].offʔ(to + 1) || (p.castles & castleKingside[our&1] == 0) || (gapKing[our&1] & p.board).anyʔ() {
					return false
				}
				return (castleKing[our&1] & p.allAttacks(their)).noneʔ()
			case C1, C8:
				if p.outposts[rook(our)].offʔ(to - 2) || (p.castles & castleQueenside[our&1] == 0) || (gapQueen[our&1] & p.board).anyʔ() {
					return false
				}
				return (castleQueen[our&1] & p.allAttacks(their)).noneʔ()
			}
			return false
		}
		return p.kingAttacksAt(from, our).onʔ(to)
	}

	// Anything castle-related is now non-legit.
	if m.castleʔ() {
		return false
	}

	// Check remaining pieces.
	switch piece.kind() {
	case Knight:
		return p.knightAttacksAt(from, our).onʔ(to)
	case Bishop:
		return p.bishopAttacksAt(from, our).onʔ(to)
	case Rook:
		return p.rookAttacksAt(from, our).onʔ(to)
	case Queen:
		return p.queenAttacksAt(from, our).onʔ(to)
	}

	return false
}

// Returns string representation of the move in long coordinate notation as
// expected by UCI, ex. `g1f3`, `e4d5` or `h7h8q`.
func (m Move) notation() string {
	var buffer bytes.Buffer

	from, to, _, _ := m.split()
	buffer.WriteByte(byte(col(from)) + 'a')
	buffer.WriteByte(byte(row(from)) + '1')
	buffer.WriteByte(byte(col(to)) + 'a')
	buffer.WriteByte(byte(row(to)) + '1')
	if m & isPromo != 0 {
		buffer.WriteByte(m.promo().char() + 32)
	}

	return buffer.String()
}

// Returns string representation of the move in long algebraic notation using
// ASCII characters only.
func (m Move) str() (str string) {
	if engine.fancyʔ {
		defer func() { engine.fancyʔ = true }()
		engine.fancyʔ = false
	}

	return m.String()
}

// By default the move is represented in long algebraic notation utilizing fancy
// UTF-8 engine setting. For example: `♘g1-f3` (fancy), `e4xd5` or `h7-h8Q`.
// This notation is used in tests, REPL, and when showing principal variation.
func (m Move) String() (str string) {
	var buffer bytes.Buffer

	from, to, piece, capture := m.split()
	if m.castleʔ() {
		if to > from {
			return `0-0`
		}
		return `0-0-0`
	}

	if !piece.pawnʔ() {
		buffer.WriteByte(piece.char())
	}
	buffer.WriteByte(byte(col(from)) + 'a')
	buffer.WriteByte(byte(row(from)) + '1')
	if capture == 0 {
		buffer.WriteByte('-')
	} else {
		buffer.WriteByte('x')
	}
	buffer.WriteByte(byte(col(to)) + 'a')
	buffer.WriteByte(byte(row(to)) + '1')
	if promo := m.promo(); promo.someʔ() {
		buffer.WriteByte(promo.char())
	}

	return buffer.String()
}
