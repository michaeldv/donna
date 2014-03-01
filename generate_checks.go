// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import()

// Non-capturing checks.
func (gen *MoveGen) GenerateChecks() *MoveGen {
        color := gen.p.color
        enemy := gen.p.color^1
        square := gen.p.outposts[king(enemy)].first()

        // Non-capturing Knight checks.
        checks := knightMoves[square]
        outposts := gen.p.outposts[knight(color)]
        for outposts != 0 {
                from := outposts.pop()
                gen.movePiece(from, knightMoves[from] & checks & ^gen.p.board)
        }

        // Non-capturing Bishop or Queen checks.
        checks = gen.p.targetsFor(square, bishop(enemy))
        outposts = gen.p.outposts[bishop(color)] | gen.p.outposts[queen(color)]
        for outposts != 0 {
                from := outposts.pop()
                targets := gen.p.targetsFor(from, bishop(enemy)) & checks & ^gen.p.outposts[enemy]
                for targets != 0 {
                        to := targets.pop()
                        if piece := gen.p.pieces[to]; piece == 0 {
                                //
                                // Empty square: simply move a bishop to check.
                                //
                                gen.add(gen.p.NewMove(from, to))
                        } else if piece.color() == color && maskDiagonal[from][square] != 0 {
                                //
                                // Non-empty square occupied by friendly piece on the same
                                // diagonal: moving the piece away causes discovered check.
                                //
                                switch piece.kind() {
                                case Pawn:
                                        // Block pawn promotions (since they are treated as
                                        // captures) and en-passant captures.
                                        prohibit := maskRank[0] | maskRank[7]
                                        if gen.p.flags.enpassant != 0 {
                                                prohibit.set(gen.p.flags.enpassant)
                                        }
                                        gen.movePawn(to, gen.p.targets(to) & ^gen.p.board & ^prohibit)
                                case King:
                                        // Make sure the king steps out of attack diaginal.
                                        gen.moveKing(to, gen.p.targets(to) & ^gen.p.board & ^maskBlock[from][square])
                                default:
                                        gen.movePiece(to, gen.p.targets(to) & ^gen.p.board)
                                }
                        }
                }
		if gen.p.pieces[from].isQueen() {
			//
			// Queen could move straight as a rook and check diagonally as a bishop
			// or move diagonally as a bishop and check straight as a rook.
			//
			targets = (gen.p.targetsFor(from, rook(color)) & checks) |
			          (gen.p.targetsFor(from, bishop(color)) & gen.p.targetsFor(square, rook(color)))
                        gen.movePiece(from, targets & ^gen.p.board)
		}
        }

        // Non-capturing Rook or Queen checks.
        checks = gen.p.targetsFor(square, rook(enemy))
        outposts = gen.p.outposts[rook(color)] | gen.p.outposts[queen(color)]
        for outposts != 0 {
                from := outposts.pop()
                targets := gen.p.targetsFor(from, rook(enemy)) & checks & ^gen.p.outposts[enemy]
                for targets != 0 {
                        to := targets.pop()
                        if piece := gen.p.pieces[to]; piece == 0 {
                                //
                                // Empty square: simply move a rook to check.
                                //
                                gen.add(gen.p.NewMove(from, to))
                        } else if piece.color() == color {
				if maskStraight[from][square] != 0 {
	                                //
	                                // Non-empty square occupied by friendly piece on the same
	                                // file or rank: moving the piece away causes discovered check.
	                                //
	                                switch piece.kind() {
	                                case Pawn:
	                                        // Block pawn promotions (since they are treated as
	                                        // captures) and en-passant captures.
	                                        prohibit := maskRank[0] | maskRank[7]
	                                        if gen.p.flags.enpassant != 0 {
	                                                prohibit.set(gen.p.flags.enpassant)
	                                        }
	                                        gen.movePawn(to, gen.p.targets(to) & ^gen.p.board & ^prohibit)
	                                case King:
	                                        // Make sure the king steps out of attack file or rank.
						prohibit := maskNone
						if row := Row(from); row == Row(square) {
							prohibit = maskRank[row]
						} else {
							prohibit = maskFile[Col(square)]
						}
	                                        gen.moveKing(to, gen.p.targets(to) & ^gen.p.board & ^prohibit)
	                                default:
	                                        gen.movePiece(to, gen.p.targets(to) & ^gen.p.board)
	                                }
				}
			}
		}
	}

        // Non-capturing Pawn checks.
        outposts = gen.p.outposts[pawn(color)] & maskIsolated[Col(square)]
        for outposts != 0 {
                from := outposts.pop()
                if target := maskPawn[color][square] & gen.p.targets(from); target != 0 {
                        gen.add(gen.p.pawnMove(from, target.pop()))
                }
        }

        return gen
}
