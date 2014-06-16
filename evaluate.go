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
	threats int 		// A sum of treats: each based on attacking piece type.
	attacks int 		// Number of attacks on squares adjacent to the king.
	attackers int 		// Number of pieces attacking king's fort.
}

// Helper structure used for evaluation tracking.
type Total struct {
	white Score 		// Score for white.
	black Score 		// Score for black.
}

//
type Evaluation struct {
	flags     uint8 	 // Evaluation flags.
	score     Score 	 // Current score.
	safety    [2]Safety 	 // King safety for both sides.
	attacks   [14]Bitmask 	 // Attack bitmasks for all the pieces on the board.
	pawns     *PawnEntry 	 // Pointer to the pawn cache entry.
	material  *MaterialEntry // Pointer to the matrial cache entry.
	position  *Position 	 // Pointer to the position we're evaluating.
	metrics   Metrics 	 // Evaluation metrics when tracking is on.
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

		eval.checkpoint(`Phase`, eval.material.phase)
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
	e.position = p

	// Initialize the score with incremental PST value and right to move.
	e.score = p.tally
	if e.position.color == White {
		e.score.add(rightToMove)
	} else {
		e.score.subtract(rightToMove)
	}

	// Set up king and pawn attacks for both sides.
	e.attacks[King] = p.kingAttacks(White)
	e.attacks[Pawn] = p.pawnAttacks(White)
	e.attacks[BlackKing] = p.kingAttacks(Black)
	e.attacks[BlackPawn] = p.pawnAttacks(Black)

	// Overall attacks for both sides include kings and pawns so far.
	e.attacks[White] = e.attacks[King] | e.attacks[Pawn]
	e.attacks[Black] = e.attacks[BlackKing] | e.attacks[BlackPawn]

	return e
}

func (e *Evaluation) run() int {
	e.analyzeMaterial()
	if e.material.flags & materialDraw != 0 {
		return 0
	}
	if e.material.flags & knownEndgame != 0 {
		return e.analyzeEndgame()
	}

	// Initialize king fort bitmasks only when we need them.
	if e.material.flags & whiteKingSafety != 0 {
		e.safety[White].fort = e.setupFort(White)
	}
	if e.material.flags & blackKingSafety != 0 {
		e.safety[Black].fort = e.setupFort(Black)
	}

	e.analyzePawns()
	e.analyzePieces()
	e.analyzeThreats()
	e.analyzeSafety()
	e.analyzePassers()
	if e.material.flags & lesserKnownEndgame != 0 {
		e.inspectEndgame()
	}

	if e.position.color == Black {
		e.score.midgame = -e.score.midgame
		e.score.endgame = -e.score.endgame
	}

	return e.score.blended(e.material.phase)
}

func (e *Evaluation) checkpoint(tag string, metric interface{}) {
	e.metrics[tag] = metric
}

func (e *Evaluation) setupFort(color int) (bitmask Bitmask) {
	bitmask = e.attacks[king(color)] | e.attacks[king(color)].pushed(color)
	switch e.position.king[color] {
	case A1, A8:
		bitmask |= e.attacks[king(color)] << 1
	case H1, H8:
		bitmask |= e.attacks[king(color)] >> 1
	}
	return
}
