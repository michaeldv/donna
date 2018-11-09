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
	color := piece.color()
	if piece.pawnʔ() {
		// Start with one square push, then try the second square.
		empty := ^p.board
		bitmask  = bit[square].up(color) & empty
		bitmask |= bitmask.up(color) & empty & maskRank[A4H4 + color]
		bitmask |= pawnAttacks[color][square] & p.outposts[color^1]

		// If the last move set the en-passant square and it is diagonally adjacent
		// to the current pawn, then add en-passant to the pawn's attack targets.
		if p.enpassant != 0 && maskPawn[color][p.enpassant].onʔ(square) {
			bitmask.set(p.enpassant)
		}
	} else {
		bitmask = p.attacksFor(square, piece) & ^p.outposts[color]
	}

	return bitmask
}

func (p *Position) attacks(square int) Bitmask {
	return p.attacksFor(square, p.pieces[square])
}

func (p *Position) attacksFor(square int, piece Piece) (bitmask Bitmask) {
	switch piece.kind() {
	case Pawn:
		return pawnAttacks[piece.color()][square]
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
	switch kind, color := piece.kind(), piece.color(); kind {
	case Bishop:
		board := p.board ^ p.outposts[queen(color)]
		return p.bishopMovesAt(square, board)
	case Rook:
		board := p.board ^ p.outposts[rook(color)] ^ p.outposts[queen(color)]
		return p.rookMovesAt(square, board)
	}

	return p.attacksFor(square, piece)
}

func (p *Position) allAttacks(color int) (bitmask Bitmask) {
	bitmask = p.pawnAttacks(color) | p.knightAttacks(color) | p.kingAttacks(color)

	for bm := p.outposts[bishop(color)] | p.outposts[queen(color)]; bm.anyʔ(); bm = bm.pop() {
		bitmask |= p.bishopMoves(bm.first())
	}

	for bm := p.outposts[rook(color)] | p.outposts[queen(color)]; bm.anyʔ(); bm = bm.pop() {
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
func (p *Position) attackers(color int, square int, board Bitmask) (bitmask Bitmask) {
	bitmask  = knightMoves[square] & p.outposts[knight(color)]
	bitmask |= maskPawn[color][square] & p.outposts[pawn(color)]
	bitmask |= kingMoves[square] & p.outposts[king(color)]
	bitmask |= p.rookMovesAt(square, board) & (p.outposts[rook(color)] | p.outposts[queen(color)])
	bitmask |= p.bishopMovesAt(square, board) & (p.outposts[bishop(color)] | p.outposts[queen(color)])

	return bitmask
}

func (p *Position) attackedʔ(color int, square int) bool {
	return (knightMoves[square] & p.outposts[knight(color)]).anyʔ() ||
	       (maskPawn[color][square] & p.outposts[pawn(color)]).anyʔ() ||
	       (kingMoves[square] & p.outposts[king(color)]).anyʔ() ||
	       (p.rookMoves(square) & (p.outposts[rook(color)] | p.outposts[queen(color)])).anyʔ() ||
	       (p.bishopMoves(square) & (p.outposts[bishop(color)] | p.outposts[queen(color)])).anyʔ()
}

func (p *Position) pawnTargets(color int, pawns Bitmask) Bitmask {
	if color == White {
		return ((pawns & ^maskFile[0]) << 7) | ((pawns & ^maskFile[7]) << 9)
	}

	return ((pawns & ^maskFile[0]) >> 9) | ((pawns & ^maskFile[7]) >> 7)
}

func (p *Position) pawnAttacks(color int) Bitmask {
	return p.pawnTargets(color, p.outposts[pawn(color)])
}

func (p *Position) knightAttacks(color int) (bitmask Bitmask) {
	for bm := p.outposts[knight(color)]; bm.anyʔ(); bm = bm.pop() {
		bitmask |= knightMoves[bm.first()]
	}

	return bitmask
}

func (p *Position) bishopAttacks(color int) (bitmask Bitmask) {
	for bm := p.outposts[bishop(color)]; bm.anyʔ(); bm = bm.pop() {
		bitmask |= p.bishopMoves(bm.first())
	}

	return bitmask
}

func (p *Position) rookAttacks(color int) (bitmask Bitmask) {
	for bm := p.outposts[rook(color)]; bm.anyʔ(); bm = bm.pop() {
		bitmask |= p.rookMoves(bm.first())
	}

	return bitmask
}

func (p *Position) queenAttacks(color int) (bitmask Bitmask) {
	for bm := p.outposts[queen(color)]; bm.anyʔ(); bm = bm.pop() {
		bitmask |= p.queenMoves(bm.first())
	}

	return bitmask
}

func (p *Position) kingAttacks(color int) Bitmask {
	return kingMoves[p.king[color]]
}

func (p *Position) knightAttacksAt(square int, color int) (bitmask Bitmask) {
	return knightMoves[square] & ^p.outposts[color]
}

func (p *Position) bishopAttacksAt(square int, color int) (bitmask Bitmask) {
	return p.bishopMoves(square) & ^p.outposts[color]
}

func (p *Position) rookAttacksAt(square int, color int) (bitmask Bitmask) {
	return p.rookMoves(square) & ^p.outposts[color]
}

func (p *Position) queenAttacksAt(square int, color int) (bitmask Bitmask) {
	return p.queenMoves(square) & ^p.outposts[color]
}

func (p *Position) kingAttacksAt(square int, color int) Bitmask {
	return kingMoves[square] & ^p.outposts[color]
}

