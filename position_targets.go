// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

// Returns a bitmask of possible Bishop moves from the given square.
func (p *Position) bishopMoves(sq Square) Bitmask {
	return p.bishopMovesAt(sq, p.board)
}

// Ditto for Rook.
func (p *Position) rookMoves(sq Square) Bitmask {
	return p.rookMovesAt(sq, p.board)
}

// Ditto for Queen.
func (p *Position) queenMoves(sq Square) Bitmask {
	return p.bishopMovesAt(sq, p.board) | p.rookMovesAt(sq, p.board)
}

// Returns a bitmask of possible Bishop moves from the given square whereas
// other pieces on the board are represented by the explicit parameter.
func (p *Position) bishopMovesAt(sq Square, board Bitmask) Bitmask {
	magic := ((bishopMagic[sq].mask & board) * bishopMagic[sq].magic) >> 55
	return bishopMagicMoves[sq][magic]
}

// Ditto for Rook.
func (p *Position) rookMovesAt(sq Square, board Bitmask) Bitmask {
	magic := ((rookMagic[sq].mask & board) * rookMagic[sq].magic) >> 52
	return rookMagicMoves[sq][magic]
}

// Ditto for Queen.
func (p *Position) queenMovesAt(sq Square, board Bitmask) Bitmask {
	return p.bishopMovesAt(sq, board) | p.rookMovesAt(sq, board)
}

func (p *Position) targets(sq Square) (bitmask Bitmask) {
	piece := p.pieces[sq]
	our := piece.color(); their := our^1

	if piece.pawnʔ() {
		// Start with one square push, then try the second square.
		empty := ^p.board
		bitmask  = bit(sq).up(our) & empty
		bitmask |= bitmask.up(our) & empty & maskRank[A4H4 + our]
		bitmask |= pawnAttacks[our][sq] & p.outposts[their]

		// If the last move set the en-passant square and it is diagonally adjacent
		// to the current pawn, then add en-passant to the pawn's attack targets.
		if p.enpassant != 0 && maskPawn[our][p.enpassant].onʔ(sq) {
			bitmask.set(p.enpassant)
		}
	} else {
		bitmask = p.attacksFor(sq, piece) & ^p.outposts[our]
	}

	return bitmask
}

func (p *Position) attacks(sq Square) Bitmask {
	return p.attacksFor(sq, p.pieces[sq])
}

func (p *Position) attacksFor(sq Square, piece Piece) (bitmask Bitmask) {
	switch piece.kind() {
	case Pawn:
		return pawnAttacks[piece.color()][sq]
	case Knight:
		return knightMoves[sq]
	case Bishop:
		return p.bishopMoves(sq)
	case Rook:
		return p.rookMoves(sq)
	case Queen:
		return p.queenMoves(sq)
	case King:
		return kingMoves[sq]
	}

	return bitmask
}

func (p *Position) xrayAttacks(sq Square) Bitmask {
	return p.xrayAttacksFor(sq, p.pieces[sq])
}

func (p *Position) xrayAttacksFor(sq Square, piece Piece) (bitmask Bitmask) {
	switch kind, our := piece.kind(), piece.color(); kind {
	case Bishop:
		board := p.board ^ p.outposts[queen(our)]
		return p.bishopMovesAt(sq, board)
	case Rook:
		board := p.board ^ p.outposts[rook(our)] ^ p.outposts[queen(our)]
		return p.rookMovesAt(sq, board)
	}

	return p.attacksFor(sq, piece)
}

func (p *Position) allAttacks(our int) (bitmask Bitmask) {
	bitmask = p.pawnAttacks(our) | p.knightAttacks(our) | p.kingAttacks(our)

	for bm := p.outposts[bishop(our)] | p.outposts[queen(our)]; bm.anyʔ(); bm = bm.pop() {
		bitmask |= p.bishopMoves(bm.first())
	}

	for bm := p.outposts[rook(our)] | p.outposts[queen(our)]; bm.anyʔ(); bm = bm.pop() {
		bitmask |= p.rookMoves(bm.first())
	}

	return bitmask
}

// Returns a bitmask of pieces that attack given square. The resulting bitmask
// only counts pieces of requested color.
//
// This method is used in static exchange evaluation so instead of using current
// board bitmask (p.board) we pass the one that gets continuously updated during
// the evaluation.
func (p *Position) attackers(our int, sq Square, board Bitmask) (bitmask Bitmask) {
	bitmask  = knightMoves[sq] & p.outposts[knight(our)]
	bitmask |= maskPawn[our][sq] & p.outposts[pawn(our)]
	bitmask |= kingMoves[sq] & p.outposts[king(our)]
	bitmask |= p.rookMovesAt(sq, board) & (p.outposts[rook(our)] | p.outposts[queen(our)])
	bitmask |= p.bishopMovesAt(sq, board) & (p.outposts[bishop(our)] | p.outposts[queen(our)])

	return bitmask
}

func (p *Position) attackedʔ(our int, sq Square) bool {
	return p.attackedAtʔ(our, sq, p.board)
}

func (p *Position) attackedAtʔ(our int, sq Square, board Bitmask) bool {
	return (knightMoves[sq] & p.outposts[knight(our)]).anyʔ() ||
	       (maskPawn[our][sq] & p.outposts[pawn(our)]).anyʔ() ||
	       (kingMoves[sq] & p.outposts[king(our)]).anyʔ() ||
	       (p.rookMovesAt(sq, board) & (p.outposts[rook(our)] | p.outposts[queen(our)])).anyʔ() ||
	       (p.bishopMovesAt(sq, board) & (p.outposts[bishop(our)] | p.outposts[queen(our)])).anyʔ()
}

func (p *Position) pawnTargets(our int, pawns Bitmask) Bitmask {
	if our == White {
		return ((pawns & ^maskFile[0]) << 7) | ((pawns & ^maskFile[7]) << 9)
	}

	return ((pawns & ^maskFile[0]) >> 9) | ((pawns & ^maskFile[7]) >> 7)
}

func (p *Position) pawnAttacks(our int) Bitmask {
	return p.pawnTargets(our, p.outposts[pawn(our)])
}

func (p *Position) knightAttacks(our int) (bitmask Bitmask) {
	for bm := p.outposts[knight(our)]; bm.anyʔ(); bm = bm.pop() {
		bitmask |= knightMoves[bm.first()]
	}

	return bitmask
}

func (p *Position) bishopAttacks(our int) (bitmask Bitmask) {
	for bm := p.outposts[bishop(our)]; bm.anyʔ(); bm = bm.pop() {
		bitmask |= p.bishopMoves(bm.first())
	}

	return bitmask
}

func (p *Position) rookAttacks(our int) (bitmask Bitmask) {
	for bm := p.outposts[rook(our)]; bm.anyʔ(); bm = bm.pop() {
		bitmask |= p.rookMoves(bm.first())
	}

	return bitmask
}

func (p *Position) queenAttacks(our int) (bitmask Bitmask) {
	for bm := p.outposts[queen(our)]; bm.anyʔ(); bm = bm.pop() {
		bitmask |= p.queenMoves(bm.first())
	}

	return bitmask
}

func (p *Position) kingAttacks(our int) Bitmask {
	return kingMoves[p.king[our]]
}

func (p *Position) knightAttacksAt(sq Square, our int) (bitmask Bitmask) {
	return knightMoves[sq] & ^p.outposts[our]
}

func (p *Position) bishopAttacksAt(sq Square, our int) (bitmask Bitmask) {
	return p.bishopMoves(sq) & ^p.outposts[our]
}

func (p *Position) rookAttacksAt(sq Square, our int) (bitmask Bitmask) {
	return p.rookMoves(sq) & ^p.outposts[our]
}

func (p *Position) queenAttacksAt(sq Square, our int) (bitmask Bitmask) {
	return p.queenMoves(sq) & ^p.outposts[our]
}

func (p *Position) kingAttacksAt(sq Square, our int) Bitmask {
	return kingMoves[sq] & ^p.outposts[our]
}

