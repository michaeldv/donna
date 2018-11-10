// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

// Non-capturing checks.
func (gen *MoveGen) generateChecks() *MoveGen {
	p := gen.p
	our, their := p.colors()
	square := p.king[their&1]
	r, c := coordinate(square)
	prohibit := maskNone
	empty := ^p.board
	friendly := ^p.outposts[their&1]

	// Non-capturing Knight checks.
	checks := knightMoves[square]
	for bm := p.outposts[knight(our)]; bm.anyʔ(); bm = bm.pop() {
		from := bm.first()
		gen.movePiece(from, knightMoves[from] & checks & empty)
	}

	// Non-capturing Bishop or Queen checks.
	checks = p.bishopAttacksAt(square, their)
	for bm := p.outposts[bishop(our)] | p.outposts[queen(our)]; bm.anyʔ(); bm = bm.pop() {
		from := bm.first()
		diagonal := (r != row(from) && c != col(from))
		for bm := p.bishopAttacksAt(from, their) & checks & friendly; bm.anyʔ(); bm = bm.pop() {
			to := bm.first()
			if piece := p.pieces[to]; piece == 0 {
				// Empty square: simply move a bishop to check.
				gen.add(NewMove(p, from, to))
			} else if diagonal && piece.color() == our && maskLine[from][square].anyʔ() {
				// Non-empty square occupied by friendly piece on the same
				// diagonal: moving the piece away causes discovered check.
				switch piece.kind() {
				case Pawn:
					// Block pawn promotions (since they are treated as
					// captures) and en-passant captures.
					prohibit = maskRank[0] | maskRank[7]
					if p.enpassant != 0 {
						prohibit.set(p.enpassant)
					}
					gen.movePawn(to, p.targets(to) & empty & ^prohibit)
				case King:
					// Make sure the king steps out of attack diaginal.
					gen.moveKing(to, p.targets(to) & empty & ^maskBlock[from][square])
				default:
					gen.movePiece(to, p.targets(to) & empty)
				}
			}
		}
		if p.pieces[from].queenʔ() {
			// Queen could move straight as a rook and check diagonally as a bishop
			// or move diagonally as a bishop and check straight as a rook.
			targets := (p.rookAttacksAt(from, our) & checks) |
				   (p.bishopAttacksAt(from, our) & p.rookAttacksAt(square, our))
			gen.movePiece(from, targets & empty)
		}
	}

	// Non-capturing Rook or Queen checks.
	checks = p.rookAttacksAt(square, their)
	for bm := p.outposts[rook(our)] | p.outposts[queen(our)]; bm.anyʔ(); bm = bm.pop() {
		from := bm.first()
		straight := (r == row(from) || c == col(from))
		for bm := p.rookAttacksAt(from, their) & checks & friendly; bm.anyʔ(); bm = bm.pop() {
			to := bm.first()
			if piece := p.pieces[to]; piece == 0 {
				// Empty square: simply move a rook to check.
				gen.add(NewMove(p, from, to))
			} else if straight && piece.color() == our && maskLine[from][square].anyʔ() {
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
					prohibit = maskRank[0] | maskRank[7]
					if p.enpassant != 0 {
						prohibit.set(p.enpassant)
					}
					gen.movePawn(to, p.targets(to) & empty & ^prohibit)
				case King:
					// Make sure the king steps out of attack file or rank.
					if row(from) == r {
						prohibit = maskRank[r]
					} else {
						prohibit = maskFile[c]
					}
					gen.moveKing(to, p.targets(to) & empty & ^prohibit)
				default:
					gen.movePiece(to, p.targets(to) & empty)
				}
			}
		}
	}

	// Non-capturing Pawn checks.
	for bm := p.outposts[pawn(our)] & maskIsolated[c]; bm.anyʔ(); bm = bm.pop() {
		from := bm.first()
		if bm := maskPawn[our&1][square] & p.targets(from); bm.anyʔ() {
			gen.add(NewPawnMove(p, from, bm.first()))
		}
	}

	return gen
}
