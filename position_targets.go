// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
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

// Returns a bitmask of possible Bishop moves from the given square wherees
// other pieces on the board are represented by the explicit parameter.
func (p *Position) bishopMovesAt(square int, board Bitmask) Bitmask {
	magic := ((bishopMagic[square].mask & board) * bishopMagic[square].magic) >> 55
	return bishopMagicMoves[square][magic]
}

// Returns a bitmask of possible Rook moves from the given square wherees other
// pieces on the board are represented by the explicit parameter.
func (p *Position) rookMovesAt(square int, board Bitmask) Bitmask {
	magic := ((rookMagic[square].mask & board) * rookMagic[square].magic) >> 52
	return rookMagicMoves[square][magic]
}

func (p *Position) targets(square int) Bitmask {
	return p.targetsFor(square, p.pieces[square])
}

func (p *Position) targetsFor(square int, piece Piece) (bitmask Bitmask) {
	switch kind, color := piece.kind(), piece.color(); kind {
	case Pawn:
		bitmask = pawnMoves[color][square] & p.outposts[color^1]
		//
		// If the square in front of the pawn is empty then add it as possible
		// target.
		//
		if target := square + eight[color]; p.board.isClear(target) {
			bitmask.set(target)
			//
			// If the pawn is in its initial position and two squares in front of
			// the pawn are empty then add the second square as possible target.
			//
			if RelRow(square, color) == 1 {
				if target += eight[color]; p.board.isClear(target) {
					bitmask.set(target)
				}
			}
		}
		//
		// If the last move set the en-passant square and it is diagonally adjacent
		// to the current pawn, then add en-passant to the pawn's attack targets.
		//
		if p.enpassant != 0 && maskPawn[color][p.enpassant].isSet(square) {
			bitmask.set(p.enpassant)
		}
	case Knight:
		bitmask = knightMoves[square] & ^p.outposts[color]
	case Bishop:
		bitmask = p.bishopMoves(square) & ^p.outposts[color]
	case Rook:
		bitmask = p.rookMoves(square) & ^p.outposts[color]
	case Queen:
		bitmask = (p.bishopMoves(square) | p.rookMoves(square)) & ^p.outposts[color]
	case King:
		bitmask = kingMoves[square] & ^p.outposts[color]
	}
	return
}

func (p *Position) xrayTargets(square int) Bitmask {
	return p.xrayTargetsFor(square, p.pieces[square])
}

func (p *Position) xrayTargetsFor(square int, piece Piece) (bitmask Bitmask) {
	switch kind, color := piece.kind(), piece.color(); kind {
	case Bishop:
		board := p.board ^ p.outposts[queen(color)]
		bitmask = p.bishopMovesAt(square, board) & ^p.outposts[color]
	case Rook:
		board := p.board ^ p.outposts[queen(color)]
		bitmask = p.rookMovesAt(square, board) & ^p.outposts[color]
	}
	return p.targetsFor(square, piece)
}

func (p *Position) attacks(square int) Bitmask {
	return p.attacksFor(square, p.pieces[square])
}

func (p *Position) attacksFor(square int, piece Piece) (bitmask Bitmask) {
	switch kind, color := piece.kind(), piece.color(); kind {
	case Pawn:
		return pawnMoves[color][square]
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
	return
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
		board := p.board ^ p.outposts[queen(color)]
		return p.rookMovesAt(square, board)
	}
	return p.attacksFor(square, piece)
}

func (p *Position) allAttacks(color int) (bitmask Bitmask) {
	bitmask = p.pawnAttacks(color) | p.knightAttacks(color) | p.kingAttacks(color)

	outposts := p.outposts[bishop(color)] | p.outposts[queen(color)]
	for outposts != 0 {
		bitmask |= p.bishopMoves(outposts.pop())
	}

	outposts = p.outposts[rook(color)] | p.outposts[queen(color)]
	for outposts != 0 {
		bitmask |= p.rookMoves(outposts.pop())
	}
	return
}

// Returns a bitmask of pieces that attack given square. The resulting bitmask
// only counts pieces of requested color.
//
// This method is used in static exchange evaluation so instead of using current
// board bitmask (p.board) we pass the one that gets continuously updated during
// the evaluation.
func (p *Position) attackers(square, color int, board Bitmask) Bitmask {
	attackers := knightMoves[square] & p.outposts[knight(color)]
	attackers |= maskPawn[color][square] & p.outposts[pawn(color)]
	attackers |= kingMoves[square] & p.outposts[king(color)]
	attackers |= p.rookMovesAt(square, board) & (p.outposts[rook(color)] | p.outposts[queen(color)])
	attackers |= p.bishopMovesAt(square, board) & (p.outposts[bishop(color)] | p.outposts[queen(color)])

	return attackers
}

func (p *Position) isAttacked(square, color int) bool {
	return (knightMoves[square] & p.outposts[knight(color)]) != 0 ||
		(maskPawn[color][square] & p.outposts[pawn(color)]) != 0 ||
		(kingMoves[square] & p.outposts[king(color)]) != 0 ||
		(p.rookMoves(square) & (p.outposts[rook(color)]|p.outposts[queen(color)])) != 0 ||
		(p.bishopMoves(square) & (p.outposts[bishop(color)]|p.outposts[queen(color)])) != 0
}

func (p *Position) pawnAttacks(color int) (bitmask Bitmask) {
	if color == White {
		bitmask = (p.outposts[Pawn] & ^maskFile[0]) << 7
		bitmask |= (p.outposts[Pawn] & ^maskFile[7]) << 9
	} else {
		bitmask = (p.outposts[BlackPawn] & ^maskFile[0]) >> 9
		bitmask |= (p.outposts[BlackPawn] & ^maskFile[7]) >> 7
	}
	return
}

func (p *Position) knightAttacks(color int) (bitmask Bitmask) {
	outposts := p.outposts[knight(color)]
	for outposts != 0 {
		bitmask |= knightMoves[outposts.pop()]
	}
	return
}

func (p *Position) bishopAttacks(color int) (bitmask Bitmask) {
	outposts := p.outposts[bishop(color)]
	for outposts != 0 {
		bitmask |= p.bishopMoves(outposts.pop())
	}
	return
}

func (p *Position) rookAttacks(color int) (bitmask Bitmask) {
	outposts := p.outposts[rook(color)]
	for outposts != 0 {
		bitmask |= p.rookMoves(outposts.pop())
	}
	return
}

func (p *Position) queenAttacks(color int) (bitmask Bitmask) {
	outposts := p.outposts[queen(color)]
	for outposts != 0 {
		square := outposts.pop()
		bitmask |= p.rookMoves(square) | p.bishopMoves(square)
	}
	return
}

func (p *Position) kingAttacks(color int) Bitmask {
	return kingMoves[p.king[color]]
}

func (p *Position) strongestPiece(color int, targets Bitmask) Piece {
	if targets & p.outposts[queen(color)] != 0 {
		return queen(color)
	}
	if targets & p.outposts[rook(color)] != 0 {
		return rook(color)
	}
	if targets & p.outposts[bishop(color)] != 0 {
		return bishop(color)
	}
	if targets & p.outposts[knight(color)] != 0 {
		return knight(color)
	}
	if targets & p.outposts[pawn(color)] != 0 {
		return pawn(color)
	}
	return Piece(0)
}