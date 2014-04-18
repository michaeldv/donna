// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
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
        var position *donna.Position

        setup := func() {
                if game == nil || position == nil {
                        game = donna.NewGame().InitialPosition()
                        position = game.CacheSize(32).Start(donna.White)
                        fmt.Printf("%s\n", position)
                }
        }

        think := func() {
                if move := game.Think(8, position); move != 0 {
                        position = position.MakeMove(move)
                        fmt.Printf("%s\n", position)
                }
        }

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
                        game, position = nil, nil
                        setup()
                case `go`:
                        setup()
                        think()
                default:
                        setup()
                        if move := position.NewMoveFromString(command); move != 0 {
                                if advance := position.MakeMove(move); advance != nil {
                                        position = advance
                                        think()
                                        continue
                                }
                        }
                        // Invalid move (typo) or non-evasion of check.
                        fmt.Printf("%s appears to be an invalid move.\n", command)
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

func benchmark() {
	// Sergey Kaminer, 1935
	// 1.h8Q+ Kxh8 2.Ng4+ Kg7 3.Qh7+ Kf8 4.Qg8+ Ke7 5.Qe8+ Kd6 6.Qd7+ Kc5 7.Qb5+ Kd6
	// 8.Qb6+ Ke7 9.Qc7+ Kf8 10.Qc8+ Ke7 11.Qd7+ Kf8 12.Qe8+ Kg7 13.Qg8+ Kxg8 14.Kxf6+
	game := donna.NewGame().Setup(`Kd1,Qh2,Nh6,a4,g3,h7`, `Kg7,Qe4,Bf6,b7,e6,g6`)
	fmt.Printf("%s\n", game)
	game.CacheSize(32).Think(8, game.Start(donna.White))

        // Bobby Fischer vs. James Sherwin, New Jersey Open 1957 after 16 moves.
        // http://www.chessgames.com/perl/chessgame?gid=1008366
        // Fischer played 17. h2-h4!
        game = donna.NewGame().Setup(`Kg1,Qc2,Ra1,Re1,Bc1,Bg2,Ng5,a2,b2,c3,d4,f2,g3,h2`, `Kg8,Qd6,Ra8,Rf8,Bc8,Nd5,Ng6,a7,b6,c4,e6,f7,g7,h7`)
        fmt.Printf("%s\n", game)
        game.CacheSize(32).Think(9, game.Start(donna.White))

        // Mikhail Botvinnik vs. Jose Raul Capablanca, AVRO 1936 after 29 moves.
        // Botvinnik played 30. Bb2-a3!
        game = donna.NewGame().Setup(`Kg1,Qe5,Bb2,Ng3,c3,d4,e6,g2,h2`, `Kg7,Qe7,Nb3,Nf6,a7,b6,c4,d5,g6,h7`)
        fmt.Printf("%s\n", game)
        game.CacheSize(32).Think(9, game.Start(donna.White))

}

func perft(depth int) (total int64){
        p := donna.NewGame().InitialPosition().Start(donna.White)
        start := time.Now()
	total = p.Perft(depth)
        finish := time.Since(start).Seconds()
        fmt.Printf("\n  Nodes: %d\n", total)
        fmt.Printf("Elapsed: %.2fs\n", finish)
        fmt.Printf("Nodes/s: %.2f\n", float64(total) / finish)
        return
}
