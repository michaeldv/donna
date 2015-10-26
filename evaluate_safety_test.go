// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`github.com/michaeldv/donna/expect`; `testing`)

// f2/g2/h2 and f7/g7/h7 (perfect cover).
func TestSafety000(t *testing.T) {
	NewGame(`Kg1,Qd1,Ra1,Rf1,Nf3,a2,d4,f2,g2,h2`, `M,Kg8,Qd8,Ra8,Rf8,Nf6,a7,d5,f7,g7,h7`).start().EvaluateWithTrace()
	white := eval.metrics[`-Cover`].(Total).white
	black := eval.metrics[`-Cover`].(Total).black

	expect.Eq(t, black.midgame, 167)
	expect.Eq(t, black.endgame, 0)
	expect.Eq(t, white, black)
}

// f2,g2,H3 vs f7/G6/h7 (one square pawn distance).
func TestSafety010(t *testing.T) {
	NewGame(`Kg1,Qd1,Ra1,Rf1,Nf3,a2,d4,f2,g2,h3`, `M,Kg8,Qd8,Ra8,Rf8,Nf6,a7,d5,f7,g6,h7`).start().EvaluateWithTrace()
	white := eval.metrics[`-Cover`].(Total).white
	black := eval.metrics[`-Cover`].(Total).black

	expect.Eq(t, black.midgame, 163 - penaltyCover[2])
	expect.Eq(t, white, black)
}

// F4,g2,h2 vs f7/g7/H5 (two squares pawn distance).
func TestSafety020(t *testing.T) {
	NewGame(`Kg1,Qd1,Ra1,Rf1,Nf3,a2,d4,f4,g2,h2`, `M,Kg8,Qd8,Ra8,Rf8,Nf6,a7,d5,f7,g7,h5`).start().EvaluateWithTrace()
	white := eval.metrics[`-Cover`].(Total).white
	black := eval.metrics[`-Cover`].(Total).black

	expect.Eq(t, black.midgame, 157 - penaltyCover[3])
	expect.Eq(t, white, black)
}

// F5,g2,h2 vs f7/g7/H4 (three squares pawn distance).
func TestSafety030(t *testing.T) {
	NewGame(`Kg1,Qd1,Ra1,Rf1,Nf3,a2,d4,f5,g2,h2`, `M,Kg8,Qd8,Ra8,Rf8,Nf6,a7,d5,f7,g7,h4`).start().EvaluateWithTrace()
	white := eval.metrics[`-Cover`].(Total).white
	black := eval.metrics[`-Cover`].(Total).black

	expect.Eq(t, black.midgame, 155 - penaltyCover[4])
	expect.Eq(t, white, black)
}

// F4,G3,h2 vs f7/G6/H5 (one and two squares pawn distances).
func TestSafety040(t *testing.T) {
	NewGame(`Kg1,Qd1,Ra1,Rf1,Nf3,a2,d4,f4,g3,h2`, `M,Kg8,Qd8,Ra8,Rf8,Nf6,a7,d5,f7,g6,h5`).start().EvaluateWithTrace()
	white := eval.metrics[`-Cover`].(Total).white
	black := eval.metrics[`-Cover`].(Total).black

	expect.Eq(t, black.midgame, 154 - penaltyCover[2] - penaltyCover[3])
	expect.Eq(t, white, black)
}

// F3,F4,g2,h2 vs F6/F5,g7/h7 (one square, doubled pawns).
func TestSafety100(t *testing.T) {
	NewGame(`Kg1,Qd1,Ra1,Rf1,Nf3,a2,d4,f3,f4,g2,h2`, `M,Kg8,Qd8,Ra8,Rf8,Nf6,a7,d5,f6,f5,g7,h7`).start().EvaluateWithTrace()
	white := eval.metrics[`-Cover`].(Total).white
	black := eval.metrics[`-Cover`].(Total).black

	expect.Eq(t, white.midgame, 163 - penaltyCover[2])
	expect.Eq(t, black.midgame, 163 - penaltyCover[2])
	expect.Eq(t, white, black)
}

// Kg2,f2,G3,h2 vs Kg7,f7,G7/h7 (ajacent).
func TestSafety110(t *testing.T) {
	NewGame(`Kg2,Qd1,Ra1,Rf1,Nf3,a2,d4,f2,g3,h2`, `M,Kg7,Qd8,Ra8,Rf8,Nf6,a7,d5,f7,g6,h7`).start().EvaluateWithTrace()
	white := eval.metrics[`-Cover`].(Total).white
	black := eval.metrics[`-Cover`].(Total).black

	expect.Eq(t, white.midgame, 167 - penaltyCover[0])
	expect.Eq(t, black.midgame, 167 - penaltyCover[0])
	expect.Eq(t, white, black)
}

func TestSafety120(t *testing.T) {
	game := NewGame(`Ke1,Qf3,Ra1,Rh1,Bc1,Bf1,Nc3,a2,b2,c2,d4,e3,f2,g2,h3`, `Ke8,Qd8,Ra8,Rh8,Bf8,Nc6,Nf6,a7,b7,c7,d5,e7,f7,g7,h7`)
	game.start().EvaluateWithTrace()
	white := eval.metrics[`-Cover`].(Total).white

	expect.Eq(t, white.midgame, 149)
}

func TestSafety130(t *testing.T) {
	game := NewGame(`Ke1,Qd1,Ra1,Rh1,Bc1,Bf1,Nc3,a2,b2,c2,d4,e3,f2,f3,h3`, `Ke8,Qd8,Ra8,Rh8,Bf8,Nc6,Nf6,a7,b7,c7,d5,e7,f7,g7,h7`)
	game.start().EvaluateWithTrace()
	white := eval.metrics[`-Cover`].(Total).white

	expect.Eq(t, white.midgame, 119)
}

// Friendly pawn distance.
func TestSafety200(t *testing.T) {
	NewGame(`Ke1,Qd1`, `M,Ke8,Qd8,f7`).start().EvaluateWithTrace()
	black := eval.metrics[`-Cover`].(Total).black
	expect.Eq(t, black.endgame, 0)

	// NewGame(`Ke1,Qd1`, `M,Ke8,Qd8,g7`).start().EvaluateWithTrace()
	// black = eval.metrics[`-Cover`].(Total).black
	// expect.Eq(t, black.endgame, -kingByPawn.endgame * 1)
	//
	// NewGame(`Ke1,Qd1`, `M,Ke8,Qd8,h7`).start().EvaluateWithTrace()
	// black = eval.metrics[`-Cover`].(Total).black
	// expect.Eq(t, black.endgame, -kingByPawn.endgame * 2)
	//
	// NewGame(`Ke1,Qd1`, `M,Ka8,Qd8,h2`).start().EvaluateWithTrace()
	// black = eval.metrics[`-Cover`].(Total).black
	// expect.Eq(t, black.endgame, -kingByPawn.endgame * 6)
}
