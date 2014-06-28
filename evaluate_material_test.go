// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import `testing`

// Bare kings.
func TestMaterial000(t *testing.T) {
	p := NewGame(`Ke1`, `Ke8`).Start(Black)
	p.EvaluateWithTrace()

	expect(t, eval.material.flags, uint8(materialDraw))
	expect(t, eval.material.endgame, nil)
}

// No pawns, king with a minor.
func TestMaterial010(t *testing.T) {
	p := NewGame(`Ke1,Bc1`, `Ke8`).Start(White)
	p.EvaluateWithTrace()

	expect(t, eval.material.flags, uint8(materialDraw))
	expect(t, eval.material.endgame, nil)
}

func TestMaterial015(t *testing.T) {
	p := NewGame(`Ke1,Bc1`, `Ke8,Nb8`).Start(White)
	p.EvaluateWithTrace()

	expect(t, eval.material.flags, uint8(materialDraw))
	expect(t, eval.material.endgame, nil)
}

// No pawns, king with two knights.
func TestMaterial020(t *testing.T) {
	p := NewGame(`Ke1,Ne2,Ne3`, `Ke8`).Start(White)
	p.EvaluateWithTrace()

	expect(t, eval.material.flags, uint8(materialDraw))
	expect(t, eval.material.endgame, nil)
}

// Known: king and a pawn vs. bare king.
func TestMaterial100(t *testing.T) {
	p := NewGame(`Ke1,e2`, `Ke8`).Start(White)
	p.EvaluateWithTrace()

	expect(t, eval.material.hash, uint64(0x5355F900C2A82DC7))
	expect(t, eval.material.flags, uint8(knownEndgame))
	expect(t, eval.material.endgame, (*Evaluation).kingAndPawnVsBareKing)
}

func TestMaterial110(t *testing.T) {
	p := NewGame(`Ke1`, `Ke8,e7`).Start(White)
	p.EvaluateWithTrace()

	expect(t, eval.material.hash, uint64(0x9D39247E33776D41))
	expect(t, eval.material.flags, uint8(knownEndgame))
	expect(t, eval.material.endgame, (*Evaluation).kingAndPawnVsBareKing)
}

// Known: king with a knight and a bishop vs. bare king.
func TestMaterial120(t *testing.T) {
	p := NewGame(`Ke1,Nb1,Bc1`, `Ke8`).Start(White)
	p.EvaluateWithTrace()

	expect(t, eval.material.hash, uint64(0xE6F0FBA55BF280F1))
	expect(t, eval.material.flags, uint8(knownEndgame))
	expect(t, eval.material.endgame, (*Evaluation).knightAndBishopVsBareKing)
}

func TestMaterial130(t *testing.T) {
	p := NewGame(`Ke1`, `Ke8,Nb8,Bc8`).Start(White)
	p.EvaluateWithTrace()

	expect(t, eval.material.hash, uint64(0x29D8066E0A562122))
	expect(t, eval.material.flags, uint8(knownEndgame))
	expect(t, eval.material.endgame, (*Evaluation).knightAndBishopVsBareKing)
}

// Lesser known endgame: king and two or more pawns vs. bare king.
func TestMaterial140(t *testing.T) {
	p := NewGame(`Ke1,a4,a5`, `Ka8`).Start(Black)
	p.EvaluateWithTrace()

	expect(t, eval.material.flags, uint8(lesserKnownEndgame))
	expect(t, eval.material.endgame, (*Evaluation).kingAndPawnsVsBareKing)
}
