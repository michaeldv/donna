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
                t.Logf("\r\t\x1B[32m%s line %d: %s\x1B[0m", file, line, actual)
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
                t.Logf("\r\t\x1B[32m%s line %d: %s\x1B[0m", file, line, actual)
        } else {
                t.Errorf("\r\t\x1B[31m%s line %d\nContains: %s\n  Actual: %s\x1B[0m", file, line, expected, actual)
        }
}
