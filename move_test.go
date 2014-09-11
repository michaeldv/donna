// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`github.com/michaeldv/donna/expect`; `testing`)

// PxQ, NxQ, BxQ, RxQ, QxQ, KxQ
func TestMove000(t *testing.T) {
	game := NewGame(`Kd6,Qd1,Ra5,Nc3,Bc4,e4`, `Kh8,Qd5`)
	p := game.Start(White)
	expect.Eq(t, NewMove(p, E4, D5).value(), 1258) // PxQ
	expect.Eq(t, NewMove(p, C3, D5).value(), 1256) // NxQ
	expect.Eq(t, NewMove(p, C4, D5).value(), 1254) // BxQ
	expect.Eq(t, NewMove(p, A5, D5).value(), 1252) // RxQ
	expect.Eq(t, NewMove(p, D1, D5).value(), 1250) // QxQ
	expect.Eq(t, NewMove(p, D6, D5).value(), 1248) // KxQ
}

// PxR, NxR, BxR, RxR, QxR, KxR
func TestMove010(t *testing.T) {
	game := NewGame(`Kd6,Qd1,Ra5,Nc3,Bc4,e4`, `Kh8,Rd5`)
	p := game.Start(White)
	expect.Eq(t, NewMove(p, E4, D5).value(), 633) // PxR
	expect.Eq(t, NewMove(p, C3, D5).value(), 631) // NxR
	expect.Eq(t, NewMove(p, C4, D5).value(), 629) // BxR
	expect.Eq(t, NewMove(p, A5, D5).value(), 627) // RxR
	expect.Eq(t, NewMove(p, D1, D5).value(), 625) // QxR
	expect.Eq(t, NewMove(p, D6, D5).value(), 623) // KxR
}

// PxB, NxB, BxB, RxB, QxB, KxB
func TestMove020(t *testing.T) {
	game := NewGame(`Kd6,Qd1,Ra5,Nc3,Bc4,e4`, `Kh8,Bd5`)
	p := game.Start(White)
	expect.Eq(t, NewMove(p, E4, D5).value(), 416) // PxB
	expect.Eq(t, NewMove(p, C3, D5).value(), 414) // NxB
	expect.Eq(t, NewMove(p, C4, D5).value(), 412) // BxB
	expect.Eq(t, NewMove(p, A5, D5).value(), 410) // RxB
	expect.Eq(t, NewMove(p, D1, D5).value(), 408) // QxB
	expect.Eq(t, NewMove(p, D6, D5).value(), 406) // KxB
}

// PxN, NxN, BxN, RxN, QxN, KxN
func TestMove030(t *testing.T) {
	game := NewGame(`Kd6,Qd1,Ra5,Nc3,Bc4,e4`, `Kh8,Nd5`)
	p := game.Start(White)
	expect.Eq(t, NewMove(p, E4, D5).value(), 406) // PxN
	expect.Eq(t, NewMove(p, C3, D5).value(), 404) // NxN
	expect.Eq(t, NewMove(p, C4, D5).value(), 402) // BxN
	expect.Eq(t, NewMove(p, A5, D5).value(), 400) // RxN
	expect.Eq(t, NewMove(p, D1, D5).value(), 398) // QxN
	expect.Eq(t, NewMove(p, D6, D5).value(), 396) // KxN
}

// PxP, NxP, BxP, RxP, QxP, KxP
func TestMove040(t *testing.T) {
	game := NewGame(`Kd6,Qd1,Ra5,Nc3,Bc4,e4`, `Kh8,d5`)
	p := game.Start(White)
	expect.Eq(t, NewMove(p, E4, D5).value(), 98) // PxP
	expect.Eq(t, NewMove(p, C3, D5).value(), 96) // NxP
	expect.Eq(t, NewMove(p, C4, D5).value(), 94) // BxP
	expect.Eq(t, NewMove(p, A5, D5).value(), 92) // RxP
	expect.Eq(t, NewMove(p, D1, D5).value(), 90) // QxP
	expect.Eq(t, NewMove(p, D6, D5).value(), 88) // KxP
}

// NewMoveFromString: move from algebraic notation.
func TestMove100(t *testing.T) {
	p := NewGame().Start()
	m1 := NewMove(p, E2, E4)
	m2 := NewMove(p, G1, F3)

	expect.Eq(t, NewMoveFromString(p, `e2e4`), m1)
	expect.Eq(t, NewMoveFromString(p, `e2-e4`), m1)
	expect.Eq(t, NewMoveFromString(p, `Ng1f3`), m2)
	expect.Eq(t, NewMoveFromString(p, `Ng1-f3`), m2)
	expect.Eq(t, NewMoveFromString(p, `Rg1-f3`), Move(0))
}

func TestMove110(t *testing.T) {
	p := NewGame(`Ke1,g7,a7`, `Ke8,Rh8,e2`).Start(White)
	m1 := NewMove(p, E1, E2) // Capture.
	m2 := NewMove(p, A7, A8).promote(Rook)  // Promo without capture.
	m3 := NewMove(p, G7, H8).promote(Queen) // Promo with capture.

	expect.Eq(t, NewMoveFromString(p, `Ke1e2`), m1)
	expect.Eq(t, NewMoveFromString(p, `Ke1xe2`), m1)
	expect.Eq(t, NewMoveFromString(p, `a7a8R`), m2)
	expect.Eq(t, NewMoveFromString(p, `a7-a8R`), m2)
	expect.Eq(t, NewMoveFromString(p, `g7h8Q`), m3)
	expect.Eq(t, NewMoveFromString(p, `g7xh8Q`), m3)
}

func TestMove120(t *testing.T) {
	p1 := NewGame(`Ke1`, `Ke8,Ra8`).Start(Black)
	m1 := NewCastle(p1, E8, C8)
	expect.Eq(t, NewMoveFromString(p1, `0-0-0`), m1)

	p2 := NewGame(`Ke1`, `Ke8,Rh8`).Start(Black)
	m2 := NewCastle(p2, E8, G8)
	expect.Eq(t, NewMoveFromString(p2, `0-0`), m2)
}

func TestMove130(t *testing.T) {
	p := NewGame().Start()
	p = p.MakeMove(NewPawnMove(p, E2, E4))
	p = p.MakeMove(NewPawnMove(p, E7, E6))
	p = p.MakeMove(NewPawnMove(p, E4, E5))
	move := NewPawnMove(p, D7, D5) // Causes en-passant on D6.

	expect.Eq(t, NewMoveFromString(p, `d7d5`), move)
	expect.True(t, NewMoveFromString(p, `d7-d5`).isEnpassant())
}

// Move to UCI coordinate notation.
func TestMove200(t *testing.T) {
	p := NewGame().Start()
	m1 := NewMove(p, E2, E4)
	m2 := NewMove(p, G1, F3)

	expect.Eq(t, m1.notation(), `e2e4`) // Pawn.
	expect.Eq(t, m2.notation(), `g1f3`) // Knight.
}

func TestMove210(t *testing.T) {
	p := NewGame(`Ke1,g7,a7`, `Ke8,Rh8,e2`).Start(White)
	m1 := NewMove(p, E1, E2) // Capture.
	m2 := NewMove(p, A7, A8).promote(Rook)  // Promo without capture.
	m3 := NewMove(p, G7, H8).promote(Queen) // Promo with capture.

	expect.Eq(t, m1.notation(), `e1e2`)
	expect.Eq(t, m2.notation(), `a7a8r`)
	expect.Eq(t, m3.notation(), `g7h8q`)
}

func TestMove220(t *testing.T) {
	p1 := NewGame(`Ke1`, `Ke8,Ra8`).Start(Black)
	m1 := NewCastle(p1, E8, C8) // 0-0-0
	expect.Eq(t, m1.notation(), `e8c8`)

	p2 := NewGame(`Ke1`, `Ke8,Rh8`).Start(Black)
	m2 := NewCastle(p2, E8, G8) // 0-0
	expect.Eq(t, m2.notation(), `e8g8`)
}

// Move from UCI coordinate notation.
func TestMove300(t *testing.T) {
	p := NewGame().Start()
	m1 := NewMove(p, E2, E4)
	m2 := NewMove(p, G1, F3)

	expect.Eq(t, NewMoveFromNotation(p, `e2e4`), m1) // Pawn.
	expect.Eq(t, NewMoveFromNotation(p, `g1f3`), m2) // Knight.
}

func TestMove310(t *testing.T) {
	p := NewGame(`Ke1,g7,a7`, `Ke8,Rh8,e2`).Start(White)
	m1 := NewMove(p, E1, E2) // Capture.
	m2 := NewMove(p, A7, A8).promote(Rook)  // Promo without capture.
	m3 := NewMove(p, G7, H8).promote(Queen) // Promo with capture.

	expect.Eq(t, NewMoveFromNotation(p, `e1e2`), m1)
	expect.Eq(t, NewMoveFromNotation(p, `a7a8r`), m2)
	expect.Eq(t, NewMoveFromNotation(p, `g7h8q`), m3)
}

func TestMove320(t *testing.T) {
	p1 := NewGame(`Ke1`, `Ke8,Ra8`).Start(Black)
	m1 := NewCastle(p1, E8, C8) // 0-0-0
	expect.Eq(t, NewMoveFromNotation(p1, `e8c8`), m1)

	p2 := NewGame(`Ke1`, `Ke8,Rh8`).Start(Black)
	m2 := NewCastle(p2, E8, G8) // 0-0
	expect.Eq(t, NewMoveFromNotation(p2, `e8g8`), m2)
}

func TestMove330(t *testing.T) {
	p := NewGame().Start()
	p = p.MakeMove(NewPawnMove(p, E2, E4))
	p = p.MakeMove(NewPawnMove(p, E7, E6))
	p = p.MakeMove(NewPawnMove(p, E4, E5))
	move := NewPawnMove(p, D7, D5) // Causes en-passant on D6.

	expect.Eq(t, NewMoveFromNotation(p, `d7d5`), move)
	expect.True(t, NewMoveFromNotation(p, `d7d5`).isEnpassant())
}

// Only pawns can do en-passant capture.
func TestMove340(t *testing.T) {
	p := NewGame(`Kg1,d2`, `Kc2,Qa3,Rh3,Be4,Nc1,c4`).Start(White)
	p = p.MakeMove(NewEnpassant(p, D2, D4)) // Causes en-passant on D3.
	bQ := NewMove(p, A3, D3)
	bR := NewMove(p, H3, D3)
	bB := NewMove(p, E4, D3)
	bN := NewMove(p, C1, D3)
	bK := NewMove(p, C2, D3)
	bP := NewMove(p, C4, D3)

	expect.Eq(t, bQ.capture(), Piece(0))
	expect.Eq(t, bR.capture(), Piece(0))
	expect.Eq(t, bB.capture(), Piece(0))
	expect.Eq(t, bN.capture(), Piece(0))
	expect.Eq(t, bK.capture(), Piece(0))
	expect.Eq(t, bP.capture(), Piece(Pawn))

	expect.Eq(t, bQ & isCapture, Move(0))
	expect.Eq(t, bR & isCapture, Move(0))
	expect.Eq(t, bB & isCapture, Move(0))
	expect.Eq(t, bN & isCapture, Move(0))
	expect.Eq(t, bK & isCapture, Move(0))
	expect.Ne(t, bP & isCapture, Move(0)) // Ne() for Pawn.
}
