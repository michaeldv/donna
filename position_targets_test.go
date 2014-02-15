// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna
import (`testing`)

// Pawn targets.
func TestTargets000(t *testing.T) {
        game := NewGame().Setup(`Kd1,e2`, `Ke8,d4`)
        position := game.Start(White)

        expect(t, position.targetsMask(E2), Bit(E3) | Bit(E4)) // e3,e4
        expect(t, position.targetsMask(D4), Bit(D3))           // d3
}

func TestTargets010(t *testing.T) {
        game := NewGame().Setup(`Kd1,e2,d3`, `Ke8,d4,e4`)
        position := game.Start(White)

        expect(t, position.targetsMask(E2), Bit(E3))           // e3
        expect(t, position.targetsMask(D3), Bit(E4))           // e4
        expect(t, position.targetsMask(D4), Bitmask(0))        // None.
        expect(t, position.targetsMask(E4), Bit(D3) | Bit(E3)) // d3,e3
}

func TestTargets020(t *testing.T) {
        game := NewGame().Setup(`Kd1,e2`, `Ke8,d3,f3`)
        position := game.Start(White)

        expect(t, position.targetsMask(E2), Bit(D3) | Bit(E3) | Bit(E4) | Bit(F3)) // d3,e3,e4,f3
        expect(t, position.targetsMask(D3), Bit(E2) | Bit(D2)) // e2,d2
        expect(t, position.targetsMask(F3), Bit(E2) | Bit(F2)) // e2,f2
}

func TestTargets030(t *testing.T) {
        game := NewGame().Setup(`Kd1,e2`, `Ke8,d4`)
        position := game.Start(White)
        position = position.MakeMove(position.NewEnpassant(E2, E4)) // Creates en-passant on e3.

        expect(t, position.targetsMask(E4), Bit(E5))           // e5
        expect(t, position.targetsMask(D4), Bit(D3) | Bit(E3)) // d3, e3 (en-passant).
}
