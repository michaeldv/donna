// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`)

func TestEvaluate000(t *testing.T) {
        game := NewGame().InitialPosition()
        score := game.Start(White).Evaluate()
        expect(t, score, 4) // Right to move.
}

func TestEvaluate010(t *testing.T) { // After 1. e2-e4
        game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
                                `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e7,f7,g7,h7`)
        score := game.Start(Black).Evaluate()
        expect(t, score, -89) // +89 for white.
}

func TestEvaluate020(t *testing.T) { // After 1. e2-e4 e7-e5
        game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
                                `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`)
        score := game.Start(White).Evaluate()
        expect(t, score, 4) // Right to move.
}

func TestEvaluate030(t *testing.T) { // After 1. e2-e4 e7-e5 2. Ng1-f3
        game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Nf3,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
                                `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`)
        score := game.Start(Black).Evaluate()
        expect(t, score, -47) // +47 for White.
}

func TestEvaluate035(t *testing.T) { // After 1. e2-e4 e7-e5 2. Ng1-f3 Ng8-f6
        game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Nf3,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
                                `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Nf6,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`)
        score := game.Start(White).Evaluate()
        expect(t, score, 4) // Right to move.
}

func TestEvaluate040(t *testing.T) { // After 1. e2-e4 e7-e5 2. Ng1-f3 Nb8-c6
        game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Nf3,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
                                `Ra8,Nc6,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`)
        score := game.Start(White).Evaluate()
        expect(t, score, -18) // +18 for White.
}

func TestEvaluate050(t *testing.T) { // After 1. e2-e4 e7-e5 2. Ng1-f3 Nb8-c6 3. Nb1-c3 Ng8-f6
        game := NewGame().Setup(`Ra1,Nc3,Bc1,Qd1,Ke1,Bf1,Nf3,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
                                `Ra8,Nc6,Bc8,Qd8,Ke8,Bf8,Nf6,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`)
        score := game.Start(White).Evaluate()
        expect(t, score, 4) // Right to move.
}

// Doubled pawns.
func TestEvaluate100(t *testing.T) {
        game := NewGame().Setup(`Ke1,h2,h3`, `Ke8,a7,a6`)
        score := game.Start(White).Evaluate()

        expect(t, score, 2) // Right to move for the endgame.
}

func TestEvaluate110(t *testing.T) {
        game := NewGame().Setup(`Ke1,h2,h3`, `Ke8,a7,h7`)
        score := game.Start(White).Evaluate()

        expect(t, score, -101)
}

func TestEvaluate120(t *testing.T) {
        game := NewGame().Setup(`Ke1,f4,f5`, `Ke8,f7,h7`)
        score := game.Start(White).Evaluate()

        expect(t, score, -53)
}

// Passed pawns.
func TestEvaluate200(t *testing.T) {
        game := NewGame().Setup(`Ke1,h4`, `Ke8,h5`) // Blocked.
        score := game.Start(White).Evaluate()

        expect(t, score, 2)
}

func TestEvaluate210(t *testing.T) {
        game := NewGame().Setup(`Ke1,h4`, `Ke8,g7`) // Can't pass.
        score := game.Start(White).Evaluate()

        expect(t, score, 21)
}

func TestEvaluate220(t *testing.T) {
        game := NewGame().Setup(`Ke1,e4`, `Ke8,d6`) // Can't pass.
        score := game.Start(White).Evaluate()

        expect(t, score, 7)
}

func TestEvaluate230(t *testing.T) {
        game := NewGame().Setup(`Ke1,e5`, `Ke8,e4`) // Both passing.
        score := game.Start(White).Evaluate()

        expect(t, score, 2)
}

func TestEvaluate240(t *testing.T) {
        game := NewGame().Setup(`Ke1,e5`, `Ke8,d5`) // Both passing but white is closer.
        score := game.Start(White).Evaluate()

        expect(t, score, 31)
}

func TestEvaluate250(t *testing.T) {
        game := NewGame().Setup(`Ke1,a5`, `Ke8,h7`) // Both passing but white is much closer.
        score := game.Start(White).Evaluate()

        expect(t, score, 64)
}

// Isolated pawns.
func TestEvaluate300(t *testing.T) {
        game := NewGame().Setup(`Ke1,a5,c5`, `Ke8,f4,h4`) // All pawns are isolated.
        score := game.Start(White).Evaluate()

        expect(t, score, 2)
}

func TestEvaluate310(t *testing.T) {
        game := NewGame().Setup(`Ke1,a2,c2,e2`, `Ke8,a7,b7,c7`) // White pawns are isolated.
        score := game.Start(White).Evaluate()

        expect(t, score, -80)
}

// Rooks.
func TestEvaluate400(t *testing.T) {
        game := NewGame().Setup(`Ke1,Ra7`, `Ke8,Rh3`) // White on 7th.
        score := game.Start(White).Evaluate()

        expect(t, score, 9)
}

func TestEvaluate410(t *testing.T) {
        game := NewGame().Setup(`Ke1,Rb1,Ng2,a2`, `Ke8,Rh8,Nb7,h7`) // White on open file.
        score := game.Start(White).Evaluate()

        expect(t, score, 64)
}

func TestEvaluate420(t *testing.T) {
        game := NewGame().Setup(`Ke1,Rb1,a2,g2`, `Ke8,Rh8,h7,b7`) // White on semi-open file.
        score := game.Start(White).Evaluate()

        expect(t, score, 92)
}

// King shield.
func TestEvaluate500(t *testing.T) {
        game := NewGame().Setup(`Kg1,f2,g2,h2,Qa3,Na4`, `Kg8,f7,g7,h7,Qa6,Na5`) // h2,g2,h2 == f7,g7,h7
        score := game.Start(White).Evaluate()

        expect(t, score, 2)
}
func TestEvaluate505(t *testing.T) {
        game := NewGame().Setup(`Kg1,f2,g2,h2,Qa3,Na4`, `Kg8,f7,g6,h7,Qa6,Na5`) // h2,g2,h2 vs f7,G6,h7
        score := game.Start(White).Evaluate()

        expect(t, score, 22)
}

func TestEvaluate510(t *testing.T) {
        game := NewGame().Setup(`Kg1,f2,g2,h2,Qa3,Na4`, `Kg8,f5,g6,h7,Qa6,Na5`) // h2,g2,h2 vs F5,G6,h7
        score := game.Start(White).Evaluate()

        expect(t, score, 19)
}

func TestEvaluate520(t *testing.T) {
        game := NewGame().Setup(`Kg1,f2,g2,h2,Qa3,Na4`, `Kg8,a7,f7,g7,Qa6,Na5`) // h2,g2,h2 vs A7,f7,g7
        score := game.Start(White).Evaluate()

        expect(t, score, 42)
}

func TestEvaluate530(t *testing.T) {
        game := NewGame().Setup(`Kb1,a3,b2,c2,Qh3,Nh4`, `Kb8,a7,b7,c7,Qh6,Nh5`) // A3,b2,c2 vs a7,b7,c7
        score := game.Start(White).Evaluate()

        expect(t, score, -9)
}

func TestEvaluate540(t *testing.T) {
        game := NewGame().Setup(`Kb1,a3,b4,c2,Qh3,Nh4`, `Kb8,a7,b7,c7,Qh6,Nh5`) // A3,B4,c2 vs a7,b7,c7
        score := game.Start(White).Evaluate()

        expect(t, score, 2)
}

func TestEvaluate550(t *testing.T) {
        game := NewGame().Setup(`Kb1,b2,c2,h2,Qh3,Nh4`, `Kb8,a7,b7,c7,Qh6,Nh5`) // b2,c2,H2 vs a7,b7,c7
        score := game.Start(White).Evaluate()

        expect(t, score, -36)
}

func TestEvaluate560(t *testing.T) {
        game := NewGame().Setup(`Ka1,a3,b2,Qc1,Nd2`, `Kh8,g7,h6,Qf8,Ne7`) // a3,b2 == g7,h6
        score := game.Start(White).Evaluate()

        expect(t, score, 2)
}

func TestEvaluate570(t *testing.T) {
        game := NewGame().Setup(`Kb1,a2,c2,f2,g2,h2`, `Kg8,a7,c7,f7,g7,h7`) // B2 hole but not enough power to bother.
        score := game.Start(White).Evaluate()

        expect(t, score, 11)
}
