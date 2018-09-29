// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

import(`github.com/michaeldv/donna/expect`; `testing`)

// Piece captures.
func TestGenCaptures000(t *testing.T) {
	p := NewGame(`Ka1,Qd1,Rh1,Bb3,Ne5`, `Ka8,Qd8,Rh8,Be6,Ng6`).start()
	white := NewMoveGen(p).generateCaptures()
	expect.Eq(t, white.allMoves(), `[Qd1xd8 Rh1xh8 Bb3xe6 Ne5xg6]`)

	p.color = Black
	black := NewMoveGen(p).generateCaptures()
	expect.Eq(t, black.allMoves(), `[Be6xb3 Ng6xe5 Qd8xd1 Rh8xh1]`)
}

// Pawn captures on rows 2-6 (no promotions).
func TestGenCaptures010(t *testing.T) {
	p := NewGame(`Ka1,b3,c4,d5`, `Ka8,a4,b5,c6,e6`).start()
	white := NewMoveGen(p).generateCaptures()
	expect.Eq(t, white.allMoves(), `[b3xa4 c4xb5 d5xc6 d5xe6]`)

	p.color = Black
	black := NewMoveGen(p).generateCaptures()
	expect.Eq(t, black.allMoves(), `[a4xb3 b5xc4 c6xd5 e6xd5]`)
}

// Pawn captures with promotion, rows 1-7 (queen promotions only when generating captures).
func TestGenCaptures020(t *testing.T) {
	p := NewGame(`Ka1,Bh1,Ng1,a2,b7,e7`, `Kb8,Rd8,Be8,Rf8,h2`).start()
	white := NewMoveGen(p).generateCaptures()
	expect.Eq(t, white.allMoves(), `[e7xd8Q e7xf8Q]`) 

	p.color = Black
	black := NewMoveGen(p).generateCaptures()
	expect.Eq(t, black.allMoves(), `[h2xg1Q Kb8xb7]`)
}

// Pawn promotions without capture, rows 1-7 (queen promotions only when generating captures).
func TestGenCaptures030(t *testing.T) {
	p := NewGame(`Ka1,a2,e7`, `Ka8,Rd8,a7,h2`).start()
	white := NewMoveGen(p).generateCaptures()
	expect.Eq(t, white.allMoves(), `[e7xd8Q e7-e8Q]`)

	p.color = Black
	black := NewMoveGen(p).generateCaptures()
	expect.Eq(t, black.allMoves(), `[h2-h1Q]`)
}

// Captures/promotions sort order.
func TestGenCaptures100(t *testing.T) {
	p := NewGame(`Kg1,Qh4,Rg2,Rf1,Nd6,f7`, `Kh8,Qc3,Rd8,Rd7,Ne6,a3,h7`).start()
	gen := NewMoveGen(p).generateCaptures()
	expect.Eq(t, gen.allMoves(), `[f7-f8Q Qh4xh7 Qh4xd8]`)
}

func TestGenCaptures110(t *testing.T) {
	p := NewGame(`Kg1,Qh4,Rg2,Rf1,Nd6,f7`, `Kh8,Qc3,Rd8,Rd7,Ne6,a3,h7`).start()
	gen := NewMoveGen(p).generateCaptures().rank(Move(0))
	expect.Eq(t, gen.allMoves(), `[f7-f8Q Qh4xd8 Qh4xh7]`)
}

func TestGenCaptures120(t *testing.T) {
	p := NewGame(`Kg1,Qh4,Rg2,Rf1,Nd6,f7`, `Kh8,Qc3,Rd8,Rd7,Ne6,a3,h7`).start()
	gen := NewMoveGen(p).generateCaptures().quickRank()
	expect.Eq(t, gen.allMoves(), `[f7-f8Q Qh4xd8 Qh4xh7]`)
}

// Move legality.
func TestGenCaptures130(t *testing.T) {
	p := NewGame(`Ka1,Rd8,Nb1,f6,h6`, `Kh8,Re1,Ng8,a3,c3`).start()
	white := NewMoveGen(p).generateCaptures().quickRank()
	expect.Eq(t, white.allMoves(), `[Rd8xg8 Nb1xc3 Nb1xa3]`) // Nb1 is pinned but nevertheless.

	p.color = Black
	black := NewMoveGen(p).generateCaptures()
	expect.Eq(t, black.allMoves(), `[Re1xb1 Ng8xf6 Ng8xh6]`) // Ng8 is pinned but nevertheless.
}
