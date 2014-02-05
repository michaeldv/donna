// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (gen *MoveGen) GenerateCaptures() *MoveGen {
        color := gen.p.color
        gen.pawnCaptures(color)
        gen.pieceCaptures(color)
        return gen
}

// Generates all pseudo-legal pawn captures and Queen promotions.
func (gen *MoveGen) pawnCaptures(color int) *MoveGen {
        pawns := gen.p.outposts[Pawn(color)]

        for pawns != 0 {
                square := pawns.pop()
                //
                // First check capture targets on rows 2-7 (no promotions).
                //
                targets := gen.p.targets[square] & gen.p.board[color^1] & 0x00FFFFFFFFFFFF00
                for targets != 0 {
                        gen.list[gen.tail].move = gen.p.NewMove(square, targets.pop())
                        gen.tail++
                }
                //
                // Now check promo rows. The might include capture targets as well
                // as empty promo square in front of the pawn.
                //
                if RelRow(square, color) == 6 {
                        //
                        // Select maskRank[7] for white and maskRank[0] for black.
                        //
                        targets  = gen.p.targets[square] & maskRank[7 - 7 * color]
                        targets |= gen.p.board[2] & Bit(square + eight[color])

                        for targets != 0 {
                                gen.list[gen.tail].move = gen.p.NewMove(square, targets.pop()).promote(QUEEN)
                                gen.tail++
                        }
                }
        }
        return gen
}

// Generates all pseudo-legal captures by pieces other than pawn.
func (gen *MoveGen) pieceCaptures(color int) *MoveGen {
	for _, kind := range [5]int{ KNIGHT, BISHOP, ROOK, QUEEN, KING } {
	        outposts := gen.p.outposts[Piece(kind|color)]
	        for outposts != 0 {
	                square := outposts.pop()
	                targets := gen.p.targets[square] & gen.p.board[color^1]
	                for targets != 0 {
	                        gen.list[gen.tail].move = gen.p.NewMove(square, targets.pop())
	                        gen.tail++
	                }
	        }
	}
	return gen
}
