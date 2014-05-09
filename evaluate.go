// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

type Evaluator struct {
	phase    int
	score    Score
	attacks  [2]int
	threats  [2]int
	position *Position
}

// Use single statically allocated variable to avoid garbage collection overhead.
var eval Evaluator

func (p *Position) Evaluate() int {
	eval = Evaluator{p.phase(), Score{0, 0}, [2]int{0, 0}, [2]int{0, 0}, p}
	eval.analyzeMaterial()
	eval.analyzePieces()
	eval.analyzePawns()
	eval.analyzeSafety()

	if p.color == White {
		eval.score.add(rightToMove)
	} else {
		eval.score.subtract(rightToMove)
		eval.score.midgame = -eval.score.midgame
		eval.score.endgame = -eval.score.endgame
	}

	return eval.score.blended(eval.phase)
}

func (e *Evaluator) analyzeMaterial() {
	counters := &e.position.count
	for _, piece := range []Piece{Pawn, Knight, Bishop, Rook, Queen} {
		count := counters[piece] - counters[piece|Black]
		e.score.add(piece.value().times(count))
	}
}

func (e *Evaluator) strongEnough(color int) bool {
	p := e.position
	return p.count[queen(color)] > 0 &&
		(p.count[rook(color)] > 0 || p.count[bishop(color)] > 0 || p.count[knight(color)] > 0)
}
