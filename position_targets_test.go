// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

import(`github.com/michaeldv/donna/expect`; `testing`)

// Pawn targets.
func TestTargets000(t *testing.T) {
	game := NewGame(`Kd1,e2`, `Ke8,d4`)
	position := game.start()

	expect.Eq(t, position.targets(E2), bit[E3]|bit[E4]) // e3,e4
	expect.Eq(t, position.targets(D4), bit[D3])         // d3
}

func TestTargets010(t *testing.T) {
	game := NewGame(`Kd1,e2,d3`, `Ke8,d4,e4`)
	position := game.start()

	expect.Eq(t, position.targets(E2), bit[E3])         // e3
	expect.Eq(t, position.targets(D3), bit[E4])         // e4
	expect.Eq(t, position.targets(D4), Bitmask(0))      // None.
	expect.Eq(t, position.targets(E4), bit[D3]|bit[E3]) // d3,e3
}

func TestTargets020(t *testing.T) {
	game := NewGame(`Kd1,e2`, `Ke8,d3,f3`)
	position := game.start()

	expect.Eq(t, position.targets(E2), bit[D3]|bit[E3]|bit[E4]|bit[F3]) // d3,e3,e4,f3
	expect.Eq(t, position.targets(D3), bit[E2]|bit[D2])                 // e2,d2
	expect.Eq(t, position.targets(F3), bit[E2]|bit[F2])                 // e2,f2
}

func TestTargets030(t *testing.T) {
	game := NewGame(`Kd1,e2`, `Ke8,d4`)
	position := game.start()
	position = position.makeMove(NewEnpassant(position, E2, E4)) // Creates en-passant on e3.

	expect.Eq(t, position.targets(E4), bit[E5])         // e5
	expect.Eq(t, position.targets(D4), bit[D3]|bit[E3]) // d3, e3 (en-passant).
}

// Pawn attacks.
func TestTargets040(t *testing.T) {
	game := NewGame(`Ke1,a3,b3,c3,d3,e3,f3,g3,h3`, `Ke8,a6,b6,c6,d6,e6,f6,g6,h6`)
	position := game.start()
	expect.Eq(t, position.pawnAttacks(White), bit[A4]|bit[B4]|bit[C4]|bit[D4]|bit[E4]|bit[F4]|bit[G4]|bit[H4])
	expect.Eq(t, position.pawnAttacks(Black), bit[A5]|bit[B5]|bit[C5]|bit[D5]|bit[E5]|bit[F5]|bit[G5]|bit[H5])
}
