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
}