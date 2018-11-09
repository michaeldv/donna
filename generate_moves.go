// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

func (gen *MoveGen) generateRootMoves() *MoveGen {
	gen.generateAllMoves()

	if !gen.onlyMoveʔ() {
		gen.validOnly().rank(Move(0))
	}

	return gen
}

// Copies last move returned by nextMove() to the top of the list shifting
// remaining moves down. Head/tail pointers remain unchanged.
func (gen *MoveGen) rearrangeRootMoves() *MoveGen {
	if gen.head > 0 {
		best := gen.list[gen.head - 1]
		copy(gen.list[1:], gen.list[0:gen.head - 1])
		gen.list[0] = best
	}

	return gen
}

func (gen *MoveGen) generateAllMoves() *MoveGen {
	if gen.p.inCheckʔ(gen.p.color) {
		return gen.generateEvasions()
	}

	return gen.generateMoves()
}

func (gen *MoveGen) generateMoves() *MoveGen {
	color := gen.p.color

	return gen.pawnMoves(color).pieceMoves(color).kingMoves(color)
}

func (gen *MoveGen) pawnMoves(color int) *MoveGen {
	for bm := gen.p.outposts[pawn(color)]; bm.anyʔ(); bm = bm.pop() {
		square := bm.first()
		gen.movePawn(square, gen.p.targets(square))
	}

	return gen
}

// Go over all pieces except pawns and the king.
func (gen *MoveGen) pieceMoves(color int) *MoveGen {
	for bm := gen.p.outposts[color] ^ gen.p.outposts[pawn(color)] ^ gen.p.outposts[king(color)]; bm.anyʔ(); bm = bm.pop() {
		square := bm.first()
		gen.movePiece(square, gen.p.targets(square))
	}

	return gen
}

func (gen *MoveGen) kingMoves(color int) *MoveGen {
	if gen.p.outposts[king(color)].anyʔ() {
		square := gen.p.king[color]
		gen.moveKing(square, gen.p.targets(square))

		kingside, queenside := gen.p.canCastleʔ(color)
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
	for bm := targets; bm.anyʔ(); bm = bm.pop() {
		target := bm.first()
		if target > H1 && target < A8 {
			gen.add(NewPawnMove(gen.p, square, target))
		} else { // Promotion.
			mQ, mR, mB, mN := NewPromotion(gen.p, square, target)
			gen.add(mQ).add(mN).add(mR).add(mB)
		}
	}

	return gen
}

func (gen *MoveGen) moveKing(square int, targets Bitmask) *MoveGen {
	for bm := targets; bm.anyʔ(); bm = bm.pop() {
		target := bm.first()
		if abs(square - target) == 2 {
			gen.add(NewCastle(gen.p, square, target))
		} else {
			gen.add(NewMove(gen.p, square, target))
		}
	}

	return gen
}

func (gen *MoveGen) movePiece(square int, targets Bitmask) *MoveGen {
	for bm := targets; bm.anyʔ(); bm = bm.pop() {
		gen.add(NewMove(gen.p, square, bm.first()))
	}

	return gen
}
