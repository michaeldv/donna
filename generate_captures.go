// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

func (gen *MoveGen) generateCaptures() *MoveGen {
	our, their := gen.p.colors()
	return gen.pawnCaptures(our, their).pieceCaptures(our, their)
}

// Generates all pseudo-legal pawn captures and promotions.
func (gen *MoveGen) pawnCaptures(our, their int) *MoveGen {
	opponent := gen.p.pick(their).all

	for pawns := gen.p.pick(our).pawns; pawns.anyʔ(); pawns = pawns.pop() {
		square := pawns.first()

		// For pawns on files 2-6 the moves include captures only,
		// while for pawns on the 7th file the moves include captures
		// as well as queen promotion.
		if rank(our, square) != A7H7 {
			gen.movePawn(square, gen.p.targets(square) & opponent)
		} else {
			for bm := gen.p.targets(square); bm.anyʔ(); bm = bm.pop() {
				target := bm.first()
				mQ, _, _, _ := NewPromotion(gen.p, square, target)
				gen.add(mQ)
			}
		}
	}

	return gen
}

// Generates all pseudo-legal captures by pieces other than pawn.
func (gen *MoveGen) pieceCaptures(our, their int) *MoveGen {
	side := gen.p.pick(our)
	opponent := gen.p.pick(their).all

	for bm := side.all ^ side.pawns ^ side.king; bm.anyʔ(); bm = bm.pop() {
		square := bm.first()
		gen.movePiece(square, gen.p.targets(square) & opponent)
	}

	if side.king.anyʔ() {
		square := gen.p.pick(our).home
		gen.moveKing(square, gen.p.targets(square) & opponent)
	}

	return gen
}
