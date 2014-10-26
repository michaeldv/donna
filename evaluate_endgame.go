// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
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
	markdown := e.material.endgame(e)
	if markdown == 0 {
		e.score = Score{0, 0}
	} else if markdown > 0xFFFF {
		mul, div := markdown >> 16, markdown & 0xFFFF
		e.score.endgame = e.score.endgame * mul / div
	} else if markdown > 0 {
		e.score.endgame /= markdown
	}
}

func (e *Evaluation) strongerSide() int {
	if e.score.endgame > 0 {
		return White
	}
	return Black
}

// Known endgames where we calculate the exact score.
func (e *Evaluation) winAgainstBareKing() int { // STUB.
	return e.score.blended(e.material.phase)
}

func (e *Evaluation) knightAndBishopVsBareKing() int { // STUB.
	return e.score.blended(e.material.phase)
}

func (e *Evaluation) twoBishopsVsBareKing() int { // STUB.
	return e.score.blended(e.material.phase)
}

func (e *Evaluation) kingAndPawnVsBareKing() int {
	var color, wKing, bKing, wPawn int

	if e.strongerSide() == White {
		color = e.position.color
		wKing = e.position.king[White]
		bKing = e.position.king[Black]
		wPawn = e.position.outposts[Pawn].last()
	} else {
		color = e.position.color^1
		wKing = 64 + ^e.position.king[Black]
		bKing = 64 + ^e.position.king[White]
		wPawn = 64 + ^e.position.outposts[BlackPawn].last()
	}

	index := color + (wKing << 1) + (bKing << 7) + ((wPawn - 8) << 13)
	if bitbase[index / 64] & (1 << uint(index & 0x3F)) == 0 {
		return 0
	}

	if color == Black {
		return -DecisiveAdvantage
	}
	return DecisiveAdvantage
}

// Lesser known endgames where we calculate endgame score markdown.
func (e *Evaluation) kingAndPawnsVsBareKing() int {
	color := e.strongerSide()

	pawns := e.position.outposts[pawn(color)]
	row, col := coordinate(e.position.king[color^1])

	// Pawns on A file with bare king opposing them.
	if (pawns & ^maskFile[A1] == 0) && (pawns & ^maskInFront[color^1][row*8] == 0) && col <= B1 {
		return 0
	}

	// Pawns on H file with bare king opposing them.
	if (pawns & ^maskFile[H1] == 0) && (pawns & ^maskInFront[color^1][row*8+7] == 0) && col >= G1 {
		return 0
	}

	return -1
}

// Bishop-only endgame: drop the score if we have opposite-colored bishops.
func (e *Evaluation) bishopsAndPawns() int {
	if e.oppositeBishops() {
		if pawns := abs(e.position.count[Pawn] - e.position.count[BlackPawn]); pawns == 1 {
			return 8 // --> 1/8 of original score.
		}
		return 2 // --> 1/2 of original score.
	}

	return -1
}

// Single bishops plus some minors: drop the score if we have opposite-colored bishops.
func (e *Evaluation) drawishBishops() int {
	if e.oppositeBishops() {
		return 4 // --> 1/4 of original score.
	}
	return -1
}

func (e *Evaluation) kingAndPawnVsKingAndPawn() int { // STUB.
	return -1 // 96
}

func (e *Evaluation) bishopAndPawnVsBareKing() int { // STUB.
	return -1 // 96
}

func (e *Evaluation) rookAndPawnVsRook() int { // STUB.
	return -1 // 96
}

func (e *Evaluation) queenVsRookAndPawns() int { // STUB.
	return -1 // 96
}

func (e *Evaluation) lastPawnLeft() int {
	color := e.strongerSide()
	count := &e.position.count

	if (color == White && count[Pawn] == 1) || (color == Black && count[BlackPawn] == 1) {
		return (3 << 16) | 4 // --> 3/4 of original score.
	}

	return -1
}

func (e *Evaluation) noPawnsLeft() int {
	color := e.strongerSide()
	count := &e.position.count
	whiteMinorOnly := count[Queen] + count[Rook] == 0 && count[Bishop] + count[Knight] == 1
	blackMinorOnly := count[BlackQueen] + count[BlackRook] == 0 && count[BlackBishop] + count[BlackKnight] == 1

	if color == White && count[Pawn] == 0 {
		if whiteMinorOnly {
			// There is a theoretical chance of winning if opponent's pawns are on
			// edge files (ex. some puzzles).
			if e.position.outposts[BlackPawn] & (maskFile[0] | maskFile[7]) != 0 {
				return 64 // --> 1/64 of original score.
			}
			return 0
		} else if blackMinorOnly {
			return 16 // --> 1/16 of original score.
		}
		return (3 << 16) | 16 // --> 3/16 of original score.
	}

	if color == Black && count[BlackPawn] == 0 {
		if blackMinorOnly {
			if e.position.outposts[Pawn] & (maskFile[0] | maskFile[7]) != 0 {
				return 64 // --> 1/64 of original score.
			}
			return 0
		} else if whiteMinorOnly {
			return 16 // --> 1/16 of original score.
		}
		return (3 << 16) | 16 // --> 3/16 of original score.
	}

	return -1
}
