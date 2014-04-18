// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
	`testing`
)

// Mate in 2.

// The very first chess puzzle I had solved as a kid.
func TestSearch000(t *testing.T) {
	move := NewGame().Setup(`Kf8,Rh1,g6`, `Kh8,Bg8,g7,h7`).Start(White).search(3)
	expect(t, move, `Rh1-h6`)
}

func TestSearch020(t *testing.T) {
	move := NewGame().Setup(`Kf4,Qc2,Nc5`, `Kd4`).Start(White).search(3)
	expect(t, move, `Nc5-b7`)
}

func TestSearch030(t *testing.T) {
	move := NewGame().Setup(`Kf2,Qf7,Nf3`, `Kg4`).Start(White).search(3)
	expect(t, move, `Qf7-f6`)
}

func TestSearch040(t *testing.T) {
	move := NewGame().Setup(`Kc3,Qc2,Ra4`, `Kb5`).Start(White).search(3)
	expect(t, move, `Qc2-g6`)
}

func TestSearch050(t *testing.T) {
	move := NewGame().Setup(`Ke5,Qc1,Rf3,Bg2`, `Ke2,Nd5,Nb1`).Start(White).search(3)
	expect(t, move, `Rf3-d3`)
}

func TestSearch060(t *testing.T) {
	move := NewGame().Setup(`Kf1,Qa8,Bf7,Ng2`, `Kg4`).Start(White).search(3)
	expect(t, move, `Qa8-b8`)
}

func TestSearch070(t *testing.T) {
	move := NewGame().Setup(`Ke5,Rd3,Bb1`, `Kh7`).Start(White).search(3)
	expect(t, move, `Ke5-f6`)
}

// Puzzles with pawns.

func TestSearch080(t *testing.T) {
	move := NewGame().Setup(`Kg3,Bc1,Nc3,Bg2`, `Kg1,Re1,e3`).Start(White).search(3)
	expect(t, move, `Bc1-a3`)
}

func TestSearch090(t *testing.T) {
	move := NewGame().Setup(`Kf2,Qb8,Be7,f3`, `Kh5,h6,g5`).Start(White).search(3)
	expect(t, move, `Qb8-b1`)
}

func TestSearch100(t *testing.T) {
	move := NewGame().Setup(`Ke6,Qg3,b3,c2`, `Ke4,e7,f5`).Start(White).search(3)
	expect(t, move, `b3-b4`)
}

func TestSearch110(t *testing.T) {
	move := NewGame().Setup(`Kf1,Qh6,Nd2,Nf2`, `Kc1,c2,c3`).Start(White).search(3)
	expect(t, move, `Qh6-a6`)
}

func TestSearch120(t *testing.T) { // En-passant
	move := NewGame().Setup(`Kd5,Qc8,c5,e5,g6`, `Ke7,d7`).Start(White).search(3)
	expect(t, move, `Kd5-e4`)
}

func TestSearch130(t *testing.T) { // En-passant
	move := NewGame().Setup(`Ke7,Rf8,Ba3,Bc2,e5,g5`, `Kg7,c3,h7`).Start(White).search(3)
	expect(t, move, `Ba3-c1`)
}

func TestSearch140(t *testing.T) { // En-passant, stalemate
	move := NewGame().Setup(`Kc6,Rh4,Bb5,a3,c2,d3`, `Ka5,c5,d4,h5`).Start(White).search(3)
	expect(t, move, `c2-c4`)
}

func TestSearch150(t *testing.T) { // Stalemate after Qg7-c3
	move := NewGame().Setup(`Kb4,Qg7,Nc1`, `Kb1`).Start(White).search(3)
	expect(t, move, `Kb4-c3`)
}

func TestSearch160(t *testing.T) { // Pawn promotion
	move := NewGame().Setup(`Ka8,Qc4,b7`, `Ka5`).Start(White).search(3)
	expect(t, move, `b7-b8B`)
}

func TestSearch170(t *testing.T) { // Pawn promotion
	move := NewGame().Setup(`Kf8,Rc6,Be4,Nd7,c7`, `Ke6,d6`).Start(White).search(3)
	expect(t, move, `c7-c8R`)
}

func TestSearch180(t *testing.T) { // Pawn promotion
	move := NewGame().Setup(`Kc6,c7`, `Ka7`).Start(White).search(3)
	expect(t, move, `c7-c8R`)
}

func TestSearch190(t *testing.T) { // Pawn promotion
	move := NewGame().Setup(`Kc4,a7,c7`, `Ka5`).Start(White).search(3)
	expect(t, move, `c7-c8N`)
}

func TestSearch195(t *testing.T) { // King-side castle
	move := NewGame().Setup(`Ke1,Rf1,Rh1`, `Ka1`).Start(White).search(3)
	expect(t, move, `Rf1-f2`)
}

func TestSearch196(t *testing.T) { // Queen-side castle
	move := NewGame().Setup(`Ke1,Ra1,Rb1`, `Kg1`).Start(White).search(3)
	expect(t, move, `Rb1-b2`)
}

// Mate in 3.

func TestSearch200(t *testing.T) {
	move := NewGame().Setup(`Kf8,Re7,Nd5`, `Kh8,Bh5`).Start(White).search(5)
	expect(t, move, `Re7-g7`)
}

func TestSearch210(t *testing.T) {
	move := NewGame().Setup(`Kf8,Bf7,Nf3,e5`, `Kh8,e6,h7`).Start(White).search(5)
	expect(t, move, `Bf7-g8`)
}

func TestSearch220(t *testing.T) { // Pawn promotion
	move := NewGame().Setup(`Kf3,h7`, `Kh1,h3`).Start(White).search(5)
	expect(t, move, `h7-h8R`)
}

func TestSearch230(t *testing.T) { // Pawn promotion
	move := NewGame().Setup(`Kd8,c7,e4,f7`, `Ke6,e5`).Start(White).search(5)
	expect(t, move, `f7-f8R`)
}

func TestSearch240(t *testing.T) { // Pawn promotion
	move := NewGame().Setup(`Kh3,f7,g7`, `Kh6`).Start(White).search(5)
	expect(t, move, `g7-g8Q`)
}

func TestSearch250(t *testing.T) { // Pawn promotion
	move := NewGame().Setup(`Ke4,c7,d6,e7,f6,g7`, `Ke6`).Start(White).search(5)
	expect(t, move, `e7-e8B`)
}

// Mate in 4.

func TestSearch260(t *testing.T) { // Pawn promotion
	move := NewGame().Setup(`Kf6,Nf8,Nh6`, `Kh8,f7,h7`).Start(White).search(7)
	expect(t, move, `Nf8-e6`)
}

func TestSearch270(t *testing.T) { // Pawn promotion/stalemate
	move := NewGame().Setup(`Kf2,e7`, `Kh1,d2`).Start(White).search(7)
	expect(t, move, `e7-e8R`)
}

func TestSearch280(t *testing.T) { // Stalemate
	move := NewGame().Setup(`Kc1,Nb4,a2`, `Ka1,b5`).Start(White).search(7)
	expect(t, move, `a2-a4`)
}

func TestSearch290(t *testing.T) { // Stalemate
	move := NewGame().Setup(`Kh6,Rd3,h7`, `Kh8,Bd7`).Start(White).search(7)
	expect(t, move, `Rd3-d6`)
}

func TestSearch300(t *testing.T) {
	move := NewGame().Setup(`Kc6,Bc1,Ne5`, `Kc8,Ra8,a7,a6`).Start(White).search(7)
	expect(t, move, `Ne5-f7`)
}

// // Bobby Fischer vs. James Sherwin benchmark.
// func TestSearch900(t *testing.T) {
//         move := NewGame().Setup(`Kg1,Qc2,Ra1,Re1,Bc1,Bg2,Ng5,a2,b2,c3,d4,f2,g3,h2`,
//                                 `Kg8,Qd6,Ra8,Rf8,Bc8,Nd5,Ng6,a7,b6,c4,e6,f7,g7,h7`).Start(White).search(8)
//         expect(t, move, `h2-h4`)
// }
//
// // Mikhail Botvinnik vs. Jose Raul Capablanca
// func TestSearch910(t *testing.T) {
//         move := NewGame().Setup(`Kg1,Qe5,Bb2,Ng3,c3,d4,e6,g2,h2`,
//                                 `Kg7,Qe7,Nb3,Nf6,a7,b6,c4,d5,g6,h7`).Start(White).search(10)
//         expect(t, move, `Bb2-a3`)
// }
//
