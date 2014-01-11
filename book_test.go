// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`)

// See test key values at http://hardy.uhasselt.be/Toga/book_format.html
func TestBook000(t *testing.T) {
        position := NewGame().InitialPosition().Start(White)

        expect(t, position.polyglot(), uint64(0x463B96181691FC9C))
}

func TestBook010(t *testing.T) { // 1. e4
        position := NewGame().InitialPosition().Start(White)
        position = position.MakeMove(NewMove(position, E2, E4))

        expect(t, position.polyglot(), uint64(0x823C9B50FD114196))
}

func TestBook020(t *testing.T) { // 1. e4 d5
        position := NewGame().InitialPosition().Start(White)
        position = position.MakeMove(NewMove(position, E2, E4))
        position = position.MakeMove(NewMove(position, D7, D5))

        expect(t, position.polyglot(), uint64(0x0756B94461C50FB0))
}

func TestBook030(t *testing.T) { // 1. e4 d5 2. e5
        position := NewGame().InitialPosition().Start(White)
        position = position.MakeMove(NewMove(position, E2, E4))
        position = position.MakeMove(NewMove(position, D7, D5))
        position = position.MakeMove(NewMove(position, E4, E5))

        expect(t, position.polyglot(), uint64(0x662FAFB965DB29D4))
}

func TestBook040(t *testing.T) { // 1. e4 d5 2. e5 f5 <-- Enpassant
        position := NewGame().InitialPosition().Start(White)
        position = position.MakeMove(NewMove(position, E2, E4))
        position = position.MakeMove(NewMove(position, D7, D5))
        position = position.MakeMove(NewMove(position, E4, E5))
        position = position.MakeMove(NewMove(position, F7, F5))

        expect(t, position.polyglot(), uint64(0x22A48B5A8E47FF78))
        expect(t, position.flags.enpassant, F6)
}

func TestBook050(t *testing.T) { // TODO: 1. e4 d5 2. e5 f5 3. Ke2 <-- White Castle
        position := NewGame().InitialPosition().Start(White)
        position = position.MakeMove(NewMove(position, E2, E4))
        position = position.MakeMove(NewMove(position, D7, D5))
        position = position.MakeMove(NewMove(position, E4, E5))
        position = position.MakeMove(NewMove(position, F7, F5))
        position = position.MakeMove(NewMove(position, E1, E2))

        expect(t, position.polyglot(), uint64(0x652A607CA3F242C1))
}

func TestBook060(t *testing.T) { // TODO: 1. e4 d5 2. e5 f5 3. Ke2 Kf7 <-- Black Castle
        position := NewGame().InitialPosition().Start(White)
        position = position.MakeMove(NewMove(position, E2, E4))
        position = position.MakeMove(NewMove(position, D7, D5))
        position = position.MakeMove(NewMove(position, E4, E5))
        position = position.MakeMove(NewMove(position, F7, F5))
        position = position.MakeMove(NewMove(position, E1, E2))
        position = position.MakeMove(NewMove(position, E8, F7))

        expect(t, position.polyglot(), uint64(0x00FDD303C946BDD9))
}

func TestBook070(t *testing.T) { // 1. a2a4 b7b5 2. h2h4 b5b4 3. c2c4 <-- Enpassant
        position := NewGame().InitialPosition().Start(White)
        position = position.MakeMove(NewMove(position, A2, A4))
        position = position.MakeMove(NewMove(position, B7, B5))
        position = position.MakeMove(NewMove(position, H2, H4))
        position = position.MakeMove(NewMove(position, B5, B4))
        position = position.MakeMove(NewMove(position, C2, C4))

        expect(t, position.polyglot(), uint64(0x3C8123EA7B067637))
        expect(t, position.flags.enpassant, C3)
}

func TestBook080(t *testing.T) { // TODO: 1. a2a4 b7b5 2. h2h4 b5b4 3. c2c4 b4xc3 4. Ra1a3 <-- Enpassant/Castle
        position := NewGame().InitialPosition().Start(White)
        position = position.MakeMove(NewMove(position, A2, A4))
        position = position.MakeMove(NewMove(position, B7, B5))
        position = position.MakeMove(NewMove(position, H2, H4))
        position = position.MakeMove(NewMove(position, B5, B4))
        position = position.MakeMove(NewMove(position, C2, C4))
        position = position.MakeMove(NewMove(position, B4, C3))
        position = position.MakeMove(NewMove(position, A1, A3))

        expect(t, position.polyglot(), uint64(0x5C3F9B829B279560))
}


func TestBook100(t *testing.T) { // 1. e4 e5
        position := NewGame().InitialPosition().Start(White)
        position = position.MakeMove(NewMove(position, E2, E4))
        position = position.MakeMove(NewMove(position, E7, E5))

        expect(t, position.polyglot(), uint64(0x0844931A6EF4B9A0))
}

func TestBook110(t *testing.T) { // 1. d4 d5
        position := NewGame().InitialPosition().Start(White)
        position = position.MakeMove(NewMove(position, D2, D4))
        position = position.MakeMove(NewMove(position, D7, D5))

        expect(t, position.polyglot(), uint64(0x06649BA69B8C9FF8))
}
