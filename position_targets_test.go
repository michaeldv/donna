// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna
import (`testing`)
//
// See http://cinnamonchess.altervista.org/bitboard_calculator/Calc.html
//
// Targets of white king on E1 include G1 and C1 if castles are allowed.
func TestTargets010(t *testing.T) {
        game := NewGame().Setup(`Ke1,Ra1,Rh1`, `Ke8`)
        position := game.Start(White)

        expect(t, uint64(position.targets[E1]), uint64(0x000000000000386C))
}

// Targets of white king on E1 *do not* include G1 and C1 if castles *are not* allowed.
func TestTargets011(t *testing.T) {
        game := NewGame().Setup(`Ke1`, `Ke8`)
        position := game.Start(White)

        expect(t, uint64(position.targets[E1]), uint64(0x0000000000003828))
}

// No G1 target (white bishop block).
func TestTargets012(t *testing.T) {
        game := NewGame().Setup(`Ke1,Rh1,Bf1`, `Ke8`)
        position := game.Start(White)

        expect(t, uint64(position.targets[E1]), uint64(0x0000000000003808))
}

// No G1 target (F1 under attack).
func TestTargets013(t *testing.T) {
        game := NewGame().Setup(`Ke1,Rh1`, `Ke8,Rf8`)
        position := game.Start(White)

        expect(t, uint64(position.targets[E1]), uint64(0x0000000000003828))
}

// No C1 target (D1 under attack)
func TestTargets014(t *testing.T) {
        game := NewGame().Setup(`Ke1,Ra1`, `Ke8,Rd8`)
        position := game.Start(White)

        expect(t, uint64(position.targets[E1]), uint64(0x0000000000003828))
}

// Targets of black king on E1 *do not* include G1 and C1.
func TestTargets015(t *testing.T) {
        game := NewGame().Setup(`Ke8`, `Ke1`)
        position := game.Start(Black)

        expect(t, uint64(position.targets[E1]), uint64(0x0000000000003828))
}

// Targets of black king on E8 include G8 and C8 if castles are allowed.
func TestTargets020(t *testing.T) {
        game := NewGame().Setup(`Ke1`, `Ke8,Ra8,Rh8`)
        position := game.Start(Black)

        expect(t, uint64(position.targets[E8]), uint64(0x6C38000000000000))
}

// Targets of black king on E8 *do not* include G8 and C8 if castles *are not* allowed.
func TestTargets021(t *testing.T) {
        game := NewGame().Setup(`Ke1`, `Ke8`)
        position := game.Start(Black)

        expect(t, uint64(position.targets[E8]), uint64(0x2838000000000000))
}

// No C8 target (black knight block).
func TestTargets022(t *testing.T) {
        game := NewGame().Setup(`Ke1`, `Ke8,Ra8,Nb8`)
        position := game.Start(Black)

        expect(t, uint64(position.targets[E8]), uint64(0x2838000000000000))
}

// No C8 target (D8 under attack).
func TestTargets023(t *testing.T) {
        game := NewGame().Setup(`Ke1,Rd1`, `Ke8,Ra8`)
        position := game.Start(Black)

        expect(t, uint64(position.targets[E8]), uint64(0x2838000000000000))
}

// No G8 target (F8 under attack)
func TestTargets024(t *testing.T) {
        game := NewGame().Setup(`Ke1,Rf1`, `Ke8,Rh8`)
        position := game.Start(Black)

        expect(t, uint64(position.targets[E8]), uint64(0x2838000000000000))
}

// Targets of white king on E8 *do not* include G8 and C8.
func TestTargets025(t *testing.T) {
        game := NewGame().Setup(`Ke8`, `Ke1`)
        position := game.Start(White)

        expect(t, uint64(position.targets[E8]), uint64(0x2838000000000000))
}

// King on G1.
func TestTargets030(t *testing.T) {
        game := NewGame().Setup(`Kg1`, `Ke8`)
        position := game.Start(White)

        expect(t, uint64(position.targets[G1]), uint64(0x000000000000E0A0))
}

// King on H1.
func TestTargets035(t *testing.T) {
        game := NewGame().Setup(`Kh1`, `Ke8`)
        position := game.Start(White)

        expect(t, uint64(position.targets[H1]), uint64(0x000000000000C040))
}

// King on C8.
func TestTargets040(t *testing.T) {
        game := NewGame().Setup(`Kd1`, `Kc8`)
        position := game.Start(Black)

        expect(t, uint64(position.targets[C8]), uint64(0x0A0E000000000000))
}

// King on D4.
func TestTargets050(t *testing.T) {
        game := NewGame().Setup(`Kd4`, `Ke8`)
        position := game.Start(White)

        expect(t, uint64(position.targets[D4]), uint64(0x0000001C141C0000))
}

// Pawn targets.
func TestTargets100(t *testing.T) {
        game := NewGame().Setup(`Kd1,e2`, `Ke8,d4`)
        position := game.Start(White)

        expect(t, uint64(position.targets[E2]), uint(Bit(E3)|Bit(E4)))          // e3,e4
        expect(t, uint64(position.targets[D4]), uint(Bit(D3)))                  // d3
}

func TestTargets110(t *testing.T) {
        game := NewGame().Setup(`Kd1,e2,d3`, `Ke8,d4,e4`)
        position := game.Start(White)

        expect(t, uint64(position.targets[E2]), uint64(Bit(E3)))                // e3
        expect(t, uint64(position.targets[D3]), uint64(Bit(E4)))                // e4
        expect(t, uint64(position.targets[D4]), uint64(0))                      // None.
        expect(t, uint64(position.targets[E4]), uint64(Bit(D3)|Bit(E3)))        // d3,e3
}

func TestTargets120(t *testing.T) {
        game := NewGame().Setup(`Kd1,e2`, `Ke8,d3,f3`)
        position := game.Start(White)

        expect(t, uint64(position.targets[E2]), uint64(Bit(D3)|Bit(E3)|Bit(E4)|Bit(F3))) // d3,e3,e4,f3
        expect(t, uint64(position.targets[D3]), uint64(Bit(E2)|Bit(D2)))        // e2,d2
        expect(t, uint64(position.targets[F3]), uint64(Bit(E2)|Bit(F2)))        // e2,f2
}

func TestTargets130(t *testing.T) {
        game := NewGame().Setup(`Kd1,e2`, `Ke8,d4`)
        position := game.Start(White)
        position = position.MakeMove(NewEnpassant(position, E2, E4))            // Creates en-passant on e3.

        expect(t, uint64(position.targets[E4]), uint64(Bit(E5)))                // e5
        expect(t, uint64(position.targets[D4]), uint64(Bit(D3)|Bit(E3)))        // d3, e3 (en-passant).
}
