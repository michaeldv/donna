// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

// Returns true if *non-evasion* move is valid, i.e. it is possible to make
// the move in current position without violating chess rules. If the king is
// in check the generator is expected to generate valid evasions where extra
// validation is not needed.
func (p *Position) isValid(move Move, pins Bitmask) bool {
	color := move.color() // TODO: make color part of move split.
	from, to, piece, capture := move.split()

	// For rare en-passant pawn captures we validate the move by actually
	// making it, and then taking it back.
	if p.enpassant != 0 && to == p.enpassant && capture.isPawn() {
		if position := p.MakeMove(move); position != nil {
			position.UndoLastMove()
			return true
		}
		return false
	}

	// King's move is valid when a) the move is a castle or b) the destination
	// square is not being attacked by the opponent.
	if piece.isKing() {
		return (move & isCastle != 0) || !p.isAttacked(to, color^1)
	}

	// For all other peices the move is valid when it doesn't cause a
	// check. For pinned sliders this includes moves along the pinning
	// file, rank, or diagonal.
	return pins == 0 || pins.isClear(from) || IsBetween(from, to, p.king[color])
}

// Returns a bitmask of all pinned pieces preventing a check for the king on
// given square. The color of the pieces match the color of the king.
func (p *Position) pinnedMask(square int) (mask Bitmask) {
	color := p.pieces[square].color()
	enemy := color ^ 1
	attackers := (p.outposts[bishop(enemy)] | p.outposts[queen(enemy)]) & bishopMagicMoves[square][0]
	attackers |= (p.outposts[rook(enemy)] | p.outposts[queen(enemy)]) & rookMagicMoves[square][0]

	for attackers != 0 {
		attackSquare := attackers.pop()
		blockers := maskBlock[square][attackSquare] & ^bit[attackSquare] & p.board

		if blockers.count() == 1 {
			mask |= blockers & p.outposts[color] // Only friendly pieces are pinned.
		}
	}
	return
}
