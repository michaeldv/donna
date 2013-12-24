// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna
import (`testing`; `fmt`)
//
// See http://cinnamonchess.altervista.org/bitboard_calculator/Calc.html
//
// Targets of white king on E1 include G1 and C1.
func TestTargets010(t *testing.T) {
        game := NewGame().Setup(`Ke1`, `Ke8`)
        position := game.Start()

        expect(t, fmt.Sprintf(`0x%016X`, uint64(position.targets[E1])), `0x000000000000386C`)
}

// Targets of black king on E1 *do not* include G1 and C1.
func TestTargets015(t *testing.T) {
        game := NewGame().Setup(`Ke8`, `Ke1`)
        game.current = BLACK
        position := game.Start()

        expect(t, fmt.Sprintf(`0x%016X`, uint64(position.targets[E1])), `0x0000000000003828`)
}

// Targets of black king on E8 include G8 and C8.
func TestTargets020(t *testing.T) {
        game := NewGame().Setup(`Ke1`, `Ke8`)
        game.current = BLACK
        position := game.Start()

        expect(t, fmt.Sprintf(`0x%016X`, uint64(position.targets[E8])), `0x6C38000000000000`)
}

// Targets of white king on E8 *do not* include G8 and C8.
func TestTargets025(t *testing.T) {
        game := NewGame().Setup(`Ke8`, `Ke1`)
        position := game.Start()

        expect(t, fmt.Sprintf(`0x%016X`, uint64(position.targets[E8])), `0x2838000000000000`)
}

// King on G1.
func TestTargets030(t *testing.T) {
        game := NewGame().Setup(`Kg1`, `Ke8`)
        position := game.Start()

        expect(t, fmt.Sprintf(`0x%016X`, uint64(position.targets[G1])), `0x000000000000E0A0`)
}

// King on H1.
func TestTargets035(t *testing.T) {
        game := NewGame().Setup(`Kh1`, `Ke8`)
        position := game.Start()

        expect(t, fmt.Sprintf(`0x%016X`, uint64(position.targets[H1])), `0x000000000000C040`)
}

// King on C8.
func TestTargets040(t *testing.T) {
        game := NewGame().Setup(`Kd1`, `Kc8`)
        game.current = BLACK
        position := game.Start()

        expect(t, fmt.Sprintf(`0x%016X`, uint64(position.targets[C8])), `0x0A0E000000000000`)
}

// King on D4.
func TestTargets050(t *testing.T) {
        game := NewGame().Setup(`Kd4`, `Ke8`)
        position := game.Start()

        expect(t, fmt.Sprintf(`0x%016X`, uint64(position.targets[D4])), `0x0000001C141C0000`)
}
