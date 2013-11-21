package main

import (
        `github.com/michaeldv/lape`
        `fmt`
)

func main() {
        b := new(lape.Bitboard).Initialize()
        fmt.Printf("Main... %v\n", b)

        fmt.Println(b.Rook[0].ToString())
        fmt.Println(b.Rook[9].ToString())
        fmt.Println(b.Rook[63].ToString())

        fmt.Println(b.Knight[0].ToString())
        fmt.Println(b.Knight[9].ToString())
        fmt.Println(b.Knight[63].ToString())

        fmt.Println(b.Bishop[0].ToString())
        fmt.Println(b.Bishop[9].ToString())
        fmt.Println(b.Bishop[63].ToString())

        fmt.Println(b.Queen[0].ToString())
        fmt.Println(b.Queen[9].ToString())
        fmt.Println(b.Queen[63].ToString())

        fmt.Println(b.King[0].ToString())
        fmt.Println(b.King[9].ToString())
        fmt.Println(b.King[63].ToString())

        game := new(lape.Game).Initialize().SetInitialPosition()
        fmt.Println("...Initial Position...\n")
        fmt.Println(game.ToString())

        fmt.Println("...Moves...\n")
        m := new(lape.Move).Initialize(lape.Index(1,4), lape.Index(3,4), lape.Pawn(0), 0)
        fmt.Println(m.ToString())
        m = new(lape.Move).Initialize(lape.Index(7,6), lape.Index(5,5), lape.Knight(1), 0)
        fmt.Println(m.ToString())
        m = new(lape.Move).Initialize(lape.Index(1,3), lape.Index(3,3), lape.Pawn(0), 0)
        fmt.Println(m.ToString())
        m = new(lape.Move).Initialize(lape.Index(5,5), lape.Index(3,4), lape.Knight(1), lape.Pawn(0))
        fmt.Println(m.ToString())

        fmt.Println("...Make Move...\n")

        move := game.MakeMove(1)
        fmt.Println(move.ToString())
}