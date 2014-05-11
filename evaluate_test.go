// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
	`testing`
)

// Initial position.
func TestEvaluate000(t *testing.T) {
	p := NewGame().InitialPosition().Start(White)
	score := p.Evaluate()
	expect(t, score, rightToMove.blended(p.phase())) // Right to move only.
}

// After 1. e2-e4
func TestEvaluate010(t *testing.T) {
	p := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
		`Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e7,f7,g7,h7`).Start(Black)
	score := p.Evaluate()
	expect(t, score, -43) // +43 for white.
}

// After 1. e2-e4 e7-e5
func TestEvaluate020(t *testing.T) {
	p := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
		`Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`).Start(White)
	score := p.Evaluate()
	expect(t, score, rightToMove.blended(p.phase())) // Right to move only.
}

// After 1. e2-e4 e7-e5 2. Ng1-f3
func TestEvaluate030(t *testing.T) {
	p := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Nf3,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
		`Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`).Start(Black)
	score := p.Evaluate()
	expect(t, score, -49) // +49 for White.
}

// After 1. e2-e4 e7-e5 2. Ng1-f3 Ng8-f6
func TestEvaluate040(t *testing.T) {
	p := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Nf3,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
		`Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Nf6,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`).Start(White)
	score := p.Evaluate()
	expect(t, score, rightToMove.blended(p.phase())) // Right to move only.
}

// After 1. e2-e4 e7-e5 2. Ng1-f3 Nb8-c6
func TestEvaluate050(t *testing.T) {
	p := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Nf3,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
		`Ra8,Nc6,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`).Start(White)
	score := p.Evaluate()
	expect(t, score, -3) // +3 for White.
}

// After 1. e2-e4 e7-e5 2. Ng1-f3 Nb8-c6 3. Nb1-c3 Ng8-f6
func TestEvaluate060(t *testing.T) {
	p := NewGame().Setup(`Ra1,Nc3,Bc1,Qd1,Ke1,Bf1,Nf3,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
		`Ra8,Nc6,Bc8,Qd8,Ke8,Bf8,Nf6,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`).Start(White)
	score := p.Evaluate()
	expect(t, score, rightToMove.blended(p.phase())) // Right to move only.
}

func TestEvaluate999(t *testing.T) {
	p := NewGame().Setup(`Ke2,Rg7,Be6,c3,g3,h2`,
		             `Kb8,Bb7,h6,e5,e4,b5,a5`).Start(Black)
	LogOn(); defer LogOff()
	Lop(p)
	Lop(p.material())
	score := p.Evaluate()
	expect(t, score, 1)
}
