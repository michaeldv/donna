// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

// const = brains * looks * availability
const (
	whiteKingSafety    = 0x01  // Should we worry about white king's safety?
	blackKingSafety    = 0x02  // Ditto for the black king.
	materialDraw       = 0x04  // King vs. King (with minor)
	knownEndgame       = 0x08  // Where we calculate exact score.
	lesserKnownEndgame = 0x10  // Where we set score markdown value.
	singleBishops      = 0x20  // Sides might have bishops on opposite color squares.
)

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

type Function func(*Evaluation) int
type MaterialEntry struct {
	score     Score 	// Score adjustment for the given material.
	endgame   Function 	// Function to analyze an endgame position.
	phase     int 		// Game phase based on available material.
	turf      int 		// Home turf score for the game opening.
	flags     uint8    	// Evaluation flags based on material balance.
}

type Evaluation struct {
	score     Score 	 // Current score.
	safety    [2]Safety 	 // King safety data for both sides.
	attacks   [14]Bitmask 	 // Attack bitmasks for all the pieces on the board.
	pins      [2]Bitmask     // Bitmask of pinned pieces for both sides.
	pawns     *PawnEntry 	 // Pointer to the pawn cache entry.
	material  *MaterialEntry // Pointer to the matrial base entry.
	position  *Position 	 // Pointer to the position we're evaluating.
	metrics   Metrics 	 // Evaluation metrics when tracking is on.
}

// Use single statically allocated variable to avoid garbage collection overhead.
var eval Evaluation

// The following statement is true. The previous statement is false. Main position
// evaluation method that returns single blended score.
func (p *Position) Evaluate() int {
	return eval.init(p).run()
}

// Auxiliary evaluation method that captures individual evaluation metrics. This
// is useful when we want to see evaluation summary.
func (p *Position) EvaluateWithTrace() (int, Metrics) {
	eval.init(p)
	eval.metrics = make(Metrics)

	engine.trace = true
	defer func() {
		var tempo Total
		var final Score

		if p.color == White {
			tempo.white.add(rightToMove)
			final.add(eval.score)
		} else {
			tempo.black.add(rightToMove)
			final.sub(eval.score)
		}

		eval.checkpoint(`Phase`, eval.material.phase)
		eval.checkpoint(`Imbalance`, eval.material.score)
		eval.checkpoint(`PST`, p.tally)
		eval.checkpoint(`Tempo`, tempo)
		eval.checkpoint(`Final`, final)
		engine.trace = false
	}()

	return eval.run(), eval.metrics
}

func (e *Evaluation) init(p *Position) *Evaluation {
	eval = Evaluation{}
	e.position = p

	// Initialize the score with incremental PST value and right to move.
	e.score = p.tally
	if p.color == White {
		e.score.add(rightToMove)
	} else {
		e.score.sub(rightToMove)
	}

	// Set up king and pawn attacks for both sides.
	e.attacks[King] = p.kingAttacks(White)
	e.attacks[Pawn] = p.pawnAttacks(White)
	e.attacks[BlackKing] = p.kingAttacks(Black)
	e.attacks[BlackPawn] = p.pawnAttacks(Black)

	// Overall attacks for both sides include kings and pawns so far.
	e.attacks[White] = e.attacks[King] | e.attacks[Pawn]
	e.attacks[Black] = e.attacks[BlackKing] | e.attacks[BlackPawn]

	// Pinned pieces for both sides that have restricted mobility.
	e.pins[White] = p.pins(p.king[White])
	e.pins[Black] = p.pins(p.king[Black])

	return e
}

func (e *Evaluation) run() int {
	e.material = &materialBase[e.position.balance]

	e.score.add(e.material.score)
	if e.material.flags & knownEndgame != 0 {
		return e.evaluateEndgame()
	}

	e.analyzePawns()
	e.analyzePieces()
	e.analyzeThreats()
	e.analyzeSafety()
	e.analyzePassers()
	e.wrapUp()

	return e.score.blended(e.material.phase)
}

func (e *Evaluation) wrapUp() {

	// Adjust the endgame score if we have lesser known endgame.
	if e.score.endgame != 0 && e.material.flags & lesserKnownEndgame != 0 {
		e.inspectEndgame()
	}

	// Flip the sign for black so that blended evaluation score always
	// represents the white side.
	if e.position.color == Black {
		e.score.midgame = -e.score.midgame
		e.score.endgame = -e.score.endgame
	}
}

func (e *Evaluation) checkpoint(tag string, metric interface{}) {
	e.metrics[tag] = metric
}

func (e *Evaluation) oppositeBishops() bool {
	bishops := e.position.outposts[Bishop] | e.position.outposts[BlackBishop]

	return bishops & maskDark != 0 && bishops & ^maskDark != 0
}

// Returns true if material tables indicate that the king can't defend himself.
func (e *Evaluation) isKingUnsafe(color int) bool {
	if color == White {
		return e.material.flags & whiteKingSafety != 0
	}

	return e.material.flags & blackKingSafety != 0
}
