// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`; `fmt`; `runtime`; `strings`)

func expect(t *testing.T, actual, expected interface{}) {
        var passed bool

        _, file, line, _ := runtime.Caller(1)           // Get the calling file path and line number.
        file = file[strings.LastIndex(file, `/`) + 1:]  // Keep file name only.

        switch expected.(type) {
        case bool:
                if assertion, ok := actual.(bool); ok {
                        passed = (assertion == expected)
                }
        case int:
                if assertion, ok := actual.(int); ok {
                        passed = (assertion == expected)
                }
        case uint64:
                if assertion, ok := actual.(uint64); ok {
                        passed = (assertion == expected)
                }
        default:
                passed = (fmt.Sprintf(`%v`, actual) == fmt.Sprintf(`%v`, expected))
        }

        if passed {
                t.Logf("\r\t\x1B[32m%s line %d: %v\x1B[0m", file, line, actual)
        } else {
                t.Errorf("\r\t\x1B[31m%s line %d\nExpected: %v\n  Actual: %v\x1B[0m", file, line, expected, actual)
        }
}

func contains(t *testing.T, actual interface{}, expected string) {
        containsMatcher(t, fmt.Sprintf(`%v`, actual), expected, true)
}

func doesNotContain(t *testing.T, actual interface{}, expected string) {
        containsMatcher(t, fmt.Sprintf(`%v`, actual), expected, false)
}

func containsMatcher(t *testing.T, actual, expected string, match bool) {
        var passed bool

        _, file, line, _ := runtime.Caller(2)
        file = file[strings.LastIndex(file, `/`) + 1:]

        if match {
                passed = (actual == expected)
        } else {
                passed = (actual != expected)
        }

        if passed {
                t.Logf("\r\t\x1B[32m%s line %d: %v\x1B[0m", file, line, actual)
        } else {
                t.Errorf("\r\t\x1B[31m%s line %d\nContains: %s\n  Actual: %s\x1B[0m", file, line, expected, actual)
        }
}

func TestBitmask000(t *testing.T) { // White
        passed := [8]Bitmask{0}
        for square := A2; square <= H2; square++ {
                i := square - A2
                if Col(square) > 0 {
                        passed[i].Fill(square - 1, 8, 0, 0x00FFFFFFFFFFFFFF)
                }
                passed[i].Fill(square, 8, 0, 0x00FFFFFFFFFFFFFF)
                if Col(square) < 7 {
                        passed[i].Fill(square + 1, 8, 0, 0x00FFFFFFFFFFFFFF)
                }
        }
        expect(t, passed[0], Bitmask(0x0303030303030000))
        expect(t, passed[1], Bitmask(0x0707070707070000))
        expect(t, passed[2], Bitmask(0x0E0E0E0E0E0E0000))
        expect(t, passed[3], Bitmask(0x1C1C1C1C1C1C0000))
        expect(t, passed[4], Bitmask(0x3838383838380000))
        expect(t, passed[5], Bitmask(0x7070707070700000))
        expect(t, passed[6], Bitmask(0xE0E0E0E0E0E00000))
        expect(t, passed[7], Bitmask(0xC0C0C0C0C0C00000))
}

func TestBitmask010(t *testing.T) { // Black
        passed := [8]Bitmask{0}
        for square := A7; square <= H7; square++ {
                i := square - A7
                if Col(square) > 0 {
                        passed[i].Fill(square - 1, -8, 0, 0xFFFFFFFFFFFFFF00)
                }
                passed[i].Fill(square, -8, 0, 0xFFFFFFFFFFFFFF00)
                if Col(square) < 7 {
                        passed[i].Fill(square + 1, -8, 0, 0xFFFFFFFFFFFFFF00)
                }
        }
        expect(t, passed[0], Bitmask(0x0000030303030303))
        expect(t, passed[1], Bitmask(0x0000070707070707))
        expect(t, passed[2], Bitmask(0x00000E0E0E0E0E0E))
        expect(t, passed[3], Bitmask(0x00001C1C1C1C1C1C))
        expect(t, passed[4], Bitmask(0x0000383838383838))
        expect(t, passed[5], Bitmask(0x0000707070707070))
        expect(t, passed[6], Bitmask(0x0000E0E0E0E0E0E0))
        expect(t, passed[7], Bitmask(0x0000C0C0C0C0C0C0))
}
