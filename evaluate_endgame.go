// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

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

func (e *Evaluation) strongerSide() int {
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
		color = e.position.color
		wKing = e.position.king[White]
		bKing = e.position.king[Black]
		wPawn = e.position.outposts[Pawn].last()
	} else {
		color = e.position.color ^ 1
		wKing = 64 + ^e.position.king[Black]
		bKing = 64 + ^e.position.king[White]
		wPawn = 64 + ^e.position.outposts[BlackPawn].last()
	}

	index := color + (wKing << 1) + (bKing << 7) + ((wPawn - 8) << 13)
	if (bitbase[index / 64] & bit[index & 0x3F]).noneʔ() {
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
	row, col := coordinate(e.position.king[color^1])

	// Pawns on A file with bare king opposing them.
	if (pawns & ^maskFile[A1]).noneʔ() && (pawns & ^maskInFront[color^1][row * 8]).noneʔ() && col <= B1 {
		return DrawScore
	}

	// Pawns on H file with bare king opposing them.
	if (pawns & ^maskFile[H1]).noneʔ() && (pawns & ^maskInFront[color^1][row * 8 + 7]).noneʔ() && col >= G1 {
		return DrawScore
	}

	return ExistingScore
}

// Bishop-only endgame: drop the score if we have opposite-colored bishops.
func (e *Evaluation) bishopsAndPawns() int {
	if e.oppositeBishopsʔ() {
		outposts := &e.position.outposts
		if abs(outposts[Pawn].count() - outposts[BlackPawn].count()) == 1 {
			return e.fraction(1, 8) // 1/8
		}
		return e.fraction(1, 2) // 1/2
	}

	return ExistingScore
}

// Single bishops plus some other pieces: drop the score if we have opposite-colored
// bishops but only if other minors/majors are balanced.
func (e *Evaluation) drawishBishops() int {
	if e.oppositeBishopsʔ() {
		outposts := &e.position.outposts
		wN, bN := outposts[Knight].count(), outposts[BlackKnight].count()
		wR, bR := outposts[Rook].count(), outposts[BlackRook].count()
		extraPawns := abs(outposts[Pawn].count() - outposts[BlackPawn].count())

		if wN == bN && wR == bR && extraPawns <= 2 {
			return e.fraction(1, 4) // 1/4
		}
	}

	return ExistingScore
}

func (e *Evaluation) kingAndPawnVsKingAndPawn() int {
	if e.score.endgame == 0 {
		return ExistingScore
	}

	p := e.position
	unstoppableʔ := func(color int, square int) bool {
		if (p.outposts[color^1] & maskInFront[color][square]).noneʔ() {
			mask := maskNone
			if p.color == color {
				mask = maskSquare[color][square]
			} else {
				mask = maskSquareEx[color][square]
			}
			return (mask & p.outposts[king(color^1)]).noneʔ()
		}
		return false
	}

	// Check if either side has unstoppable pawn.
	white := unstoppableʔ(White, p.outposts[pawn(White)].first())
	black := unstoppableʔ(Black, p.outposts[pawn(Black)].first())
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
	if rank(our, pawns.first()) < A5H5 || (pawns & (maskFile[0] | maskFile[7])).anyʔ() {

		// Temporarily remove opponent's pawn.
		piece = pawn(our^1)
		pawns = p.outposts[piece]		// -> Save: opponent's pawn bitmask.
		square = pawns.first()			// -> Save: opponent's pawn square.

		p.outposts[piece] = maskNone
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

// One side has 1 pawn and the other side has 1 or more pawns: reduce
// score if both sides have exactly 1 pawn.
func (e *Evaluation) lastPawnLeft() int {
	outposts := &e.position.outposts

	if outposts[Pawn].singleʔ() && outposts[BlackPawn].singleʔ() {
		return e.fraction(3, 4) // 3/4
	}

	return ExistingScore
}

// One side has 0 pawn and the other side has 0 or more pawns.
func (e *Evaluation) noPawnsLeft() int {
	color := e.strongerSide()
	outposts := &e.position.outposts
	whiteMinorOnly := outposts[Queen].noneʔ() && outposts[Rook].noneʔ() && (outposts[Bishop] | outposts[Knight]).singleʔ()
	blackMinorOnly := outposts[BlackQueen].noneʔ() && outposts[BlackRook].noneʔ() && (outposts[BlackBishop] | outposts[BlackKnight]).singleʔ()

	// Check for opposite bishops first.
	if whiteMinorOnly && blackMinorOnly && outposts[Knight].noneʔ() && outposts[BlackKnight].noneʔ() && e.oppositeBishopsʔ() {
		pawn := outposts[pawn(color)]			// The passer.
		king := outposts[king(color^1)]			// Defending king.
		path := maskInFront[color][pawn.first()]	// Path in front of the passer.
		safe := maskDark				// Safe squares for king to block on.
		if (outposts[bishop(color)] & maskDark).anyʔ() {
			safe = ^maskDark
		}

		// Draw if king blocks the passer on safe square.
		if (king & safe).anyʔ() && (king & path).anyʔ() {
			return DrawScore
		}

		// Draw if bishop attacks a square in front of the passer.
		if (e.position.bishopAttacks(color^1) & path).anyʔ() {
			return DrawScore
		}
	}

	if color == White && outposts[Pawn].noneʔ() {
		if whiteMinorOnly {
			// There is a theoretical chance of winning if opponent's pawns are on
			// edge files (ex. some puzzles).
			if (outposts[BlackPawn] & (maskFile[0] | maskFile[7])).anyʔ() {
				return e.fraction(1, 64) // 1/64
			}
			return DrawScore
		} else if blackMinorOnly {
			return e.fraction(1, 16) // 1/16
		}
		return e.fraction(3, 16) // 3/16
	}

	if color == Black && outposts[BlackPawn].noneʔ() {
		if blackMinorOnly {
			if (outposts[Pawn] & (maskFile[0] | maskFile[7])).anyʔ() {
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
