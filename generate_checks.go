// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import()

// Non-capturing checks.
func (gen *MoveGen) GenerateChecks() *MoveGen {
        color := gen.p.color
        enemy := gen.p.color^1
        square := gen.p.outposts[King(enemy)].first()

        // Non-capturing Knight checks.
        checks := knightMoves[square]
        outposts := gen.p.outposts[Knight(color)]
        for outposts != 0 {
                from := outposts.pop()
                gen.movePiece(from, knightMoves[from] & checks & ^gen.p.board[2])
        }

        // Non-capturing Bishop or Queen checks.
        checks = gen.p.Targets(square, Bishop(enemy))
        outposts = gen.p.outposts[Bishop(color)] | gen.p.outposts[Queen(color)]
        for outposts != 0 {
                from := outposts.pop()
                targets := gen.p.Targets(from, Bishop(enemy)) & checks & ^gen.p.board[enemy]
                for targets != 0 {
                        to := targets.pop()
                        if piece := gen.p.pieces[to]; piece == 0 {
                                //
                                // Empty square: simply move a bishop to check.
                                //
                                gen.list[gen.tail].move = gen.p.NewMove(from, to)
                                gen.tail++
                        } else if piece.color() == color && sameDiagonal[from][square] {
                                //
                                // Non-empty square occupied by friendly piece on the same
                                // diagonal: moving the piece away causes discovered check.
                                //
                                switch piece.kind() {
                                case PAWN:
                                        // Block pawn promotions (since they are treated as
                                        // captures) and en-passant captures.
                                        prohibit := maskRank[0] | maskRank[7]
                                        if gen.p.flags.enpassant != 0 {
                                                prohibit.set(gen.p.flags.enpassant)
                                        }
                                        gen.movePawn(to, gen.p.targets[to] & ^gen.p.board[2] & ^prohibit)
                                case KING:
                                        // Make sure the king steps out of attack diaginal.
                                        gen.moveKing(to, gen.p.targets[to] & ^gen.p.board[2] & ^maskBlock[from][square])
                                default:
                                        gen.movePiece(to, gen.p.targets[to] & ^gen.p.board[2])
                                }
                        }
                }
		if gen.p.pieces[from].isQueen() {
			//
			// Queen could move straight as a rook and check diagonally as a bishop
			// or move diagonally as a bishop and check straight as a rook.
			//
			targets = (gen.p.Targets(from, Rook(color)) & checks) |
			          (gen.p.Targets(from, Bishop(color)) & gen.p.Targets(square, Rook(color)))
                        gen.movePiece(from, targets & ^gen.p.board[2])
		}
        }

        // Non-capturing Rook or Queen checks.
        checks = gen.p.Targets(square, Rook(enemy))
        outposts = gen.p.outposts[Rook(color)] | gen.p.outposts[Queen(color)]
        for outposts != 0 {
                from := outposts.pop()
                targets := gen.p.Targets(from, Rook(enemy)) & checks & ^gen.p.board[enemy]
                for targets != 0 {
                        to := targets.pop()
                        if piece := gen.p.pieces[to]; piece == 0 {
                                //
                                // Empty square: simply move a rook to check.
                                //
                                gen.list[gen.tail].move = gen.p.NewMove(from, to)
                                gen.tail++
                        } else if piece.color() == color {
				sameRow := Row(from) == Row(square)
				sameCol := Col(from) == Col(square)
				if sameRow || sameCol {
	                                //
	                                // Non-empty square occupied by friendly piece on the same
	                                // file or rank: moving the piece away causes discovered check.
	                                //
	                                switch piece.kind() {
	                                case PAWN:
	                                        // Block pawn promotions (since they are treated as
	                                        // captures) and en-passant captures.
	                                        prohibit := maskRank[0] | maskRank[7]
	                                        if gen.p.flags.enpassant != 0 {
	                                                prohibit.set(gen.p.flags.enpassant)
	                                        }
	                                        gen.movePawn(to, gen.p.targets[to] & ^gen.p.board[2] & ^prohibit)
	                                case KING:
	                                        // Make sure the king steps out of attack file or rank.
						prohibit := maskNone
						if sameRow {
							prohibit = maskRank[Row(square)]
						} else {
							prohibit = maskFile[Col(square)]
						}
	                                        gen.moveKing(to, gen.p.targets[to] & ^gen.p.board[2] & ^prohibit)
	                                default:
	                                        gen.movePiece(to, gen.p.targets[to] & ^gen.p.board[2])
	                                }
				}
			}
		}
	}

        // Non-capturing Pawn checks.
        return gen
}
