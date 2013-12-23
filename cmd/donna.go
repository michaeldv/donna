// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package main

import (
        `fmt`
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
                case `exit`, `quit`:
                        return
                case `help`:
                        fmt.Println(`help: not implemented yet.`)
                case `new`:
                        game = donna.NewGame().InitialPosition()
                        position = game.Start()
                        fmt.Printf("%s\n", position)
                case `go`:
                        if game == nil || position == nil {
                                game = donna.NewGame().InitialPosition()
                                position = game.Start()
                        }
                        move = game.Think(3, position)
                        position = position.MakeMove(move)
                        fmt.Printf("%s\n", position)
                default:
                        if game == nil || position == nil {
                                game = donna.NewGame().InitialPosition()
                                position = game.Start()
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
