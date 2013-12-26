// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`)

func TestPosition010(t *testing.T) {
        game := NewGame().Setup(`Ke1,e2`, `Kg8,d7`)
        initial := game.Start()
        expect(t, initial.enpassant, 0)

        position := initial.MakeMove(NewMove(E2, E4, Pawn(WHITE), Piece(0)))
        expect(t, position.enpassant, E3)

        initial.color = BLACK
        position = initial.MakeMove(NewMove(D7, D5, Pawn(BLACK), Piece(0)))
        expect(t, position.enpassant, D6)
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
        game := NewGame().Setup(`Ke1,Ra1,Rh1`, `Ke8`)
        p := game.Start()

        expect(t, p.isKingSideCastleAllowed(), true)
        expect(t, p.isQueenSideCastleAllowed(), true)

        game = NewGame().Setup(`Ke1`, `Ke8,Ra8,Rh8`)
        p = NewPosition(game, game.pieces, BLACK, 0)

        expect(t, p.isKingSideCastleAllowed(), true)
        expect(t, p.isQueenSideCastleAllowed(), true)
}


func TestPosition040(t *testing.T) { // King checked.
        game := NewGame().Setup(`Ke1,Ra1,Rh1`, `Ke8,Bg3`)
        p := game.Start()

        expect(t, p.isKingSideCastleAllowed(), false)
        expect(t, p.isQueenSideCastleAllowed(), false)

        game = NewGame().Setup(`Ke1,Bg6`, `Ke8,Ra8,Rh8`)
        p = NewPosition(game, game.pieces, BLACK, 0)

        expect(t, p.isKingSideCastleAllowed(), false)
        expect(t, p.isQueenSideCastleAllowed(), false)
}

func TestPosition050(t *testing.T) { // Attacked square.
        game := NewGame().Setup(`Ke1,Ra1,Rh1`, `Ke8,Bb3,Bh3`)
        p := game.Start()

        expect(t, p.isKingSideCastleAllowed(), false)
        expect(t, p.isQueenSideCastleAllowed(), false)

        game = NewGame().Setup(`Ke1,Bb6,Bh6`, `Ke8,Ra8,Rh8`)
        p = NewPosition(game, game.pieces, BLACK, 0)

        expect(t, p.isKingSideCastleAllowed(), false)
        expect(t, p.isQueenSideCastleAllowed(), false)
}

func TestPosition060(t *testing.T) { // Wrong square.
        game := NewGame().Setup(`Ke1,Ra8,Rh8`, `Ke5`)
        p := game.Start()

        expect(t, p.isKingSideCastleAllowed(), false)
        expect(t, p.isQueenSideCastleAllowed(), false)

        game = NewGame().Setup(`Ke2,Ra1,Rh1`, `Ke8`)
        p = game.Start()

        expect(t, p.isKingSideCastleAllowed(), false)
        expect(t, p.isQueenSideCastleAllowed(), false)

        game = NewGame().Setup(`Ke4`, `Ke8,Ra1,Rh1`)
        p = NewPosition(game, game.pieces, BLACK, 0)

        expect(t, p.isKingSideCastleAllowed(), false)
        expect(t, p.isQueenSideCastleAllowed(), false)

        game = NewGame().Setup(`Ke4`, `Ke7,Ra8,Rh8`)
        p = NewPosition(game, game.pieces, BLACK, 0)

        expect(t, p.isKingSideCastleAllowed(), false)
        expect(t, p.isQueenSideCastleAllowed(), false)
}

func TestPosition070(t *testing.T) { // Missing rooks.
        game := NewGame().Setup(`Ke1`, `Ke8`)
        p := game.Start()

        expect(t, p.isKingSideCastleAllowed(), false)
        expect(t, p.isQueenSideCastleAllowed(), false)

        game = NewGame().Setup(`Ke1`, `Ke8`)
        p = NewPosition(game, game.pieces, BLACK, 0)

        expect(t, p.isKingSideCastleAllowed(), false)
        expect(t, p.isQueenSideCastleAllowed(), false)
}
