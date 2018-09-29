// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

import(`github.com/michaeldv/donna/expect`; `testing`)

func TestGenerateMoves000(t *testing.T) {
	gen := NewMoveGen(NewGame().start()).generateMoves()

	// All possible moves in the initial position, pawn-to-queen, left-to right, unsorted.
	expect.Eq(t, gen.allMoves(), `[a2-a3 a2-a4 b2-b3 b2-b4 c2-c3 c2-c4 d2-d3 d2-d4 e2-e3 e2-e4 f2-f3 f2-f4 g2-g3 g2-g4 h2-h3 h2-h4 Nb1-a3 Nb1-c3 Ng1-f3 Ng1-h3]`)
}

func TestGenerateMoves020(t *testing.T) {
	game := NewGame(`a2,b3,c4,d2,e6,f5,g4,h3`, `a3,b4,c5,e7,f6,g5,h4,Kg8`)
	gen := NewMoveGen(game.start()).generateMoves()

	// All possible moves, left-to right, unsorted.
	expect.Eq(t, gen.allMoves(), `[d2-d3 d2-d4]`)
}

func TestGenerateMoves030(t *testing.T) {
	game := NewGame(`a2,e4,g2`, `b3,f5,f3,h3,Kg8`)
	gen := NewMoveGen(game.start()).generateMoves()

	// All possible moves, left-to right, unsorted.
	expect.Eq(t, gen.allMoves(), `[a2-a3 a2xb3 a2-a4 g2xf3 g2-g3 g2xh3 g2-g4 e4-e5 e4xf5]`)
}

// Castles.
func TestGenerateMoves031(t *testing.T) {
	p := NewGame(`Ke1,Rh1,h2`, `Ke8,Ra8,a7`).start()
	white := NewMoveGen(p).generateMoves()
	expect.Contain(t, white.allMoves(), `0-0`)

	p.color = Black
	black := NewMoveGen(p).generateMoves()
	expect.Contain(t, black.allMoves(), `0-0-0`)
}

// Should not include castles when rook has moved.
func TestGenerateMoves040(t *testing.T) {
	p := NewGame(`Ke1,Rf1,g2`, `Ke8`).start()
	white := NewMoveGen(p).generateMoves()
	expect.NotContain(t, white.allMoves(), `0-0`)
}

func TestGenerateMoves050(t *testing.T) {
	p := NewGame(`Ke1,Rb1,b2`, `Ke8`).start()
	white := NewMoveGen(p).generateMoves()
	expect.NotContain(t, white.allMoves(), `0-0`)
}

// Should not include castles when king has moved.
func TestGenerateMoves060(t *testing.T) {
	p := NewGame(`Kf1,Ra1,a2,Rh1,h2`, `Ke8`).start()
	white := NewMoveGen(p).generateMoves()
	expect.NotContain(t, white.allMoves(), `0-0`)
}

// Should not include castles when rooks are not there.
func TestGenerateMoves070(t *testing.T) {
	p := NewGame(`Ke1`, `Ke8`).start()
	white := NewMoveGen(p).generateMoves()
	expect.NotContain(t, white.allMoves(), `0-0`)
}

// Should not include castles when king is in check.
func TestGenerateMoves080(t *testing.T) {
	p := NewGame(`Ke1,Ra1,Rf1`, `Ke8,Re7`).start()
	white := NewMoveGen(p).generateMoves()
	expect.NotContain(t, white.allMoves(), `0-0`)
}

// Should not include castles when target square is a capture.
func TestGenerateMoves090(t *testing.T) {
	p := NewGame(`Ke1,Ra1,Rf1`, `Ke8,Nc1,Ng1`).start()
	white := NewMoveGen(p).generateMoves()
	expect.NotContain(t, white.allMoves(), `0-0`)
}

// Should not include castles when king is to jump over attacked square.
func TestGenerateMoves100(t *testing.T) {
	p := NewGame(`Ke1,Ra1,Rf1`, `Ke8,Bc4,Bf4`).start()
	white := NewMoveGen(p).generateMoves()
	expect.NotContain(t, white.allMoves(), `0-0`)
}

// Pawn moves that include promotions.
func TestGenerateMoves200(t *testing.T) {
	p := NewGame(`Ka1,a6,b7`, `Kh8,g3,h2`).start()
	white := NewMoveGen(p).pawnMoves(White)
	expect.Eq(t, white.allMoves(), `[a6-a7 b7-b8Q b7-b8N b7-b8R b7-b8B]`)

	p.color = Black
	black := NewMoveGen(p).pawnMoves(Black)
	expect.Eq(t, black.allMoves(), `[h2-h1Q h2-h1N h2-h1R h2-h1B g3-g2]`)
}

// Pawn moves that include jumps.
func TestGenerateMoves210(t *testing.T) {
	p := NewGame(`Ka1,a6`, `Kh8,a7,g7,h6`).start()
	white := NewMoveGen(p).pawnMoves(White)
	expect.Eq(t, white.allMoves(), `[]`)

	p.color = Black
	black := NewMoveGen(p).pawnMoves(Black)
	expect.Eq(t, black.allMoves(), `[h6-h5 g7-g5 g7-g6]`)
}

// Pawn captures without promotions.
func TestGenerateMoves220(t *testing.T) {
	p := NewGame(`Ka1,a6,f6,g5`, `Kh8,b7,g7,h6`).start()
	white := NewMoveGen(p).pawnCaptures(White)
	expect.Eq(t, white.allMoves(), `[g5xh6 a6xb7 f6xg7]`)

	p.color = Black
	black := NewMoveGen(p).pawnCaptures(Black)
	expect.Eq(t, black.allMoves(), `[h6xg5 b7xa6 g7xf6]`)
}

// Pawn captures with Queen promotion (captures generate queen promotions only).
func TestGenerateMoves230(t *testing.T) {
	p := NewGame(`Ka1,Rh1,Bf1,c7`, `Kh8,Nb8,Qd8,g2`).start()
	white := NewMoveGen(p).pawnCaptures(White)
	expect.Eq(t, white.allMoves(), `[c7xb8Q c7-c8Q c7xd8Q]`)

	p.color = Black
	black := NewMoveGen(p).pawnCaptures(Black)
	expect.Eq(t, black.allMoves(), `[g2xf1Q g2-g1Q g2xh1Q]`)
}

// Rearrange root moves.
func TestGenerateMoves300(t *testing.T) {
	p := NewGame().start()
	gen := NewMoveGen(p).generateMoves().validOnly()

	// Pick from the middle.
	gen.head = 10 // e2-e4
	gen.rearrangeRootMoves().reset()
	expect.Eq(t, gen.allMoves(), `[e2-e4 a2-a3 a2-a4 b2-b3 b2-b4 c2-c3 c2-c4 d2-d3 d2-d4 e2-e3 f2-f3 f2-f4 g2-g3 g2-g4 h2-h3 h2-h4 Nb1-a3 Nb1-c3 Ng1-f3 Ng1-h3]`)

	// Pick first one.
	gen.head = 1
	gen.rearrangeRootMoves().reset()
	expect.Eq(t, gen.allMoves(), `[e2-e4 a2-a3 a2-a4 b2-b3 b2-b4 c2-c3 c2-c4 d2-d3 d2-d4 e2-e3 f2-f3 f2-f4 g2-g3 g2-g4 h2-h3 h2-h4 Nb1-a3 Nb1-c3 Ng1-f3 Ng1-h3]`)

	// Pick last one.
	gen.head = gen.tail
	gen.rearrangeRootMoves().reset()
	expect.Eq(t, gen.allMoves(), `[Ng1-h3 e2-e4 a2-a3 a2-a4 b2-b3 b2-b4 c2-c3 c2-c4 d2-d3 d2-d4 e2-e3 f2-f3 f2-f4 g2-g3 g2-g4 h2-h3 h2-h4 Nb1-a3 Nb1-c3 Ng1-f3]`)
}
