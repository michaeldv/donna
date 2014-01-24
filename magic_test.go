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
        expect(t, maskPawnAttack[White][A3], Bitmask( Bit(B2) ))
        expect(t, maskPawnAttack[White][D5], Bitmask( Bit(C4) | Bit(E4) ))
        expect(t, maskPawnAttack[White][F8], Bitmask( Bit(E7) | Bit(G7) ))
        expect(t, maskPawnAttack[Black][H4], Bitmask( Bit(G5) ))
        expect(t, maskPawnAttack[Black][C5], Bitmask( Bit(B6) | Bit(D6) ))
        expect(t, maskPawnAttack[Black][B1], Bitmask( Bit(A2) | Bit(C2) ))
}
