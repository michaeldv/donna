// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna
import (`testing`)

// Pawn targets.
func TestTargets000(t *testing.T) {
        game := NewGame().Setup(`Kd1,e2`, `Ke8,d4`)
        position := game.Start(White)

        expect(t, position.targets(E2), bit[E3] | bit[E4]) // e3,e4
        expect(t, position.targets(D4), bit[D3])           // d3
}

func TestTargets010(t *testing.T) {
        game := NewGame().Setup(`Kd1,e2,d3`, `Ke8,d4,e4`)
        position := game.Start(White)

        expect(t, position.targets(E2), bit[E3])           // e3
        expect(t, position.targets(D3), bit[E4])           // e4
        expect(t, position.targets(D4), maskNone)          // None.
        expect(t, position.targets(E4), bit[D3] | bit[E3]) // d3,e3
}

func TestTargets020(t *testing.T) {
        game := NewGame().Setup(`Kd1,e2`, `Ke8,d3,f3`)
        position := game.Start(White)

        expect(t, position.targets(E2), bit[D3] | bit[E3] | bit[E4] | bit[F3]) // d3,e3,e4,f3
        expect(t, position.targets(D3), bit[E2] | bit[D2]) // e2,d2
        expect(t, position.targets(F3), bit[E2] | bit[F2]) // e2,f2
}

func TestTargets030(t *testing.T) {
        game := NewGame().Setup(`Kd1,e2`, `Ke8,d4`)
        position := game.Start(White)
        position = position.MakeMove(position.NewEnpassant(E2, E4)) // Creates en-passant on e3.

        expect(t, position.targets(E4), bit[E5])           // e5
        expect(t, position.targets(D4), bit[D3] | bit[E3]) // d3, e3 (en-passant).
}
