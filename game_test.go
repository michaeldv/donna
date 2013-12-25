// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`)

// The very first chess puzzle I had solved as a kid.
func TestGame000(t *testing.T) {
        move := NewGame().Setup(`Kf8,Rh1,g6`, `Kh8,Bg8,g7,h7`).Search(2)
        expect(t, move, `Rh1-h6`)
}

// Puzzle samples were taken from
// http://chessproblem.ru/index.php?kind=2&f_ot=0&f_do=8&lev=0

func TestGame010(t *testing.T) {
        move := NewGame().Setup(`Kf4,Qc2,Nc5`, `Kd4`).Search(2)
        expect(t, move, `Nc5-b7`)
}

func TestGame020(t *testing.T) {
        move := NewGame().Setup(`Kf8,Qf6`, `Kh7,Nf5`).Search(2)
        expect(t, move, `Qf6-g5`)
}

func TestGame030(t *testing.T) {
        move := NewGame().Setup(`Kf2,Qf7,Nf3`, `Kg4`).Search(2)
        expect(t, move, `Qf7-f6`)
}

func TestGame040(t *testing.T) {
        move := NewGame().Setup(`Kc3,Qc2,Ra4`, `Kb5`).Search(2)
        expect(t, move, `Qc2-g6`)
}

func TestGame050(t *testing.T) {
        move := NewGame().Setup(`Ke5,Qc1,Rf3,Bg2`, `Ke2,Nd5,Nb1`).Search(2)
        expect(t, move, `Rf3-d3`)
}

func TestGame060(t *testing.T) {
        move := NewGame().Setup(`Kf1,Qa8,Bf7,Ng2`, `Kg4`).Search(2)
        expect(t, move, `Qa8-b8`)
}

func TestGame070(t *testing.T) {
        move := NewGame().Setup(`Ke5,Rd3,Bb1`, `Kh7`).Search(2)
        expect(t, move, `Ke5-f6`)
}

// Puzzles with pawns.

func TestGame080(t *testing.T) {
        move := NewGame().Setup(`Kg3,Bc1,Nc3,Bg2`, `Kg1,Re1,e3`).Search(2)
        expect(t, move, `Bc1-a3`)
}

func TestGame090(t *testing.T) {
        move := NewGame().Setup(`Kf2,Qb8,Be7,f3`, `Kh5,h6,g5`).Search(2)
        expect(t, move, `Qb8-b1`)
}

func TestGame100(t *testing.T) {
        move := NewGame().Setup(`Ke6,Qg3,b3,c2`, `Ke4,e7,f5`).Search(2)
        expect(t, move, `b3-b4`)
}

func TestGame110(t *testing.T) {
        move := NewGame().Setup(`Kf1,Qh6,Nd2,Nf2`, `Kc1,c2,c3`).Search(2)
        expect(t, move, `Qh6-a6`)
}

func TestGame120(t *testing.T) { // En-passant
        move := NewGame().Setup(`Kd5,Qc8,c5,e5,g6`, `Ke7,d7`).Search(2)
        expect(t, move, `Kd5-e4`)
}

func TestGame130(t *testing.T) { // En-passant
        move := NewGame().Setup(`Ke7,Rf8,Ba3,Bc2,e5,g5`, `Kg7,c3,h7`).Search(2)
        expect(t, move, `Ba3-c1`)
}

func TestGame140(t *testing.T) { // En-passant, stalemate
        move := NewGame().Setup(`Kc6,Rh4,Bb5,a3,c2,d3`, `Ka5,c5,d4,h5`).Search(2)
        expect(t, move, `c2-c4`)
}

func TestGame150(t *testing.T) { // Stalemate after Qg7-c3
        move := NewGame().Setup(`Kb4,Qg7,Nc1`, `Kb1`).Search(2)
        expect(t, move, `Kb4-c3`)
}

func TestGame160(t *testing.T) { // Pawn promotion
        move := NewGame().Setup(`Ka8,Qc4,b7`, `Ka5`).Search(2)
        expect(t, move, `b7-b8B`)
}

func TestGame170(t *testing.T) { // Pawn promotion
        move := NewGame().Setup(`Kf8,Rc6,Be4,Nd7,c7`, `Ke6,d6`).Search(2)
        expect(t, move, `c7-c8R`)
}

func TestGame180(t *testing.T) { // Pawn promotion
        move := NewGame().Setup(`Kc6,c7`, `Ka7`).Search(2)
        expect(t, move, `c7-c8R`)
}

func TestGame190(t *testing.T) { // Pawn promotion
        move := NewGame().Setup(`Kc4,a7,c7`, `Ka5`).Search(2)
        expect(t, move, `c7-c8N`)
}

func TestGame195(t *testing.T) { // King-side castle
        move := NewGame().Setup(`Ke1,Rf1,Rh1`, `Ka1`).Search(2)
        expect(t, move, `Rf1-f2`)
}

func TestGame196(t *testing.T) { // Queen-side castle
        move := NewGame().Setup(`Ke1,Ra1,Rb1`, `Kg1`).Search(2)
        expect(t, move, `Rb1-b2`)
}


// Mate in 3

func TestGame200(t *testing.T) {
        move := NewGame().Setup(`Kf8,Re7,Nd5`, `Kh8,Bh5`).Search(3)
        expect(t, move, `Re7-g7`)
}

func TestGame210(t *testing.T) {
        move := NewGame().Setup(`Kf8,Bf7,Nf3,e5`, `Kh8,e6,h7`).Search(3)
        expect(t, move, `Bf7-g8`)
}

func TestGame220(t *testing.T) { // Pawn promotion
        move := NewGame().Setup(`Kf3,h7`, `Kh1,h3`).Search(3)
        expect(t, move, `h7-h8R`)
}

func TestGame230(t *testing.T) { // Pawn promotion
        move := NewGame().Setup(`Kd8,c7,e4,f7`, `Ke6,e5`).Search(3)
        expect(t, move, `f7-f8R`)
}

func TestGame240(t *testing.T) { // Pawn promotion
        move := NewGame().Setup(`Kh3,f7,g7`, `Kh6`).Search(3)
        expect(t, move, `g7-g8Q`)
}

func TestGame250(t *testing.T) { // Pawn promotion
        move := NewGame().Setup(`Ke4,c7,d6,e7,f6,g7`, `Ke6`).Search(3)
        expect(t, move, `e7-e8B`)
}

// Mate in 4

func TestGame260(t *testing.T) { // Pawn promotion
        move := NewGame().Setup(`Kf6,Nf8,Nh6`, `Kh8,f7,h7`).Search(4)
        expect(t, move, `Nf8-e6`)
}
