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

func (e *Evaluation) strongerSide() (int, int) {
	if e.score.endgame > 0 {
		return White, Black
	}
	return Black, White
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

	stronger, _ := e.strongerSide()
	if stronger == White {
		color = e.position.color
		wKing = e.position.white.home
		bKing = e.position.black.home
		wPawn = e.position.white.pawns.last()
	} else {
		color = e.position.color ^ 1
		wKing = 64 + ^e.position.black.home
		bKing = 64 + ^e.position.white.home
		wPawn = 64 + ^e.position.black.pawns.last()
	}

	index := color + (wKing << 1) + (bKing << 7) + ((wPawn - 8) << 13)
	if (bitbase[index / 64] & bit(index)).noneʔ() {
		return DrawScore
	}

	if stronger == Black {
		return BlackWinning
	}

	return WhiteWinning
}

// Lesser known endgames where we calculate endgame score markdown.
func (e *Evaluation) kingAndPawnsVsBareKing() int {
	our, their := e.strongerSide()

	pawns := e.position.pick(our).pawns
	row, col := coordinate(e.position.pick(their).home)

	// Pawns on A file with bare king opposing them.
	if (pawns & ^maskFile[A1]).noneʔ() && (pawns & ^maskInFront[their&1][row * 8]).noneʔ() && col <= B1 {
		return DrawScore
	}

	// Pawns on H file with bare king opposing them.
	if (pawns & ^maskFile[H1]).noneʔ() && (pawns & ^maskInFront[their&1][row * 8 + 7]).noneʔ() && col >= G1 {
		return DrawScore
	}

	return ExistingScore
}

// Bishop-only endgame: drop the score if we have opposite-colored bishops.
func (e *Evaluation) bishopsAndPawns() int {
	if e.oppositeBishopsʔ() {
		if abs(e.position.white.pawns.count() - e.position.black.pawns.count()) == 1 {
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
		wN, bN := e.position.white.knights.count(), e.position.black.knights.count()
		wR, bR := e.position.white.rooks.count(), e.position.black.rooks.count()
		extraPawns := abs(e.position.white.pawns.count() - e.position.black.pawns.count())

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

	unstoppableʔ := func(our int, square int) bool {
		their := our^1

		if (p.pick(their).all & maskInFront[our&1][square]).noneʔ() {
			mask := maskNone
			if p.color == our {
				mask = maskSquare[our&1][square]
			} else {
				mask = maskSquareEx[our&1][square]
			}
			return (mask & p.pick(their).king).noneʔ()
		}

		return false
	}

	// Check if either side has unstoppable pawn.
	white := unstoppableʔ(White, p.white.pawns.first())
	black := unstoppableʔ(Black, p.black.pawns.first())
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
	our, their := p.colors()
	piece := pawn(our)
	pawns := p.pick(our).pawns
	square := pawns.first()
	if rank(our, pawns.first()) < A5H5 || (pawns & (maskFile[0] | maskFile[7])).anyʔ() {

		// Temporarily remove opponent's pawn.
		piece = pawn(their)
		pawns = p.pick(their).pawns		// -> Save: opponent's pawn bitmask.
		square = pawns.first()			// -> Save: opponent's pawn square.

		p.pick(their).pawns = maskNone
		p.pieces[square] = Piece(0)

		// Temporarily adjust score so that when e.strongerSide() gets called
		// by kingAndPawnVsBareKing() it returns side to move.
		score := e.score.endgame		// -> Save: endgame score.
		e.score.endgame = let(our == White, 1, -1)

		// When we're done restore original endgame score and opponent's pawn.
		defer func() {
			e.score.endgame = score		// <- Restore: endgame score.
			p.pieces[square] = piece	// <- Restore: opponent's pawn square.
			p.pick(their).pawns = pawns	// <- Restore: opponent's pawn bitmask.
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
	if e.position.white.pawns.singleʔ() && e.position.black.pawns.singleʔ() {
		return e.fraction(3, 4) // 3/4
	}

	return ExistingScore
}

// One side has 0 pawn and the other side has 0 or more pawns.
func (e *Evaluation) noPawnsLeft() int {
	p := e.position
	our, their := e.strongerSide()
	whiteMinorOnly := p.white.queens.noneʔ() && p.white.rooks.noneʔ() && (p.white.bishops | p.white.knights).singleʔ()
	blackMinorOnly := p.black.queens.noneʔ() && p.black.rooks.noneʔ() && (p.black.bishops | p.black.knights).singleʔ()

	// Check for opposite bishops first.
	if whiteMinorOnly && blackMinorOnly && p.white.knights.noneʔ() && p.black.knights.noneʔ() && e.oppositeBishopsʔ() {
		pawn := p.pick(our).pawns			// The passer.
		king := p.pick(their).king			// Defending king.
		path := maskInFront[our&1][pawn.first()]	// Path in front of the passer.
		safe := maskDark				// Safe squares for king to block on.
		if (p.pick(our).bishops & maskDark).anyʔ() {
			safe = ^maskDark
		}

		// Draw if king blocks the passer on safe square.
		if (king & safe).anyʔ() && (king & path).anyʔ() {
			return DrawScore
		}

		// Draw if bishop attacks a square in front of the passer.
		if (e.position.bishopAttacks(their) & path).anyʔ() {
			return DrawScore
		}
	}

	if our == White && p.white.pawns.noneʔ() {
		if whiteMinorOnly {
			// There is a theoretical chance of winning if opponent's pawns are on
			// edge files (ex. some puzzles).
			if (p.black.pawns & (maskFile[0] | maskFile[7])).anyʔ() {
				return e.fraction(1, 64) // 1/64
			}
			return DrawScore
		} else if blackMinorOnly {
			return e.fraction(1, 16) // 1/16
		}
		return e.fraction(3, 16) // 3/16
	}

	if our == Black && p.black.pawns.noneʔ() {
		if blackMinorOnly {
			if (p.white.pawns & (maskFile[0] | maskFile[7])).anyʔ() {
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

