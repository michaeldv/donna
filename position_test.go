// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`)

func TestPosition010(t *testing.T) {
        p := NewGame().Setup(`Ke1,e2`, `Kg8,d7,f7`).Start(White)
        expect(t, p.enpassant, 0)

        p = p.MakeMove(p.NewMove(E2, E4))
        expect(t, p.enpassant, 0)

        p = p.MakeMove(p.NewMove(D7, D5))
        expect(t, p.enpassant, 0)

        p = p.MakeMove(p.NewMove(E4, E5))
        expect(t, p.enpassant, 0)

        p = p.MakeMove(p.NewEnpassant(F7, F5))
        expect(t, p.enpassant, F6)
}

// Castle tests.
func TestPosition030(t *testing.T) { // Everything is OK.
        p := NewGame().Setup(`Ke1,Ra1,Rh1`, `Ke8`).Start(White)
        kingside, queenside := p.canCastle(p.color)
        expect(t, kingside, true)
        expect(t, queenside, true)

        p = NewGame().Setup(`Ke1`, `Ke8,Ra8,Rh8`).Start(Black)
        kingside, queenside = p.canCastle(p.color)
        expect(t, kingside, true)
        expect(t, queenside, true)
}


func TestPosition040(t *testing.T) { // King checked.
        p := NewGame().Setup(`Ke1,Ra1,Rh1`, `Ke8,Bg3`).Start(White)
        kingside, queenside := p.canCastle(p.color)
        expect(t, kingside, false)
        expect(t, queenside, false)

        p = NewGame().Setup(`Ke1,Bg6`, `Ke8,Ra8,Rh8`).Start(Black)
        kingside, queenside = p.canCastle(p.color)
        expect(t, kingside, false)
        expect(t, queenside, false)
}

func TestPosition050(t *testing.T) { // Attacked square.
        p := NewGame().Setup(`Ke1,Ra1,Rh1`, `Ke8,Bb3,Bh3`).Start(White)
        kingside, queenside := p.canCastle(p.color)
        expect(t, kingside, false)
        expect(t, queenside, false)

        p = NewGame().Setup(`Ke1,Bb6,Bh6`, `Ke8,Ra8,Rh8`).Start(Black)
        kingside, queenside = p.canCastle(p.color)
        expect(t, kingside, false)
        expect(t, queenside, false)
}

func TestPosition060(t *testing.T) { // Wrong square.
        p := NewGame().Setup(`Ke1,Ra8,Rh8`, `Ke5`).Start(White)
        kingside, queenside := p.canCastle(p.color)
        expect(t, kingside, false)
        expect(t, queenside, false)

        p = NewGame().Setup(`Ke2,Ra1,Rh1`, `Ke8`).Start(White)
        kingside, queenside = p.canCastle(p.color)
        expect(t, kingside, false)
        expect(t, queenside, false)

        p = NewGame().Setup(`Ke4`, `Ke8,Ra1,Rh1`).Start(Black)
        kingside, queenside = p.canCastle(p.color)
        expect(t, kingside, false)
        expect(t, queenside, false)

        p = NewGame().Setup(`Ke4`, `Ke7,Ra8,Rh8`).Start(Black)
        kingside, queenside = p.canCastle(p.color)
        expect(t, kingside, false)
        expect(t, queenside, false)
}

func TestPosition070(t *testing.T) { // Missing rooks.
        p := NewGame().Setup(`Ke1`, `Ke8`).Start(White)
        kingside, queenside := p.canCastle(p.color)
        expect(t, kingside, false)
        expect(t, queenside, false)

        p = NewGame().Setup(`Ke1`, `Ke8`).Start(Black)
        kingside, queenside = p.canCastle(p.color)
        expect(t, kingside, false)
        expect(t, queenside, false)
}

func TestPosition080(t *testing.T) { // Rooks on wrong squares.
        p := NewGame().Setup(`Ke1,Rb1`, `Ke8`).Start(White)
        kingside, queenside := p.canCastle(p.color)
        expect(t, kingside, false)
        expect(t, queenside, false)

        p = NewGame().Setup(`Ke1,Rb1,Rh1`, `Ke8`).Start(White)
        kingside, queenside = p.canCastle(p.color)
        expect(t, kingside, true)
        expect(t, queenside, false)

        p = NewGame().Setup(`Ke1,Ra1,Rf1`, `Ke8`).Start(White)
        kingside, queenside = p.canCastle(p.color)
        expect(t, kingside, false)
        expect(t, queenside, true)
}

func TestPosition081(t *testing.T) { // Rook has moved.
        p := NewGame().Setup(`Ke1,Ra1,Rh1`, `Ke8`).Start(White)
        p = p.MakeMove(p.NewMove(A1, A2))
        p = p.MakeMove(p.NewMove(E8, E7))
        p = p.MakeMove(p.NewMove(A2, A1))

        kingside, queenside := p.canCastle(White)
        expect(t, kingside, true)
        expect(t, queenside, false)
}

func TestPosition082(t *testing.T) { // King has moved.
        p := NewGame().Setup(`Ke1`, `Ke8,Ra8,Rh8`).Start(Black)
        p = p.MakeMove(p.NewMove(E8, E7))
        p = p.MakeMove(p.NewMove(E1, E2))
        p = p.MakeMove(p.NewMove(E7, E8))

        kingside, queenside := p.canCastle(Black)
        expect(t, kingside, false)
        expect(t, queenside, false)
}

func TestPosition083(t *testing.T) { // Rook is taken.
        p := NewGame().Setup(`Ke1,Nb6`, `Ke8,Ra8,Rh8`).Start(White)
        p = p.MakeMove(p.NewMove(B6, A8))

        kingside, queenside := p.canCastle(Black)
        expect(t, kingside, true)
        expect(t, queenside, false)
}

// Straight repetition.
func TestPosition100(t *testing.T) {
        p := NewGame().InitialPosition().Start(White) // Initial 1.
        p = p.MakeMove(p.NewMove(G1, F3));  p = p.MakeMove(p.NewMove(G8, F6)) // 1.
        expect(t, p.isRepetition(), false)
        p = p.MakeMove(p.NewMove(F3, G1));  p = p.MakeMove(p.NewMove(F6, G8)) // Initial 2.
        expect(t, p.isRepetition(), false)
        p = p.MakeMove(p.NewMove(G1, F3));  p = p.MakeMove(p.NewMove(G8, F6)) // 2.
        expect(t, p.isRepetition(), false)
        p = p.MakeMove(p.NewMove(F3, G1));  p = p.MakeMove(p.NewMove(F6, G8)) // Initial 3.
        expect(t, p.isRepetition(), true)
        p = p.MakeMove(p.NewMove(G1, F3));  p = p.MakeMove(p.NewMove(G8, F6)) // 3.
        expect(t, p.isRepetition(), true)
}

// Repetition with some moves in between.
func TestPosition110(t *testing.T) {
        p := NewGame().InitialPosition().Start(White)
        p = p.MakeMove(p.NewMove(E2, E4));  p = p.MakeMove(p.NewMove(E7, E5))

        p = p.MakeMove(p.NewMove(G1, F3));  p = p.MakeMove(p.NewMove(G8, F6)) // 1.
        p = p.MakeMove(p.NewMove(B1, C3));  p = p.MakeMove(p.NewMove(B8, C6))
        p = p.MakeMove(p.NewMove(F1, C4));  p = p.MakeMove(p.NewMove(F8, C5))
        p = p.MakeMove(p.NewMove(C3, B1));  p = p.MakeMove(p.NewMove(C6, B8))
        p = p.MakeMove(p.NewMove(C4, F1));  p = p.MakeMove(p.NewMove(C5, F8)) // 2.

        expect(t, p.isRepetition(), false)

        p = p.MakeMove(p.NewMove(F1, C4));  p = p.MakeMove(p.NewMove(F8, C5))
        p = p.MakeMove(p.NewMove(B1, C3));  p = p.MakeMove(p.NewMove(B8, C6))
        p = p.MakeMove(p.NewMove(C4, F1));  p = p.MakeMove(p.NewMove(C5, F8))
        p = p.MakeMove(p.NewMove(C3, B1));  p = p.MakeMove(p.NewMove(C6, B8)) // 3.

        expect(t, p.isRepetition(), true)
}

// Irreversible 0-0.
func TestPosition120(t *testing.T) {
        p := NewGame().Setup(`Ke1,Rh1,h2`, `Ke8,Ra8,a7`).Start(White)
        p = p.MakeMove(p.NewMove(H2, H4));  p = p.MakeMove(p.NewMove(A7, A5)) // 1.
        p = p.MakeMove(p.NewMove(E1, E2));  p = p.MakeMove(p.NewMove(E8, E7)) // King has moved.
        p = p.MakeMove(p.NewMove(E2, E1));  p = p.MakeMove(p.NewMove(E7, E8)) // 2.
        p = p.MakeMove(p.NewMove(E1, E2));  p = p.MakeMove(p.NewMove(E8, E7)) // King has moved again.
        p = p.MakeMove(p.NewMove(E2, E1));  p = p.MakeMove(p.NewMove(E7, E8)) // 3.
        expect(t, p.isRepetition(), false) // <-- Lost 0-0 right.

        p = p.MakeMove(p.NewMove(E1, E2));  p = p.MakeMove(p.NewMove(E8, E7)) // King has moved again.
        p = p.MakeMove(p.NewMove(E2, E1));  p = p.MakeMove(p.NewMove(E7, E8)) // 4.
        expect(t, p.isRepetition(), true) // <-- 3 time repetioion with lost 0-0 right.
}

// Incremental hash recalculation tests (see book_test.go).
func TestPosition200(t *testing.T) { // 1. e4
        p := NewGame().InitialPosition().Start(White)
        p = p.MakeMove(p.NewMove(E2, E4))

        expect(t, p.polyglot(), uint64(0x823C9B50FD114196))
        expect(t, p.enpassant, 0)
        expect(t, p.castles, uint8(0x0F))
}

func TestPosition210(t *testing.T) { // 1. e4 d5
        p := NewGame().InitialPosition().Start(White)
        p = p.MakeMove(p.NewMove(E2, E4))
        p = p.MakeMove(p.NewMove(D7, D5))

        expect(t, p.polyglot(), uint64(0x0756B94461C50FB0))
        expect(t, p.enpassant, 0)
        expect(t, p.castles, uint8(0x0F))
}

func TestPosition220(t *testing.T) { // 1. e4 d5 2. e5
        p := NewGame().InitialPosition().Start(White)
        p = p.MakeMove(p.NewMove(E2, E4))
        p = p.MakeMove(p.NewMove(D7, D5))
        p = p.MakeMove(p.NewMove(E4, E5))

        expect(t, p.polyglot(), uint64(0x662FAFB965DB29D4))
        expect(t, p.enpassant, 0)
        expect(t, p.castles, uint8(0x0F))
}

func TestPosition230(t *testing.T) { // 1. e4 d5 2. e5 f5 <-- Enpassant
        p := NewGame().InitialPosition().Start(White)
        p = p.MakeMove(p.NewMove(E2, E4))
        p = p.MakeMove(p.NewMove(D7, D5))
        p = p.MakeMove(p.NewMove(E4, E5))
        p = p.MakeMove(p.NewEnpassant(F7, F5))

        expect(t, p.polyglot(), uint64(0x22A48B5A8E47FF78))
        expect(t, p.enpassant, F6)
        expect(t, p.castles, uint8(0x0F))
}

func TestPosition240(t *testing.T) { // 1. e4 d5 2. e5 f5 3. Ke2 <-- White Castle
        p := NewGame().InitialPosition().Start(White)
        p = p.MakeMove(p.NewMove(E2, E4))
        p = p.MakeMove(p.NewMove(D7, D5))
        p = p.MakeMove(p.NewMove(E4, E5))
        p = p.MakeMove(p.NewMove(F7, F5))
        p = p.MakeMove(p.NewMove(E1, E2))

        expect(t, p.polyglot(), uint64(0x652A607CA3F242C1))
        expect(t, p.enpassant, 0)
        expect(t, p.castles, castleKingside[Black] | castleQueenside[Black])
}

func TestPosition250(t *testing.T) { // 1. e4 d5 2. e5 f5 3. Ke2 Kf7 <-- Black Castle
        p := NewGame().InitialPosition().Start(White)
        p = p.MakeMove(p.NewMove(E2, E4))
        p = p.MakeMove(p.NewMove(D7, D5))
        p = p.MakeMove(p.NewMove(E4, E5))
        p = p.MakeMove(p.NewMove(F7, F5))
        p = p.MakeMove(p.NewMove(E1, E2))
        p = p.MakeMove(p.NewMove(E8, F7))

        expect(t, p.polyglot(), uint64(0x00FDD303C946BDD9))
        expect(t, p.enpassant, 0)
        expect(t, p.castles, uint8(0))
}

func TestPosition260(t *testing.T) { // 1. a2a4 b7b5 2. h2h4 b5b4 3. c2c4 <-- Enpassant
        p := NewGame().InitialPosition().Start(White)
        p = p.MakeMove(p.NewMove(A2, A4))
        p = p.MakeMove(p.NewMove(B7, B5))
        p = p.MakeMove(p.NewMove(H2, H4))
        p = p.MakeMove(p.NewMove(B5, B4))
        p = p.MakeMove(p.NewEnpassant(C2, C4))

        expect(t, p.polyglot(), uint64(0x3C8123EA7B067637))
        expect(t, p.enpassant, C3)
        expect(t, p.castles, uint8(0x0F))
}

func TestPosition270(t *testing.T) { // 1. a2a4 b7b5 2. h2h4 b5b4 3. c2c4 b4xc3 4. Ra1a3 <-- Enpassant/Castle
        p := NewGame().InitialPosition().Start(White)
        p = p.MakeMove(p.NewMove(A2, A4))
        p = p.MakeMove(p.NewMove(B7, B5))
        p = p.MakeMove(p.NewMove(H2, H4))
        p = p.MakeMove(p.NewMove(B5, B4))
        p = p.MakeMove(p.NewEnpassant(C2, C4))
        p = p.MakeMove(p.NewMove(B4, C3))
        p = p.MakeMove(p.NewMove(A1, A3))

        expect(t, p.polyglot(), uint64(0x5C3F9B829B279560))
        expect(t, p.enpassant, 0)
        expect(t, p.castles, castleKingside[White] | castleKingside[Black] | castleQueenside[Black])
}

// Position status.
func TestPosition300(t *testing.T) {
        p := NewGame().InitialPosition().Start(White)
        expect(t, p.status(p.NewMove(E2, E4), p.Evaluate()), InProgress)
}

// Mate in 1 move.
func TestPosition310(t *testing.T) {
        p := NewGame().Setup(`Kf8,Rh1,g6`, `Kh8,Bg8,g7,h7`).Start(White)
        rootNode = node // Reset Ply().
        expect(t, p.status(p.NewMove(H1, H6), Checkmate - Ply()), WhiteWinning)
}

// Forced stalemate.
func TestPosition320(t *testing.T) {
        p := NewGame().Setup(`Kf7,b2,b4,h6`, `Kh8,Ba4,b3,b5,h7`).Start(White)
        expect(t, p.status(p.NewMove(F7, F8), 0), Stalemate)
}

// Self-imposed stalemate.
func TestPosition330(t *testing.T) {
        p := NewGame().Setup(`Ka1,g3,h2`, `Kh5,h3,g4,g5,g6,h7`).Start(Black)
        p = p.MakeMove(p.NewMove(H7, H6))
        expect(t, p.status(p.NewMove(A1, B2), 0), Stalemate)
}

// Draw by repetition.
func TestPosition340(t *testing.T) {
        p := NewGame().Setup(`Ka1,g3,h2`, `Kh5,h3,g4,g5,g6,h7`).Start(Black) // Initial.

        p = p.MakeMove(p.NewMove(H5, H6)); p = p.MakeMove(p.NewMove(A1, A2))
        p = p.MakeMove(p.NewMove(H6, H5)); p = p.MakeMove(p.NewMove(A2, A1)) // Rep #2.
        expect(t, p.status(p.NewMove(H5, H6), 0), InProgress)

        p = p.MakeMove(p.NewMove(H5, H6)); p = p.MakeMove(p.NewMove(A1, A2))
        p = p.MakeMove(p.NewMove(H6, H5)) // -- No p.NewMove(A2, A1) here --

        rootNode = node // Reset Ply().
        expect(t, p.status(p.NewMove(A2, A1), 0), Repetition) // <-- Ka2-a1 causes rep #3.
}

// Position after null move.
func TestPosition350(t *testing.T) {
        p := NewGame().Setup(`Ke1,Qd1,d2,e2`, `Kg8,Qf8,f7,g7`).Start(White)

        p = p.MakeNullMove()
        expect(t, p.isNull(), true)

        p = p.TakeBackNullMove()
        p = p.MakeMove(p.NewMove(E2, E4))
        expect(t, p.isNull(), false)
}
