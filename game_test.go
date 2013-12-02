package lape

import (`testing`)

func expect(t *testing.T, move *Move, answer string) {
        right := move.String()
        if answer != right {
                t.Error(`Expected: ` + answer + `, got: ` + right)
        } else {
                t.Log(right)
        }
}

func Test01(t *testing.T) {
        Settings.Log = false
        Settings.Fancy = false
        move := new(Game).Initialize().Setup(`Qc2,Nc5,Kf4`, `Kd4`).Search(2)
        expect(t, move, `Nc5-b7`)
}

func Test02(t *testing.T) {
        move := new(Game).Initialize().Setup(`Qf6,Kf8`, `Kh7,Nf5`).Search(2)
        expect(t, move, `Qf6-g5`)
}

func Test03(t *testing.T) {
        move := new(Game).Initialize().Setup(`Qf7,Nf3,Kf2`, `Kg4`).Search(2)
        expect(t, move, `Qf7-f6`)
}

func Test04(t *testing.T) {
        move := new(Game).Initialize().Setup(`Qc2,Kc3,Ra4`, `Kb5`).Search(2)
        expect(t, move, `Qc2-g6`)
}

func Test05(t *testing.T) {
        move := new(Game).Initialize().Setup(`Ke5,Qc1,Rf3,Bg2`, `Ke2,Nd5,Nb1`).Search(2)
        expect(t, move, `Rf3-d3`)
}

func Test06(t *testing.T) {
        move := new(Game).Initialize().Setup(`Qa8,Bf7,Ng2,Kf1`, `Kg4`).Search(2)
        expect(t, move, `Qa8-b8`)
}

func Test07(t *testing.T) {
        move := new(Game).Initialize().Setup(`Bb1,Rd3,Ke5`, `Kh7`).Search(2)
        expect(t, move, `Ke5-f6`)
}

// game.Setup(`Kb4,Nc1,Qg7`, `Kb1`) // TODO: Qc3 stalemate vs. Kc3+-
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
