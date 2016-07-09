// Copyright (c) 2014-2016 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (p *Position) movePiece(piece Piece, from, to int) *Position {
	p.pieces[from], p.pieces[to] = 0, piece
	p.outposts[piece] ^= bit[from] | bit[to]
	p.outposts[piece.color()] ^= bit[from] | bit[to]

	// Update position's hash values.
	random := piece.polyglot(from) ^ piece.polyglot(to)
	p.id ^= random
	if piece.isPawn() {
		p.pawnId ^= random
	}

	// Update positional score.
	p.tally.sub(pst[piece][from]).add(pst[piece][to])

	return p
}

func (p *Position) promotePawn(pawn Piece, from, to int, promo Piece) *Position {
	p.pieces[from], p.pieces[to] = 0, promo
	p.outposts[pawn] ^= bit[from]
	p.outposts[promo] ^= bit[to]
	p.outposts[pawn.color()] ^= bit[from] | bit[to]

	// Update position's hash values and material balance.
	random := pawn.polyglot(from)
	p.id ^= random ^ promo.polyglot(to)
	p.pawnId ^= random
	p.balance += materialBalance[promo] - materialBalance[pawn]

	// Update positional score.
	p.tally.sub(pst[pawn][from]).add(pst[promo][to])

	return p
}

func (p *Position) capturePiece(capture Piece, from, to int) *Position {
	p.outposts[capture] ^= bit[to]
	p.outposts[capture.color()] ^= bit[to]

	// Update position's hash values and material balance.
	random := capture.polyglot(to)
	p.id ^= random
	if capture.isPawn() {
		p.pawnId ^= random
	}
	p.balance -= materialBalance[capture]

	// Update positional score.
	p.tally.sub(pst[capture][to])

	return p
}

func (p *Position) captureEnpassant(capture Piece, from, to int) *Position {
	enpassant := to - up[capture.color()^1]

	p.pieces[enpassant] = 0
	p.outposts[capture] ^= bit[enpassant]
	p.outposts[capture.color()] ^= bit[enpassant]

	// Update position's hash values and material balance.
	random := capture.polyglot(enpassant)
	p.id ^= random
	p.pawnId ^= random
	p.balance -= materialBalance[capture]

	// Update positional score.
	p.tally.sub(pst[capture][enpassant])

	return p
}

func (p *Position) makeMove(move Move) *Position {
	color := move.color()
	from, to, piece, capture := move.split()

	// Copy over the contents of previous tree node to the current one.
	node++
	tree[node] = *p // => tree[node] = tree[node - 1]
	pp := &tree[node]

	pp.enpassant, pp.reversible = 0, true
	pp.count50++

	if capture != 0 {
		pp.count50, pp.reversible = 0, false
		if to != 0 && to == int(p.enpassant) {
			pp.captureEnpassant(pawn(color^1), from, to)
			pp.id ^= hashEnpassant[p.enpassant & 7] // p.enpassant column.
		} else {
			pp.capturePiece(capture, from, to)
		}
	}

	if promo := move.promo(); promo == 0 {
		pp.movePiece(piece, from, to)

		if piece.isKing() {
			pp.king[color] = uint8(to)
			if move.isCastle() {
				pp.reversible = false
				switch to {
				case G1:
					pp.movePiece(Rook, H1, F1)
				case C1:
					pp.movePiece(Rook, A1, D1)
				case G8:
					pp.movePiece(BlackRook, H8, F8)
				case C8:
					pp.movePiece(BlackRook, A8, D8)
				}
			}
		} else if piece.isPawn() {
			pp.count50, pp.reversible = 0, false
			if move.isEnpassant() {
				pp.enpassant = uint8(from + up[color]) // Save the en-passant square.
				pp.id ^= hashEnpassant[pp.enpassant & 7]
			}
		}
	} else {
		pp.count50, pp.reversible = 0, false
		pp.promotePawn(piece, from, to, promo)
	}

	// Set up the board bitmask, update castle rights, finish off incremental
	// hash value, and flip the color.
	pp.board = pp.outposts[White] | pp.outposts[Black]
	pp.castles &= castleRights[from] & castleRights[to]
	pp.id ^= hashCastle[p.castles] ^ hashCastle[pp.castles]
	pp.id ^= polyglotRandomWhite
	pp.color ^= 1 // <-- Flip side to move.
	pp.score = Unknown

	return &tree[node] // pp
}

// Makes "null" move by copying over previous node position (i.e. preserving all pieces
// intact) and flipping the color.
func (p *Position) makeNullMove() *Position {
	node++
	tree[node] = *p // => tree[node] = tree[node - 1]
	pp := &tree[node]

	// Flipping side to move obviously invalidates the enpassant square.
	if pp.enpassant != 0 {
		pp.id ^= hashEnpassant[pp.enpassant & 7]
		pp.enpassant = 0
	}
	pp.id ^= polyglotRandomWhite
	pp.color ^= 1 // <-- Flip side to move.
	pp.count50++

	return &tree[node] // pp
}

// Restores previous position effectively taking back the last move made.
func (p *Position) undoLastMove() *Position {
	if node > 0 {
		node--
	}
	return &tree[node]
}

func (p *Position) isInCheck(color uint8) bool {
	return p.isAttacked(color^1, int(p.king[color]))
}

func (p *Position) isNull() bool {
	return node > 0 && tree[node].board == tree[node-1].board
}

func (p *Position) fifty() bool {
	return p.count50 >= 100
}

func (p *Position) repetition() bool {
	if !p.reversible || node < 1 {
		return false
	}

	for previous := node - 1; previous >= 0; previous-- {
		if !tree[previous].reversible {
			return false
		}
		if tree[previous].id == p.id {
			return true
		}
	}

	return false
}

func (p *Position) thirdRepetition() bool {
	if !p.reversible || node < 4 {
		return false
	}

	for previous, repetitions := node - 2, 1; previous >= 0; previous -= 2 {
		if !tree[previous].reversible || !tree[previous + 1].reversible {
			return false
		}
		if tree[previous].id == p.id {
			repetitions++
			if repetitions == 3 {
				return true
			}
		}
	}

	return false
}

// Returns a pair of booleans that indicate whether given side is allowed to
// castle kingside and queenside.
func (p *Position) canCastle(color uint8) (kingside, queenside bool) {

	// Start off with simple checks.
	kingside = (p.castles & castleKingside[color] != 0) && (gapKing[color] & p.board == 0)
	queenside = (p.castles & castleQueenside[color] != 0) && (gapQueen[color] & p.board == 0)

	// If it still looks like the castles are possible perform more expensive
	// final check.
	if kingside || queenside {
		attacks := p.allAttacks(color^1)
		kingside = kingside && (castleKing[color] & attacks == 0)
		queenside = queenside && (castleQueen[color] & attacks == 0)
	}

	return kingside, queenside
}

// Returns a bitmask of all pinned pieces preventing a check for the king on
// given square. The color of the pieces match the color of the king.
func (p *Position) pins(square uint8) (bitmask Bitmask) {
	our := p.pieces[square].color()
	their := our^1

	attackers := (p.outposts[bishop(their)] | p.outposts[queen(their)]) & bishopMagicMoves[square][0]
	attackers |= (p.outposts[rook(their)] | p.outposts[queen(their)]) & rookMagicMoves[square][0]

	for bm := attackers; bm.any(); bm = bm.pop() {
		attackSquare := bm.first()
		blockers := maskBlock[square][attackSquare] & ^bit[attackSquare] & p.board

		if blockers.single() {
			bitmask |= blockers & p.outposts[our] // Only friendly pieces are pinned.
		}
	}

	return bitmask
}
