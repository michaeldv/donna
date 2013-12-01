package lape

import (`testing`)

func Test01(t *testing.T) {
        Settings.Log = false
        Settings.Fancy = false
        game := new(Game).Initialize().Setup(`Qc2,Nc5,Kf4`, `Kd4`) // Nb7
        move := game.Search(2)
        if move.String() != `♘ c5-b7` {
                t.Error(`01`)
        }
}

func Test02(t *testing.T) {
        game := new(Game).Initialize().Setup(`Qf6,Kf8`, `Kh7,Nf5`) // Qg5
        move := game.Search(2)
        if move.String() != `♕ f6-g5` {
                t.Error(`02`)
        }
}