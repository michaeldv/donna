// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import `testing`

// Attacks by minor piece.
func TestEvaluateThreats000(t *testing.T) {
	// Baseline: bishop defended by pawn.
	NewGame(`Kh1,Ne4`, `Kg7,Bf6,e7`).Start(White).EvaluateWithTrace()
	baseline := eval.metrics[`Threats`].(Total).white

	// Bishop not defended by pawn.
	NewGame(`Kh1,Ne4,a2`, `Ke7,Bf6,a7`).Start(White).EvaluateWithTrace()
	score := eval.metrics[`Threats`].(Total).white
	expect(t, score.minus(baseline), bonusMinorThreat[Bishop/2])

	// Bishop and rook not defended by pawn (rook is stronger).
	NewGame(`Kh1,Ne4,a2`, `Ke7,Bf6,Rd6,a7`).Start(White).EvaluateWithTrace()
	score = eval.metrics[`Threats`].(Total).white
	expect(t, score.minus(baseline), bonusMinorThreat[Rook/2])

	// Hanging bishop with extra bonus for the right to move.
	NewGame(`Kh1,Ne4,a2`, `Ka8,Bf6,a7`).Start(White).EvaluateWithTrace()
	score = eval.metrics[`Threats`].(Total).white
	expect(t, score.minus(baseline), bonusMinorThreat[Bishop/2].plus(hangingAttack.times(2)))
}

// Attacks by major piece.
func TestEvaluateThreats010(t *testing.T) {
	// Baseline: bishop defended by pawn.
	NewGame(`Kh1,Rf1`, `Kg7,Bf6,e7`).Start(White).EvaluateWithTrace()
	baseline := eval.metrics[`Threats`].(Total).white

	// Bishop not defended by pawn.
	NewGame(`Kh1,Rf1`, `Ke7,Bf6`).Start(White).EvaluateWithTrace()
	score := eval.metrics[`Threats`].(Total).white
	expect(t, score.minus(baseline), bonusMajorThreat[Bishop/2])

	// Bishop and queen not defended by pawn (queen is stronger).
	NewGame(`Kh1,Rf1`, `Ke7,Qa1,Bf6`).Start(White).EvaluateWithTrace()
	score = eval.metrics[`Threats`].(Total).white
	expect(t, score.minus(baseline), bonusMajorThreat[Queen/2])

	// Hanging bishop with extra bonus for the right to move.
	NewGame(`Kh1,Rf1`, `Kh8,Bf6`).Start(White).EvaluateWithTrace()
	score = eval.metrics[`Threats`].(Total).white
	expect(t, score.minus(baseline), bonusMajorThreat[Bishop/2].plus(hangingAttack.times(2)))
}
