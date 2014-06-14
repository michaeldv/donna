// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (e *Evaluation) analyzeEndgame() int {
	return e.material.endgame(e)
}

func (e *Evaluation) inspectEndgame() {
	if e.score.endgame != 0 {
		markdown := e.material.endgame(e)
		if markdown > 0 {
			e.score.endgame *= markdown / 128
		}
	}
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
	return -1 // 96
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


