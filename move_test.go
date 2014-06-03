// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
	`testing`
)

// PxQ, NxQ, BxQ, RxQ, QxQ, KxQ
func TestMove000(t *testing.T) {
	game := NewGame().Setup(`Kd6,Qd1,Ra5,Nc3,Bc4,e4`, `Kh8,Qd5`)
	p := game.Start(White)
	expect(t, p.NewMove(E4, D5).value(), 1258) // PxQ
	expect(t, p.NewMove(C3, D5).value(), 1256) // NxQ
	expect(t, p.NewMove(C4, D5).value(), 1254) // BxQ
	expect(t, p.NewMove(A5, D5).value(), 1252) // RxQ
	expect(t, p.NewMove(D1, D5).value(), 1250) // QxQ
	expect(t, p.NewMove(D6, D5).value(), 1248) // KxQ
}

// PxR, NxR, BxR, RxR, QxR, KxR
func TestMove010(t *testing.T) {
	game := NewGame().Setup(`Kd6,Qd1,Ra5,Nc3,Bc4,e4`, `Kh8,Rd5`)
	p := game.Start(White)
	expect(t, p.NewMove(E4, D5).value(), 633) // PxR
	expect(t, p.NewMove(C3, D5).value(), 631) // NxR
	expect(t, p.NewMove(C4, D5).value(), 629) // BxR
	expect(t, p.NewMove(A5, D5).value(), 627) // RxR
	expect(t, p.NewMove(D1, D5).value(), 625) // QxR
	expect(t, p.NewMove(D6, D5).value(), 623) // KxR
}

// PxB, NxB, BxB, RxB, QxB, KxB
func TestMove020(t *testing.T) {
	game := NewGame().Setup(`Kd6,Qd1,Ra5,Nc3,Bc4,e4`, `Kh8,Bd5`)
	p := game.Start(White)
	expect(t, p.NewMove(E4, D5).value(), 416) // PxB
	expect(t, p.NewMove(C3, D5).value(), 414) // NxB
	expect(t, p.NewMove(C4, D5).value(), 412) // BxB
	expect(t, p.NewMove(A5, D5).value(), 410) // RxB
	expect(t, p.NewMove(D1, D5).value(), 408) // QxB
	expect(t, p.NewMove(D6, D5).value(), 406) // KxB
}

// PxN, NxN, BxN, RxN, QxN, KxN
func TestMove030(t *testing.T) {
	game := NewGame().Setup(`Kd6,Qd1,Ra5,Nc3,Bc4,e4`, `Kh8,Nd5`)
	p := game.Start(White)
	expect(t, p.NewMove(E4, D5).value(), 406) // PxN
	expect(t, p.NewMove(C3, D5).value(), 404) // NxN
	expect(t, p.NewMove(C4, D5).value(), 402) // BxN
	expect(t, p.NewMove(A5, D5).value(), 400) // RxN
	expect(t, p.NewMove(D1, D5).value(), 398) // QxN
	expect(t, p.NewMove(D6, D5).value(), 396) // KxN
}

// PxP, NxP, BxP, RxP, QxP, KxP
func TestMove040(t *testing.T) {
	game := NewGame().Setup(`Kd6,Qd1,Ra5,Nc3,Bc4,e4`, `Kh8,d5`)
	p := game.Start(White)
	expect(t, p.NewMove(E4, D5).value(), 98) // PxP
	expect(t, p.NewMove(C3, D5).value(), 96) // NxP
	expect(t, p.NewMove(C4, D5).value(), 94) // BxP
	expect(t, p.NewMove(A5, D5).value(), 92) // RxP
	expect(t, p.NewMove(D1, D5).value(), 90) // QxP
	expect(t, p.NewMove(D6, D5).value(), 88) // KxP
}
