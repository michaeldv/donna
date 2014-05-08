// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

type Evaluator struct {
	score    Score
	attacks  [2]int
	threats  [2]int
	position *Position
}

// Use single statically allocated variable to avoid garbage collection overhead.
var evaluator Evaluator

func (p *Position) Evaluate() int {
	evaluator = Evaluator{Score{0, 0}, [2]int{0, 0}, [2]int{0, 0}, p}
	evaluator.analyzeMaterial()
	evaluator.analyzePieces()
	evaluator.analyzePawns()
	evaluator.analyzeSafety()

	if p.color == White {
		evaluator.score.add(rightToMove)
		return p.blended(evaluator.score)
	} else {
		evaluator.score.subtract(rightToMove)
		evaluator.score.midgame = -evaluator.score.midgame
		evaluator.score.endgame = -evaluator.score.endgame
	}

	return p.blended(evaluator.score)
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
