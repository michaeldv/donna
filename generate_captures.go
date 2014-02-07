// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (gen *MoveGen) GenerateCaptures() *MoveGen {
        color := gen.p.color
        return gen.pawnCaptures(color).pieceCaptures(color)
}

// Generates all pseudo-legal pawn captures and Queen promotions.
func (gen *MoveGen) pawnCaptures(color int) *MoveGen {
        pawns := gen.p.outposts[Pawn(color)]

        for pawns != 0 {
                //
                // First check capture targets on rows 2-7 (no promotions).
                //
                square := pawns.pop()
                targets := gen.p.targets[square] & gen.p.board[color^1] & 0x00FFFFFFFFFFFF00
                for targets != 0 {
                        gen.add(gen.p.NewMove(square, targets.pop()))
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
                                gen.add(gen.p.NewMove(square, targets.pop()).promote(QUEEN))
                        }
                }
        }
        return gen
}

// Generates all pseudo-legal captures by pieces other than pawn.
func (gen *MoveGen) pieceCaptures(color int) *MoveGen {
        enemy := color^1
        for _, kind := range [4]int{ KNIGHT, BISHOP, ROOK, QUEEN } {
	        outposts := gen.p.outposts[Piece(kind|color)]
	        for outposts != 0 {
	                square := outposts.pop()
	                gen.movePiece(square, gen.p.targets[square] & gen.p.board[enemy])
	        }
        }
        if king := gen.p.outposts[King(color)]; king != 0 {
                square := king.pop()
                gen.moveKing(square, gen.p.targets[square] & gen.p.board[enemy])
        }
        return gen
}
