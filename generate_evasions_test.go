// Copyright (c) 2014-2016 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`github.com/michaeldv/donna/expect`; `testing`)

// Check evasions (king retreats).
func TestGenEvasions270(t *testing.T) {
	game := NewGame(`Kh6,g7`, `M,Kf8`)
	black := NewMoveGen(game.start()).generateEvasions()
	expect.Eq(t, black.allMoves(), `[Kf8-e7 Kf8-f7 Kf8-e8 Kf8-g8]`)
}

// Check evasions (king retreats).
func TestGenEvasions280(t *testing.T) {
	game := NewGame(`Ka1`, `Kf8,Nc4,b2`)
	white := NewMoveGen(game.start()).generateEvasions()
	expect.Eq(t, white.allMoves(), `[Ka1-b1 Ka1-a2]`)
}

// Check evasions (king retreats).
func TestGenEvasions290(t *testing.T) {
	game := NewGame(`Ka1,h6,g7`, `M,Kf8`)
	black := NewMoveGen(game.start()).generateEvasions()
	expect.Eq(t, black.allMoves(), `[Kf8-e7 Kf8-f7 Kf8-e8 Kf8-g8]`)
}

// Check evasions (captures/blocks by major pieces).
func TestGenEvasions300(t *testing.T) {
	game := NewGame(`Ka5,Ra8`, `M,Kg8,Qf6,Re6,Bc6,Bg7,Na6,Nb6,f7,g6,h5`)
	black := NewMoveGen(game.start()).generateEvasions()
	expect.Eq(t, black.allMoves(), `[Kg8-h7 Na6-b8 Nb6xa8 Nb6-c8 Bc6xa8 Bc6-e8 Re6-e8 Qf6-d8 Bg7-f8]`)
}

// Check evasions (double check).
func TestGenEvasions310(t *testing.T) {
	game := NewGame(`Ka1,Ra8,Nf6`, `M,Kg8,Qc6,Bb7,Nb6,f7,g6,h7`)
	black := NewMoveGen(game.start()).generateEvasions()
	expect.Eq(t, black.allMoves(), `[Kg8-g7]`)
}

// Check evasions (double check).
func TestGenEvasions320(t *testing.T) {
	game := NewGame(`Ka1,Ra8,Be5`, `M,Kh8,Qd5,Bb7,Nb6`)
	black := NewMoveGen(game.start()).generateEvasions()
	expect.Eq(t, black.allMoves(), `[Kh8-h7]`)
}

// Check evasions (pawn captures).
func TestGenEvasions330(t *testing.T) {
	game := NewGame(`Kh6,Be5`, `M,Kh8,d6`)
	black := NewMoveGen(game.start()).generateEvasions()
	expect.Eq(t, black.allMoves(), `[Kh8-g8 d6xe5]`)
}

// Check evasions (pawn captures).
func TestGenEvasions340(t *testing.T) {
	game := NewGame(`Ke1,e2,f2,g2`, `Kc3,Nf3`)
	white := NewMoveGen(game.start()).generateEvasions()
	expect.Eq(t, white.allMoves(), `[Ke1-d1 Ke1-f1 e2xf3 g2xf3]`)
}

// Check evasions (pawn blocks).
func TestGenEvasions350(t *testing.T) {
	game := NewGame(`Kf8,Ba1`, `M,Kh8,b3,c4,d5,e6,f7`)
	black := NewMoveGen(game.start()).generateEvasions()
	expect.Eq(t, black.allMoves(), `[Kh8-h7 b3-b2 c4-c3 d5-d4 e6-e5 f7-f6]`)
}

// Check evasions (pawn blocks).
func TestGenEvasions360(t *testing.T) {
	game := NewGame(`Ka5,a4,b4,c4,d4,e4,f4,h4`, `Kh8,Qh5`)
	white := NewMoveGen(game.start()).generateEvasions()
	expect.Eq(t, white.allMoves(), `[Ka5-a6 Ka5-b6 b4-b5 c4-c5 d4-d5 e4-e5 f4-f5]`)
}

// Check evasions (pawn jumps).
func TestGenEvasions365(t *testing.T) {
	game := NewGame(`Ke1,Rh4`, `M,Kh6,h7`)
	white := NewMoveGen(game.start()).generateEvasions()
	expect.Eq(t, white.allMoves(), `[Kh6-g5 Kh6-g6 Kh6-g7]`) // No h7-h5 jump.
}

// Check evasions (en-passant pawn capture).
func TestGenEvasions370(t *testing.T) {
	game := NewGame(`Kd4,d5,f5`, `M,Kd8,e7`)
	black := game.start()
	white := NewMoveGen(black.makeMove(NewEnpassant(black, E7, E5))).generateEvasions()
	expect.Eq(t, white.allMoves(), `[Kd4-c3 Kd4-d3 Kd4-e3 Kd4-c4 Kd4-e4 Kd4-c5 Kd4xe5 d5xe6 f5xe6]`)
}

// Check evasions (en-passant pawn capture).
func TestGenEvasions380(t *testing.T) {
	game := NewGame(`Kb1,b2`, `Ka5,a4,c5,c4`)
	white := game.start()
	black := NewMoveGen(white.makeMove(NewEnpassant(white, B2, B4))).generateEvasions()
	expect.Eq(t, black.allMoves(), `[Ka5xb4 Ka5-b5 Ka5-a6 Ka5-b6 c5xb4 a4xb3 c4xb3]`)
}

// Check evasions (en-passant pawn capture).
func TestGenEvasions385(t *testing.T) {
	game := NewGame(`Ke4,c5,e5`, `M,Ke7,d7`)
	black := game.start()
	white := NewMoveGen(black.makeMove(NewEnpassant(black, D7, D5))).generateEvasions()
	for move := white.nextMove(); move != 0; move = white.nextMove() {
		if move.piece() == Pawn {
			expect.Eq(t, move.to(), D6)
			expect.Eq(t, move.color(), uint8(White))
			expect.Eq(t, Piece(move.capture()), Piece(BlackPawn))
			expect.Eq(t, Piece(move.promo()), Piece(0))
		}
	}
}

// Check evasions (en-passant pawn capture).
func TestGenEvasions386(t *testing.T) {
	game := NewGame(`Kf1,Qa2`, `M,Kf7,b2`)
	black := NewMoveGen(game.start()).generateEvasions()
	expect.NotContain(t, black.allMoves(), `b2-a1`)
}

// Check evasions (pawn jumps).
func TestGenEvasions390(t *testing.T) {
	game := NewGame(`Kh4,a2,b2,c2,d2,e2,f2,g2`, `Kd8,Ra4`)
	white := NewMoveGen(game.start()).generateEvasions()
	expect.Eq(t, white.allMoves(), `[Kh4-g3 Kh4-h3 Kh4-g5 Kh4-h5 b2-b4 c2-c4 d2-d4 e2-e4 f2-f4 g2-g4]`)
}

// Check evasions (pawn jumps).
func TestGenEvasions400(t *testing.T) {
	game := NewGame(`Kd8,Rh5`, `M,Ka5,b7,c7,d7,e7,f7,g7,h7`)
	black := NewMoveGen(game.start()).generateEvasions()
	expect.Eq(t, black.allMoves(), `[Ka5-a4 Ka5-b4 Ka5-a6 Ka5-b6 b7-b5 c7-c5 d7-d5 e7-e5 f7-f5 g7-g5]`)
}

// Check evasions (pawn jump, sets en-passant).
func TestGenEvasions405(t *testing.T) {
	p := NewGame(`Ke1,Qd4,d5`, `M,Kg7,e7,g6,h7`).start()
	black := NewMoveGen(p).generateEvasions()
	e7e5 := black.list[black.head + 4].move
	expect.Eq(t, black.allMoves(), `[Kg7-h6 Kg7-f7 Kg7-f8 Kg7-g8 e7-e5]`)
	p = p.makeMove(e7e5)
	expect.Eq(t, p.enpassant, uint8(E6))
}

// Check evasions (piece x piece captures).
func TestGenEvasions410(t *testing.T) {
	game := NewGame(`Ke1,Qd1,Rc7,Bd3,Nd4,c2,f2`, `M,Kh8,Nb4`)
	black := game.start()
	white := NewMoveGen(black.makeMove(NewMove(black, B4, C2))).generateEvasions()
	expect.Eq(t, white.allMoves(), `[Ke1-f1 Ke1-d2 Ke1-e2 Qd1xc2 Bd3xc2 Nd4xc2 Rc7xc2]`)
}

// Check evasions (pawn x piece captures).
func TestGenEvasions420(t *testing.T) {
	game := NewGame(`Ke1,Qd1,Rd7,Bf1,Nc1,c2,d3,f2`, `M,Kh8,Nb4`)
	black := game.start()
	white := NewMoveGen(black.makeMove(NewMove(black, B4, D3))).generateEvasions()
	expect.Eq(t, white.allMoves(), `[Ke1-d2 Ke1-e2 c2xd3 Nc1xd3 Qd1xd3 Bf1xd3 Rd7xd3]`)
}

// Check evasions (king x piece captures).
func TestGenEvasions430(t *testing.T) {
	game := NewGame(`Ke1,Qf7,f2`, `M,Kh8,Qh4`)
	black := game.start()
	white := NewMoveGen(black.makeMove(NewMove(black, H4, F2))).generateEvasions()
	expect.Eq(t, white.allMoves(), `[Ke1-d1 Ke1xf2 Qf7xf2]`)
}

// Pawn promotion to block.
func TestGenEvasions440(t *testing.T) {
	game := NewGame(`Kf1,Qf3,Nf2`, `Ka1,b2`)
	white := game.start()
	black := NewMoveGen(white.makeMove(NewMove(white, F3, D1))).generateEvasions()
	expect.Eq(t, black.allMoves(), `[Ka1-a2 b2-b1Q]`)
}

// Pawn promotion to block or capture.
func TestGenEvasions450(t *testing.T) {
	game := NewGame(`Kf1,Qf3,Nf2`, `Ka1,b2,c2`)
	white := game.start()
	black := NewMoveGen(white.makeMove(NewMove(white, F3, D1))).generateEvasions()
	expect.Eq(t, black.allMoves(), `[Ka1-a2 c2xd1Q b2-b1Q c2-c1Q]`)
}

// Pawn promotion to capture.
func TestGenEvasions460(t *testing.T) {
	game := NewGame(`Kf1,Qf3,Nf2`, `Kc1,c2,d2`)
	white := game.start()
	black := NewMoveGen(white.makeMove(NewMove(white, F3, D1))).generateEvasions()
	expect.Eq(t, black.allMoves(), `[Kc1-b2 c2xd1Q]`)
}
