// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

import(`github.com/michaeldv/donna/expect`; `testing`)

func TestBitmask000(t *testing.T) { // maskPassed White
	expect.Eq(t, maskPassed[White][A2], Bitmask(0x0303030303030000))
	expect.Eq(t, maskPassed[White][B2], Bitmask(0x0707070707070000))
	expect.Eq(t, maskPassed[White][C2], Bitmask(0x0E0E0E0E0E0E0000))
	expect.Eq(t, maskPassed[White][D2], Bitmask(0x1C1C1C1C1C1C0000))
	expect.Eq(t, maskPassed[White][E2], Bitmask(0x3838383838380000))
	expect.Eq(t, maskPassed[White][F2], Bitmask(0x7070707070700000))
	expect.Eq(t, maskPassed[White][G2], Bitmask(0xE0E0E0E0E0E00000))
	expect.Eq(t, maskPassed[White][H2], Bitmask(0xC0C0C0C0C0C00000))

	expect.Eq(t, maskPassed[White][A1], Bitmask(0x0303030303030300))
	expect.Eq(t, maskPassed[White][H8], Bitmask(0x0000000000000000))
	expect.Eq(t, maskPassed[White][C6], Bitmask(0x0E0E000000000000))
}

func TestBitmask010(t *testing.T) { // maskPassed Black
	expect.Eq(t, maskPassed[Black][A7], Bitmask(0x0000030303030303))
	expect.Eq(t, maskPassed[Black][B7], Bitmask(0x0000070707070707))
	expect.Eq(t, maskPassed[Black][C7], Bitmask(0x00000E0E0E0E0E0E))
	expect.Eq(t, maskPassed[Black][D7], Bitmask(0x00001C1C1C1C1C1C))
	expect.Eq(t, maskPassed[Black][E7], Bitmask(0x0000383838383838))
	expect.Eq(t, maskPassed[Black][F7], Bitmask(0x0000707070707070))
	expect.Eq(t, maskPassed[Black][G7], Bitmask(0x0000E0E0E0E0E0E0))
	expect.Eq(t, maskPassed[Black][H7], Bitmask(0x0000C0C0C0C0C0C0))

	expect.Eq(t, maskPassed[Black][A1], Bitmask(0x0000000000000000))
	expect.Eq(t, maskPassed[Black][H8], Bitmask(0x00C0C0C0C0C0C0C0))
	expect.Eq(t, maskPassed[Black][C6], Bitmask(0x0000000E0E0E0E0E))
}

func TestBitmask020(t *testing.T) { // maskInFront White
	expect.Eq(t, maskInFront[0][A4], Bitmask(0x0101010100000000))
	expect.Eq(t, maskInFront[0][B4], Bitmask(0x0202020200000000))
	expect.Eq(t, maskInFront[0][C4], Bitmask(0x0404040400000000))
	expect.Eq(t, maskInFront[0][D4], Bitmask(0x0808080800000000))
	expect.Eq(t, maskInFront[0][E4], Bitmask(0x1010101000000000))
	expect.Eq(t, maskInFront[0][F4], Bitmask(0x2020202000000000))
	expect.Eq(t, maskInFront[0][G4], Bitmask(0x4040404000000000))
	expect.Eq(t, maskInFront[0][H4], Bitmask(0x8080808000000000))
}

func TestBitmask030(t *testing.T) { // maskInFront Black
	expect.Eq(t, maskInFront[1][A7], Bitmask(0x0000010101010101))
	expect.Eq(t, maskInFront[1][B7], Bitmask(0x0000020202020202))
	expect.Eq(t, maskInFront[1][C7], Bitmask(0x0000040404040404))
	expect.Eq(t, maskInFront[1][D7], Bitmask(0x0000080808080808))
	expect.Eq(t, maskInFront[1][E7], Bitmask(0x0000101010101010))
	expect.Eq(t, maskInFront[1][F7], Bitmask(0x0000202020202020))
	expect.Eq(t, maskInFront[1][G7], Bitmask(0x0000404040404040))
	expect.Eq(t, maskInFront[1][H7], Bitmask(0x0000808080808080))
}

func TestBitmask040(t *testing.T) { 	// maskBlock A1->H8 (fill)
	expect.Eq(t, maskBlock[A1][A1], Bitmask(0x0000000000000000))
	expect.Eq(t, maskBlock[A1][B2], Bitmask(0x0000000000000200))
	expect.Eq(t, maskBlock[A1][C3], Bitmask(0x0000000000040200))
	expect.Eq(t, maskBlock[A1][D4], Bitmask(0x0000000008040200))
	expect.Eq(t, maskBlock[A1][E5], Bitmask(0x0000001008040200))
	expect.Eq(t, maskBlock[A1][F6], Bitmask(0x0000201008040200))
	expect.Eq(t, maskBlock[A1][G7], Bitmask(0x0040201008040200))
	expect.Eq(t, maskBlock[A1][H8], Bitmask(0x8040201008040200))
}

func TestBitmask050(t *testing.T) { // maskBlock H1->A8 (fill)
	expect.Eq(t, maskBlock[H1][H1], Bitmask(0x0000000000000000))
	expect.Eq(t, maskBlock[H1][G2], Bitmask(0x0000000000004000))
	expect.Eq(t, maskBlock[H1][F3], Bitmask(0x0000000000204000))
	expect.Eq(t, maskBlock[H1][E4], Bitmask(0x0000000010204000))
	expect.Eq(t, maskBlock[H1][D5], Bitmask(0x0000000810204000))
	expect.Eq(t, maskBlock[H1][C6], Bitmask(0x0000040810204000))
	expect.Eq(t, maskBlock[H1][B7], Bitmask(0x0002040810204000))
	expect.Eq(t, maskBlock[H1][A8], Bitmask(0x0102040810204000))
}

func TestBitmask060(t *testing.T) { // maskEvade A1->H8 (spot)
	expect.Eq(t, maskEvade[A1][A1], Bitmask(0xFFFFFFFFFFFFFFFF))
	expect.Eq(t, maskEvade[B2][A1], Bitmask(0xFFFFFFFFFFFBFFFF))
	expect.Eq(t, maskEvade[C3][A1], Bitmask(0xFFFFFFFFF7FFFFFF))
	expect.Eq(t, maskEvade[D4][A1], Bitmask(0xFFFFFFEFFFFFFFFF))
	expect.Eq(t, maskEvade[E5][A1], Bitmask(0xFFFFDFFFFFFFFFFF))
	expect.Eq(t, maskEvade[F6][A1], Bitmask(0xFFBFFFFFFFFFFFFF))
	expect.Eq(t, maskEvade[G7][A1], Bitmask(0x7FFFFFFFFFFFFFFF))
	expect.Eq(t, maskEvade[H8][A1], Bitmask(0xFFFFFFFFFFFFFFFF))
}

func TestBitmask070(t *testing.T) { // maskEvade H1->A8 (spot)
	expect.Eq(t, maskEvade[H1][H1], Bitmask(0xFFFFFFFFFFFFFFFF))
	expect.Eq(t, maskEvade[G2][H1], Bitmask(0xFFFFFFFFFFDFFFFF))
	expect.Eq(t, maskEvade[F3][H1], Bitmask(0xFFFFFFFFEFFFFFFF))
	expect.Eq(t, maskEvade[E4][H1], Bitmask(0xFFFFFFF7FFFFFFFF))
	expect.Eq(t, maskEvade[D5][H1], Bitmask(0xFFFFFBFFFFFFFFFF))
	expect.Eq(t, maskEvade[C6][H1], Bitmask(0xFFFDFFFFFFFFFFFF))
	expect.Eq(t, maskEvade[B7][H1], Bitmask(0xFEFFFFFFFFFFFFFF))
	expect.Eq(t, maskEvade[A8][H1], Bitmask(0xFFFFFFFFFFFFFFFF))
}

func TestBitmask100(t *testing.T) {
	mask := Bitmask(0x0000000000000001)
	expect.Eq(t, mask.pop(), Bitmask(0x0000000000000000))

	mask = Bitmask(0xFFFFFFFFFFFFFFF0)
	expect.Eq(t, mask.pop(), Bitmask(0xFFFFFFFFFFFFFFE0))

	mask = Bitmask(0x8000000000000000)
	expect.Eq(t, mask.pop(), Bitmask(0x0000000000000000))
}
