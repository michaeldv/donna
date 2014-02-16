// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`)

func TestPosition010(t *testing.T) {
        p := NewGame().Setup(`Ke1,e2`, `Kg8,d7,f7`).Start(White)
        expect(t, p.flags.enpassant, 0)

        p = p.MakeMove(p.NewMove(E2, E4))
        expect(t, p.flags.enpassant, 0)

        p = p.MakeMove(p.NewMove(D7, D5))
        expect(t, p.flags.enpassant, 0)

        p = p.MakeMove(p.NewMove(E4, E5))
        expect(t, p.flags.enpassant, 0)

        p = p.MakeMove(p.NewEnpassant(F7, F5))
        expect(t, p.flags.enpassant, F6)
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
        p = p.MakeMove(p.NewMove(F3, G1));  p = p.MakeMove(p.NewMove(F6, G8)) // Initial 2.
        p = p.MakeMove(p.NewMove(G1, F3));  p = p.MakeMove(p.NewMove(G8, F6)) // 2.
        p = p.MakeMove(p.NewMove(F3, G1));  p = p.MakeMove(p.NewMove(F6, G8)) // Initial 3.
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

