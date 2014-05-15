// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

type Total struct {
	white Score
	black Score
}

type Evaluator struct {
	phase     int
	score     Score
	attacks   [2]int
	threats   [2]int
	summary   map[string]interface{}
	position  *Position
}

// Use single statically allocated variable to avoid garbage collection overhead.
var eval Evaluator

func (p *Position) Evaluate() int {
	eval = Evaluator{p.phase(), p.tally, [2]int{0, 0}, [2]int{0, 0}, nil, p}
	return eval.run()
}

func (p *Position) EvaluateWithTrace() (int, map[string]interface{}) {
	eval = Evaluator{p.phase(), p.tally, [2]int{0, 0}, [2]int{0, 0}, make(map[string]interface{}), p}

	Settings.Trace = true
	defer func() {
		var tempo Total
		var final Score

		if p.color == White {
			tempo.white.add(rightToMove)
			final.add(eval.score)
		} else {
			tempo.black.add(rightToMove)
			final.subtract(eval.score)
		}

		eval.checkpoint(`Phase`, eval.phase)
		eval.checkpoint(`PST`, p.tally)
		eval.checkpoint(`Tempo`, tempo)
		eval.checkpoint(`Final`, final)
		Settings.Trace = false
	}()

	return eval.run(), eval.summary
}

func (e *Evaluator) run() int {
	e.analyzePieces()
	e.analyzePawns()
	e.analyzeSafety()

	if e.position.color == White {
		e.score.add(rightToMove)
	} else {
		e.score.subtract(rightToMove)
		e.score.midgame = -e.score.midgame
		e.score.endgame = -e.score.endgame
	}

	return e.score.blended(e.phase)
}

func (e *Evaluator) checkpoint(tag string, total interface{}) {
	e.summary[tag] = total
}

func (e *Evaluator) strongEnough(color int) bool {
	p := e.position
	return p.count[queen(color)] > 0 &&
		(p.count[rook(color)] > 0 || p.count[bishop(color)] > 0 || p.count[knight(color)] > 0)
}
