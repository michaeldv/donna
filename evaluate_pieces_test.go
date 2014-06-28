// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import `testing`

// Rooks.
func TestEvaluatePieces000(t *testing.T) {
	NewGame(`Ke4,Ne1,Ra1,Rh1,a2,b2,c2,f2,g2,h2`, `Ke8,e7`).Start(White).EvaluateWithTrace()
	baseline := eval.metrics[`-Rooks`].(Total).white

	// H1 rook boxed (can castle): 1x penalty.
	NewGame(`Ke1,Ra1,Rh1,a2,b2,c2,f2,g2,h2`, `Ke8,e7`).Start(White).EvaluateWithTrace()
	rooks := eval.metrics[`-Rooks`].(Total).white
	expect(t, rooks.minus(baseline), rookBoxed.times(-1))

	// H1 rook boxed (can't castle): 2x penalty.
	NewGame(`Ke1,Ra1,Rf1,a2,b2,c2,f2,g2,h2`, `Ke8,e7`).Start(White).EvaluateWithTrace()
	boxed := eval.metrics[`-Rooks`].(Total).white
	expect(t, boxed.minus(baseline), rookBoxed.times(-2))
}
