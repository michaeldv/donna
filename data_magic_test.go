// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`)

func TestMagic000(t *testing.T) {
        expect(t, maskBlock[C3][H8], Bitmask( bit[D4] | bit[E5] | bit[F6] | bit[G7] | bit[H8] ))
        expect(t, maskBlock[C3][C8], Bitmask( bit[C4] | bit[C5] | bit[C6] | bit[C7] | bit[C8] ))
        expect(t, maskBlock[C3][A5], Bitmask( bit[B4] | bit[A5]                               ))
        expect(t, maskBlock[C3][A3], Bitmask( bit[B3] | bit[A3]                               ))
        expect(t, maskBlock[C3][A1], Bitmask( bit[B2] | bit[A1]                               ))
        expect(t, maskBlock[C3][C1], Bitmask( bit[C2] | bit[C1]                               ))
        expect(t, maskBlock[C3][E1], Bitmask( bit[D2] | bit[E1]                               ))
        expect(t, maskBlock[C3][H3], Bitmask( bit[D3] | bit[E3] | bit[F3] | bit[G3] | bit[H3] ))
        expect(t, maskBlock[C3][E7], Bitmask(0))
}

func TestMagic010(t *testing.T) {
        expect(t, maskEvade[C3][H8], Bitmask( ^bit[B2] ))
        expect(t, maskEvade[C3][C8], Bitmask( ^bit[C2] ))
        expect(t, maskEvade[C3][A5], Bitmask( ^bit[D2] ))
        expect(t, maskEvade[C3][A3], Bitmask( ^bit[D3] ))
        expect(t, maskEvade[C3][A1], Bitmask( ^bit[D4] ))
        expect(t, maskEvade[C3][C1], Bitmask( ^bit[C4] ))
        expect(t, maskEvade[C3][E1], Bitmask( ^bit[B4] ))
        expect(t, maskEvade[C3][H3], Bitmask( ^bit[B3] ))
        expect(t, maskEvade[C3][E7], Bitmask(maskFull))
}

func TestMagic020(t *testing.T) {
        expect(t, maskPawn[White][A3], Bitmask( bit[B2] ))
        expect(t, maskPawn[White][D5], Bitmask( bit[C4] | bit[E4] ))
        expect(t, maskPawn[White][F8], Bitmask( bit[E7] | bit[G7] ))
        expect(t, maskPawn[Black][H4], Bitmask( bit[G5] ))
        expect(t, maskPawn[Black][C5], Bitmask( bit[B6] | bit[D6] ))
        expect(t, maskPawn[Black][B1], Bitmask( bit[A2] | bit[C2] ))
}

func TestMagic030(t *testing.T) {
        // Same file.
        expect(t, maskStraight[A2][A5], maskFile[0])
        expect(t, maskStraight[H6][H1], maskFile[7])
        // Same rank.
        expect(t, maskStraight[A2][F2], maskRank[1])
        expect(t, maskStraight[H6][B6], maskRank[5])
        // Edge cases.
        expect(t, maskStraight[A1][C5], maskNone) // Random squares.
        expect(t, maskStraight[E4][E4], maskNone) // Same square.
}

func TestMagic040(t *testing.T) {
        // Same diagonal.
        expect(t, maskDiagonal[C4][F7], bit[A2] | bit[B3] | bit[C4] | bit[D5] | bit[E6] | bit[F7] | bit[G8])
        expect(t, maskDiagonal[F6][H8], maskA1H8)
        expect(t, maskDiagonal[F1][H3], bit[F1] | bit[G2] | bit[H3])
        // Same anti-diagonal.
        expect(t, maskDiagonal[C2][B3], bit[D1] | bit[C2] | bit[B3] | bit[A4])
        expect(t, maskDiagonal[F3][B7], maskH1A8)
        expect(t, maskDiagonal[H3][D7], bit[H3] | bit[G4] | bit[F5] | bit[E6] | bit[D7] | bit[C8])
        // Edge cases.
        expect(t, maskDiagonal[A2][G4], maskNone) // Random squares.
        expect(t, maskDiagonal[E4][E4], maskNone) // Same square.
}
