// Copyright (c) 2014-2016 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

// Non-capturing checks.
func (gen *MoveGen) generateChecks() *MoveGen {
	p := gen.p
	color, enemy := p.color, p.color^1
	square := int(p.king[enemy])
	r, c := coordinate(square)

	// Non-capturing Knight checks.
	checks := knightMoves[square]
	for bm := p.outposts[knight(color)]; bm.any(); bm = bm.pop() {
		from := bm.first()
		gen.movePiece(from, knightMoves[from] & checks & ^p.board)
	}

	// Non-capturing Bishop or Queen checks.
	checks = p.bishopAttacksAt(square, enemy)
	for bm := p.outposts[bishop(color)] | p.outposts[queen(color)]; bm.any(); bm = bm.pop() {
		from := bm.first()
		diagonal := (r != row(from) && c != col(from))
		for bm := p.bishopAttacksAt(from, enemy) & checks & ^p.outposts[enemy]; bm.any(); bm = bm.pop() {
			to := bm.first()
			if piece := p.pieces[to]; piece == 0 {
				// Empty square: simply move a bishop to check.
				gen.add(NewMove(p, from, to))
			} else if diagonal && piece.color() == color && maskLine[from][square].any() {
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
			targets := (p.rookAttacksAt(from, color) & checks) |
				   (p.bishopAttacksAt(from, color) & p.rookAttacksAt(square, color))
			gen.movePiece(from, targets & ^p.board)
		}
	}

	// Non-capturing Rook or Queen checks.
	checks = p.rookAttacksAt(square, enemy)
	for bm := p.outposts[rook(color)] | p.outposts[queen(color)]; bm.any(); bm = bm.pop() {
		from := bm.first()
		straight := (r == row(from) || c == col(from))
		for bm := p.rookAttacksAt(from, enemy) & checks & ^p.outposts[enemy]; bm.any(); bm = bm.pop() {
			to := bm.first()
			if piece := p.pieces[to]; piece == 0 {
				// Empty square: simply move a rook to check.
				gen.add(NewMove(p, from, to))
			} else if straight && piece.color() == color && maskLine[from][square].any() {
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
					prohibit := Bitmask(0)
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

	// Non-capturing Pawn checks.
	for bm := p.outposts[pawn(color)] & maskIsolated[col(square)]; bm.any(); bm = bm.pop() {
		from := bm.first()
		if bm := maskPawn[color][square] & p.targets(from); bm.any() {
			gen.add(NewPawnMove(p, from, bm.first()))
		}
	}

	return gen
}
