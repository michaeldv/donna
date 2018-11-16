// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

func (p *Position) movePiece(piece Piece, from, to Square) *Position {
	bm := bit(from) | bit(to)
	p.pieces[from], p.pieces[to] = 0, piece
	p.outposts[piece] ^= bm
	p.outposts[piece.color()] ^= bm

	// Update position's hash values.
	random := piece.polyglot(from) ^ piece.polyglot(to)
	p.id ^= random
	if piece.pawnʔ() {
		p.pid ^= random
	}

	// Update positional score.
	p.tally.sub(pst[piece][from]).add(pst[piece][to])

	return p
}

func (p *Position) promotePawn(pawn Piece, from, to Square, promo Piece) *Position {
	p.pieces[from], p.pieces[to] = 0, promo
	p.outposts[pawn] ^= bit(from)
	p.outposts[promo] ^= bit(to)
	p.outposts[pawn.color()] ^= bit(from) | bit(to)

	// Update position's hash values and material balance.
	random := pawn.polyglot(from)
	p.id ^= random ^ promo.polyglot(to)
	p.pid ^= random
	p.balance += materialBalance[promo] - materialBalance[pawn]

	// Update positional score.
	p.tally.sub(pst[pawn][from]).add(pst[promo][to])

	return p
}

func (p *Position) capturePiece(capture Piece, from, to Square) *Position {
	p.outposts[capture] ^= bit(to)
	p.outposts[capture.color()] ^= bit(to)

	// Update position's hash values and material balance.
	random := capture.polyglot(to)
	p.id ^= random
	if capture.pawnʔ() {
		p.pid ^= random
	}
	p.balance -= materialBalance[capture]

	// Update positional score.
	p.tally.sub(pst[capture][to])

	return p
}

func (p *Position) captureEnpassant(capture Piece, from, to Square) *Position {
	enpassant := to.push(capture.color())

	p.pieces[enpassant] = 0
	p.outposts[capture] ^= bit(enpassant)
	p.outposts[capture.color()] ^= bit(enpassant)

	// Update position's hash values and material balance.
	random := capture.polyglot(enpassant)
	p.id ^= random
	p.pid ^= random
	p.balance -= materialBalance[capture]

	// Update positional score.
	p.tally.sub(pst[capture][enpassant])

	return p
}

func (p *Position) makeMove(move Move) *Position {
	our := move.color(); their := our^1
	from, to, piece, capture := move.split()

	// Copy over the contents of previous tree node to the current one.
	node++
	tree[node] = *p // => tree[node] = tree[node - 1]
	pp := &tree[node]

	pp.enpassant, pp.reversibleʔ = 0, true
	if capture == 0 {
		if p.enpassant != 0 {
			pp.id ^= polyglotRandomEp[p.enpassant & 7]
		}
	} else if p.enpassant == 0 || to != p.enpassant {
		pp.count50, pp.reversibleʔ = 0, false
		pp.capturePiece(capture, from, to)
	}

	if piece.pawnʔ() {
		pp.count50, pp.reversibleʔ = 0, false
		if p.enpassant != 0 && to == p.enpassant {
			pp.captureEnpassant(pawn(their), from, to)
			pp.id ^= polyglotRandomEp[p.enpassant & 7] // p.enpassant column.
		}
		if promo := move.promo(); promo != 0 {
			pp.promotePawn(piece, from, to, promo)
		} else {
			pp.movePiece(piece, from, to)
			if from ^ to == 16 && (pawnAttacks[our][from.push(our)] & p.outposts[pawn(their)]).anyʔ() {
				pp.enpassant = from.push(our) // Save the en-passant square.
				pp.id ^= polyglotRandomEp[pp.enpassant & 7]
			}
		}
	} else if piece.kingʔ() {
		pp.movePiece(piece, from, to)
		pp.count50++
		pp.king[our] = to
		pp.castles &= ^maskRank[7*our] // Both castling rights are gone.
		if move.castleʔ() {
			pp.reversibleʔ = false
			if to == from + 2 {
				pp.movePiece(rook(our), to + 1, to - 1)
			} else if to == from - 2 {
				pp.movePiece(rook(our), to - 2, to + 1)
			}
		}
	} else {
		pp.movePiece(piece, from, to)
		pp.count50++
	}

	// Set up the board bitmask, update castle rights, finish off incremental
	// hash value, and flip the color.
	pp.board = pp.outposts[White] | pp.outposts[Black]
	pp.castles &= ^(bit(from) | bit(to))
	pp.id ^= p.polycast() ^ pp.polycast()
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
		pp.id ^= polyglotRandomEp[pp.enpassant & 7]
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

func (p *Position) inCheckʔ(our int) bool {
	return p.attackedʔ(our^1, p.king[our])
}

func (p *Position) nlNodeʔ() bool {
	return node > 0 && tree[node].board == tree[node-1].board
}

func (p *Position) fiftyʔ() bool {
	return p.count50 >= 100
}

func (p *Position) repetitionʔ() bool {
	if !p.reversibleʔ || node < 1 {
		return false
	}

	for previous := node - 1; previous >= 0; previous-- {
		if !tree[previous].reversibleʔ {
			return false
		}
		if tree[previous].id == p.id {
			return true
		}
	}

	return false
}

func (p *Position) thirdRepetitionʔ() bool {
	if !p.reversibleʔ || node < 4 {
		return false
	}

	for previous, repetitions := node - 2, 1; previous >= 0; previous -= 2 {
		if !tree[previous].reversibleʔ || !tree[previous + 1].reversibleʔ {
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
func (p *Position) canCastleʔ(our int) (kingside, queenside bool) {
	their := our^1

	// Start off with simple checks.
	kingside = p.castles.onʔ(Square(H1).flip(their)) && (gapKing[our] & p.board).noneʔ()
	queenside = p.castles.onʔ(Square(A1).flip(their)) && (gapQueen[our] & p.board).noneʔ()

	// If it still looks like the castles are possible perform more expensive
	// final check.
	if kingside || queenside {
		attacks := p.allAttacks(their)
		kingside = kingside && (castleKing[our] & attacks).noneʔ()
		queenside = queenside && (castleQueen[our] & attacks).noneʔ()
	}

	return kingside, queenside
}

// Returns a bitmask of all pinned pieces preventing a check for the king on
// given square. The color of the pieces match the color of the king.
func (p *Position) pins(sq Square) (bitmask Bitmask) {
	our := p.pieces[sq].color()
	their := our^1

	attackers := (p.outposts[bishop(their)] | p.outposts[queen(their)]) & bishopMagicMoves[sq][0]
	attackers |= (p.outposts[rook(their)] | p.outposts[queen(their)]) & rookMagicMoves[sq][0]

	for bm := attackers; bm.anyʔ(); bm = bm.pop() {
		target := bm.first()
		blockers := maskBlock[sq][target] & ^bit(target) & p.board

		if blockers.singleʔ() {
			bitmask |= blockers & p.outposts[our] // Only friendly pieces are pinned.
		}
	}

	return bitmask
}
