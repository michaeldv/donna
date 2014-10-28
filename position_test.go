// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`github.com/michaeldv/donna/expect`; `testing`)

// Initial position: castles, no en-passant.
func TestPosition000(t *testing.T) {
	p := NewGame(`rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1`).Start()
	expect.Eq(t, p.color, White)
	expect.Eq(t, p.castles, uint(0x0F))
	expect.Eq(t, p.enpassant, 0)
	expect.Eq(t, p.king[White], E1)
	expect.Eq(t, p.king[Black], E8)
	expect.Eq(t, p.outposts[Pawn], bit[A2]|bit[B2]|bit[C2]|bit[D2]|bit[E2]|bit[F2]|bit[G2]|bit[H2])
	expect.Eq(t, p.outposts[Knight], bit[B1]|bit[G1])
	expect.Eq(t, p.outposts[Bishop], bit[C1]|bit[F1])
	expect.Eq(t, p.outposts[Rook], bit[A1]|bit[H1])
	expect.Eq(t, p.outposts[Queen], bit[D1])
	expect.Eq(t, p.outposts[King], bit[E1])
	expect.Eq(t, p.outposts[BlackPawn], bit[A7]|bit[B7]|bit[C7]|bit[D7]|bit[E7]|bit[F7]|bit[G7]|bit[H7])
	expect.Eq(t, p.outposts[BlackKnight], bit[B8]|bit[G8])
	expect.Eq(t, p.outposts[BlackBishop], bit[C8]|bit[F8])
	expect.Eq(t, p.outposts[BlackRook], bit[A8]|bit[H8])
	expect.Eq(t, p.outposts[BlackQueen], bit[D8])
	expect.Eq(t, p.outposts[BlackKing], bit[E8])
}

// Castles, no en-passant.
func TestPosition010(t *testing.T) {
	p := NewGame(`2r1kb1r/pp3ppp/2n1b3/1q1N2B1/1P2Q3/8/P4PPP/3RK1NR w Kk - 42 42`).Start()
	expect.Eq(t, p.color, White)
	expect.Eq(t, p.castles, castleKingside[White] | castleKingside[Black])
	expect.Eq(t, p.enpassant, 0)
	expect.Eq(t, p.king[White], E1)
	expect.Eq(t, p.king[Black], E8)
	expect.Eq(t, p.outposts[Pawn], bit[A2]|bit[B4]|bit[F2]|bit[G2]|bit[H2])
	expect.Eq(t, p.outposts[Knight], bit[D5]|bit[G1])
	expect.Eq(t, p.outposts[Bishop], bit[G5])
	expect.Eq(t, p.outposts[Rook], bit[D1]|bit[H1])
	expect.Eq(t, p.outposts[Queen], bit[E4])
	expect.Eq(t, p.outposts[King], bit[E1])
	expect.Eq(t, p.outposts[BlackPawn], bit[A7]|bit[B7]|bit[F7]|bit[G7]|bit[H7])
	expect.Eq(t, p.outposts[BlackKnight], bit[C6])
	expect.Eq(t, p.outposts[BlackBishop], bit[E6]|bit[F8])
	expect.Eq(t, p.outposts[BlackRook], bit[C8]|bit[H8])
	expect.Eq(t, p.outposts[BlackQueen], bit[B5])
	expect.Eq(t, p.outposts[BlackKing], bit[E8])
}

// No castles, en-passant.
func TestPosition020(t *testing.T) {
	p := NewGame(`1rr2k2/p1q5/3p2Q1/3Pp2p/8/1P3P2/1KPRN3/8 w - e6 42 42`).Start()
	expect.Eq(t, p.color, White)
	expect.Eq(t, p.castles, uint8(0))
	expect.Eq(t, p.enpassant, E6)
	expect.Eq(t, p.king[White], B2)
	expect.Eq(t, p.king[Black], F8)
	expect.Eq(t, p.outposts[Pawn], bit[B3]|bit[C2]|bit[D5]|bit[F3])
	expect.Eq(t, p.outposts[Knight], bit[E2])
	expect.Eq(t, p.outposts[Bishop], Bitmask(0))
	expect.Eq(t, p.outposts[Rook], bit[D2])
	expect.Eq(t, p.outposts[Queen], bit[G6])
	expect.Eq(t, p.outposts[King], bit[B2])
	expect.Eq(t, p.outposts[BlackPawn], bit[A7]|bit[D6]|bit[E5]|bit[H5])
	expect.Eq(t, p.outposts[BlackKnight], Bitmask(0))
	expect.Eq(t, p.outposts[BlackBishop], Bitmask(0))
	expect.Eq(t, p.outposts[BlackRook], bit[B8]|bit[C8])
	expect.Eq(t, p.outposts[BlackQueen], bit[C7])
	expect.Eq(t, p.outposts[BlackKing], bit[F8])
}

// Position to FEN tests.

// Initial position: castles, no en-passant.
func TestPosition100(t *testing.T) {
	p := NewGame(`rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1`).Start()
	expect.Eq(t, p.fen(), `rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq -`)
}

// Castles, no en-passant.
func TestPosition110(t *testing.T) {
	p := NewGame(`2r1kb1r/pp3ppp/2n1b3/1q1N2B1/1P2Q3/8/P4PPP/3RK1NR w Kk - 42 42`).Start()
	expect.Eq(t, p.fen(), `2r1kb1r/pp3ppp/2n1b3/1q1N2B1/1P2Q3/8/P4PPP/3RK1NR w Kk -`)
}

// No castles, en-passant.
func TestPosition120(t *testing.T) {
	p := NewGame(`1rr2k2/p1q5/3p2Q1/3Pp2p/8/1P3P2/1KPRN3/8 w - e6 42 42`).Start()
	expect.Eq(t, p.fen(), `1rr2k2/p1q5/3p2Q1/3Pp2p/8/1P3P2/1KPRN3/8 w - e6`)
}

// Position status.
func TestPosition200(t *testing.T) {
	p := NewGame().Start()
	expect.Eq(t, p.status(NewMove(p, E2, E4), p.Evaluate()), InProgress)
}

// Mate in 1 move.
func TestPosition210(t *testing.T) {
	p := NewGame(`Kf8,Rh1,g6`, `Kh8,Bg8,g7,h7`).Start(White)
	rootNode = node // Reset ply().
	expect.Eq(t, p.status(NewMove(p, H1, H6), Checkmate - ply()), WhiteWinning)
}

// Forced stalemate.
func TestPosition220(t *testing.T) {
	p := NewGame(`Kf7,b2,b4,h6`, `Kh8,Ba4,b3,b5,h7`).Start(White)
	expect.Eq(t, p.status(NewMove(p, F7, F8), 0), Stalemate)
}

// Self-imposed stalemate.
func TestPosition230(t *testing.T) {
	p := NewGame(`Ka1,g3,h2`, `Kh5,h3,g4,g5,g6,h7`).Start(Black)
	p = p.MakeMove(NewMove(p, H7, H6))
	expect.Eq(t, p.status(NewMove(p, A1, B2), 0), Stalemate)
}

// Draw by repetition.
func TestPosition240(t *testing.T) {
	p := NewGame(`Ka1,g3,h2`, `Kh5,h3,g4,g5,g6,h7`).Start(Black) // Initial.

	p = p.MakeMove(NewMove(p, H5, H6))
	p = p.MakeMove(NewMove(p, A1, A2))
	p = p.MakeMove(NewMove(p, H6, H5))
	p = p.MakeMove(NewMove(p, A2, A1)) // Rep #2.
	expect.Eq(t, p.status(NewMove(p, H5, H6), 0), InProgress)

	p = p.MakeMove(NewMove(p, H5, H6))
	p = p.MakeMove(NewMove(p, A1, A2))
	p = p.MakeMove(NewMove(p, H6, H5)) // -- No NewMove(p, A2, A1) here --

	rootNode = node // Reset ply().
	expect.Eq(t, p.status(NewMove(p, A2, A1), 0), Repetition) // <-- Ka2-a1 causes rep #3.
}

// Insufficient material.
func TestPosition250(t *testing.T) {
	p := NewGame(`Ka1,Bb2`, `Kh5`).Start(White)
	expect.True(t, p.insufficient())
}
