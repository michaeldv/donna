package donna

import (`testing`; `runtime`; `strings`)

func expect(t *testing.T, actual, expected string) {
        _, file, line, _ := runtime.Caller(1)           // Get the calling file path and line number.
        file = file[strings.LastIndex(file, `/`) + 1:]  // Keep file name only.

        if expected != actual {
                t.Errorf("\r\t\x1B[31m%s line %d\n%s\x1B[0m", file, line, `Expected: ` + expected + "\n  Actual: " + actual)
        } else {
                t.Logf("\r\t\x1B[32m%s line %d: %s\x1B[0m", file, line, actual)
        }
}

