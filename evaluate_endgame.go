// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (e *Evaluation) evaluateEndgame() int {
	score := e.material.endgame(e)

	if e.position.color == Black {
		return -score
	}

	return score
}

func (e *Evaluation) inspectEndgame() {
	switch markdown := e.material.endgame(e); markdown {
	case ExistingScore:
		return
	case DrawScore:
		e.score = Score{0, 0}
	case WhiteWinning:
		e.score = Score{WhiteWinning, WhiteWinning}
	case BlackWinning:
		e.score = Score{-BlackWinning, -BlackWinning}
	default:
		mul, div := markdown >> 16, markdown & 0xFFFF
		if div != 0 {
			if mul != 0 {
				e.score.endgame = e.score.endgame * mul
			}
			e.score.endgame /= div
		}
	}
}

// Packs fractional markdown value as expected by inspectEndgame().
func (e *Evaluation) fraction(mul, div int) int {
	if mul == 1 {
		return div
	}
	return (mul << 16) | div
}

func (e *Evaluation) strongerSide() uint8 {
	if e.score.endgame > 0 {
		return White
	}
	return Black
}

// Known endgames where we calculate the exact score.
func (e *Evaluation) winAgainstBareKing() int { 	// STUB.
	return e.score.blended(e.material.phase)
}

func (e *Evaluation) knightAndBishopVsBareKing() int {	// STUB.
	return e.score.blended(e.material.phase)
}

func (e *Evaluation) twoBishopsVsBareKing() int { 	// STUB.
	return e.score.blended(e.material.phase)
}

func (e *Evaluation) kingAndPawnVsBareKing() int {
	var color, wKing, bKing, wPawn int

	stronger := e.strongerSide()
	if stronger == White {
		color = int(e.position.color)
		wKing = int(e.position.king[White])
		bKing = int(e.position.king[Black])
		wPawn = e.position.outposts[Pawn].last()
	} else {
		color = int(e.position.color)^1
		wKing = 64 + ^int(e.position.king[Black])
		bKing = 64 + ^int(e.position.king[White])
		wPawn = 64 + ^e.position.outposts[BlackPawn].last()
	}

	index := color + (wKing << 1) + (bKing << 7) + ((wPawn - 8) << 13)
	if bitbase[index / 64] & (1 << uint(index & 0x3F)) == 0 {
		return DrawScore
	}

	if stronger == Black {
		return BlackWinning
	}

	return WhiteWinning
}

// Lesser known endgames where we calculate endgame score markdown.
func (e *Evaluation) kingAndPawnsVsBareKing() int {
	color := e.strongerSide()

	pawns := e.position.outposts[pawn(color)]
	row, col := coordinate(int(e.position.king[color^1]))

	// Pawns on A file with bare king opposing them.
	if (pawns & ^maskFile[A1]).empty() && (pawns & ^maskInFront[color^1][row * 8]).empty() && col <= B1 {
		return DrawScore
	}

	// Pawns on H file with bare king opposing them.
	if (pawns & ^maskFile[H1]).empty() && (pawns & ^maskInFront[color^1][row * 8 + 7]).empty() && col >= G1 {
		return DrawScore
	}

	return ExistingScore
}

// Bishop-only endgame: drop the score if we have opposite-colored bishops.
func (e *Evaluation) bishopsAndPawns() int {
	if e.oppositeBishops() {
		outposts := &e.position.outposts
		if (outposts[Pawn] | outposts[BlackPawn]).count() == 2 {
			return e.fraction(1, 8) // 1/8
		}
		return e.fraction(1, 2) // 1/2
	}

	return ExistingScore
}

// Single bishops plus some minors: drop the score if we have opposite-colored bishops.
func (e *Evaluation) drawishBishops() int {
	if e.oppositeBishops() {
		return e.fraction(1, 4) // 1/4
	}

	return ExistingScore
}

func (e *Evaluation) kingAndPawnVsKingAndPawn() int {
	if e.score.endgame == 0 {
		return ExistingScore
	}

	p := e.position
	unstoppable := func(color uint8, square int) bool {
		if (p.outposts[color^1] & maskInFront[color][square]).empty() {
			mask := Bitmask(0)
			if p.color == color {
				mask = maskSquare[color][square]
			} else {
				mask = maskSquareEx[color][square]
			}
			return (mask & p.outposts[king(color^1)]).empty()
		}
		return false
	}

	// Check if either side has unstoppable pawn.
	white := unstoppable(White, p.outposts[pawn(White)].first())
	black := unstoppable(Black, p.outposts[pawn(Black)].first())
	if white {
		e.score.endgame = WhiteWinning
	}
	if black {
		e.score.endgame = BlackWinning
	}
	if white || black {
		return ExistingScore
	}

	// Try to evaluate the endgame using KPK bitbase. If the opposite side is not loosing
	// without the pawn it's unlikely the game is lost with the pawn present.
	our := p.color
	piece := pawn(our)
	pawns := p.outposts[piece]
	square := pawns.first()
	if rank(our, pawns.first()) < A5H5 || (pawns & (maskFile[0] | maskFile[7])).any() {

		// Temporarily remove opponent's pawn.
		piece = pawn(our^1)
		pawns = p.outposts[piece]		// -> Save: opponent's pawn bitmask.
		square = pawns.first()			// -> Save: opponent's pawn square.

		p.outposts[piece] = Bitmask(0)
		p.pieces[square] = Piece(0)

		// Temporarily adjust score so that when e.strongerSide() gets called
		// by kingAndPawnVsBareKing() it returns side to move.
		score := e.score.endgame		// -> Save: endgame score.
		e.score.endgame = let(our == White, 1, -1)

		// When we're done restore original endgame score and opponent's pawn.
		defer func() {
			e.score.endgame = score		// <- Restore: endgame score.
			p.pieces[square] = piece	// <- Restore: opponent's pawn square.
			p.outposts[piece] = pawns	// <- Restore: opponent's pawn bitmask.
		}()

		if e.kingAndPawnVsBareKing() == DrawScore {
			return DrawScore
		}
	}

	return ExistingScore
}

func (e *Evaluation) bishopAndPawnVsBareKing() int {	// STUB.
	return ExistingScore
}

func (e *Evaluation) rookAndPawnVsRook() int {		// STUB.
	return ExistingScore
}

func (e *Evaluation) queenVsRookAndPawns() int { 	// STUB.
	return ExistingScore
}

func (e *Evaluation) lastPawnLeft() int {
	color := e.strongerSide()
	outposts := &e.position.outposts

	if (color == White && outposts[Pawn].count() == 1) || (color == Black && outposts[BlackPawn].count() == 1) {
		return e.fraction(3, 4) // 3/4
	}

	return ExistingScore
}

func (e *Evaluation) noPawnsLeft() int {
	color := e.strongerSide()
	outposts := &e.position.outposts
	whiteMinorOnly := outposts[Queen].empty() && outposts[Rook].empty() && (outposts[Bishop] | outposts[Knight]).count() == 1
	blackMinorOnly := outposts[BlackQueen].empty() && outposts[BlackRook].empty() && (outposts[BlackBishop] | outposts[BlackKnight]).count() == 1

	if color == White && outposts[Pawn].empty() {
		if whiteMinorOnly {
			// There is a theoretical chance of winning if opponent's pawns are on
			// edge files (ex. some puzzles).
			if (outposts[BlackPawn] & (maskFile[0] | maskFile[7])).any() {
				return e.fraction(1, 64) // 1/64
			}
			return DrawScore
		} else if blackMinorOnly {
			return e.fraction(1, 16) // 1/16
		}
		return e.fraction(3, 16) // 3/16
	}

	if color == Black && outposts[BlackPawn].empty() {
		if blackMinorOnly {
			if (outposts[Pawn] & (maskFile[0] | maskFile[7])).any() {
				return e.fraction(1, 64) // 1/64
			}
			return DrawScore
		} else if whiteMinorOnly {
			return e.fraction(1, 16) // 1/16
		}
		return e.fraction(3, 16) // 3/16
	}

	return ExistingScore
}
