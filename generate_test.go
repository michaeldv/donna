package donna

import (`testing`; `fmt`)

func TestGenerate000(t *testing.T) {
        game := NewGame().InitialPosition()
        moves := game.Start().Moves()

        // Possible moves includes illegal move that causes a check.
        expect(t, moves, `[Nb1-c3 Ng1-f3 e2-e4 d2-d4 c2-c4 e2-e3 d2-d3 f2-f4 f2-f3 c2-c3 b2-b3 g2-g4 g2-g3 b2-b4 h2-h4 a2-a3 a2-a4 h2-h3 Nb1-a3 Ng1-h3]`)
}

func TestGenerate010(t *testing.T) {
        Settings.Log = false
        Settings.Fancy = false
        game := NewGame().Setup(`a2,b3,c4,d2,e6,f5,g4,h3`, `Kg8`)
        moves := game.Start().Moves()

        // Moves should be sorted by relative strength.
        expect(t, moves, `[h3-h4 a2-a4 e6-e7 a2-a3 g4-g5 f5-f6 b3-b4 c4-c5 d2-d4 d2-d3]`)
}

func TestGenerate020(t *testing.T) {
        game := NewGame().Setup(`a2,b3,c4,d2,e6,f5,g4,h3`, `a3,b4,c5,e7,f6,g5,h4,Kg8`)
        moves := game.Start().Moves()

        expect(t, moves, `[d2-d4 d2-d3]`)
}

func TestGenerate030(t *testing.T) {
        game := NewGame().Setup(`a2,e4,g2`, `b3,f5,f3,h3,Kg8`)
        moves := game.Start().Moves()

        expect(t, moves, `[a2xb3 g2xf3 g2xh3 e4xf5 a2-a4 a2-a3 g2-g4 g2-g3 e4-e5]`)
}

// Should not include castles when rook has moved.
func TestGenerate040(t *testing.T) {
        game := NewGame().Setup(`Ke1,Rf1,g2`, `Ke8`)
        moves := fmt.Sprintf(`%v`, game.Start().Moves())

        doesNotContain(t, moves, `0-0`)
}

func TestGenerate050(t *testing.T) {
        game := NewGame().Setup(`Ke1,Rb1,b2`, `Ke8`)
        moves := fmt.Sprintf(`%v`, game.Start().Moves())

        doesNotContain(t, moves, `0-0`)
}

// Should not include castles when king has moved.
func TestGenerate060(t *testing.T) {
        game := NewGame().Setup(`Kf1,Ra1,a2,Rh1,h2`, `Ke8`)
        moves := fmt.Sprintf(`%v`, game.Start().Moves())

        doesNotContain(t, moves, `0-0`)
}

// Should not include castles when rooks are not there.
func TestGenerate070(t *testing.T) {
        game := NewGame().Setup(`Ke1`, `Ke8`)
        moves := fmt.Sprintf(`%v`, game.Start().Moves())

        doesNotContain(t, moves, `0-0`)
}

// Should not include castles when king is in check.
func TestGenerate080(t *testing.T) {
        game := NewGame().Setup(`Ke1,Ra1,Rf1`, `Ke8,Re7`)
        moves := fmt.Sprintf(`%v`, game.Start().Moves())

        doesNotContain(t, moves, `0-0`)
}

// Should not include castles when target square is a capture.
func TestGenerate090(t *testing.T) {
        game := NewGame().Setup(`Ke1,Ra1,Rf1`, `Ke8,Nc1,Ng1`)
        moves := fmt.Sprintf(`%v`, game.Start().Moves())

        doesNotContain(t, moves, `0-0`)
}

// Should not include castles when king is to jump over attacked square.
func TestGenerate100(t *testing.T) {
        game := NewGame().Setup(`Ke1,Ra1,Rf1`, `Ke8,Bc4,Bf4`)
        moves := fmt.Sprintf(`%v`, game.Start().Moves())

        doesNotContain(t, moves, `0-0`)
}
