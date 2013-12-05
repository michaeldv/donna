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

        game := lape.NewGame().InitialPosition()
        // fmt.Println("...Initial Position...\n")
        // fmt.Println(game.String())

        // fmt.Println("...Moves...\n")
        // m := lape.NewMove(lape.Index(1,4), lape.Index(3,4), lape.Pawn(0), 0)
        // fmt.Println(m.String())
        // m = lape.NewMove(lape.Index(7,6), lape.Index(5,5), lape.Knight(1), 0)
        // fmt.Println(m.String())
        // m = lape.NewMove(lape.Index(1,3), lape.Index(3,3), lape.Pawn(0), 0)
        // fmt.Println(m.String())
        // m = lape.NewMove(lape.Index(5,5), lape.Index(3,4), lape.Knight(1), lape.Pawn(0))
        // fmt.Println(m.String())

        // http://chessproblem.ru/index.php?kind=2&f_ot=0&f_do=8&lev=0
        move := game.Search(2)
        fmt.Printf("Best move: %s\n", move)
}