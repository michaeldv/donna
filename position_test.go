// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`)

func TestPosition010(t *testing.T) {
        p := NewGame().Setup(`Ke1,e2`, `Kg8,d7,f7`).Start()
        expect(t, p.enpassant, 0)

        p = p.MakeMove(NewMove(p, E2, E4))
        expect(t, p.enpassant, 0)

        p = p.MakeMove(NewMove(p, D7, D5))
        expect(t, p.enpassant, 0)

        p = p.MakeMove(NewMove(p, E4, E5))
        expect(t, p.enpassant, 0)

        p = p.MakeMove(NewMove(p, F7, F5))
        expect(t, p.enpassant, F6)
}

func TestPosition020(t *testing.T) {
        game := NewGame().Setup(`Ke1,b7,e4`, `Kg8,Ra8,h2`)
        p := game.Start()

        expect(t, p.isPawnPromotion(Pawn(WHITE), A8), true)
        expect(t, p.isPawnPromotion(Pawn(WHITE), B8), true)
        expect(t, p.isPawnPromotion(Pawn(WHITE), E5), false)
        expect(t, p.isPawnPromotion(Pawn(BLACK), H1), true)
}

// Castle tests.

func TestPosition030(t *testing.T) { // Everything is OK.
        p := NewGame().Setup(`Ke1,Ra1,Rh1`, `Ke8`).Start() // White to move.

        expect(t, p.isKingSideCastleAllowed(p.color), true)
        expect(t, p.isQueenSideCastleAllowed(p.color), true)

        p = NewGame().Setup(`Ke1`, `Ke8,Ra8,Rh8`).Start()  // White to move.
        p.color ^= 1                                       // Black to move.

        expect(t, p.isKingSideCastleAllowed(p.color), true)
        expect(t, p.isQueenSideCastleAllowed(p.color), true)
}


func TestPosition040(t *testing.T) { // King checked.
        p := NewGame().Setup(`Ke1,Ra1,Rh1`, `Ke8,Bg3`).Start()

        expect(t, p.isKingSideCastleAllowed(p.color), false)
        expect(t, p.isQueenSideCastleAllowed(p.color), false)

        p = NewGame().Setup(`Ke1,Bg6`, `Ke8,Ra8,Rh8`).Start()
        p.color ^= 1

        expect(t, p.isKingSideCastleAllowed(p.color), false)
        expect(t, p.isQueenSideCastleAllowed(p.color), false)
}

func TestPosition050(t *testing.T) { // Attacked square.
        p := NewGame().Setup(`Ke1,Ra1,Rh1`, `Ke8,Bb3,Bh3`).Start()

        expect(t, p.isKingSideCastleAllowed(p.color), false)
        expect(t, p.isQueenSideCastleAllowed(p.color), false)

        p = NewGame().Setup(`Ke1,Bb6,Bh6`, `Ke8,Ra8,Rh8`).Start()
        p.color ^= 1

        expect(t, p.isKingSideCastleAllowed(p.color), false)
        expect(t, p.isQueenSideCastleAllowed(p.color), false)
}

func TestPosition060(t *testing.T) { // Wrong square.
        p := NewGame().Setup(`Ke1,Ra8,Rh8`, `Ke5`).Start()

        expect(t, p.isKingSideCastleAllowed(p.color), false)
        expect(t, p.isQueenSideCastleAllowed(p.color), false)

        p = NewGame().Setup(`Ke2,Ra1,Rh1`, `Ke8`).Start()

        expect(t, p.isKingSideCastleAllowed(p.color), false)
        expect(t, p.isQueenSideCastleAllowed(p.color), false)

        p = NewGame().Setup(`Ke4`, `Ke8,Ra1,Rh1`).Start()
        p.color ^= 1

        expect(t, p.isKingSideCastleAllowed(p.color), false)
        expect(t, p.isQueenSideCastleAllowed(p.color), false)

        p = NewGame().Setup(`Ke4`, `Ke7,Ra8,Rh8`).Start()
        p.color ^= 1

        expect(t, p.isKingSideCastleAllowed(p.color), false)
        expect(t, p.isQueenSideCastleAllowed(p.color), false)
}

func TestPosition070(t *testing.T) { // Missing rooks.
        p := NewGame().Setup(`Ke1`, `Ke8`).Start()

        expect(t, p.isKingSideCastleAllowed(p.color), false)
        expect(t, p.isQueenSideCastleAllowed(p.color), false)

        p = NewGame().Setup(`Ke1`, `Ke8`).Start()
        p.color ^= 1

        expect(t, p.isKingSideCastleAllowed(p.color), false)
        expect(t, p.isQueenSideCastleAllowed(p.color), false)
}

func TestPosition080(t *testing.T) { // Rooks on wrong squares.
        p := NewGame().Setup(`Ke1,Rb1`, `Ke8`).Start()

        expect(t, p.isKingSideCastleAllowed(p.color), false)
        expect(t, p.isQueenSideCastleAllowed(p.color), false)

        p = NewGame().Setup(`Ke1,Rb1,Rh1`, `Ke8`).Start()

        expect(t, p.isKingSideCastleAllowed(p.color), true)
        expect(t, p.isQueenSideCastleAllowed(p.color), false)

        p = NewGame().Setup(`Ke1,Ra1,Rf1`, `Ke8`).Start()

        expect(t, p.isKingSideCastleAllowed(p.color), false)
        expect(t, p.isQueenSideCastleAllowed(p.color), true)
}
