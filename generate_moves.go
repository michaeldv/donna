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
	our := gen.p.color

	return gen.pawnMoves(our).pieceMoves(our).kingMoves(our)
}

func (gen *MoveGen) pawnMoves(our int) *MoveGen {
	for bm := gen.p.outposts[pawn(our)]; bm.anyʔ(); bm = bm.pop() {
		square := bm.first()
		gen.movePawn(square, gen.p.targets(square))
	}

	return gen
}

// Go over all pieces except pawns and the king.
func (gen *MoveGen) pieceMoves(our int) *MoveGen {
	for bm := gen.p.outposts[our&1] ^ gen.p.outposts[pawn(our)] ^ gen.p.outposts[king(our)]; bm.anyʔ(); bm = bm.pop() {
		square := bm.first()
		gen.movePiece(square, gen.p.targets(square))
	}

	return gen
}

func (gen *MoveGen) kingMoves(our int) *MoveGen {
	if gen.p.outposts[king(our)].anyʔ() {
		square := gen.p.king[our&1]
		gen.moveKing(square, gen.p.targets(square))

		kingside, queenside := gen.p.canCastleʔ(our)
		if kingside {
			gen.moveKing(square, bit(Square(G1 + 56 * our)))
		}
		if queenside {
			gen.moveKing(square, bit(Square(C1 + 56 * our)))
		}
	}

	return gen
}

func (gen *MoveGen) movePawn(sq Square, targets Bitmask) *MoveGen {
	for bm := targets; bm.anyʔ(); bm = bm.pop() {
		target := bm.first()
		if target > H1 && target < A8 {
			gen.add(NewMove(gen.p, sq, target))
		} else { // Promotion.
			mQ, mR, mB, mN := NewPromotion(gen.p, sq, target)
			gen.add(mQ).add(mN).add(mR).add(mB)
		}
	}

	return gen
}

func (gen *MoveGen) moveKing(sq Square, targets Bitmask) *MoveGen {
	for bm := targets; bm.anyʔ(); bm = bm.pop() {
		target := bm.first()
		if abs(int(sq) - int(target)) == 2 {
			gen.add(NewCastle(gen.p, sq, target))
		} else {
			gen.add(NewMove(gen.p, sq, target))
		}
	}

	return gen
}

func (gen *MoveGen) movePiece(sq Square, targets Bitmask) *MoveGen {
	for bm := targets; bm.anyʔ(); bm = bm.pop() {
		gen.add(NewMove(gen.p, sq, bm.first()))
	}

	return gen
}
