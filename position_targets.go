// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

// Returns a bitmask of possible Bishop moves from the given square.
func (p *Position) bishopMoves(square int) Bitmask {
	return p.bishopMovesAt(square, p.board)
}

// Returns a bitmask of possible Rook moves from the given square.
func (p *Position) rookMoves(square int) Bitmask {
	return p.rookMovesAt(square, p.board)
}

// Returns a bitmask of possible Bishop moves from the given square whereas
// other pieces on the board are represented by the explicit parameter.
func (p *Position) bishopMovesAt(square int, board Bitmask) Bitmask {
	magic := ((bishopMagic[square].mask & board) * bishopMagic[square].magic) >> 55
	return bishopMagicMoves[square][magic]
}

// Returns a bitmask of possible Rook moves from the given square whereas
// other pieces on the board are represented by the explicit parameter.
func (p *Position) rookMovesAt(square int, board Bitmask) Bitmask {
	magic := ((rookMagic[square].mask & board) * rookMagic[square].magic) >> 52
	return rookMagicMoves[square][magic]
}

func (p *Position) targets(square int) Bitmask {
	return p.targetsFor(square, p.pieces[square])
}

func (p *Position) targetsFor(square int, piece Piece) (bitmask Bitmask) {
	color := piece.color()
	if piece.isPawn() {
		// Start with one square push, then try the second square.
		empty := ^p.board
		if color == White {
			bitmask |= (bit[square] << 8) & empty
			bitmask |= (bitmask << 8) & empty & maskRank[3]
		} else {
			bitmask |= (bit[square] >> 8) & empty
			bitmask |= (bitmask >> 8) & empty & maskRank[4]
		}
		bitmask |= pawnAttacks[color][square] & p.outposts[color^1]

		// If the last move set the en-passant square and it is diagonally adjacent
		// to the current pawn, then add en-passant to the pawn's attack targets.
		if p.enpassant != 0 && maskPawn[color][p.enpassant].on(square) {
			bitmask |= bit[p.enpassant]
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
	switch kind, color := piece.kind(), piece.color(); kind {
	case Pawn:
		return pawnAttacks[color][square]
	case Knight:
		return knightMoves[square]
	case Bishop:
		return p.bishopMoves(square)
	case Rook:
		return p.rookMoves(square)
	case Queen:
		return p.bishopMoves(square) | p.rookMoves(square)
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

func (p *Position) allAttacks(color uint8) (bitmask Bitmask) {
	bitmask = p.pawnAttacks(color) | p.knightAttacks(color) | p.kingAttacks(color)

	outposts := p.outposts[bishop(color)] | p.outposts[queen(color)]
	for outposts.any() {
		bitmask |= p.bishopMoves(outposts.pop())
	}

	outposts = p.outposts[rook(color)] | p.outposts[queen(color)]
	for outposts.any() {
		bitmask |= p.rookMoves(outposts.pop())
	}

	return bitmask
}

// Returns a bitmask of pieces that attack given square. The resulting bitmask
// only counts pieces of requested color.
//
// This method is used in static exchange evaluation so instead of using current
// board bitmask (p.board) we pass the one that gets continuously updated during
// the evaluation.
func (p *Position) attackers(color uint8, square int, board Bitmask) (bitmask Bitmask) {
	bitmask  = knightMoves[square] & p.outposts[knight(color)]
	bitmask |= maskPawn[color][square] & p.outposts[pawn(color)]
	bitmask |= kingMoves[square] & p.outposts[king(color)]
	bitmask |= p.rookMovesAt(square, board) & (p.outposts[rook(color)] | p.outposts[queen(color)])
	bitmask |= p.bishopMovesAt(square, board) & (p.outposts[bishop(color)] | p.outposts[queen(color)])

	return bitmask
}

func (p *Position) isAttacked(color uint8, square int) bool {
	return (knightMoves[square] & p.outposts[knight(color)]).any() ||
	       (maskPawn[color][square] & p.outposts[pawn(color)]).any() ||
	       (kingMoves[square] & p.outposts[king(color)]).any() ||
	       (p.rookMoves(square) & (p.outposts[rook(color)] | p.outposts[queen(color)])).any() ||
	       (p.bishopMoves(square) & (p.outposts[bishop(color)] | p.outposts[queen(color)])).any()
}

func (p *Position) pawnTargets(color uint8, pawns Bitmask) Bitmask {
	if color == White {
		return ((pawns & ^maskFile[0]) << 7) | ((pawns & ^maskFile[7]) << 9)
	}

	return ((pawns & ^maskFile[0]) >> 9) | ((pawns & ^maskFile[7]) >> 7)
}

func (p *Position) pawnAttacks(color uint8) Bitmask {
	return p.pawnTargets(color, p.outposts[pawn(color)])
}

func (p *Position) knightAttacks(color uint8) (bitmask Bitmask) {
	outposts := p.outposts[knight(color)]
	for outposts.any() {
		bitmask |= knightMoves[outposts.pop()]
	}

	return bitmask
}

func (p *Position) bishopAttacks(color uint8) (bitmask Bitmask) {
	outposts := p.outposts[bishop(color)]
	for outposts.any() {
		bitmask |= p.bishopMoves(outposts.pop())
	}

	return bitmask
}

func (p *Position) rookAttacks(color uint8) (bitmask Bitmask) {
	outposts := p.outposts[rook(color)]
	for outposts.any() {
		bitmask |= p.rookMoves(outposts.pop())
	}

	return bitmask
}

func (p *Position) queenAttacks(color uint8) (bitmask Bitmask) {
	outposts := p.outposts[queen(color)]
	for outposts.any() {
		square := outposts.pop()
		bitmask |= p.rookMoves(square) | p.bishopMoves(square)
	}

	return bitmask
}

func (p *Position) kingAttacks(color uint8) Bitmask {
	return kingMoves[p.king[color]]
}
