// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`github.com/michaeldv/donna/expect`; `testing`)

// Doubled pawns.
func TestEvaluatePawns100(t *testing.T) {
	p := NewGame(`Ke1,h2,h3`, `Kd8,a7,a6`).start()
	score := p.Evaluate()
	expect.Eq(t, score, rightToMove.endgame) // Right to move only.
}

func TestEvaluatePawns110(t *testing.T) {
	game := NewGame(`Ke1,h2,h3`, `Ke8,a7,h7`)
	score := game.start().Evaluate()

	expect.Eq(t, score, -10)
}

func TestEvaluatePawns120(t *testing.T) {
	game := NewGame(`Ke1,f4,f5`, `Ke8,f7,h7`)
	score := game.start().Evaluate()

	expect.Eq(t, score, -31)
}

// Passed pawns.
func TestEvaluatePawns200(t *testing.T) {
	game := NewGame(`Ke1,a4,h4`, `Ke8,a5,h5`) // Blocked.
	score := game.start().Evaluate()

	expect.Eq(t, score, 10)
}

func TestEvaluatePawns210(t *testing.T) {
	game := NewGame(`Ke1,a4,h4`, `Ke8,a5,g7`) // Can't pass.
	score := game.start().Evaluate()

	expect.Eq(t, score, 8)
}

func TestEvaluatePawns220(t *testing.T) {
	game := NewGame(`Ke1,a4,e4`, `Ke8,a5,d6`) // Can't pass.
	score := game.start().Evaluate()

	expect.Eq(t, score, -2)
}

func TestEvaluatePawns230(t *testing.T) {
	game := NewGame(`Ke1,a4,e5`, `Ke8,a5,e4`) // Both passing.
	score := game.start().Evaluate()

	expect.Eq(t, score, 10)
}

func TestEvaluatePawns240(t *testing.T) {
	game := NewGame(`Kd1,e5`, `Ke8,d5`) // Both passing but white is closer.
	score := game.start().Evaluate()    // No KPKP since the pawn crossed A4H4.

	expect.Eq(t, score, 32)
}

func TestEvaluatePawns250(t *testing.T) {
	game := NewGame(`Ke1,a5,b2`, `Kd8,g7,h7`) // Both passing but white is much closer.
	score := game.start().Evaluate()

	expect.Eq(t, score, 101)
}

// Isolated pawns.
func TestEvaluatePawns300(t *testing.T) {
	game := NewGame(`Ke1,a5,c5`, `Kd8,f4,h4`) // All pawns are isolated.
	score := game.start().Evaluate()

	expect.Eq(t, score, 10)
}

func TestEvaluatePawns310(t *testing.T) {
	game := NewGame(`Ke1,a2,c2,e2`, `Ke8,a7,b7,c7`) // White pawns are isolated.
	score := game.start().Evaluate()

	expect.Eq(t, score, -28)
}

// Rooks.
func TestEvaluatePawns400(t *testing.T) {
	game := NewGame(`Ke1,Ra7`, `Ke8,Rh3`) // White on 7th.
	score := game.start().Evaluate()

	expect.Eq(t, score, 5)
}

func TestEvaluatePawns410(t *testing.T) {
	game := NewGame(`Ke1,Rb1,Ng2,a2`, `Ke8,Rh8,Nb7,h7`) // White on open file.
	score := game.start().Evaluate()

	expect.Eq(t, score, 74)
}

func TestEvaluatePawns420(t *testing.T) {
	game := NewGame(`Ke1,Rb1,a2,g2`, `Ke8,Rh8,h7,b7`) // White on semi-open file.
	score := game.start().Evaluate()

	expect.Eq(t, score, 84)
}

// King shield.
func TestEvaluatePawns500(t *testing.T) {
	game := NewGame(`Kg1,f2,g2,h2,Qa3,Na4`, `Kg8,f7,g7,h7,Qa6,Na5`) // h2,g2,h2 == f7,g7,h7
	score := game.start().Evaluate()

	expect.Eq(t, score, 10)
}
func TestEvaluatePawns505(t *testing.T) {
	game := NewGame(`Kg1,f2,g2,h2,Qa3,Na4`, `Kg8,f7,g6,h7,Qa6,Na5`) // h2,g2,h2 vs f7,G6,h7
	score := game.start().Evaluate()

	expect.Eq(t, score, 30)
}

func TestEvaluatePawns510(t *testing.T) {
	game := NewGame(`Kg1,f2,g2,h2,Qa3,Na4`, `Kg8,f5,g6,h7,Qa6,Na5`) // h2,g2,h2 vs F5,G6,h7
	score := game.start().Evaluate()

	expect.Eq(t, score, 40)
}

func TestEvaluatePawns520(t *testing.T) {
	game := NewGame(`Kg1,f2,g2,h2,Qa3,Na4`, `Kg8,a7,f7,g7,Qa6,Na5`) // h2,g2,h2 vs A7,f7,g7
	score := game.start().Evaluate()

	expect.Eq(t, score, 45)
}

func TestEvaluatePawns530(t *testing.T) {
	game := NewGame(`Kb1,a3,b2,c2,Qh3,Nh4`, `Kb8,a7,b7,c7,Qh6,Nh5`) // A3,b2,c2 vs a7,b7,c7
	score := game.start().Evaluate()

	expect.Eq(t, score, 0)
}

func TestEvaluatePawns540(t *testing.T) {
	game := NewGame(`Kb1,a3,b4,c2,Qh3,Nh4`, `Kb8,a7,b7,c7,Qh6,Nh5`) // A3,B4,c2 vs a7,b7,c7
	score := game.start().Evaluate()

	expect.Eq(t, score, -25)
}

func TestEvaluatePawns550(t *testing.T) {
	game := NewGame(`Kb1,b2,c2,h2,Qh3,Nh4`, `Kb8,a7,b7,c7,Qh6,Nh5`) // b2,c2,H2 vs a7,b7,c7
	score := game.start().Evaluate()

	expect.Eq(t, score, -25)
}

func TestEvaluatePawns560(t *testing.T) {
	game := NewGame(`Ka1,a3,b2,Qc1,Nd2`, `Kh8,g7,h6,Qf8,Ne7`) // a3,b2 == g7,h6
	score := game.start().Evaluate()

	expect.Eq(t, score, 10)
}

func TestEvaluatePawns570(t *testing.T) {
	game := NewGame(`Kb1,a2,c2,f2,g2,h2`, `Kg8,a7,c7,f7,g7,h7`) // B2 hole but not enough power to bother.
	score := game.start().Evaluate()

	expect.Eq(t, score, 10)
}

// Unstoppable passers.
func TestEvaluatePawns600(t *testing.T) {
	game := NewGame(`Kd4,f2,g2,h2`, `Kg8,g7,h7,a3`) // Kd4-c3 stops A3 pawn.
	score := game.start().Evaluate()

	expect.Eq(t, score, -88)
}

func TestEvaluatePawns610(t *testing.T) {
	game := NewGame(`Kd4,f2,g2,h2`, `M99,Kg8,g7,h7,a3`) // a3-a2 makes the pawn unstoppable.
	score := game.start().Evaluate()

	expect.Eq(t, score, 1178)
}

func TestEvaluatePawns620(t *testing.T) {
	game := NewGame(`Ka1,b4,g2`, `Kg8,g7,h7`) // b4-b5 is unstoppable.
	score := game.start().Evaluate()

	expect.Eq(t, score, 1053)
}

func TestEvaluatePawns630(t *testing.T) {
	game := NewGame(`Ka1,b4,h2`, `M99,Kg8,g7,h7`) // Kg8-f8 stops B4 pawn.
	score := game.start().Evaluate()

	expect.Eq(t, score, 31)
}
