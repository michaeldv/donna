package main

import (
        `github.com/michaeldv/lape`
        `fmt`
)

func main() {
        lape.Settings.Log = true
        lape.Settings.Fancy = true
        fmt.Printf("Main...\n")

        // fmt.Println(b.Rook[0].String())
        // fmt.Println(b.Rook[9].String())
        // fmt.Println(b.Rook[63].String())
        // 
        // fmt.Println(b.Knight[0].String())
        // fmt.Println(b.Knight[9].String())
        // fmt.Println(b.Knight[63].String())
        // 
        // fmt.Println(b.Bishop[0].String())
        // fmt.Println(b.Bishop[9].String())
        // fmt.Println(b.Bishop[63].String())
        // 
        // fmt.Println(b.Queen[0].String())
        // fmt.Println(b.Queen[9].String())
        // fmt.Println(b.Queen[63].String())
        // 
        // fmt.Println(b.King[0].String())
        // fmt.Println(b.King[9].String())
        // fmt.Println(b.King[63].String())

        game := new(lape.Game).Initialize()
        // fmt.Println("...Initial Position...\n")
        // fmt.Println(game.String())

        // fmt.Println("...Moves...\n")
        // m := new(lape.Move).Initialize(lape.Index(1,4), lape.Index(3,4), lape.Pawn(0), 0)
        // fmt.Println(m.String())
        // m = new(lape.Move).Initialize(lape.Index(7,6), lape.Index(5,5), lape.Knight(1), 0)
        // fmt.Println(m.String())
        // m = new(lape.Move).Initialize(lape.Index(1,3), lape.Index(3,3), lape.Pawn(0), 0)
        // fmt.Println(m.String())
        // m = new(lape.Move).Initialize(lape.Index(5,5), lape.Index(3,4), lape.Knight(1), lape.Pawn(0))
        // fmt.Println(m.String())

        // http://chessproblem.ru/index.php?kind=2&f_ot=0&f_do=8&lev=0
        game.Setup(`Qc2,Nc5,Kf4`, `Kd4`) // Nb7
        // game.Setup(`Qf6,Kf8`, `Kh7,Nf5`) // Qg5
        // game.Setup(`Qf7,Nf3,Kf2`, `Kg4`) // Qf6
        // game.Setup(`Qc2,Kc3,Ra4`, `Kb5`) // Qg6
        // game.Setup(`Kb4,Nc1,Qg7`, `Kb1`) // TODO: Qc3 stalemate vs. Kc3+-
        // game.Setup(`Ke5,Qc1,Rf3,Bg2`, `Ke2,Nd5,Nb1`) // Rd3
        // game.Setup(`Qa8,Bf7,Ng2,Kf1`, `Kg4`) // Qb8
        // game.Setup(`Bb1,Rd3,Ke5`, `Kh7`) // Kf1
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
        move := game.Search(2)
        fmt.Printf("Best move: %s\n", move)
}