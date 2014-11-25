// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

// Non-capturing checks.
func (gen *MoveGen) generateChecks() *MoveGen {
	p := gen.p
	color, enemy := p.color, p.color^1
	square := int(p.king[enemy])

	// Non-capturing Knight checks.
	checks := knightMoves[square]
	outposts := p.outposts[knight(color)]
	for outposts != 0 {
		from := outposts.pop()
		gen.movePiece(from, knightMoves[from] & checks & ^p.board)
	}

	// Non-capturing Bishop or Queen checks.
	checks = p.targetsFor(square, bishop(enemy))
	outposts = p.outposts[bishop(color)] | p.outposts[queen(color)]
	for outposts != 0 {
		from := outposts.pop()
		targets := p.targetsFor(from, bishop(enemy)) & checks & ^p.outposts[enemy]
		for targets != 0 {
			to := targets.pop()
			if piece := p.pieces[to]; piece == 0 {
				// Empty square: simply move a bishop to check.
				gen.add(NewMove(p, from, to))
			} else if piece.color() == color && maskDiagonal[from][square] != 0 {
				// Non-empty square occupied by friendly piece on the same
				// diagonal: moving the piece away causes discovered check.
				switch piece.kind() {
				case Pawn:
					// Block pawn promotions (since they are treated as
					// captures) and en-passant captures.
					prohibit := maskRank[0] | maskRank[7]
					if p.enpassant != 0 {
						prohibit |= bit[p.enpassant]
					}
					gen.movePawn(to, p.targets(to) & ^p.board & ^prohibit)
				case King:
					// Make sure the king steps out of attack diaginal.
					gen.moveKing(to, p.targets(to) & ^p.board & ^maskBlock[from][square])
				default:
					gen.movePiece(to, p.targets(to) & ^p.board)
				}
			}
		}
		if p.pieces[from].isQueen() {
			// Queen could move straight as a rook and check diagonally as a bishop
			// or move diagonally as a bishop and check straight as a rook.
			targets = (p.targetsFor(from, rook(color)) & checks) |
				  (p.targetsFor(from, bishop(color)) & p.targetsFor(square, rook(color)))
			gen.movePiece(from, targets & ^p.board)
		}
	}

	// Non-capturing Rook or Queen checks.
	checks = p.targetsFor(square, rook(enemy))
	outposts = p.outposts[rook(color)] | p.outposts[queen(color)]
	for outposts != 0 {
		from := outposts.pop()
		targets := p.targetsFor(from, rook(enemy)) & checks & ^p.outposts[enemy]
		for targets != 0 {
			to := targets.pop()
			if piece := p.pieces[to]; piece == 0 {
				// Empty square: simply move a rook to check.
				gen.add(NewMove(p, from, to))
			} else if piece.color() == color {
				if maskStraight[from][square] != 0 {
					// Non-empty square occupied by friendly piece on the same
					// file or rank: moving the piece away causes discovered check.
					switch piece.kind() {
					case Pawn:
						// If pawn and rook share the same file then non-capturing
						// discovered check is not possible since the pawn is going
						// to stay on the same file no matter what.
						if col(from) == col(to) {
							continue
						}
						// Block pawn promotions (since they are treated as captures)
						// and en-passant captures.
						prohibit := maskRank[0] | maskRank[7]
						if p.enpassant != 0 {
							prohibit |= bit[p.enpassant]
						}
						gen.movePawn(to, p.targets(to) & ^p.board & ^prohibit)
					case King:
						// Make sure the king steps out of attack file or rank.
						prohibit := maskNone
						if r := row(from); r == row(square) {
							prohibit = maskRank[r]
						} else {
							prohibit = maskFile[col(square)]
						}
						gen.moveKing(to, p.targets(to) & ^p.board & ^prohibit)
					default:
						gen.movePiece(to, p.targets(to) & ^p.board)
					}
				}
			}
		}
	}

	// Non-capturing Pawn checks.
	outposts = p.outposts[pawn(color)] & maskIsolated[col(square)]
	for outposts != 0 {
		from := outposts.pop()
		if target := maskPawn[color][square] & p.targets(from); target != 0 {
			gen.add(NewPawnMove(p, from, target.pop()))
		}
	}

	return gen
}
