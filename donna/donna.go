package main

import (
        `github.com/michaeldv/donna`
)

func main() {
        donna.Settings.Log = false//true
        donna.Settings.Fancy = true

        // game := donna.NewGame().Setup(`Ke5,Rd2,Rf2,Ne8`, `Ke3,e6,e4`) // Rb2
        // game := donna.NewGame().Setup(`Ka1,Rh6,Na6,d3`, `Kb5,a5`) // Kb2
        game := donna.NewGame().Setup(`Kf6,Nf8,Nh6`, `Kh8,f7,h7`) // Ne6
        game.Think(4)
}