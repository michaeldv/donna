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

// Straight repetition.
func TestPosition100(t *testing.T) {
        p := NewGame().InitialPosition().Start() // Initial 1.
        p = p.MakeMove(NewMove(p, G1, F3));  p = p.MakeMove(NewMove(p, G8, F6)) // 1.
        p = p.MakeMove(NewMove(p, F3, G1));  p = p.MakeMove(NewMove(p, F6, G8)) // Initial 2.
        p = p.MakeMove(NewMove(p, G1, F3));  p = p.MakeMove(NewMove(p, G8, F6)) // 2.
        p = p.MakeMove(NewMove(p, F3, G1));  p = p.MakeMove(NewMove(p, F6, G8)) // Initial 3.
        p = p.MakeMove(NewMove(p, G1, F3));  p = p.MakeMove(NewMove(p, G8, F6)) // 3.

        expect(t, p.isRepetition(), true)
}

// Repetition with some moves in between.
func TestPosition110(t *testing.T) {
        p := NewGame().InitialPosition().Start()
        p = p.MakeMove(NewMove(p, E2, E4));  p = p.MakeMove(NewMove(p, E7, E5))

        p = p.MakeMove(NewMove(p, G1, F3));  p = p.MakeMove(NewMove(p, G8, F6)) // 1.
        p = p.MakeMove(NewMove(p, B1, C3));  p = p.MakeMove(NewMove(p, B8, C6))
        p = p.MakeMove(NewMove(p, F1, C4));  p = p.MakeMove(NewMove(p, F8, C5))
        p = p.MakeMove(NewMove(p, C3, B1));  p = p.MakeMove(NewMove(p, C6, B8))
        p = p.MakeMove(NewMove(p, C4, F1));  p = p.MakeMove(NewMove(p, C5, F8)) // 2.

        p = p.MakeMove(NewMove(p, F1, C4));  p = p.MakeMove(NewMove(p, F8, C5))
        p = p.MakeMove(NewMove(p, B1, C3));  p = p.MakeMove(NewMove(p, B8, C6))
        p = p.MakeMove(NewMove(p, C4, F1));  p = p.MakeMove(NewMove(p, C5, F8))
        p = p.MakeMove(NewMove(p, C3, B1));  p = p.MakeMove(NewMove(p, C6, B8)) // 3.

        expect(t, p.isRepetition(), true)
}

// Irreversible 0-0.
func TestPosition120(t *testing.T) {
        p := NewGame().Setup(`Ke1,Rh1,h2`, `Ke8,Ra8,a7`).Start()
        p = p.MakeMove(NewMove(p, H2, H4));  p = p.MakeMove(NewMove(p, A7, A5)) // 1.
        p = p.MakeMove(NewMove(p, E1, E2));  p = p.MakeMove(NewMove(p, E8, E7)) // King has moved.
        p = p.MakeMove(NewMove(p, E2, E1));  p = p.MakeMove(NewMove(p, E7, E8)) // 2.
        p = p.MakeMove(NewMove(p, E1, E2));  p = p.MakeMove(NewMove(p, E8, E7)) // King has moved again.
        p = p.MakeMove(NewMove(p, E2, E1));  p = p.MakeMove(NewMove(p, E7, E8)) // 3.
        expect(t, p.isRepetition(), false) // <-- Lost 0-0 right.

        p = p.MakeMove(NewMove(p, E1, E2));  p = p.MakeMove(NewMove(p, E8, E7)) // King has moved again.
        p = p.MakeMove(NewMove(p, E2, E1));  p = p.MakeMove(NewMove(p, E7, E8)) // 4.
        expect(t, p.isRepetition(), true) // <-- 3 time repetioion with lost 0-0 right.
}

// Undo the move.
// func TestPosition200(t *testing.T) {
//         Settings.Log = true
//         p := NewGame().InitialPosition().Start()
//         move :=  NewMove(p, E2, E4)
//         p.MakeMove(move)
//         Log("MOVE %s\n%s", move, p)
//         p.takeBack(move)
//         Log("TAKE BACK %s\n%s", move, p)
//         Settings.Log = false
// }
// 
// func TestPosition210(t *testing.T) {
//         Settings.Log = true
//         p := NewGame().Setup(`Ke1,Rh1`, `Ke8,Re7`).Start()
//         move :=  NewMove(p, H1, H2)
//         pos := p.MakeMove(move)
//         Log("MOVE %s\n%s", move, pos)
//         // p.takeBack(move)
//         // Log("TAKE BACK %s\n%s", move, p)
//         Settings.Log = false
// }

// func TestPosition230(t *testing.T) {
//         Settings.Log, Settings.Fancy = true, true
// 
//         p := NewGame().InitialPosition().Start()
//         p.MakeMove(NewMove(p, E2, E4));
//         p.MakeMove(NewMove(p, E7, E6))
//         p.MakeMove(NewMove(p, E4, E5));
// 
//         ep := NewMove(p, F7, F5)
//         p.MakeMove(ep)
//         Log("EP1 %s, ep %d\n%s\n\n", ep, p.enpassant, p)
// 
//         ep = NewMove(p, E5, F6)
//         p.MakeMove(ep)
//         Log("EP2 %s, ep %d\n%s\n\n", ep, p.enpassant, p)
// 
//         Log("Taking back %s, move.enpassant %v\n", ep, ep.enpassant)
//         p.takeBack(ep)
//         Log("EP3 %s, ep %d\n%s\n\n", ep, p.enpassant, p)
// 
//         Settings.Log, Settings.Fancy = false, false
// }

func TestPosition240(t *testing.T) {
        Settings.Log, Settings.Fancy = true, true

        p := NewGame().Setup(`Ke1,h7`, `Ke8,a2`).Start()
        Log("Initial p.inCheck = %v, p.stage = %d\n%s\n\n", p.inCheck, p.stage, p)
        Log("%v", p.count)

        game      := p.game
        pieces    := p.pieces
        targets   := p.targets
        board     := p.board
        attacks   := p.attacks
        outposts  := p.outposts
        count     := p.count
        enpassant := p.enpassant
        color     := p.color
        stage     := p.stage
        history   := p.history
        inCheck   := p.inCheck
        can00     := p.can00
        can000    := p.can000

        move := NewMove(p, H7, H8).promote(QUEEN)
        p.MakeMove(move);
        Log("Move %s p.inCheck = %v, p.stage = %d\n%s\n\n", move, p.inCheck, p.stage, p)

        Log("Taking back %s\n", move)
        p.takeBack(move)
        Log("Move %s p.inCheck = %v, p.stage = %d\n%s\n\n", move, p.inCheck, p.stage, p)

        Lop("game      : ", game      == p.game)        //t
        Lop("pieces    : ", pieces    == p.pieces)
        Lop("targets   : ", targets   == p.targets)
        Lop("board     : ", board     == p.board)
        Lop("attacks   : ", attacks   == p.attacks)
        Lop("outposts  : ", outposts  == p.outposts)
        Lop("count     : ", count     == p.count)
        Lop("enpassant : ", enpassant == p.enpassant)   //t
        Lop("color     : ", color     == p.color)
        Lop("stage     : ", stage     == p.stage)
        Lop("history   : ", history   == p.history)     //t
        Lop("inCheck   : ", inCheck   == p.inCheck)
        Lop("can00     : ", can00     == p.can00)       //t
        Lop("can000    : ", can000    == p.can000)      //t
        Log("%v", p.count)

        Settings.Log, Settings.Fancy = false, false
}

// PAWN   = 1 << 1 // 2
// KNIGHT = 2 << 1 // 4
// BISHOP = 3 << 1 // 6
// ROOK   = 4 << 1 // 8
// QUEEN  = 5 << 1 // 10
// KING   = 6 << 1 // 12
// 
// [0 0  1 1 0 0 0 0 0 0 0 0 1 1 0 0]
// [0 0 -2 1 0 0 0 0 0 0 1 0 1 1 0 0]
// - 2 white pawns and 1 queen
