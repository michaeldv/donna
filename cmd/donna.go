// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package main

import (
        `fmt`
        `time`
        `github.com/michaeldv/donna`
)

func repl() {
        var game *donna.Game
        var move *donna.Move
        var position *donna.Position

        for command := ``; ; command = `` {
                fmt.Print(`donna> `)
                fmt.Scanf(`%s`, &command)
                switch command {
                case ``:
                case `bench`:
                        benchmark()
                case `perft`:
                        perft(5)
                case `exit`, `quit`:
                        return
                case `help`:
                        fmt.Println(`help: not implemented yet.`)
                case `new`:
                        game = donna.NewGame().InitialPosition()
                        position = game.Start(donna.White)
                        fmt.Printf("%s\n", position)
                case `go`:
                        if game == nil || position == nil {
                                game = donna.NewGame().InitialPosition()
                                position = game.Start(donna.White)
                        }
                        move = game.Think(3, position)
                        position = position.MakeMove(move)
                        fmt.Printf("%s\n", position)
                default:
                        if game == nil || position == nil {
                                game = donna.NewGame().InitialPosition()
                                position = game.Start(donna.White)
                        }
                        move = donna.NewMoveFromString(command, position)
                        if move != nil {
                                position = position.MakeMove(move)
                                fmt.Printf("%s\n", position)
                                move = game.Think(3, position)
                                position = position.MakeMove(move)
                                fmt.Printf("%s\n", position)
                        } else {
                                fmt.Printf("%s appears to be an invalid move.\n", command)
                        }
                }
        }
}

func main() {
        donna.Settings.Log = false//true
        donna.Settings.Fancy = true

        // donna.NewGame().Setup(`Ka7,Qb1,Bg2`, `Ka5,b3,g3`).Think(4, nil) // Qb2
        // donna.NewGame().Setup(`Kh5,Qg7,Be5,f2,f3`, `Kh1`).Think(4, nil) // Bh2
        // donna.NewGame().Setup(`Kd3,Rd8,a5,b2,f2,g5`, `Kd1`).Think(4, nil) // Rd4
        repl()
}

//
// Bobby Fischer vs. James Sherwin, New Jersey Open 1957 after 16 moves.
// http://www.chessgames.com/perl/chessgame?gid=1008366
// Fischer played 17. h2-h4
//
func benchmark() {
        game := donna.NewGame().Setup(`Kg1,Qc2,Ra1,Re1,Bc1,Bg2,Ng5,a2,b2,c3,d4,f2,g3,h2`,
                                      `Kg8,Qd6,Ra8,Rf8,Bc8,Nd5,Ng6,a7,b6,c4,e6,f7,g7,h7`)

        fmt.Printf("%s\n", game)
        for i := 0;  i < 3; i++ {
                game.Think(3, nil)
        }
}

func perft(depth int) (total int64){
        p := donna.NewGame().InitialPosition().Start(donna.White)

        start := time.Now()
        moves := p.Moves(0)
        for _, move := range moves {
                if position := p.MakeMove(move); position != nil {
                        delta := position.Perft(depth - 1)
                        total += delta
                        position.TakeBack(move)
                        fmt.Printf("%7s - %d\n", move, delta)
                }
        }
        finish := time.Since(start).Seconds()
        fmt.Printf("\n  Nodes: %d\n", total)
        fmt.Printf("Elapsed: %.2fs\n", finish)
        fmt.Printf("Nodes/s: %.2f\n", float64(total) / finish)
        return
}
