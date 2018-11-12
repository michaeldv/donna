// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

// Returns a bitmask of possible Bishop moves from the given square.
func (p *Position) bishopMoves(square int) Bitmask {
	return p.bishopMovesAt(square, p.board)
}

// Ditto for Rook.
func (p *Position) rookMoves(square int) Bitmask {
	return p.rookMovesAt(square, p.board)
}

// Ditto for Queen.
func (p *Position) queenMoves(square int) Bitmask {
	return p.bishopMovesAt(square, p.board) | p.rookMovesAt(square, p.board)
}

// Returns a bitmask of possible Bishop moves from the given square whereas
// other pieces on the board are represented by the explicit parameter.
func (p *Position) bishopMovesAt(square int, board Bitmask) Bitmask {
	magic := ((bishopMagic[square].mask & board) * bishopMagic[square].magic) >> 55
	return bishopMagicMoves[square][magic]
}

// Ditto for Rook.
func (p *Position) rookMovesAt(square int, board Bitmask) Bitmask {
	magic := ((rookMagic[square].mask & board) * rookMagic[square].magic) >> 52
	return rookMagicMoves[square][magic]
}

// Ditto for Queen.
func (p *Position) queenMovesAt(square int, board Bitmask) Bitmask {
	return p.bishopMovesAt(square, board) | p.rookMovesAt(square, board)
}

func (p *Position) targets(square int) (bitmask Bitmask) {
	piece := p.pieces[square]
	our := piece.color()
	their := our^1

	if piece.pawnʔ() {
		// Start with one square push, then try the second square.
		empty := ^p.board
		bitmask  = bit(square).up(our) & empty
		bitmask |= bitmask.up(our) & empty & maskRank[A4H4 + our]
		bitmask |= pawnAttacks[our&1][square] & p.pick(their).all

		// If the last move set the en-passant square and it is diagonally adjacent
		// to the current pawn, then add en-passant to the pawn's attack targets.
		if p.enpassant != 0 && maskPawn[our&1][p.enpassant].onʔ(square) {
			bitmask.set(p.enpassant)
		}
	} else {
		bitmask = p.attacksFor(square, piece) & ^p.pick(our).all
	}

	return bitmask
}

func (p *Position) attacks(square int) Bitmask {
	return p.attacksFor(square, p.pieces[square])
}

func (p *Position) attacksFor(square int, piece Piece) (bitmask Bitmask) {
	switch piece.kind() {
	case Pawn:
		return pawnAttacks[piece.color()&1][square]
	case Knight:
		return knightMoves[square]
	case Bishop:
		return p.bishopMoves(square)
	case Rook:
		return p.rookMoves(square)
	case Queen:
		return p.queenMoves(square)
	case King:
		return kingMoves[square]
	}

	return bitmask
}

func (p *Position) xrayAttacks(square int) Bitmask {
	return p.xrayAttacksFor(square, p.pieces[square])
}

func (p *Position) xrayAttacksFor(square int, piece Piece) (bitmask Bitmask) {
	switch kind, our := piece.kind(), piece.color(); kind {
	case Bishop:
		board := p.board ^ p.pick(our).queens
		return p.bishopMovesAt(square, board)
	case Rook:
		board := p.board ^ p.pick(our).rooks ^ p.pick(our).queens
		return p.rookMovesAt(square, board)
	}

	return p.attacksFor(square, piece)
}

func (p *Position) allAttacks(our int) (bitmask Bitmask) {
	side := p.pick(our)
	bitmask = p.pawnAttacks(our) | p.knightAttacks(our) | p.kingAttacks(our)

	for bm := side.bishops | side.queens; bm.anyʔ(); bm = bm.pop() {
		bitmask |= p.bishopMoves(bm.first())
	}

	for bm := side.rooks | side.queens; bm.anyʔ(); bm = bm.pop() {
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
func (p *Position) attackers(our int, square int, board Bitmask) (bitmask Bitmask) {
	side := p.pick(our)
	bitmask  = knightMoves[square] & side.knights
	bitmask |= maskPawn[our&1][square] & side.pawns
	bitmask |= kingMoves[square] & side.king
	bitmask |= p.rookMovesAt(square, board) & (side.rooks | side.queens)
	bitmask |= p.bishopMovesAt(square, board) & (side.bishops | side.queens)
	return bitmask
}

func (p *Position) attackedʔ(our int, square int) bool {
	side := p.pick(our)
	return (maskPawn[our&1][square] & side.pawns).anyʔ() ||
	       (knightMoves[square] & side.knights).anyʔ() ||
	       (kingMoves[square] & side.king).anyʔ() ||
	       (p.rookMoves(square) & (side.rooks | side.queens)).anyʔ() ||
	       (p.bishopMoves(square) & (side.bishops | side.queens)).anyʔ()

}

func (p *Position) pawnTargets(our int, pawns Bitmask) Bitmask {
	if our == White {
		return ((pawns & ^maskFile[0]) << 7) | ((pawns & ^maskFile[7]) << 9)
	}

	return ((pawns & ^maskFile[0]) >> 9) | ((pawns & ^maskFile[7]) >> 7)
}

func (p *Position) pawnAttacks(our int) Bitmask {
	return p.pawnTargets(our, p.pick(our).pawns)
}

func (p *Position) knightAttacks(our int) (bitmask Bitmask) {
	for bm := p.pick(our).knights; bm.anyʔ(); bm = bm.pop() {
		bitmask |= knightMoves[bm.first()]
	}

	return bitmask
}

func (p *Position) bishopAttacks(our int) (bitmask Bitmask) {
	for bm := p.pick(our).bishops; bm.anyʔ(); bm = bm.pop() {
		bitmask |= p.bishopMoves(bm.first())
	}

	return bitmask
}

func (p *Position) rookAttacks(our int) (bitmask Bitmask) {
	for bm := p.pick(our).rooks; bm.anyʔ(); bm = bm.pop() {
		bitmask |= p.rookMoves(bm.first())
	}

	return bitmask
}

func (p *Position) queenAttacks(our int) (bitmask Bitmask) {
	for bm := p.pick(our).queens; bm.anyʔ(); bm = bm.pop() {
		bitmask |= p.queenMoves(bm.first())
	}

	return bitmask
}

func (p *Position) kingAttacks(our int) Bitmask {
	return kingMoves[p.pick(our).home]
}

func (p *Position) knightAttacksAt(square int, our int) (bitmask Bitmask) {
	return knightMoves[square] & ^p.pick(our).all
}

func (p *Position) bishopAttacksAt(square int, our int) (bitmask Bitmask) {
	return p.bishopMoves(square) & ^p.pick(our).all
}

func (p *Position) rookAttacksAt(square int, our int) (bitmask Bitmask) {
	return p.rookMoves(square) & ^p.pick(our).all
}

func (p *Position) queenAttacksAt(square int, our int) (bitmask Bitmask) {
	return p.queenMoves(square) & ^p.pick(our).all
}

func (p *Position) kingAttacksAt(square int, our int) Bitmask {
	return kingMoves[square] & ^p.pick(our).all
}

