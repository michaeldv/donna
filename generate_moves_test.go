// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`)

// Old tests for naive move genarator modernized to test the incremental one.
///////////////////////////////////////////////////////////////////////////////
func TestGenerate000(t *testing.T) {
        game := NewGame().InitialPosition()
        gen := game.Start(White).StartMoveGen(1).GenerateMoves()

        // All possible moves in the initial position, pawn-to-queen, left-to right, unsorted.
        expect(t, gen.allMoves(), `[a2-a3 a2-a4 b2-b3 b2-b4 c2-c3 c2-c4 d2-d3 d2-d4 e2-e3 e2-e4 f2-f3 f2-f4 g2-g3 g2-g4 h2-h3 h2-h4 Nb1-a3 Nb1-c3 Ng1-f3 Ng1-h3]`)
}

func TestGenerate020(t *testing.T) {
        game := NewGame().Setup(`a2,b3,c4,d2,e6,f5,g4,h3`, `a3,b4,c5,e7,f6,g5,h4,Kg8`)
        gen := game.Start(White).StartMoveGen(1).GenerateMoves()

        // All possible moves, left-to right, unsorted.
        expect(t, gen.allMoves(), `[d2-d3 d2-d4]`)
}

func TestGenerate030(t *testing.T) {
        game := NewGame().Setup(`a2,e4,g2`, `b3,f5,f3,h3,Kg8`)
        gen := game.Start(White).StartMoveGen(1).GenerateMoves()

        // All possible moves, left-to right, unsorted.
        expect(t, gen.allMoves(), `[a2-a3 a2xb3 a2-a4 g2xf3 g2-g3 g2xh3 g2-g4 e4-e5 e4xf5]`)
}

// Should not include castles when rook has moved.
func TestGenerate040(t *testing.T) {
        game := NewGame().Setup(`Ke1,Rf1,g2`, `Ke8`)
        gen := game.Start(White).StartMoveGen(1).GenerateMoves()

        doesNotContain(t, gen.allMoves(), `0-0`)
}

func TestGenerate050(t *testing.T) {
        game := NewGame().Setup(`Ke1,Rb1,b2`, `Ke8`)
        gen := game.Start(White).StartMoveGen(1).GenerateMoves()

        doesNotContain(t, gen.allMoves(), `0-0`)
}

// Should not include castles when king has moved.
func TestGenerate060(t *testing.T) {
        game := NewGame().Setup(`Kf1,Ra1,a2,Rh1,h2`, `Ke8`)
        gen := game.Start(White).StartMoveGen(1).GenerateMoves()

        doesNotContain(t, gen.allMoves(), `0-0`)
}

// Should not include castles when rooks are not there.
func TestGenerate070(t *testing.T) {
        game := NewGame().Setup(`Ke1`, `Ke8`)
        gen := game.Start(White).StartMoveGen(1).GenerateMoves()

        doesNotContain(t, gen.allMoves(), `0-0`)
}

// Should not include castles when king is in check.
func TestGenerate080(t *testing.T) {
        game := NewGame().Setup(`Ke1,Ra1,Rf1`, `Ke8,Re7`)
        gen := game.Start(White).StartMoveGen(1).GenerateMoves()

        doesNotContain(t, gen.allMoves(), `0-0`)
}

// Should not include castles when target square is a capture.
func TestGenerate090(t *testing.T) {
        game := NewGame().Setup(`Ke1,Ra1,Rf1`, `Ke8,Nc1,Ng1`)
        gen := game.Start(White).StartMoveGen(1).GenerateMoves()

        doesNotContain(t, gen.allMoves(), `0-0`)
}

// Should not include castles when king is to jump over attacked square.
func TestGenerate100(t *testing.T) {
        game := NewGame().Setup(`Ke1,Ra1,Rf1`, `Ke8,Bc4,Bf4`)
        gen := game.Start(White).StartMoveGen(1).GenerateMoves()

        doesNotContain(t, gen.allMoves(), `0-0`)
}

// Pawn moves that include promotions.
func TestGenerate200(t *testing.T) {
        game := NewGame().Setup(`Ka1,a6,b7`, `Kh8,g3,h2`)
        white := game.Start(White).StartMoveGen(1).pawnMoves(White)
        expect(t, white.allMoves(), `[a6-a7 b7-b8Q b7-b8R b7-b8B b7-b8N]`)

        black := game.Start(Black).StartMoveGen(1).pawnMoves(Black)
        expect(t, black.allMoves(), `[h2-h1Q h2-h1R h2-h1B h2-h1N g3-g2]`)
}

// Pawn moves that include jumps.
func TestGenerate210(t *testing.T) {
        game := NewGame().Setup(`Ka1,a6`, `Kh8,a7,g7,h6`)
        white := game.Start(White).StartMoveGen(1).pawnMoves(White)
        expect(t, white.allMoves(), `[]`)

        black := game.Start(Black).StartMoveGen(1).pawnMoves(Black)
        expect(t, black.allMoves(), `[h6-h5 g7-g5 g7-g6]`)
}

// Pawn captures without promotions.
func TestGenerate220(t *testing.T) {
        game := NewGame().Setup(`Ka1,a6,f6,g5`, `Kh8,b7,g7,h6`)
        white := game.Start(White).StartMoveGen(1).pawnCaptures(White)
        expect(t, white.allMoves(), `[g5xh6 a6xb7 f6xg7]`)

        black := game.Start(Black).StartMoveGen(1).pawnCaptures(Black)
        expect(t, black.allMoves(), `[h6xg5 b7xa6 g7xf6]`)
}

// Pawn captures with Queen promotion.
func TestGenerate230(t *testing.T) {
        game := NewGame().Setup(`Ka1,Rh1,Bf1,c7`, `Kh8,Nb8,Qd8,g2`)
        white := game.Start(White).StartMoveGen(1).pawnCaptures(White)
        expect(t, white.allMoves(), `[c7xb8Q c7xb8R c7xb8B c7xb8N c7-c8Q c7-c8R c7-c8B c7-c8N c7xd8Q c7xd8R c7xd8B c7xd8N]`)

        black := game.Start(Black).StartMoveGen(1).pawnCaptures(Black)
        expect(t, black.allMoves(), `[g2xf1Q g2xf1R g2xf1B g2xf1N g2-g1Q g2-g1R g2-g1B g2-g1N g2xh1Q g2xh1R g2xh1B g2xh1N]`)
}
