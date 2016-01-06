// Copyright (c) 2014-2016 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`github.com/michaeldv/donna/expect`; `testing`)

func TestMagic000(t *testing.T) {
	expect.Eq(t, maskBlock[C3][H8], Bitmask( bit[D4] | bit[E5] | bit[F6] | bit[G7] | bit[H8] ))
	expect.Eq(t, maskBlock[C3][C8], Bitmask( bit[C4] | bit[C5] | bit[C6] | bit[C7] | bit[C8] ))
	expect.Eq(t, maskBlock[C3][A5], Bitmask( bit[B4] | bit[A5]                               ))
	expect.Eq(t, maskBlock[C3][A3], Bitmask( bit[B3] | bit[A3]                               ))
	expect.Eq(t, maskBlock[C3][A1], Bitmask( bit[B2] | bit[A1]                               ))
	expect.Eq(t, maskBlock[C3][C1], Bitmask( bit[C2] | bit[C1]                               ))
	expect.Eq(t, maskBlock[C3][E1], Bitmask( bit[D2] | bit[E1]                               ))
	expect.Eq(t, maskBlock[C3][H3], Bitmask( bit[D3] | bit[E3] | bit[F3] | bit[G3] | bit[H3] ))
	expect.Eq(t, maskBlock[C3][E7], Bitmask(0))
}

func TestMagic010(t *testing.T) {
	expect.Eq(t, maskEvade[C3][H8], Bitmask( ^bit[B2] ))
	expect.Eq(t, maskEvade[C3][C8], Bitmask( ^bit[C2] ))
	expect.Eq(t, maskEvade[C3][A5], Bitmask( ^bit[D2] ))
	expect.Eq(t, maskEvade[C3][A3], Bitmask( ^bit[D3] ))
	expect.Eq(t, maskEvade[C3][A1], Bitmask( ^bit[D4] ))
	expect.Eq(t, maskEvade[C3][C1], Bitmask( ^bit[C4] ))
	expect.Eq(t, maskEvade[C3][E1], Bitmask( ^bit[B4] ))
	expect.Eq(t, maskEvade[C3][H3], Bitmask( ^bit[B3] ))
	expect.Eq(t, maskEvade[C3][E7], Bitmask(maskFull))
}

func TestMagic020(t *testing.T) {
	expect.Eq(t, maskPawn[White][A3], Bitmask( bit[B2] ))
	expect.Eq(t, maskPawn[White][D5], Bitmask( bit[C4] | bit[E4] ))
	expect.Eq(t, maskPawn[White][F8], Bitmask( bit[E7] | bit[G7] ))
	expect.Eq(t, maskPawn[Black][H4], Bitmask( bit[G5] ))
	expect.Eq(t, maskPawn[Black][C5], Bitmask( bit[B6] | bit[D6] ))
	expect.Eq(t, maskPawn[Black][B1], Bitmask( bit[A2] | bit[C2] ))
}

func TestMagic030(t *testing.T) {
	// Same file.
	expect.Eq(t, maskLine[A2][A5], maskFile[0])
	expect.Eq(t, maskLine[H6][H1], maskFile[7])
	// Same rank.
	expect.Eq(t, maskLine[A2][F2], maskRank[1])
	expect.Eq(t, maskLine[H6][B6], maskRank[5])
	// Edge cases.
	expect.Eq(t, maskLine[A1][C5], maskNone) // Random squares.
	expect.Eq(t, maskLine[E4][E4], maskNone) // Same square.
}

func TestMagic040(t *testing.T) {
	// Same diagonal.
	expect.Eq(t, maskLine[C4][F7], bit[A2] | bit[B3] | bit[C4] | bit[D5] | bit[E6] | bit[F7] | bit[G8])
	expect.Eq(t, maskLine[F6][H8], maskA1H8)
	expect.Eq(t, maskLine[F1][H3], bit[F1] | bit[G2] | bit[H3])
	// Same anti-diagonal.
	expect.Eq(t, maskLine[C2][B3], bit[D1] | bit[C2] | bit[B3] | bit[A4])
	expect.Eq(t, maskLine[F3][B7], maskH1A8)
	expect.Eq(t, maskLine[H3][D7], bit[H3] | bit[G4] | bit[F5] | bit[E6] | bit[D7] | bit[C8])
	// Edge cases.
	expect.Eq(t, maskLine[A2][G4], maskNone) // Random squares.
	expect.Eq(t, maskLine[E4][E4], maskNone) // Same square.
}

// Material base tests.

// Bare kings.
func TestMaterial000(t *testing.T) {
	balance := materialBalance[King] + materialBalance[BlackKing]
	expect.Eq(t, balance, 0)
	expect.Eq(t, materialBase[balance].flags, uint8(materialDraw))
	expect.Eq(t, materialBase[balance].endgame, nil)

	p := NewGame(`Ke1`, `Ke8`).start()
	expect.Eq(t, p.balance, balance)
}

// No pawns, king with a minor.
func TestMaterial010(t *testing.T) {
	balance := materialBalance[Bishop]
	expect.Eq(t, materialBase[balance].flags, uint8(materialDraw))
	expect.Eq(t, materialBase[balance].endgame, nil)

	p := NewGame(`Ke1,Bc1`, `Ke8`).start()
	expect.Eq(t, p.balance, balance)
}

func TestMaterial015(t *testing.T) {
	balance := materialBalance[Bishop] + materialBalance[BlackKnight]
	expect.Eq(t, materialBase[balance].flags, uint8(materialDraw))
	expect.Eq(t, materialBase[balance].endgame, nil)

	p := NewGame(`Ke1,Bc1`, `Ke8,Nb8`).start()
	expect.Eq(t, p.balance, balance)
}

// No pawns, king with two knights.
func TestMaterial020(t *testing.T) {
	balance := 2 * materialBalance[Knight]
	expect.Eq(t, materialBase[balance].flags, uint8(materialDraw))
	expect.Eq(t, materialBase[balance].endgame, nil)

	p := NewGame(`Ke1,Ne2,Ne3`, `Ke8`).start()
	expect.Eq(t, p.balance, balance)
}

// Known: king and a pawn vs. bare king.
func TestMaterial030(t *testing.T) {
	balance := materialBalance[Pawn]
	expect.Eq(t, materialBase[balance].flags, uint8(knownEndgame))
	expect.Eq(t, materialBase[balance].endgame, (*Evaluation).kingAndPawnVsBareKing)

	p := NewGame(`Ke1,e2`, `Ke8`).start()
	expect.Eq(t, p.balance, balance)
}

func TestMaterial040(t *testing.T) {
	balance := materialBalance[BlackPawn]
	expect.Eq(t, materialBase[balance].flags, uint8(knownEndgame))
	expect.Eq(t, materialBase[balance].endgame, (*Evaluation).kingAndPawnVsBareKing)

	p := NewGame(`Ke1`, `M,Ke8,e7`).start()
	expect.Eq(t, p.balance, balance)
}

// Known: king with a knight and a bishop vs. bare king.
func TestMaterial050(t *testing.T) {
	balance := materialBalance[Knight] + materialBalance[Bishop]
	expect.Eq(t, materialBase[balance].flags, uint8(knownEndgame))
	expect.Eq(t, materialBase[balance].endgame, (*Evaluation).knightAndBishopVsBareKing)

	p := NewGame(`Ke1,Nb1,Bc1`, `Ke8`).start()
	expect.Eq(t, p.balance, balance)
}

func TestMaterial060(t *testing.T) {
	balance := materialBalance[BlackKnight] + materialBalance[BlackBishop]
	expect.Eq(t, materialBase[balance].flags, uint8(knownEndgame))
	expect.Eq(t, materialBase[balance].endgame, (*Evaluation).knightAndBishopVsBareKing)

	p := NewGame(`Ke1`, `M,Ke8,Nb8,Bc8`).start()
	expect.Eq(t, p.balance, balance)
}

// Known endgame: two bishops vs. bare king.
func TestMaterial070(t *testing.T) {
	balance := 2 * materialBalance[BlackBishop]
	expect.Eq(t, materialBase[balance].flags, uint8(knownEndgame))
	expect.Eq(t, materialBase[balance].endgame, (*Evaluation).twoBishopsVsBareKing)

	p := NewGame(`Ke1`, `M,Ka8,Bg8,Bh8`).start()
	expect.Eq(t, p.balance, balance)
}

// Known endgame: king with some winning material vs. bare king.
func TestMaterial080(t *testing.T) {
	balance := materialBalance[BlackRook]
	expect.Eq(t, materialBase[balance].flags, uint8(knownEndgame))
	expect.Eq(t, materialBase[balance].endgame, (*Evaluation).winAgainstBareKing)

	p := NewGame(`Ke1`, `M,Ka8,Rh8`).start()
	expect.Eq(t, p.balance, balance)
}

// Lesser known endgame: king and two or more pawns vs. bare king.
func TestMaterial090(t *testing.T) {
	balance := 2 * materialBalance[Pawn]
	expect.Eq(t, materialBase[balance].flags, uint8(lesserKnownEndgame))
	expect.Eq(t, materialBase[balance].endgame, (*Evaluation).kingAndPawnsVsBareKing)

	p := NewGame(`Ke1,a4,a5`, `M,Ka8`).start()
	expect.Eq(t, p.balance, balance)
}

// Lesser known endgame: queen vs. rook with pawn(s)
func TestMaterial100(t *testing.T) {
	balance := materialBalance[Rook] + materialBalance[Pawn] + materialBalance[BlackQueen]
	expect.Eq(t, materialBase[balance].flags, uint8(lesserKnownEndgame))
	expect.Eq(t, materialBase[balance].endgame, (*Evaluation).queenVsRookAndPawns)

	p := NewGame(`Ke1,Re4,e5`, `M,Ka8,Qh8`).start()
	expect.Eq(t, p.balance, balance)
}

// Lesser known endgame: king and pawn vs. king and pawn.
func TestMaterial110(t *testing.T) {
	balance := materialBalance[Pawn] + materialBalance[BlackPawn]
	expect.Eq(t, materialBase[balance].flags, uint8(lesserKnownEndgame))
	expect.Eq(t, materialBase[balance].endgame, (*Evaluation).kingAndPawnVsKingAndPawn)

	p := NewGame(`Ke1,a4`, `M,Ka8,h5`).start()
	expect.Eq(t, p.balance, balance)
}

// Lesser known endgame: bishop and pawn vs. bare king.
func TestMaterial120(t *testing.T) {
	balance := materialBalance[Pawn] + materialBalance[Bishop]
	expect.Eq(t, materialBase[balance].flags, uint8(lesserKnownEndgame))
	expect.Eq(t, materialBase[balance].endgame, (*Evaluation).bishopAndPawnVsBareKing)

	p := NewGame(`Ke1,Be2,a4`, `Ka8`).start()
	expect.Eq(t, p.balance, balance)
}

// Lesser known endgame: rook and pawn vs. rook.
func TestMaterial130(t *testing.T) {
	balance := materialBalance[Rook] + materialBalance[Pawn] + materialBalance[BlackRook]
	expect.Eq(t, materialBase[balance].flags, uint8(lesserKnownEndgame))
	expect.Eq(t, materialBase[balance].endgame, (*Evaluation).rookAndPawnVsRook)

	p := NewGame(`Ke1,Re2,a4`, `Ka8,Rh8`).start()
	expect.Eq(t, p.balance, balance)
}

// Single bishops (midgame).
func TestMaterial140(t *testing.T) {
	balance := materialBalance[Pawn] * 2 + materialBalance[Bishop] + materialBalance[Knight] + materialBalance[Rook] + materialBalance[BlackPawn] * 2 + materialBalance[BlackBishop] + materialBalance[BlackKnight] + materialBalance[BlackRook]
	expect.Eq(t, materialBase[balance].flags, uint8(singleBishops | lesserKnownEndgame))
	expect.Eq(t, materialBase[balance].endgame, (*Evaluation).drawishBishops)

	p := NewGame(`Ke1,Ra1,Bc1,Nb1,d2,e2`, `Ke8,Rh8,Bf8,Ng8,d7,e7`).start()
	expect.Eq(t, p.balance, balance)
}

// Single bishops (endgame).
func TestMaterial150(t *testing.T) {
	balance := materialBalance[Bishop] + 4 * materialBalance[Pawn] + materialBalance[BlackBishop] + 3 * materialBalance[BlackPawn]
	expect.Eq(t, materialBase[balance].flags, uint8(singleBishops | lesserKnownEndgame))
	expect.Eq(t, materialBase[balance].endgame, (*Evaluation).bishopsAndPawns)

	p := NewGame(`Ke1,Bc1,a2,b2,c2,d4`, `Ke8,Bf8,f7,g7,h7`).start()
	expect.Eq(t, p.balance, balance)
}
