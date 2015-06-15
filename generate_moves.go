// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (gen *MoveGen) generateRootMoves() *MoveGen {
	gen.generateAllMoves()

	if !gen.onlyMove() {
		gen.validOnly().rank(Move(0))
	}

	return gen
}

// Copies last move returned by NextMove() to the top of the list shifting
// remaining moves down. Head/tail pointers remain unchanged.
func (gen *MoveGen) rearrangeRootMoves() *MoveGen {
	if gen.head > 0 {
		best := gen.list[gen.head - 1]
		copy(gen.list[1:], gen.list[0:gen.head - 1])
		gen.list[0] = best
	}

	return gen
}

func (gen *MoveGen) cleanupRootMoves(depth int) *MoveGen {
	if gen.size() < 2 {
		return gen
	}

	// Always preserve first higest ranking move.
	for i := 1; i < gen.tail; i++ {
		if gen.list[i].score == -depth + 1 {
			gen.tail = i
			break
		}
	}

	return gen.reset()
}

func (gen *MoveGen) generateAllMoves() *MoveGen {
	if gen.p.isInCheck(gen.p.color) {
		return gen.generateEvasions()
	}
	return gen.generateMoves()
}

func (gen *MoveGen) generateMoves() *MoveGen {
	color := gen.p.color
	return gen.pawnMoves(color).pieceMoves(color).kingMoves(color)
}

func (gen *MoveGen) pawnMoves(color uint8) *MoveGen {
	for pawns := gen.p.outposts[pawn(color)]; pawns != 0; {
		square := pawns.pop()
		gen.movePawn(square, gen.p.targets(square))
	}
	return gen
}

// Go over all pieces except pawns and the king.
func (gen *MoveGen) pieceMoves(color uint8) *MoveGen {
	outposts := gen.p.outposts[color] & ^gen.p.outposts[pawn(color)] & ^gen.p.outposts[king(color)]
	for outposts != 0 {
		square := outposts.pop()
		gen.movePiece(square, gen.p.targets(square))
	}
	return gen
}

func (gen *MoveGen) kingMoves(color uint8) *MoveGen {
	if gen.p.outposts[king(color)] != 0 {
		square := int(gen.p.king[color])
		gen.moveKing(square, gen.p.targets(square))

		kingside, queenside := gen.p.canCastle(color)
		if kingside {
			gen.moveKing(square, bit[G1 + 56 * color])
		}
		if queenside {
			gen.moveKing(square, bit[C1 + 56 * color])
		}
	}
	return gen
}

func (gen *MoveGen) movePawn(square int, targets Bitmask) *MoveGen {
	for targets != 0 {
		target := targets.pop()
		if target > H1 && target < A8 {
			gen.add(NewPawnMove(gen.p, square, target))
		} else { // Promotion.
			mQ, mR, mB, mN := NewPromotion(gen.p, square, target)
			gen.add(mQ).add(mR).add(mB).add(mN)
		}
	}
	return gen
}

func (gen *MoveGen) moveKing(square int, targets Bitmask) *MoveGen {
	for targets != 0 {
		target := targets.pop()
		if abs(square - target) == 2 {
			gen.add(NewCastle(gen.p, square, target))
		} else {
			gen.add(NewMove(gen.p, square, target))
		}
	}
	return gen
}

func (gen *MoveGen) movePiece(square int, targets Bitmask) *MoveGen {
	for targets != 0 {
		gen.add(NewMove(gen.p, square, targets.pop()))
	}
	return gen
}
