package main

import (
        `github.com/michaeldv/donna`
)

func main() {
        donna.Settings.Log = false//true
        donna.Settings.Fancy = true

        game := donna.NewGame().Setup(`Kf8,Re7,Nd5`, `Kh8,Bh5`) // .InitialPosition()
        game.Think(3)
}