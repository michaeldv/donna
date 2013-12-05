package main

import (
        `github.com/michaeldv/donna`
        `fmt`
)

func main() {
        donna.Settings.Log = true
        donna.Settings.Fancy = true
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

        game := donna.NewGame().InitialPosition()
        // fmt.Println("...Initial Position...\n")
        // fmt.Println(game.String())

        // fmt.Println("...Moves...\n")
        // m := donna.NewMove(donna.D2, donna.D4, donna.Pawn(donna.WHITE), 0)
        // fmt.Println(m.String())
        // m = donna.NewMove(donna.G8, donna.F6, donna.Knight(donna.BLACK), 0)
        // fmt.Println(m.String())
        // m = donna.NewMove(donna.E2, donna.E4, donna.Pawn(donna.WHITE), 0)
        // fmt.Println(m.String())
        // m = donna.NewMove(donna.F6, donna.E4, donna.Knight(donna.BLACK), donna.Pawn(donna.WHITE))
        // fmt.Println(m.String())

        // http://chessproblem.ru/index.php?kind=2&f_ot=0&f_do=8&lev=0
        move := game.Search(2)
        fmt.Printf("Best move: %s\n", move)
}