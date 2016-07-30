// Copyright (c) 2014-2016 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`github.com/michaeldv/donna/expect`; `testing`)

// Pawns.
func TestGenQuiets000(t *testing.T) {
	p := NewGame(`Ka1,Rb1,a2,b2,c4,d7`, `Kh8,Rc1,Ne8,a3,b5`).start()
	gen := NewMoveGen(p).generateQuiets()
	expect.Eq(t, gen.allMoves(), `[b2-b4 b2-b3 c4-c5]`)
}

func TestGenQuiets001(t *testing.T) {
	p := NewGame(`Kh1,Rc8,Ne1,a6,b4`, `M,Ka8,Rb8,a7,b7,c5,d2`).start()
	gen := NewMoveGen(p).generateQuiets()
	expect.Eq(t, gen.allMoves(), `[c5-c4 b7-b5 b7-b6]`)
}

// Knights.
func TestGenQuiets010(t *testing.T) {
	p := NewGame(`Ka1,Nb1,a2`, `Kh8,Rb2,a3,d2`).start()
	gen := NewMoveGen(p).generateQuiets()
	expect.Eq(t, gen.allMoves(), `[Nb1-c3]`)
}

func TestGenQuiets011(t *testing.T) {
	p := NewGame(`Kh1,Rb7,a6,d7`, `M,Ka8,Nb8,a7`).start()
	gen := NewMoveGen(p).generateQuiets()
	expect.Eq(t, gen.allMoves(), `[Nb8-c6]`)
}

// Bishops.
func TestGenQuiets020(t *testing.T) {
	p := NewGame(`Ka1,Bb1`, `Kh8,Rb2,a2,e4`).start()
	gen := NewMoveGen(p).generateQuiets()
	expect.Eq(t, gen.allMoves(), `[Bb1-c2 Bb1-d3]`)
}

func TestGenQuiets021(t *testing.T) {
	p := NewGame(`Kh1,Rb7,a7,e5`, `M,Ka8,Bb8`).start()
	gen := NewMoveGen(p).generateQuiets()
	expect.Eq(t, gen.allMoves(), `[Bb8-d6 Bb8-c7]`)
}

// Rooks.
func TestGenQuiets030(t *testing.T) {
	p := NewGame(`Ka1,Rb1`, `Kh8,Rb2,Re1,a2`).start()
	gen := NewMoveGen(p).generateQuiets()
	expect.Eq(t, gen.allMoves(), `[Rb1-c1 Rb1-d1]`)
}

func TestGenQuiets031(t *testing.T) {
	p := NewGame(`Kh1,Rb7,Re8,a7`, `M,Ka8,Rb8`).start()
	gen := NewMoveGen(p).generateQuiets()
	expect.Eq(t, gen.allMoves(), `[Rb8-c8 Rb8-d8]`)
}

// Queens.
func TestGenQuiets040(t *testing.T) {
	p := NewGame(`Ka1,Qb1`, `Kh8,Rb2,Re1,a2,e4`).start()
	gen := NewMoveGen(p).generateQuiets()
	expect.Eq(t, gen.allMoves(), `[Qb1-c1 Qb1-d1 Qb1-c2 Qb1-d3]`)
}

func TestGenQuiets041(t *testing.T) {
	p := NewGame(`Kh1,Rb7,Re8,a7,e5`, `M,Ka8,Qb8`).start()
	gen := NewMoveGen(p).generateQuiets()
	expect.Eq(t, gen.allMoves(), `[Qb8-d6 Qb8-c7 Qb8-c8 Qb8-d8]`)
}

// 0-0, King.
func TestGenQuiets050(t *testing.T) {
	p := NewGame(`Ke1,Rh1,h2`, `Ka8,h3`).start()
	gen := NewMoveGen(p).generateQuiets()
	expect.Eq(t, gen.allMoves(), `[0-0 Rh1-f1 Rh1-g1 Ke1-d1 Ke1-f1 Ke1-d2 Ke1-e2 Ke1-f2]`)
}

func TestGenQuiets051(t *testing.T) {
	p := NewGame(`Ke1,Rh1,h2`, `Ka8,e2,h3`).start()
	gen := NewMoveGen(p).generateQuiets()
	expect.Eq(t, gen.allMoves(), `[Rh1-f1 Rh1-g1 Ke1-d1 Ke1-f1 Ke1-d2 Ke1-f2]`)
}

func TestGenQuiets052(t *testing.T) {
	p := NewGame(`Ka1,h6`, `M,Ke8,Rh8,h7`).start()
	gen := NewMoveGen(p).generateQuiets()
	expect.Eq(t, gen.allMoves(), `[0-0 Rh8-f8 Rh8-g8 Ke8-d7 Ke8-e7 Ke8-f7 Ke8-d8 Ke8-f8]`)
}

func TestGenQuiets053(t *testing.T) {
	p := NewGame(`Ka1,e7,h6`, `M,Ke8,Rh8,h7`).start()
	gen := NewMoveGen(p).generateQuiets()
	expect.Eq(t, gen.allMoves(), `[Rh8-f8 Rh8-g8 Ke8-d7 Ke8-f7 Ke8-d8 Ke8-f8]`)
}

// 0-0-0, King.
func TestGenQuiets060(t *testing.T) {
	p := NewGame(`Ke1,Ra1,a2`, `Ka8,a3`).start()
	gen := NewMoveGen(p).generateQuiets()
	expect.Eq(t, gen.allMoves(), `[0-0-0 Ra1-b1 Ra1-c1 Ra1-d1 Ke1-d1 Ke1-f1 Ke1-d2 Ke1-e2 Ke1-f2]`)
}

func TestGenQuiets061(t *testing.T) {
	p := NewGame(`Ke1,Ra1,a2`, `Ka8,a3,e2`).start()
	gen := NewMoveGen(p).generateQuiets()
	expect.Eq(t, gen.allMoves(), `[Ra1-b1 Ra1-c1 Ra1-d1 Ke1-d1 Ke1-f1 Ke1-d2 Ke1-f2]`)
}

func TestGenQuiets62(t *testing.T) {
	p := NewGame(`Ka1,a6`, `M,Ke8,Ra8,a7`).start()
	gen := NewMoveGen(p).generateQuiets()
	expect.Eq(t, gen.allMoves(), `[0-0-0 Ra8-b8 Ra8-c8 Ra8-d8 Ke8-d7 Ke8-e7 Ke8-f7 Ke8-d8 Ke8-f8]`)
}

func TestGenQuiets063(t *testing.T) {
	p := NewGame(`Ka1,a6,e7`, `M,Ke8,Ra8,a7`).start()
	gen := NewMoveGen(p).generateQuiets()
	expect.Eq(t, gen.allMoves(), `[Ra8-b8 Ra8-c8 Ra8-d8 Ke8-d7 Ke8-f7 Ke8-d8 Ke8-f8]`)
}

// Move legality.
func TestGenQuiets070(t *testing.T) {
	p := NewGame(`Ka1,Bb1,Nb2`, `Ka8,Rd1,Bd4,a4,b3,c4,e4`).start() // Bb1,Nb2 both pinned, a2 under attack.
	gen := NewMoveGen(p).generateQuiets()
	expect.Eq(t, gen.allMoves(), `[Nb2-d3 Bb1-a2 Bb1-c2 Bb1-d3 Ka1-a2]`)
}

func TestGenQuiets071(t *testing.T) {
	p := NewGame(`Ka1,Rd8,Bd5,a5,b6,c5,e5`, `M,Ka8,Bb8,Nb7`).start() // Bb8,Nb7 both pinned, a7 under attack.
	gen := NewMoveGen(p).generateQuiets()
	expect.Eq(t, gen.allMoves(), `[Nb7-d6 Bb8-d6 Bb8-a7 Bb8-c7 Ka8-a7]`)
}
