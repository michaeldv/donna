// Copyright (c) 2014-2016 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (gen *MoveGen) generateCaptures() *MoveGen {
	color := gen.p.color
	return gen.pawnCaptures(color).pieceCaptures(color)
}

// Generates all pseudo-legal pawn captures and promotions.
func (gen *MoveGen) pawnCaptures(color uint8) *MoveGen {
	enemy := gen.p.outposts[color^1]
	pawns := gen.p.outposts[pawn(color)]

	for pawns.any() {
		square := pawns.pop()

		// For pawns on files 2-6 the moves include captures only,
		// while for pawns on the 7th file the moves include captures
		// as well as queen promotion.
		if rank(color, square) != A7H7 {
			gen.movePawn(square, gen.p.targets(square) & enemy)
		} else {
			targets := gen.p.targets(square)
			for targets.any() {
				target := targets.pop()
				mQ, _, _, _ := NewPromotion(gen.p, square, target)
				gen.add(mQ)
			}
		}
	}
	return gen
}

// Generates all pseudo-legal captures by pieces other than pawn.
func (gen *MoveGen) pieceCaptures(color uint8) *MoveGen {
	enemy := gen.p.outposts[color^1]
	outposts := gen.p.outposts[color] ^ gen.p.outposts[pawn(color)] ^ gen.p.outposts[king(color)]

	for outposts.any() {
		square := outposts.pop()
		gen.movePiece(square, gen.p.targets(square) & enemy)
	}
	if gen.p.outposts[king(color)].any() {
		square := int(gen.p.king[color])
		gen.moveKing(square, gen.p.targets(square) & enemy)
	}
	return gen
}
