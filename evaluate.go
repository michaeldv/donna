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
	position *Position
}

func (p *Position) Evaluate() int {
	evaluator := &Evaluator{0, rightToMove.midgame, rightToMove.endgame, p}
	evaluator.analyzeMaterial()
	evaluator.analyzeCoordination()
	evaluator.analyzePawnStructure()
	evaluator.analyzeKingSafety()

	if p.color == Black {
		evaluator.midgame = -evaluator.midgame
		evaluator.endgame = -evaluator.endgame
	}
	return p.score(evaluator.midgame, evaluator.endgame)
}

func (e *Evaluator) strongEnough(color int) bool {
	p := e.position
	return p.count[queen(color)] > 0 &&
		(p.count[rook(color)] > 0 || p.count[bishop(color)] > 0 || p.count[knight(color)] > 0)
}
