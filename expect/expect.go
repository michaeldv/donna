// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package expect

import(`fmt`; `runtime`; `strings`; `testing`)

func Eq(t *testing.T, actual, expected interface{}) {
	log(t, actual, expected, equal(actual, expected))
}

func True(t *testing.T, actual interface{}) {
	log(t, actual, true, equal(actual, true))
}

func False(t *testing.T, actual interface{}) {
	log(t, actual, false, equal(actual, false))
}

func Contain(t *testing.T, actual interface{}, expected string) {
	match(t, fmt.Sprintf(`%v`, actual), expected, true)
}

func NotContain(t *testing.T, actual interface{}, expected string) {
	match(t, fmt.Sprintf(`%v`, actual), expected, false)
}

func equal(actual, expected interface{}) (passed bool) {
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
	return
}

// Simple success/failure logger that assumes source test file is at Caller(2).
func log(t *testing.T, actual, expected interface{}, passed bool) {
	_, file, line, _ := runtime.Caller(2)          // Get the calling file path and line number.
	file = file[strings.LastIndex(file, `/`) + 1:] // Keep file name only.

	if !passed {
		t.Errorf("\r\t\x1B[31m%s line %d\nExpected: %v\n  Actual: %v\x1B[0m", file, line, expected, actual)
	} else if (testing.Verbose()) {
		t.Logf("\r\t\x1B[32m%s line %d: %v\x1B[0m", file, line, actual)
	}
}

func match(t *testing.T, actual, expected string, contains bool) {
	passed := (contains == strings.Contains(actual, expected))

	_, file, line, _ := runtime.Caller(2)
	file = file[strings.LastIndex(file, `/`) + 1:]

	if !passed {
		t.Errorf("\r\t\x1B[31m%s line %d\nContains: %s\n  Actual: %s\x1B[0m", file, line, expected, actual)
	} else if (testing.Verbose()) {
		t.Logf("\r\t\x1B[32m%s line %d: %v\x1B[0m", file, line, actual)
	}
}
