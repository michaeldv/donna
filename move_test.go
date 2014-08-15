// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
	`testing`
)

// PxQ, NxQ, BxQ, RxQ, QxQ, KxQ
func TestMove000(t *testing.T) {
	game := NewGame(`Kd6,Qd1,Ra5,Nc3,Bc4,e4`, `Kh8,Qd5`)
	p := game.Start(White)
	expect(t, NewMove(p, E4, D5).value(), 1258) // PxQ
	expect(t, NewMove(p, C3, D5).value(), 1256) // NxQ
	expect(t, NewMove(p, C4, D5).value(), 1254) // BxQ
	expect(t, NewMove(p, A5, D5).value(), 1252) // RxQ
	expect(t, NewMove(p, D1, D5).value(), 1250) // QxQ
	expect(t, NewMove(p, D6, D5).value(), 1248) // KxQ
}

// PxR, NxR, BxR, RxR, QxR, KxR
func TestMove010(t *testing.T) {
	game := NewGame(`Kd6,Qd1,Ra5,Nc3,Bc4,e4`, `Kh8,Rd5`)
	p := game.Start(White)
	expect(t, NewMove(p, E4, D5).value(), 633) // PxR
	expect(t, NewMove(p, C3, D5).value(), 631) // NxR
	expect(t, NewMove(p, C4, D5).value(), 629) // BxR
	expect(t, NewMove(p, A5, D5).value(), 627) // RxR
	expect(t, NewMove(p, D1, D5).value(), 625) // QxR
	expect(t, NewMove(p, D6, D5).value(), 623) // KxR
}

// PxB, NxB, BxB, RxB, QxB, KxB
func TestMove020(t *testing.T) {
	game := NewGame(`Kd6,Qd1,Ra5,Nc3,Bc4,e4`, `Kh8,Bd5`)
	p := game.Start(White)
	expect(t, NewMove(p, E4, D5).value(), 416) // PxB
	expect(t, NewMove(p, C3, D5).value(), 414) // NxB
	expect(t, NewMove(p, C4, D5).value(), 412) // BxB
	expect(t, NewMove(p, A5, D5).value(), 410) // RxB
	expect(t, NewMove(p, D1, D5).value(), 408) // QxB
	expect(t, NewMove(p, D6, D5).value(), 406) // KxB
}

// PxN, NxN, BxN, RxN, QxN, KxN
func TestMove030(t *testing.T) {
	game := NewGame(`Kd6,Qd1,Ra5,Nc3,Bc4,e4`, `Kh8,Nd5`)
	p := game.Start(White)
	expect(t, NewMove(p, E4, D5).value(), 406) // PxN
	expect(t, NewMove(p, C3, D5).value(), 404) // NxN
	expect(t, NewMove(p, C4, D5).value(), 402) // BxN
	expect(t, NewMove(p, A5, D5).value(), 400) // RxN
	expect(t, NewMove(p, D1, D5).value(), 398) // QxN
	expect(t, NewMove(p, D6, D5).value(), 396) // KxN
}

// PxP, NxP, BxP, RxP, QxP, KxP
func TestMove040(t *testing.T) {
	game := NewGame(`Kd6,Qd1,Ra5,Nc3,Bc4,e4`, `Kh8,d5`)
	p := game.Start(White)
	expect(t, NewMove(p, E4, D5).value(), 98) // PxP
	expect(t, NewMove(p, C3, D5).value(), 96) // NxP
	expect(t, NewMove(p, C4, D5).value(), 94) // BxP
	expect(t, NewMove(p, A5, D5).value(), 92) // RxP
	expect(t, NewMove(p, D1, D5).value(), 90) // QxP
	expect(t, NewMove(p, D6, D5).value(), 88) // KxP
}

// Move to UCI coordinate notation.
func TestMove200(t *testing.T) {
	p := NewGame().Start()
	m1 := NewMove(p, E2, E4)
	m2 := NewMove(p, G1, F3)

	expect(t, m1.notation(), `e2e4`) // Pawn.
	expect(t, m2.notation(), `g1f3`) // Knight.
}

func TestMove210(t *testing.T) {
	p := NewGame(`Ke1,g7,a7`, `Ke8,Rh8,e2`).Start(White)
	m1 := NewMove(p, E1, E2) // Capture.
	m2 := NewMove(p, A7, A8).promote(Rook)  // Promo without capture.
	m3 := NewMove(p, G7, H8).promote(Queen) // Promo with capture.

	expect(t, m1.notation(), `e1e2`)
	expect(t, m2.notation(), `a7a8r`)
	expect(t, m3.notation(), `g7h8q`)
}

func TestMove220(t *testing.T) {
	p1 := NewGame(`Ke1`, `Ke8,Ra8`).Start(Black)
	m1 := NewCastle(p1, E8, C8) // 0-0-0
	expect(t, m1.notation(), `e8c8`)

	p2 := NewGame(`Ke1`, `Ke8,Rh8`).Start(Black)
	m2 := NewCastle(p2, E8, G8) // 0-0
	expect(t, m2.notation(), `e8g8`)
}

// Move from UCI coordinate notation.
func TestMove300(t *testing.T) {
	p := NewGame().Start()
	m1 := NewMove(p, E2, E4)
	m2 := NewMove(p, G1, F3)

	expect(t, NewMoveFromNotation(p, `e2e4`), m1) // Pawn.
	expect(t, NewMoveFromNotation(p, `g1f3`), m2) // Knight.
}

func TestMove310(t *testing.T) {
	p := NewGame(`Ke1,g7,a7`, `Ke8,Rh8,e2`).Start(White)
	m1 := NewMove(p, E1, E2) // Capture.
	m2 := NewMove(p, A7, A8).promote(Rook)  // Promo without capture.
	m3 := NewMove(p, G7, H8).promote(Queen) // Promo with capture.

	expect(t, NewMoveFromNotation(p, `e1e2`), m1)
	expect(t, NewMoveFromNotation(p, `a7a8r`), m2)
	expect(t, NewMoveFromNotation(p, `g7h8q`), m3)
}

func TestMove320(t *testing.T) {
	p1 := NewGame(`Ke1`, `Ke8,Ra8`).Start(Black)
	m1 := NewCastle(p1, E8, C8) // 0-0-0
	expect(t, NewMoveFromNotation(p1, `e8c8`), m1)

	p2 := NewGame(`Ke1`, `Ke8,Rh8`).Start(Black)
	m2 := NewCastle(p2, E8, G8) // 0-0
	expect(t, NewMoveFromNotation(p2, `e8g8`), m2)
}

func TestMove330(t *testing.T) {
	p := NewGame().Start()
	p = p.MakeMove(NewPawnMove(p, E2, E4))
	p = p.MakeMove(NewPawnMove(p, E7, E6))
	p = p.MakeMove(NewPawnMove(p, E4, E5))
	move := NewPawnMove(p, D7, D5) // Causes en-passant on D6.

	expect(t, NewMoveFromNotation(p, `d7d5`), move)
	expect(t, NewMoveFromNotation(p, `d7d5`).isEnpassant(), true)
}
