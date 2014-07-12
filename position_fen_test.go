// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import `testing`

// Initial position: castles, no en-passant.
func TestFen100(t *testing.T) {
	p := NewGame(`rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1`).Start()
	expect(t, p.color, White)
	expect(t, p.castles, uint(0x0F))
	expect(t, p.enpassant, 0)
	expect(t, p.king[White], E1)
	expect(t, p.king[Black], E8)
	expect(t, p.outposts[Pawn], bit[A2]|bit[B2]|bit[C2]|bit[D2]|bit[E2]|bit[F2]|bit[G2]|bit[H2])
	expect(t, p.outposts[Knight], bit[B1]|bit[G1])
	expect(t, p.outposts[Bishop], bit[C1]|bit[F1])
	expect(t, p.outposts[Rook], bit[A1]|bit[H1])
	expect(t, p.outposts[Queen], bit[D1])
	expect(t, p.outposts[King], bit[E1])
	expect(t, p.outposts[BlackPawn], bit[A7]|bit[B7]|bit[C7]|bit[D7]|bit[E7]|bit[F7]|bit[G7]|bit[H7])
	expect(t, p.outposts[BlackKnight], bit[B8]|bit[G8])
	expect(t, p.outposts[BlackBishop], bit[C8]|bit[F8])
	expect(t, p.outposts[BlackRook], bit[A8]|bit[H8])
	expect(t, p.outposts[BlackQueen], bit[D8])
	expect(t, p.outposts[BlackKing], bit[E8])
}

// Castles, no en-passant.
func TestFen110(t *testing.T) {
	p := NewGame(`2r1kb1r/pp3ppp/2n1b3/1q1N2B1/1P2Q3/8/P4PPP/3RK1NR w Kk - 42 42`).Start()
	expect(t, p.color, White)
	expect(t, p.castles, castleKingside[White] | castleKingside[Black])
	expect(t, p.enpassant, 0)
	expect(t, p.king[White], E1)
	expect(t, p.king[Black], E8)
	expect(t, p.outposts[Pawn], bit[A2]|bit[B4]|bit[F2]|bit[G2]|bit[H2])
	expect(t, p.outposts[Knight], bit[D5]|bit[G1])
	expect(t, p.outposts[Bishop], bit[G5])
	expect(t, p.outposts[Rook], bit[D1]|bit[H1])
	expect(t, p.outposts[Queen], bit[E4])
	expect(t, p.outposts[King], bit[E1])
	expect(t, p.outposts[BlackPawn], bit[A7]|bit[B7]|bit[F7]|bit[G7]|bit[H7])
	expect(t, p.outposts[BlackKnight], bit[C6])
	expect(t, p.outposts[BlackBishop], bit[E6]|bit[F8])
	expect(t, p.outposts[BlackRook], bit[C8]|bit[H8])
	expect(t, p.outposts[BlackQueen], bit[B5])
	expect(t, p.outposts[BlackKing], bit[E8])
}

// No castles, en-passant.
func TestFen120(t *testing.T) {
	p := NewGame(`1rr2k2/p1q5/3p2Q1/3Pp2p/8/1P3P2/1KPRN3/8 w - e6 42 42`).Start()
	expect(t, p.color, White)
	expect(t, p.castles, uint8(0))
	expect(t, p.enpassant, E6)
	expect(t, p.king[White], B2)
	expect(t, p.king[Black], F8)
	expect(t, p.outposts[Pawn], bit[B3]|bit[C2]|bit[D5]|bit[F3])
	expect(t, p.outposts[Knight], bit[E2])
	expect(t, p.outposts[Bishop], Bitmask(0))
	expect(t, p.outposts[Rook], bit[D2])
	expect(t, p.outposts[Queen], bit[G6])
	expect(t, p.outposts[King], bit[B2])
	expect(t, p.outposts[BlackPawn], bit[A7]|bit[D6]|bit[E5]|bit[H5])
	expect(t, p.outposts[BlackKnight], Bitmask(0))
	expect(t, p.outposts[BlackBishop], Bitmask(0))
	expect(t, p.outposts[BlackRook], bit[B8]|bit[C8])
	expect(t, p.outposts[BlackQueen], bit[C7])
	expect(t, p.outposts[BlackKing], bit[F8])
}

// Position to FEN tests.

// Initial position: castles, no en-passant.
func TestFen200(t *testing.T) {
	p := NewGame(`rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1`).Start()
	expect(t, p.fen(), `rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq -`)
}

// Castles, no en-passant.
func TestFen210(t *testing.T) {
	p := NewGame(`2r1kb1r/pp3ppp/2n1b3/1q1N2B1/1P2Q3/8/P4PPP/3RK1NR w Kk - 42 42`).Start()
	expect(t, p.fen(), `2r1kb1r/pp3ppp/2n1b3/1q1N2B1/1P2Q3/8/P4PPP/3RK1NR w Kk -`)
}

// No castles, en-passant.
func TestFen220(t *testing.T) {
	p := NewGame(`1rr2k2/p1q5/3p2Q1/3Pp2p/8/1P3P2/1KPRN3/8 w - e6 42 42`).Start()
	expect(t, p.fen(), `1rr2k2/p1q5/3p2Q1/3Pp2p/8/1P3P2/1KPRN3/8 w - e6`)
}
