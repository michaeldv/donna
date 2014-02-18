// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`)

// Check evasions (king retreats).
func TestGenerate270(t *testing.T) {
        game := NewGame().Setup(`Kh6,g7`, `Kf8`)
        black := game.Start(Black).StartMoveGen(1).GenerateEvasions()
        expect(t, black.allMoves(), `[Kf8-e7 Kf8-f7 Kf8-e8 Kf8-g8]`)
}

// Check evasions (king retreats).
func TestGenerate280(t *testing.T) {
        game := NewGame().Setup(`Ka1`, `Kf8,Nc4,b2`)
        white := game.Start(White).StartMoveGen(1).GenerateEvasions()
        expect(t, white.allMoves(), `[Ka1-b1 Ka1-a2]`)
}

// Check evasions (king retreats).
func TestGenerate290(t *testing.T) {
        game := NewGame().Setup(`Ka1,h6,g7`, `Kf8`)
        black := game.Start(Black).StartMoveGen(1).GenerateEvasions()
        expect(t, black.allMoves(), `[Kf8-e7 Kf8-f7 Kf8-e8 Kf8-g8]`)
}

// Check evasions (captures/blocks by major pieces).
func TestGenerate300(t *testing.T) {
        game := NewGame().Setup(`Ka5,Ra8`, `Kg8,Qf6,Re6,Bc6,Bg7,Na6,Nb6,f7,g6,h5`)
        black := game.Start(Black).StartMoveGen(1).GenerateEvasions()
        expect(t, black.allMoves(), `[Kg8-h7 Na6-b8 Nb6xa8 Nb6-c8 Bc6xa8 Bc6-e8 Re6-e8 Qf6-d8 Bg7-f8]`)
}

// Check evasions (double check).
func TestGenerate310(t *testing.T) {
        game := NewGame().Setup(`Ka1,Ra8,Nf6`, `Kg8,Qc6,Bb7,Nb6,f7,g6,h7`)
        black := game.Start(Black).StartMoveGen(1).GenerateEvasions()
        expect(t, black.allMoves(), `[Kg8-g7]`)
}

// Check evasions (double check).
func TestGenerate320(t *testing.T) {
        game := NewGame().Setup(`Ka1,Ra8,Be5`, `Kh8,Qd5,Bb7,Nb6`)
        black := game.Start(Black).StartMoveGen(1).GenerateEvasions()
        expect(t, black.allMoves(), `[Kh8-h7]`)
}

// Check evasions (pawn captures).
func TestGenerate330(t *testing.T) {
        game := NewGame().Setup(`Kh6,Be5`, `Kh8,d6`)
        black := game.Start(Black).StartMoveGen(1).GenerateEvasions()
        expect(t, black.allMoves(), `[Kh8-g8 d6xe5]`)
}

// Check evasions (pawn captures).
func TestGenerate340(t *testing.T) {
        game := NewGame().Setup(`Ke1,e2,f2,g2`, `Kc3,Nf3`)
        white := game.Start(White).StartMoveGen(1).GenerateEvasions()
        expect(t, white.allMoves(), `[Ke1-d1 Ke1-f1 e2xf3 g2xf3]`)
}

// Check evasions (pawn blocks).
func TestGenerate350(t *testing.T) {
        game := NewGame().Setup(`Kf8,Ba1`, `Kh8,b3,c4,d5,e6,f7`)
        black := game.Start(Black).StartMoveGen(1).GenerateEvasions()
        expect(t, black.allMoves(), `[Kh8-h7 b3-b2 c4-c3 d5-d4 e6-e5 f7-f6]`)
}

// Check evasions (pawn blocks).
func TestGenerate360(t *testing.T) {
        game := NewGame().Setup(`Ka5,a4,b4,c4,d4,e4,f4,h4`, `Kh8,Qh5`)
        white := game.Start(White).StartMoveGen(1).GenerateEvasions()
        expect(t, white.allMoves(), `[Ka5-a6 Ka5-b6 b4-b5 c4-c5 d4-d5 e4-e5 f4-f5]`)
}

// Check evasions (en-passant pawn capture).
func TestGenerate370(t *testing.T) {
        game := NewGame().Setup(`Kd4,d5,f5`, `Kd8,e7`)
        black := game.Start(Black)
        white := black.MakeMove(black.NewEnpassant(E7, E5)).StartMoveGen(1).GenerateEvasions()
        expect(t, white.allMoves(), `[Kd4-c3 Kd4-d3 Kd4-e3 Kd4-c4 Kd4-e4 Kd4-c5 Kd4xe5 d5xe6 f5xe6]`)
}

// Check evasions (en-passant pawn capture).
func TestGenerate380(t *testing.T) {
        game := NewGame().Setup(`Kb1,b2`, `Ka5,a4,c5,c4`)
        white := game.Start(White)
        black := white.MakeMove(white.NewEnpassant(B2, B4)).StartMoveGen(1).GenerateEvasions()
        expect(t, black.allMoves(), `[Ka5xb4 Ka5-b5 Ka5-a6 Ka5-b6 c5xb4 a4xb3 c4xb3]`)
}

// Check evasions (pawn jumps).
func TestGenerate390(t *testing.T) {
        game := NewGame().Setup(`Kh4,a2,b2,c2,d2,e2,f2,g2`, `Kd8,Ra4`)
        white := game.Start(White).StartMoveGen(1).GenerateEvasions()
        expect(t, white.allMoves(), `[Kh4-g3 Kh4-h3 Kh4-g5 Kh4-h5 b2-b4 c2-c4 d2-d4 e2-e4 f2-f4 g2-g4]`)
}

// Check evasions (pawn jumps).
func TestGenerate400(t *testing.T) {
        game := NewGame().Setup(`Kd8,Rh5`, `Ka5,b7,c7,d7,e7,f7,g7,h7`)
        black := game.Start(Black).StartMoveGen(1).GenerateEvasions()
        expect(t, black.allMoves(), `[Ka5-a4 Ka5-b4 Ka5-a6 Ka5-b6 b7-b5 c7-c5 d7-d5 e7-e5 f7-f5 g7-g5]`)
}

// Check evasions (piece x piece captures).
func TestGenerate410(t *testing.T) {
        game := NewGame().Setup(`Ke1,Qd1,Rc7,Bd3,Nd4,c2,f2`, `Kh8,Nb4`)
        black := game.Start(Black)
        white := black.MakeMove(black.NewMove(B4, C2)).StartMoveGen(1).GenerateEvasions()
        expect(t, white.allMoves(), `[Ke1-f1 Ke1-d2 Ke1-e2 Qd1xc2 Bd3xc2 Nd4xc2 Rc7xc2]`)
}

// Check evasions (pawn x piece captures).
func TestGenerate420(t *testing.T) {
        game := NewGame().Setup(`Ke1,Qd1,Rd7,Bf1,Nc1,c2,d3,f2`, `Kh8,Nb4`)
        black := game.Start(Black)
        white := black.MakeMove(black.NewMove(B4, D3)).StartMoveGen(1).GenerateEvasions()
        expect(t, white.allMoves(), `[Ke1-d2 Ke1-e2 c2xd3 Nc1xd3 Qd1xd3 Bf1xd3 Rd7xd3]`)
}

// Check evasions (king x piece captures).
func TestGenerate430(t *testing.T) {
        game := NewGame().Setup(`Ke1,Qf7,f2`, `Kh8,Qh4`)
        black := game.Start(Black)
        white := black.MakeMove(black.NewMove(H4, F2)).StartMoveGen(1).GenerateEvasions()
        expect(t, white.allMoves(), `[Ke1-d1 Ke1xf2 Qf7xf2]`)
}