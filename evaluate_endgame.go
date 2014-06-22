// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (e *Evaluation) analyzeEndgame() int {
	score := e.material.endgame(e)
	if e.position.color == Black {
		return -score
	}
	return score
}

func (e *Evaluation) inspectEndgame() {
	if e.score.endgame != 0 {
		markdown := e.material.endgame(e)
		if markdown >= 0 {
			e.score.endgame *= markdown / 128
		}
	}
}

func (e *Evaluation) strongerSide() int {
	if e.score.endgame > 0 {
		return White
	}
	return Black
}

// Known endgames where we calculate the exact score.
func (e *Evaluation) winAgainstBareKing() int {
	return e.score.blended(e.material.phase)
}

func (e *Evaluation) knightAndBishopVsBareKing() int {
	return e.score.blended(e.material.phase)
}

func (e *Evaluation) kingAndPawnVsBareKing() int {
	return e.score.blended(e.material.phase)
}

// Lesser known endgames where we calculate endgame score markdown.
func (e *Evaluation) kingAndPawnsVsBareKing() int {
	color := e.strongerSide()

	pawns := e.position.outposts[pawn(color)]
	row, col := Coordinate(e.position.king[color^1])

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

func (e *Evaluation) kingAndPawnVsKingAndPawn() int {
	return -1 // 96
}

func (e *Evaluation) bishopAndPawnVsBareKing() int {
	return -1 // 96
}

func (e *Evaluation) rookAndPawnVsRook() int {
	return -1 // 96
}

func (e *Evaluation) queenVsRookAndPawns() int {
	return -1 // 96
}


