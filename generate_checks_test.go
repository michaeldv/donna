// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`)

// Knight.
func TestGenChecks000(t *testing.T) {
        game := NewGame().Setup(`Ka1,Nd7,Nf3,b3`, `Kh7,Nd4,f6`)
        white := game.Start(White).StartMoveGen(1).GenerateChecks()
        expect(t, white.allMoves(), `[Nf3-g5 Nd7-f8]`)

        black := game.Start(Black).StartMoveGen(1).GenerateChecks()
        expect(t, black.allMoves(), `[Nd4-c2]`)
}

// Bishop.
func TestGenChecks010(t *testing.T) {
        game := NewGame().Setup(`Kh2,Ba2`, `Kh7,Ba7`)
        white := game.Start(White).StartMoveGen(1).GenerateChecks()
        expect(t, white.allMoves(), `[Ba2-b1 Ba2-g8]`)

        black := game.Start(Black).StartMoveGen(1).GenerateChecks()
        expect(t, black.allMoves(), `[Ba7-g1 Ba7-b8]`)
}

func TestGenChecks020(t *testing.T) {
        game := NewGame().Setup(`Kf4,Bc1`, `Kc5,Bf8,h6,e3`)
        white := game.Start(White).StartMoveGen(1).GenerateChecks()
        expect(t, white.allMoves(), `[Bc1-a3]`)

        black := game.Start(Black).StartMoveGen(1).GenerateChecks()
        expect(t, black.allMoves(), `[Bf8-d6]`)
}

// Bishop: discovered non-capturing check with blocked diaginal.
func TestGenChecks030(t *testing.T) {
        game := NewGame().Setup(`Ka8,Ba1,Nb2,c3,f3`, `Kh8,Bh1,Ng2`)
        white := game.Start(White).StartMoveGen(1).GenerateChecks()
        expect(t, white.allMoves(), `[]`)

        black := game.Start(Black).StartMoveGen(1).GenerateChecks()
        expect(t, black.allMoves(), `[]`)
}

// Bishop: discovered non-capturing check: Knight.
func TestGenChecks040(t *testing.T) {
        game := NewGame().Setup(`Ka8,Ba1,Nb2,a4,h4`, `Kh8,Bh1,Ng2,c4,f4`)
        white := game.Start(White).StartMoveGen(1).GenerateChecks()
        expect(t, white.allMoves(), `[Nb2-d1 Nb2-d3]`)

        black := game.Start(Black).StartMoveGen(1).GenerateChecks()
        expect(t, black.allMoves(), `[Ng2-e1 Ng2-e3]`)
}

// Bishop: discovered non-capturing check: Rook.
func TestGenChecks050(t *testing.T) {
        game := NewGame().Setup(`Ka8,Qa1,Nb1,Rb2,b4,d2,e2`, `Kh8,Qh1,Rg2,g4`)
        white := game.Start(White).StartMoveGen(1).GenerateChecks()
        expect(t, white.allMoves(), `[Rb2-a2 Rb2-c2 Rb2-b3]`)

        black := game.Start(Black).StartMoveGen(1).GenerateChecks()
        expect(t, black.allMoves(), `[Rg2-g1 Rg2-f2 Rg2-h2 Rg2-g3]`)
}

// Bishop: discovered non-capturing check: King.
func TestGenChecks060(t *testing.T) {
        game := NewGame().Setup(`Ke5,Qc3,c4,d3,e4`, `Kh8,e6`)
        white := game.Start(White).StartMoveGen(1).GenerateChecks()
        expect(t, white.allMoves(), `[Ke5-f4 Ke5-d5 Ke5-f5 Ke5-d6]`)
}

// Bishop: discovered non-capturing check: Pawn move.
func TestGenChecks070(t *testing.T) {
        game := NewGame().Setup(`Ka8,Ba1,c3`, `Kh8,Bg2,e4`)
        white := game.Start(White).StartMoveGen(1).GenerateChecks()
        expect(t, white.allMoves(), `[c3-c4]`)

        black := game.Start(Black).StartMoveGen(1).GenerateChecks()
        expect(t, black.allMoves(), `[e4-e3]`)
}

// Bishop: discovered non-capturing check: Pawn jump.
func TestGenChecks080(t *testing.T) {
        game := NewGame().Setup(`Kh2,Bb1,c2`, `Kh7,Bb8,c7`)
        white := game.Start(White).StartMoveGen(1).GenerateChecks()
        expect(t, white.allMoves(), `[c2-c3 c2-c4]`)

        black := game.Start(Black).StartMoveGen(1).GenerateChecks()
        expect(t, black.allMoves(), `[c7-c5 c7-c6]`)
}

// Bishop: discovered non-capturing check: no pawn promotions.
func TestGenChecks090(t *testing.T) {
        game := NewGame().Setup(`Kh7,Bb8,c7`, `Kh2,Bb1,c2`)
        white := game.Start(White).StartMoveGen(1).GenerateChecks()
        expect(t, white.allMoves(), `[]`)

        black := game.Start(Black).StartMoveGen(1).GenerateChecks()
        expect(t, black.allMoves(), `[]`)
}

// Bishop: discovered non-capturing check: no enpassant captures.
func TestGenChecks100(t *testing.T) {
        p := NewGame().Setup(`Ka1,Bf4,e5`, `Kb8,f7`).Start(Black)
        white := p.MakeMove(p.NewEnpassant(F7, F5)).StartMoveGen(1).GenerateChecks()
        expect(t, white.allMoves(), `[e5-e6]`)

        p = NewGame().Setup(`Ka1,e2`, `Kb8,Be5,d4`).Start(White)
        black := p.MakeMove(p.NewEnpassant(E2, E4)).StartMoveGen(1).GenerateChecks()
        expect(t, black.allMoves(), `[d4-d3]`)
}

// Bishop: extra Rook moves for Queen.
func TestGenChecks110(t *testing.T) {
        game := NewGame().Setup(`Kb1,Qa1,f2,a2`, `Kh1,Qa7,Nb8,a6`)
        white := game.Start(White).StartMoveGen(1).GenerateChecks()
        expect(t, white.allMoves(), `[Qa1-h8 Kb1-b2 Kb1-c2]`)

        black := game.Start(Black).StartMoveGen(1).GenerateChecks()
        expect(t, black.allMoves(), `[Qa7-b6 Qa7-h7 Qa7-b7]`)
}

// Pawns.
func TestGenChecks120(t *testing.T) {
        game := NewGame().Setup(`Kb5,f2,g2,h2`, `Kg4,a7,b7,c7`)
        white := game.Start(White).StartMoveGen(1).GenerateChecks()
        expect(t, white.allMoves(), `[f2-f3 h2-h3]`)

        black := game.Start(Black).StartMoveGen(1).GenerateChecks()
        expect(t, black.allMoves(), `[a7-a6 c7-c6]`)
}

func TestGenChecks130(t *testing.T) {
        game := NewGame().Setup(`Kb4,f2,g2,h2`, `Kg5,a7,b7,c7`)
        white := game.Start(White).StartMoveGen(1).GenerateChecks()
        expect(t, white.allMoves(), `[f2-f4 h2-h4]`)

        black := game.Start(Black).StartMoveGen(1).GenerateChecks()
        expect(t, black.allMoves(), `[a7-a5 c7-c5]`)
}

func TestGenChecks140(t *testing.T) {
        game := NewGame().Setup(`Kb4,c5,f2,g2,h2`, `Kg5,a7,b7,c7,h4`)
        white := game.Start(White).StartMoveGen(1).GenerateChecks()
        expect(t, white.allMoves(), `[f2-f4]`)

        black := game.Start(Black).StartMoveGen(1).GenerateChecks()
        expect(t, black.allMoves(), `[a7-a5]`)
}

// Rook with pawn on the same rank (discovered check).
func TestGenChecks150(t *testing.T) {
        game := NewGame().Setup(`Ka4,Ra5,e5`, `Kh5,Rh4,c4`)
        white := game.Start(White).StartMoveGen(1).GenerateChecks()
        expect(t, white.allMoves(), `[e5-e6]`)

        black := game.Start(Black).StartMoveGen(1).GenerateChecks()
        expect(t, black.allMoves(), `[c4-c3]`)
}

// Rook with pawn on the same file (no check).
func TestGenChecks160(t *testing.T) {
        game := NewGame().Setup(`Kh8,Ra8,a6`, `Ka3,Rh1,h5`)
        white := game.Start(White).StartMoveGen(1).GenerateChecks()
        expect(t, white.allMoves(), `[]`)

        black := game.Start(Black).StartMoveGen(1).GenerateChecks()
        expect(t, black.allMoves(), `[]`)
}

// Rook with king on the same rank (discovered check).
func TestGenChecks170(t *testing.T) {
        game := NewGame().Setup(`Ke5,Ra5,d4,e4,f4`, `Kh5`)
        white := game.Start(White).StartMoveGen(1).GenerateChecks()
        expect(t, white.allMoves(), `[Ke5-d6 Ke5-e6 Ke5-f6]`)
}

// Rook with king on the same file (discovered check).
func TestGenChecks180(t *testing.T) {
        game := NewGame().Setup(`Kb5,Rb8,c4,c5,c6`, `Kb1`)
        white := game.Start(White).StartMoveGen(1).GenerateChecks()
        expect(t, white.allMoves(), `[Kb5-a4 Kb5-a5 Kb5-a6]`)
}
