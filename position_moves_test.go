// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna
import (`testing`)

// Unobstructed pins.
func TestMoves000(t *testing.T) {
        position := NewGame().Setup(`Ka1,Qe1,Ra8,Rh8,Bb5`, `Ke8,Re7,Bc8,Bf8,Nc6`).Start(White)
        pinned := position.pinnedMask(E8)

        expect(t, pinned, Bit(C6) | Bit(C8) | Bit(E7) | Bit(F8))
}

func TestMoves010(t *testing.T) {
        position := NewGame().Setup(`Ke4,Qe5,Rd5,Nd4,Nf4`, `Ka7,Qe8,Ra4,Rh4,Ba8`).Start(Black)
        pinned := position.pinnedMask(E4)

        expect(t, pinned, Bit(D5) | Bit(E5) | Bit(D4) | Bit(F4))
}

// Not a pin (friendly blockers).
func TestMoves020(t *testing.T) {
        position := NewGame().Setup(`Ka1,Qe1,Ra8,Rh8,Bb5,Nb8,Ng8,e4`, `Ke8,Re7,Bc8,Bf8,Nc6`).Start(White)
        pinned := position.pinnedMask(E8)

        expect(t, pinned, Bit(C6))
}

func TestMoves030(t *testing.T) {
        position := NewGame().Setup(`Ke4,Qe7,Rc6,Nb4,Ng4`, `Ka7,Qe8,Ra4,Rh4,Ba8,c4,e6,f4`).Start(Black)
        pinned := position.pinnedMask(E4)

        expect(t, pinned, Bit(C6))
}

// Not a pin (enemy blockers).
func TestMoves040(t *testing.T) {
        position := NewGame().Setup(`Ka1,Qe1,Ra8,Rh8,Bb5`, `Ke8,Re7,Rg8,Bc8,Bf8,Nc6,Nb8,e4`).Start(White)
        pinned := position.pinnedMask(E8)

        expect(t, pinned, Bit(C6))
}

func TestMoves050(t *testing.T) {
        position := NewGame().Setup(`Ke4,Qe7,Rc6,Nb4,Ng4,c4,e5,f4`, `Ka7,Qe8,Ra4,Rh4,Ba8`).Start(Black)
        pinned := position.pinnedMask(E4)

        expect(t, pinned, Bit(C6))
}
