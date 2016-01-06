// Copyright (c) 2014-2016 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`github.com/michaeldv/donna/expect`; `testing`)

// Mate in 2.

// The very first chess puzzle I had solved as a kid.
func TestSearch000(t *testing.T) {
	move := NewGame(`Kf8,Rh1,g6`, `Kh8,Bg8,g7,h7`).start().solve(3)
	expect.Eq(t, move, `Rh1-h6`)
}

func TestSearch020(t *testing.T) {
	move := NewGame(`Kf4,Qc2,Nc5`, `Kd4`).start().solve(3)
	expect.Eq(t, move, `Nc5-b7`)
}

func TestSearch030(t *testing.T) {
	move := NewGame(`Kf2,Qf7,Nf3`, `Kg4`).start().solve(3)
	expect.Eq(t, move, `Qf7-f6`)
}

func TestSearch040(t *testing.T) {
	move := NewGame(`Kc3,Qc2,Ra4`, `Kb5`).start().solve(3)
	expect.Eq(t, move, `Qc2-g6`)
}

func TestSearch050(t *testing.T) {
	move := NewGame(`Ke5,Qc1,Rf3,Bg2`, `Ke2,Nd5,Nb1`).start().solve(3)
	expect.Eq(t, move, `Rf3-d3`)
}

func TestSearch060(t *testing.T) {
	move := NewGame(`Kf1,Qa8,Bf7,Ng2`, `Kg4`).start().solve(3)
	expect.Eq(t, move, `Qa8-b8`)
}

func TestSearch070(t *testing.T) {
	move := NewGame(`Ke5,Rd3,Bb1`, `Kh7`).start().solve(3)
	expect.Eq(t, move, `Ke5-f6`)
}

// Puzzles with pawns.

func TestSearch080(t *testing.T) {
	move := NewGame(`Kg3,Bc1,Nc3,Bg2`, `Kg1,Re1,e3`).start().solve(3)
	expect.Eq(t, move, `Bc1-a3`)
}

func TestSearch090(t *testing.T) {
	move := NewGame(`Kf2,Qb8,Be7,f3`, `Kh5,h6,g5`).start().solve(3)
	expect.Eq(t, move, `Qb8-b1`)
}

func TestSearch100(t *testing.T) {
	move := NewGame(`Ke6,Qg3,b3,c2`, `Ke4,e7,f5`).start().solve(3)
	expect.Eq(t, move, `b3-b4`)
}

func TestSearch110(t *testing.T) {
	move := NewGame(`Kf1,Qh6,Nd2,Nf2`, `Kc1,c2,c3`).start().solve(3)
	expect.Eq(t, move, `Qh6-a6`)
}

func TestSearch120(t *testing.T) { // En-passant
	move := NewGame(`Kd5,Qc8,c5,e5,g6`, `Ke7,d7`).start().solve(3)
	expect.Eq(t, move, `Kd5-e4`)
}

func TestSearch130(t *testing.T) { // En-passant
	move := NewGame(`Ke7,Rf8,Ba3,Bc2,e5,g5`, `Kg7,c3,h7`).start().solve(3)
	expect.Eq(t, move, `Ba3-c1`)
}

func TestSearch140(t *testing.T) { // En-passant, stalemate
	move := NewGame(`Kc6,Rh4,Bb5,a3,c2,d3`, `Ka5,c5,d4,h5`).start().solve(3)
	expect.Eq(t, move, `c2-c4`)
}

func TestSearch150(t *testing.T) { // Stalemate after Qg7-c3
	move := NewGame(`Kb4,Qg7,Nc1`, `Kb1`).start().solve(3)
	expect.Eq(t, move, `Kb4-c3`)
}

func TestSearch160(t *testing.T) { // Pawn promotion
	move := NewGame(`Ka8,Qc4,b7`, `Ka5`).start().solve(3)
	expect.Eq(t, move, `b7-b8B`)
}

func TestSearch170(t *testing.T) { // Pawn promotion
	move := NewGame(`Kf8,Rc6,Be4,Nd7,c7`, `Ke6,d6`).start().solve(3)
	expect.Eq(t, move, `c7-c8R`)
}

func TestSearch180(t *testing.T) { // Pawn promotion
	move := NewGame(`Kc6,c7`, `Ka7`).start().solve(3)
	expect.Eq(t, move, `c7-c8R`)
}

func TestSearch190(t *testing.T) { // Pawn promotion
	move := NewGame(`Kc4,a7,c7`, `Ka5`).start().solve(3)
	expect.Eq(t, move, `c7-c8N`)
}

func TestSearch195(t *testing.T) { // King-side castle
	move := NewGame(`Ke1,Rf1,Rh1`, `Ka1`).start().solve(3)
	expect.Eq(t, move, `Rf1-f2`)
}

func TestSearch196(t *testing.T) { // Queen-side castle
	move := NewGame(`Ke1,Ra1,Rb1`, `Kg1`).start().solve(3)
	expect.Eq(t, move, `Rb1-b2`)
}

// Mate in 3.

func TestSearch200(t *testing.T) {
	move := NewGame(`Kf8,Re7,Nd5`, `Kh8,Bh5`).start().solve(5)
	expect.Eq(t, move, `Re7-g7`)
}

func TestSearch210(t *testing.T) {
	move := NewGame(`Kf8,Bf7,Nf3,e5`, `Kh8,e6,h7`).start().solve(5)
	expect.Eq(t, move, `Bf7-g8`)
}

func TestSearch220(t *testing.T) { // Pawn promotion
	move := NewGame(`Kf3,h7`, `Kh1,h3`).start().solve(5)
	expect.Eq(t, move, `h7-h8R`)
}

func TestSearch230(t *testing.T) { // Pawn promotion
	move := NewGame(`Kd8,c7,e4,f7`, `Ke6,e5`).start().solve(5)
	expect.Eq(t, move, `f7-f8R`)
}

func TestSearch240(t *testing.T) { // Pawn promotion
	move := NewGame(`Kh3,f7,g7`, `Kh6`).start().solve(5)
	expect.Eq(t, move, `g7-g8Q`)
}

func TestSearch250(t *testing.T) { // Pawn promotion
	move := NewGame(`Ke4,c7,d6,e7,f6,g7`, `Ke6`).start().solve(5)
	expect.Eq(t, move, `e7-e8B`)
}

// Mate in 4.

func TestSearch260(t *testing.T) { // Pawn promotion
	move := NewGame(`Kf6,Nf8,Nh6`, `Kh8,f7,h7`).start().solve(7)
	expect.Eq(t, move, `Nf8-e6`)
}

func TestSearch270(t *testing.T) { // Pawn promotion/stalemate
	move := NewGame(`Kf2,e7`, `Kh1,d2`).start().solve(7)
	expect.Eq(t, move, `e7-e8R`)
}

func TestSearch280(t *testing.T) { // Stalemate
	move := NewGame(`Kc1,Nb4,a2`, `Ka1,b5`).start().solve(7)
	expect.Eq(t, move, `a2-a4`)
}

func TestSearch290(t *testing.T) { // Stalemate
	move := NewGame(`Kh6,Rd3,h7`, `Kh8,Bd7`).start().solve(7)
	expect.Eq(t, move, `Rd3-d6`)
}

func TestSearch300(t *testing.T) {
	move := NewGame(`Kc6,Bc1,Ne5`, `Kc8,Ra8,a7,a6`).start().solve(7)
	expect.Eq(t, move, `Ne5-f7`)
}

// Perft.
func TestSearch400(t *testing.T) {
	position := NewGame().start()
	expect.Eq(t, position.Perft(0), int64(1))
}

func TestSearch410(t *testing.T) {
	position := NewGame().start()
	expect.Eq(t, position.Perft(1), int64(20))
}

func TestSearch420(t *testing.T) {
	position := NewGame().start()
	expect.Eq(t, position.Perft(2), int64(400))
}

func TestSearch430(t *testing.T) {
	position := NewGame().start()
	expect.Eq(t, position.Perft(3), int64(8902))
}

func TestSearch440(t *testing.T) {
	position := NewGame().start()
	expect.Eq(t, position.Perft(4), int64(197281))
}

func TestSearch450(t *testing.T) {
	position := NewGame().start()
	expect.Eq(t, position.Perft(5), int64(4865609))
}
