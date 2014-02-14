// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`)

func TestMagic000(t *testing.T) {
        expect(t, maskBlock[C3][H8], Bitmask( Bit(D4) | Bit(E5) | Bit(F6) | Bit(G7) | Bit(H8) ))
        expect(t, maskBlock[C3][C8], Bitmask( Bit(C4) | Bit(C5) | Bit(C6) | Bit(C7) | Bit(C8) ))
        expect(t, maskBlock[C3][A5], Bitmask( Bit(B4) | Bit(A5)                               ))
        expect(t, maskBlock[C3][A3], Bitmask( Bit(B3) | Bit(A3)                               ))
        expect(t, maskBlock[C3][A1], Bitmask( Bit(B2) | Bit(A1)                               ))
        expect(t, maskBlock[C3][C1], Bitmask( Bit(C2) | Bit(C1)                               ))
        expect(t, maskBlock[C3][E1], Bitmask( Bit(D2) | Bit(E1)                               ))
        expect(t, maskBlock[C3][H3], Bitmask( Bit(D3) | Bit(E3) | Bit(F3) | Bit(G3) | Bit(H3) ))
        expect(t, maskBlock[C3][E7], Bitmask(0))
}

func TestMagic010(t *testing.T) {
        expect(t, maskEvade[C3][H8], Bitmask( ^Bit(B2) ))
        expect(t, maskEvade[C3][C8], Bitmask( ^Bit(C2) ))
        expect(t, maskEvade[C3][A5], Bitmask( ^Bit(D2) ))
        expect(t, maskEvade[C3][A3], Bitmask( ^Bit(D3) ))
        expect(t, maskEvade[C3][A1], Bitmask( ^Bit(D4) ))
        expect(t, maskEvade[C3][C1], Bitmask( ^Bit(C4) ))
        expect(t, maskEvade[C3][E1], Bitmask( ^Bit(B4) ))
        expect(t, maskEvade[C3][H3], Bitmask( ^Bit(B3) ))
        expect(t, maskEvade[C3][E7], Bitmask(maskFull))
}

func TestMagic020(t *testing.T) {
        expect(t, maskPawn[White][A3], Bitmask( Bit(B2) ))
        expect(t, maskPawn[White][D5], Bitmask( Bit(C4) | Bit(E4) ))
        expect(t, maskPawn[White][F8], Bitmask( Bit(E7) | Bit(G7) ))
        expect(t, maskPawn[Black][H4], Bitmask( Bit(G5) ))
        expect(t, maskPawn[Black][C5], Bitmask( Bit(B6) | Bit(D6) ))
        expect(t, maskPawn[Black][B1], Bitmask( Bit(A2) | Bit(C2) ))
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
        expect(t, maskDiagonal[C4][F7], Bit(A2) | Bit(B3) | Bit(C4) |Bit(D5) | Bit(E6) | Bit(F7) | Bit(G8))
        expect(t, maskDiagonal[F6][H8], maskA1H8)
        expect(t, maskDiagonal[F1][H3], Bit(F1) | Bit(G2) | Bit(H3))
        // Same anti-diagonal.
        expect(t, maskDiagonal[C2][B3], Bit(D1) | Bit(C2) | Bit(B3) | Bit(A4))
        expect(t, maskDiagonal[F3][B7], maskH1A8)
        expect(t, maskDiagonal[H3][D7], Bit(H3) | Bit(G4) | Bit(F5) | Bit(E6) | Bit(D7) | Bit(C8))
        // Edge cases.
        expect(t, maskDiagonal[A2][G4], maskNone) // Random squares.
        expect(t, maskDiagonal[E4][E4], maskNone) // Same square.
}
