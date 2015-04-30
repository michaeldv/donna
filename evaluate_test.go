// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`github.com/michaeldv/donna/expect`; `testing`)

// Initial position.
func TestEvaluate000(t *testing.T) {
	p := NewGame().start()
	score := p.Evaluate()
	expect.Eq(t, score, rightToMove.midgame) // Right to move only.
}

// After 1. e2-e4
func TestEvaluate010(t *testing.T) {
	p := NewGame(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
		`M1,Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e7,f7,g7,h7`).start()
	score := p.Evaluate()
	expect.Eq(t, score, -79) // +79 for white.
}

// After 1. e2-e4 e7-e5
func TestEvaluate020(t *testing.T) {
	p := NewGame(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
		`Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`).start()
	score := p.Evaluate()
	expect.Eq(t, score, rightToMove.midgame) // Right to move only.
}

// After 1. e2-e4 e7-e5 2. Ng1-f3
func TestEvaluate030(t *testing.T) {
	p := NewGame(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Nf3,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
		`M2,Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`).start()
	score := p.Evaluate()
	expect.Eq(t, score, -78)
}

// After 1. e2-e4 e7-e5 2. Ng1-f3 Ng8-f6
func TestEvaluate040(t *testing.T) {
	p := NewGame(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Nf3,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
		`Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Nf6,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`).start()
	score := p.Evaluate()
	expect.Eq(t, score, 20)
}

// After 1. e2-e4 e7-e5 2. Ng1-f3 Nb8-c6
func TestEvaluate050(t *testing.T) {
	p := NewGame(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Nf3,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
		`Ra8,Nc6,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`).start()
	score := p.Evaluate()
	expect.Eq(t, score, 9)
}

// After 1. e2-e4 e7-e5 2. Ng1-f3 Nb8-c6 3. Nb1-c3 Ng8-f6
func TestEvaluate060(t *testing.T) {
	p := NewGame(`Ra1,Nc3,Bc1,Qd1,Ke1,Bf1,Nf3,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
		`Ra8,Nc6,Bc8,Qd8,Ke8,Bf8,Nf6,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`).start()
	score := p.Evaluate()
	expect.Eq(t, score, rightToMove.midgame) // Right to move only.
}

// Opposite-colored bishops.
func TestEvaluate070(t *testing.T) {
	p := NewGame(`Ke1,Bc1`, `Ke8,Bc8`).start()
	eval.init(p)
	expect.True(t, eval.oppositeBishops())
}

func TestEvaluate071(t *testing.T) {
	p := NewGame(`Kc4,Bd4`, `Ke8,Bd5`).start()
	eval.init(p)
	expect.True(t, eval.oppositeBishops())
}

func TestEvaluate072(t *testing.T) {
	p := NewGame(`Kc4,Bd4`, `Ke8,Be5`).start()
	eval.init(p)
	expect.False(t, eval.oppositeBishops())
}

func TestEvaluate073(t *testing.T) {
	p := NewGame(`Ke1,Bc1`, `Ke8,Bf8`).start()
	eval.init(p)
	expect.False(t, eval.oppositeBishops())
}
