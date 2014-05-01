// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

type Score struct {
	midgame int
	endgame int
}

type Evaluator struct {
	stage    int
	midgame  int
	endgame  int
	attacks  [2]int
	threats  [2]int
	position *Position
}

func (p *Position) Evaluate() int {
	evaluator := &Evaluator{0, 0, 0, [2]int{0, 0}, [2]int{0, 0}, p}
	evaluator.analyzeMaterial()
	evaluator.analyzePieces()
	evaluator.analyzePawns()
	evaluator.analyzeSafety()

	if p.color == White {
		evaluator.midgame += rightToMove.midgame
		evaluator.endgame += rightToMove.endgame
		return p.score(evaluator.midgame, evaluator.endgame)
	}
	evaluator.midgame -= rightToMove.midgame
	evaluator.endgame -= rightToMove.endgame
	return p.score(-evaluator.midgame, -evaluator.endgame)
}

func (e *Evaluator) analyzeMaterial() {
	counters := &e.position.count
	for _, piece := range []Piece{Pawn, Knight, Bishop, Rook, Queen} {
		count := counters[piece] - counters[piece|Black]
		midgame, endgame := piece.value()
		e.midgame += midgame * count
		e.endgame += endgame * count
	}
}

func (e *Evaluator) strongEnough(color int) bool {
	p := e.position
	return p.count[queen(color)] > 0 &&
		(p.count[rook(color)] > 0 || p.count[bishop(color)] > 0 || p.count[knight(color)] > 0)
}
