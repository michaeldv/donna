// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`github.com/michaeldv/donna/expect`; `testing`)

// Default move ordering.
func TestGenerate010(t *testing.T) {
	game := NewGame(`Ka1,a2,b3,c4,d2,e6,f5,g4,h3`, `Kc1`)
	gen := NewMoveGen(game.start()).generateMoves().rank(Move(0))

	expect.Eq(t, gen.allMoves(), `[a2-a3 a2-a4 d2-d3 d2-d4 b3-b4 h3-h4 c4-c5 g4-g5 f5-f6 e6-e7 Ka1-b1 Ka1-b2]`)
}

// LVA/MVV capture ordering.
func TestGenerate110(t *testing.T) {
	game := NewGame(`Kd4,e4,Nf4,Bc4,Ra5,Qh5`, `Kd8,Qd5`)
	gen := NewMoveGen(game.start()).generateCaptures().rank(Move(0))

	expect.Eq(t, gen.allMoves(), `[e4xd5 Nf4xd5 Bc4xd5 Ra5xd5 Qh5xd5 Kd4xd5]`)
}

func TestGenerate120(t *testing.T) {
	game := NewGame(`Kd4,e4,Nf4,Bc4,Ra5,Qh5`, `Kd8,Qd5,Rf5`)
	gen := NewMoveGen(game.start()).generateCaptures().rank(Move(0))

	expect.Eq(t, gen.allMoves(), `[e4xd5 Nf4xd5 Bc4xd5 Ra5xd5 Kd4xd5 e4xf5 Qh5xf5]`)
}

func TestGenerate130(t *testing.T) {
	game := NewGame(`Kd4,e4,Nf4,Bc4,Ra5,Qh5`, `Kd8,Qd5,Rf5,Bg6`)
	gen := NewMoveGen(game.start()).generateCaptures().rank(Move(0))

	expect.Eq(t, gen.allMoves(), `[e4xd5 Nf4xd5 Bc4xd5 Ra5xd5 Kd4xd5 e4xf5 Qh5xf5 Nf4xg6 Qh5xg6]`)
}

func TestGenerate140(t *testing.T) {
	game := NewGame(`Kd4,e4,Nf4,Bc4,Ra5,Qh5`, `Kd8,Qd5,Rf5,Bg6,Nh3`)
	gen := NewMoveGen(game.start()).generateCaptures().rank(Move(0))

	expect.Eq(t, gen.allMoves(), `[e4xd5 Nf4xd5 Bc4xd5 Ra5xd5 Kd4xd5 e4xf5 Qh5xf5 Nf4xg6 Qh5xg6 Nf4xh3 Qh5xh3]`)
}

func TestGenerate150(t *testing.T) {
	game := NewGame(`Kd4,e4,Nf4,Bc4,Ra5,Qh5`, `Kd8,Qd5,Rf5,Bg6,Nh3,e2`)
	gen := NewMoveGen(game.start()).generateCaptures().rank(Move(0))

	expect.Eq(t, gen.allMoves(), `[e4xd5 Nf4xd5 Bc4xd5 Ra5xd5 Kd4xd5 e4xf5 Qh5xf5 Nf4xg6 Qh5xg6 Nf4xh3 Qh5xh3 Nf4xe2 Bc4xe2 Qh5xe2]`)
}

// Vaditaing generated moves.
func TestGenerate200(t *testing.T) {
	p := NewGame(`Ke1,Qe2,d2`, `Ke8,e4`).start()
	p = p.makeMove(NewEnpassant(p, D2, D4))

	// No e4xd3 en-passant capture.
	black := NewMoveGen(p).generateMoves().validOnly()
	expect.Eq(t, black.allMoves(), `[e4-e3 Ke8-d7 Ke8-e7 Ke8-f7 Ke8-d8 Ke8-f8]`)
}

func TestGenerate210(t *testing.T) {
	p := NewGame(`Ke1,Qg2,d2`, `Ka8,e4`).start()
	p = p.makeMove(NewEnpassant(p, D2, D4))

	// Neither e4-e3 nor e4xd3 en-passant capture.
	black := NewMoveGen(p).generateMoves().validOnly()
	expect.Eq(t, black.allMoves(), `[Ka8-a7 Ka8-b7 Ka8-b8]`)
}
