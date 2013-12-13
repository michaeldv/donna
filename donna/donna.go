package main

import (
        `github.com/michaeldv/donna`
)

func main() {
        donna.Settings.Log = false//true
        donna.Settings.Fancy = true

        donna.NewGame().Setup(`Ka7,Qb1,Bg2`, `Ka5,b3,g3`).Think(4) // Qb2
        // donna.NewGame().Setup(`Kh5,Qg7,Be5,f2,f3`, `Kh1`).Think(4) // Bh2
        // donna.NewGame().Setup(`Kd3,Rd8,a5,b2,f2,g5`, `Kd1`).Think(4) // Rd4
}