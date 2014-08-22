// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`github.com/michaeldv/donna/expect`; `testing`)

func TestBitmask000(t *testing.T) { // White
	passed := [8]Bitmask{0}
	for square := A2; square <= H2; square++ {
		i := square - A2
		if Col(square) > 0 {
			passed[i].fill(square-1, 8, 0, 0x00FFFFFFFFFFFFFF)
		}
		passed[i].fill(square, 8, 0, 0x00FFFFFFFFFFFFFF)
		if Col(square) < 7 {
			passed[i].fill(square+1, 8, 0, 0x00FFFFFFFFFFFFFF)
		}
	}
	expect.Eq(t, passed[0], Bitmask(0x0303030303030000))
	expect.Eq(t, passed[1], Bitmask(0x0707070707070000))
	expect.Eq(t, passed[2], Bitmask(0x0E0E0E0E0E0E0000))
	expect.Eq(t, passed[3], Bitmask(0x1C1C1C1C1C1C0000))
	expect.Eq(t, passed[4], Bitmask(0x3838383838380000))
	expect.Eq(t, passed[5], Bitmask(0x7070707070700000))
	expect.Eq(t, passed[6], Bitmask(0xE0E0E0E0E0E00000))
	expect.Eq(t, passed[7], Bitmask(0xC0C0C0C0C0C00000))
}

func TestBitmask010(t *testing.T) { // Black
	passed := [8]Bitmask{0}
	for square := A7; square <= H7; square++ {
		i := square - A7
		if Col(square) > 0 {
			passed[i].fill(square-1, -8, 0, 0xFFFFFFFFFFFFFF00)
		}
		passed[i].fill(square, -8, 0, 0xFFFFFFFFFFFFFF00)
		if Col(square) < 7 {
			passed[i].fill(square+1, -8, 0, 0xFFFFFFFFFFFFFF00)
		}
	}
	expect.Eq(t, passed[0], Bitmask(0x0000030303030303))
	expect.Eq(t, passed[1], Bitmask(0x0000070707070707))
	expect.Eq(t, passed[2], Bitmask(0x00000E0E0E0E0E0E))
	expect.Eq(t, passed[3], Bitmask(0x00001C1C1C1C1C1C))
	expect.Eq(t, passed[4], Bitmask(0x0000383838383838))
	expect.Eq(t, passed[5], Bitmask(0x0000707070707070))
	expect.Eq(t, passed[6], Bitmask(0x0000E0E0E0E0E0E0))
	expect.Eq(t, passed[7], Bitmask(0x0000C0C0C0C0C0C0))
}

func TestBitmask030(t *testing.T) { // White
	forward := [8]Bitmask{0}
	for square := A4; square <= H4; square++ {
		i := square - A4
		forward[i].fill(square, 8, 0, 0x00FFFFFFFFFFFFFF)
	}
	expect.Eq(t, forward[0], Bitmask(0x0101010100000000))
	expect.Eq(t, forward[1], Bitmask(0x0202020200000000))
	expect.Eq(t, forward[2], Bitmask(0x0404040400000000))
	expect.Eq(t, forward[3], Bitmask(0x0808080800000000))
	expect.Eq(t, forward[4], Bitmask(0x1010101000000000))
	expect.Eq(t, forward[5], Bitmask(0x2020202000000000))
	expect.Eq(t, forward[6], Bitmask(0x4040404000000000))
	expect.Eq(t, forward[7], Bitmask(0x8080808000000000))
}

func TestBitmask040(t *testing.T) { // Black
	forward := [8]Bitmask{0}
	for square := A7; square <= H7; square++ {
		i := square - A7
		forward[i].fill(square, -8, 0, 0xFFFFFFFFFFFFFF00)
	}
	expect.Eq(t, forward[0], Bitmask(0x0000010101010101))
	expect.Eq(t, forward[1], Bitmask(0x0000020202020202))
	expect.Eq(t, forward[2], Bitmask(0x0000040404040404))
	expect.Eq(t, forward[3], Bitmask(0x0000080808080808))
	expect.Eq(t, forward[4], Bitmask(0x0000101010101010))
	expect.Eq(t, forward[5], Bitmask(0x0000202020202020))
	expect.Eq(t, forward[6], Bitmask(0x0000404040404040))
	expect.Eq(t, forward[7], Bitmask(0x0000808080808080))
}

func TestBitmask050(t *testing.T) {
	mask := Bitmask(0x0000000000000001)
	bit := mask.pop()
	expect.Eq(t, bit, 0)
	expect.Eq(t, mask, Bitmask(0x0000000000000000))

	mask = Bitmask(0x8000000000000000)
	bit = mask.pop()
	expect.Eq(t, bit, 63)
	expect.Eq(t, mask, Bitmask(0x0000000000000000))
}
