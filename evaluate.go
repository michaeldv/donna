// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

// Hash containing various evaluation metrics; used only when evaluation tracing
// is enabled.
type Metrics map[string]interface{}

// King safety information; used only in the middle game when there is enough
// material to worry about the king safety.
type Safety struct {
	fort Bitmask 		// Squares around the king plus one extra row in front.
	threat int 		// A sum of treats: each based on attacking piece type.
	homeAttacks int 	// Number of attacks on squares adjacent to the king.
	fortAttackers int 	// Number of pieces attacking king's fort.
}

// Helper structure used for evaluation tracking.
type Total struct {
	white Score 		// Score for white.
	black Score 		// Score for black.
}

//
type Evaluation struct {
	phase     int 		// Game phase based on available material.
	flags     uint8 	// Evaluation flags.
	score     Score 	// Current score.
	king      [2]Safety 	// King safety for both sides.
	targets   [14]Bitmask 	// Attack targets for all the pieces on board.
	metrics   Metrics 	// Evaluation metrics when tracking is on.
	position  *Position 	// Position we're evaluating.
}

// Use single statically allocated variable to avoid garbage collection overhead.
var eval Evaluation

// Main position evaluation method that returns single blended score.
func (p *Position) Evaluate() int {
	return eval.init(p).run()
}

// Auxiliary evaluation method that captures individual evaluation metrics. This
// is useful when we want to see evaluation summary.
func (p *Position) EvaluateWithTrace() (int, Metrics) {
	eval.init(p)
	eval.metrics = make(Metrics)

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

	return eval.run(), eval.metrics
}

// Evaluation method for use in tests. It invokes evaluation that captures the
// metrics, and returns the requested metric score.
func (p *Position) EvaluateTest(tag string) (score Score, metrics Metrics) {
	_, metrics = p.EvaluateWithTrace()

	switch metrics[tag].(type) {
	case Score:
		score = metrics[tag].(Score)
	case Total:
		if p.color == White {
			score = metrics[tag].(Total).white
		} else {
			score = metrics[tag].(Total).black
		}
	}
	return
}

func (e *Evaluation) init(p *Position) *Evaluation {
	eval = Evaluation{}
	e.phase = p.phase()
	e.score = p.tally
	e.position = p

	// Set up king and pawn attack targets for both sides.
	e.targets[King] = p.kingAttacks(White)
	e.targets[Pawn] = p.pawnAttacks(White)
	e.targets[BlackKing] = p.kingAttacks(Black)
	e.targets[BlackPawn] = p.pawnAttacks(Black)

	// Overall attack targets for both sides include kings and pawns so far.
	e.targets[White] = e.targets[King] | e.targets[Pawn]
	e.targets[Black] = e.targets[BlackKing] | e.targets[BlackPawn]

	// TODO: initialize only if we are going to evaluate king's safety.
	e.king[White].fort = e.targets[King].pushed(White)
	e.king[Black].fort = e.targets[BlackKing].pushed(Black)

	return e
}

func (e *Evaluation) run() int {
	e.analyzePawns()
	e.analyzePieces()
	e.analyzeThreats()
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

func (e *Evaluation) checkpoint(tag string, metric interface{}) {
	e.metrics[tag] = metric
}

func (e *Evaluation) strongEnough(color int) bool {
	p := e.position
	return p.count[queen(color)] > 0 &&
		(p.count[rook(color)] > 0 || p.count[bishop(color)] > 0 || p.count[knight(color)] > 0)
}
