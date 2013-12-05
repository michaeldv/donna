package donna

import (`testing`)

// The very first chess puzzle I had solved as a kid.
func TestGame000(t *testing.T) {
        move := NewGame().Setup(`Kf8,Rh1,g6`, `Kh8,Bg8,g7,h7`).Search(2)
        expect(t, move.String(), `Rh1-h6`)
}

// Puzzle samples were taken from
// http://chessproblem.ru/index.php?kind=2&f_ot=0&f_do=8&lev=0

func TestGame010(t *testing.T) {
        move := NewGame().Setup(`Kf4,Qc2,Nc5`, `Kd4`).Search(2)
        expect(t, move.String(), `Nc5-b7`)
}

func TestGame020(t *testing.T) {
        move := NewGame().Setup(`Kf8,Qf6`, `Kh7,Nf5`).Search(2)
        expect(t, move.String(), `Qf6-g5`)
}

func TestGame030(t *testing.T) {
        move := NewGame().Setup(`Kf2,Qf7,Nf3`, `Kg4`).Search(2)
        expect(t, move.String(), `Qf7-f6`)
}

func TestGame040(t *testing.T) {
        move := NewGame().Setup(`Kc3,Qc2,Ra4`, `Kb5`).Search(2)
        expect(t, move.String(), `Qc2-g6`)
}

func TestGame050(t *testing.T) {
        move := NewGame().Setup(`Ke5,Qc1,Rf3,Bg2`, `Ke2,Nd5,Nb1`).Search(2)
        expect(t, move.String(), `Rf3-d3`)
}

func TestGame060(t *testing.T) {
        move := NewGame().Setup(`Kf1,Qa8,Bf7,Ng2`, `Kg4`).Search(2)
        expect(t, move.String(), `Qa8-b8`)
}

func TestGame070(t *testing.T) {
        move := NewGame().Setup(`Ke5,Rd3,Bb1`, `Kh7`).Search(2)
        expect(t, move.String(), `Ke5-f6`)
}

// Puzzles with pawns.

func TestGame080(t *testing.T) {
        move := NewGame().Setup(`Kg3,Bc1,Nc3,Bg2`, `Kg1,Re1,e3`).Search(2)
        expect(t, move.String(), `Bc1-a3`)
}

func TestGame090(t *testing.T) {
        move := NewGame().Setup(`Kf2,Qb8,Be7,f3`, `Kh5,h6,g5`).Search(2)
        expect(t, move.String(), `Qb8-b1`)
}

func TestGame100(t *testing.T) {
        move := NewGame().Setup(`Ke6,Qg3,b3,c2`, `Ke4,e7,f5`).Search(2)
        expect(t, move.String(), `b3-b4`)
}

func TestGame110(t *testing.T) { // EN-PASSANT
        move := NewGame().Setup(`Kd5,Qc8,c5,e5,g6`, `Ke7,d7`).Search(2)
        expect(t, move.String(), `Kd5-e4`)
}

func TestGame120(t *testing.T) { // EN-PASSANT
        move := NewGame().Setup(`Ke7,Rf8,Ba3,Bc2,e5,g5`, `Kg7,c3,h7`).Search(2)
        expect(t, move.String(), `Ba3-c1`)
}

func TestGame130(t *testing.T) { // EN-PASSANT, STALEMATE
        move := NewGame().Setup(`Kc6,Rh4,Bb5,a3,c2,d3`, `Ka5,c5,d4,h5`).Search(2)
        expect(t, move.String(), `c2-c4`)
}

func TestGame140(t *testing.T) { // TODO: stalemate after Qg7-c3
        move := NewGame().Setup(`Kb4,Qg7,Nc1`, `Kb1`).Search(2)
        expect(t, move.String(), `Kb4-c3`)
}
//
// func TestGame???(t *testing.T) { // TODO: stalemate (both Qa6 and Qb6 have 32767.00 score)
//         move := NewGame().Setup(`Kf1,Qh6,Nd2,Nf2`, `Kc1,c2,c3`).Search(2)
//         expect(t, move.String(), `Qh6-a6`)
// }
//

// game.Setup(`Kg1,Qh1,Bh8,g2`, `Kg8,Rf8,f7,g6,h7`)
// game.Setup(`Kh1,Ra7,Rc7,Ba8`, `Kh8`)
// game.Setup(`Kh1,h2,g2,Qh4,Bf6,g5,g4,d4`, `Kg8,Rf8,f7,g6,h7,c8`)
// game.Setup(`Kh1,g2,h2,Nh6,Qe6`, `Kh8,Rf8,g7,h7`)
// game.Setup(`Kh1,Ra6,Rb5`, `Kh7`)
// game.Setup(`Kh1,Ra1`, `Kg8,f7,g7,h7`)
//
// game.Setup(`Kg1,f2,g2,h2`, `Kg8,Ra1`)
// game.Setup(`Kg1,f3,e2,e3`, `Kh3,Ra1`)
// game.Setup(`d2,f3,g2,Rf2,Kg1`, `Kg3,Ra1`)
// game.Setup(`a3,Bb4,a5,c3,e7,Kh2`, `a7,a5,b6,Bc7,Kg8`)
// game.Setup(`a2,Ra3,b3,a7,Kg1`, `d4,Rc4,c3,c5,Bb6,Kg8`)
