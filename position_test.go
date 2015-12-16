// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`github.com/michaeldv/donna/expect`; `testing`)

// Initial position: castles, no en-passant.
func TestPosition000(t *testing.T) {
	p := NewGame(`rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1`).start()
	expect.Eq(t, p.color, uint8(White))
	expect.Eq(t, p.castles, uint(0x0F))
	expect.Eq(t, p.enpassant, uint8(0))
	expect.Eq(t, p.king[White], uint8(E1))
	expect.Eq(t, p.king[Black], uint8(E8))
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
	p := NewGame(`2r1kb1r/pp3ppp/2n1b3/1q1N2B1/1P2Q3/8/P4PPP/3RK1NR w Kk - 42 42`).start()
	expect.Eq(t, p.color, uint8(White))
	expect.Eq(t, p.castles, castleKingside[White] | castleKingside[Black])
	expect.Eq(t, p.enpassant, uint8(0))
	expect.Eq(t, p.king[White], uint8(E1))
	expect.Eq(t, p.king[Black], uint8(E8))
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
	p := NewGame(`1rr2k2/p1q5/3p2Q1/3Pp2p/8/1P3P2/1KPRN3/8 w - e6 42 42`).start()
	expect.Eq(t, p.color, uint8(White))
	expect.Eq(t, p.castles, uint8(0))
	expect.Eq(t, p.enpassant, uint8(E6))
	expect.Eq(t, p.king[White], uint8(B2))
	expect.Eq(t, p.king[Black], uint8(F8))
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

//\\ Position to FEN tests.
// Initial position: castles, no en-passant.
func TestPosition100(t *testing.T) {
	p := NewGame(`rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1`).start()
	expect.Eq(t, p.fen(), `rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1`)
}

// Castles, no en-passant.
func TestPosition110(t *testing.T) {
	p := NewGame(`2r1kb1r/pp3ppp/2n1b3/1q1N2B1/1P2Q3/8/P4PPP/3RK1NR w Kk - 42 42`).start()
	expect.Eq(t, p.fen(), `2r1kb1r/pp3ppp/2n1b3/1q1N2B1/1P2Q3/8/P4PPP/3RK1NR w Kk - 42 1`)
}

// No castles, en-passant.
func TestPosition120(t *testing.T) {
	p := NewGame(`1rr2k2/p1q5/3p2Q1/3Pp2p/8/1P3P2/1KPRN3/8 w - e6 42 42`).start()
	expect.Eq(t, p.fen(), `1rr2k2/p1q5/3p2Q1/3Pp2p/8/1P3P2/1KPRN3/8 w - e6 42 1`)
}

//\\ Donna Chess Format (DCF) tests.
// Initial position: castles, no en-passant.
func TestPosition130(t *testing.T) {
	p := NewGame(`rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1`).start()
	expect.Eq(t, p.dcf(), `Ke1,Qd1,Ra1,Rh1,Bc1,Bf1,Nb1,Ng1,a2,b2,c2,d2,e2,f2,g2,h2 : Ke8,Qd8,Ra8,Rh8,Bc8,Bf8,Nb8,Ng8,a7,b7,c7,d7,e7,f7,g7,h7`)
}

// Castles, no en-passant.
func TestPosition140(t *testing.T) {
	p := NewGame(`2r1kb1r/pp3ppp/2n1b3/1q1N2B1/1P2Q3/8/P4PPP/3RK1NR w Kk - 42 42`).start()
	expect.Eq(t, p.dcf(), `Ke1,Qe4,Rd1,Rh1,Bg5,Ng1,Nd5,Cg1,a2,f2,g2,h2,b4 : Ke8,Qb5,Rc8,Rh8,Be6,Bf8,Nc6,Cg8,a7,b7,f7,g7,h7`)
}

// No castles, en-passant.
func TestPosition150(t *testing.T) {
	p := NewGame(`1rr2k2/p1q5/3p2Q1/3Pp2p/8/1P3P2/1KPRN3/8 w - e6 42 42`).start()
	expect.Eq(t, p.dcf(), `Kb2,Qg6,Rd2,Ne2,Ee6,c2,b3,f3,d5 : Kf8,Qc7,Rb8,Rc8,e5,h5,d6,a7`)

	pp := NewGame(`M,Kb2,Qg6,Rd2,Ne2,Ee6,c2,b3,f3,d5`, `Kf8,Qc7,Rb8,Rc8,e5,h5,d6,a7`).start()
	expect.Eq(t, pp.fen(), `1rr2k2/p1q5/3p2Q1/3Pp2p/8/1P3P2/1KPRN3/8 w - e6 0 1`)
}

// Position status.
func TestPosition200(t *testing.T) {
	p := NewGame().start()
	expect.Eq(t, p.status(NewMove(p, E2, E4), p.Evaluate()), InProgress)
}

// Mate in 1 move.
func TestPosition210(t *testing.T) {
	p := NewGame(`Kf8,Rh1,g6`, `Kh8,Bg8,g7,h7`).start()
	rootNode = node // Reset ply().
	expect.Eq(t, p.status(NewMove(p, H1, H6), Checkmate - ply()), WhiteWinning)
}

// Forced stalemate.
func TestPosition220(t *testing.T) {
	p := NewGame(`Kf7,b2,b4,h6`, `Kh8,Ba4,b3,b5,h7`).start()
	expect.Eq(t, p.status(NewMove(p, F7, F8), 0), Stalemate)
}

// Self-imposed stalemate.
func TestPosition230(t *testing.T) {
	p := NewGame(`Ka1,g3,h2`, `M,Kh5,h3,g4,g5,g6,h7`).start()
	p = p.makeMove(NewMove(p, H7, H6))
	expect.Eq(t, p.status(NewMove(p, A1, B2), 0), Stalemate)
}

// Draw by repetition.
func TestPosition240(t *testing.T) {
	p := NewGame(`Ka1,g3,h2`, `M,Kh5,h3,g4,g5,g6,h7`).start() // Initial.

	p = p.makeMove(NewMove(p, H5, H6))
	p = p.makeMove(NewMove(p, A1, A2))
	p = p.makeMove(NewMove(p, H6, H5))
	p = p.makeMove(NewMove(p, A2, A1)) // Rep #2.
	expect.Eq(t, p.status(NewMove(p, H5, H6), 0), InProgress)

	p = p.makeMove(NewMove(p, H5, H6))
	p = p.makeMove(NewMove(p, A1, A2))
	p = p.makeMove(NewMove(p, H6, H5)) // -- No NewMove(p, A2, A1) here --

	rootNode = node // Reset ply().
	expect.Eq(t, p.status(NewMove(p, A2, A1), 0), Repetition) // <-- Ka2-a1 causes rep #3.
}

// Insufficient material.
func TestPosition250(t *testing.T) {
	p := NewGame(`Ka1,Bb2`, `Kh5`).start()
	expect.True(t, p.insufficient())
}

// Restricted mobility for pinned pieces.
func TestPosition300(t *testing.T) {
	p := NewGame(`Ka1,a2,Nc3`, `Kh8,h7,Bg8`).start() // Nc3 vs Bishop, no pin.
	expect.Eq(t, p.Evaluate(), -3)
	p = NewGame(`Ka1,a2,Nc3`, `Kh8,h7,Bg7`).start() // Nc3 vs Bishop, pin on C3-G7 diagonal.
	expect.Eq(t, p.Evaluate(), -62)

}

func TestPosition310(t *testing.T) {
	p := NewGame(`Ka1,a2,Bc3`, `Kg8,h7,Bg6`).start() // Bc3 vs Bishop, no pin.
	expect.Eq(t, p.Evaluate(), 0)
	p = NewGame(`Ka1,a2,Bc3`, `Kg8,h7,Bg7`).start() // Bc3 vs Bishop, pin on C3-G7 diagonal.
	expect.Eq(t, p.Evaluate(), -23)
	p = NewGame(`Ka3,a2,Bc3`, `Kh8,h7,Rh1`).start() // Bc3 vs Rook, no pin.
	expect.Eq(t, p.Evaluate(), -206)
	p = NewGame(`Ka3,a2,Bc3`, `Kh8,h7,Rh3`).start() // Bc3 vs Rook, pin on C3-H3 file.
	expect.Eq(t, p.Evaluate(), -319)

}
