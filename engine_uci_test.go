// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
	`github.com/michaeldv/donna/expect`
	`testing`
	`io/ioutil`
	`os`
	`syscall`
)

// Mocks os.Stdin by redirecting standard input to read data from a temporary
// file we create.
func mockStdin(input string) (string, error) {
	// Create temporary file with read/write access.
	f, err := ioutil.TempFile(``, `donna`)
	if err != nil {
		return ``, err
	}

	// Save the file name and write input string to the file.
	mock := f.Name()
	f.WriteString(input)
	f.Close()

	// Reopen the file in read-only mode.
	f, err = os.Open(mock)
	if err != nil {
		return mock, err
	}
	defer f.Close()

	// Redirect os.Stdin (fd=0) to read from the file.
	syscall.Dup2(int(f.Fd()), int(os.Stdin.Fd()))

	return mock, nil
}

func restoreStdin(mock string) {
	os.Stdin = os.NewFile(uintptr(syscall.Stdin), `/dev/stdin`)
	if mock != `` {
		os.Remove(mock)
	}
}

func TestUci000(t *testing.T) {
	mock, err := mockStdin("go movetime 1234\nquit\n")
	defer restoreStdin(mock)

	if err != nil {
		t.Errorf(err.Error())
	} else {
		engine := NewEngine().Uci()
		expect.Eq(t, engine.options.moveTime, int64(1234))
		expect.Eq(t, engine.options.timeLeft, int64(0))
		expect.Eq(t, engine.options.timeInc, int64(0))
	}
}

func TestUci010(t *testing.T) {
	mock, err := mockStdin("position startpos\ngo wtime 12345 btime 98765 movestogo 42\nquit\n")
	defer restoreStdin(mock)

	if err != nil {
		t.Errorf(err.Error())
	} else {
		engine := NewEngine().Uci()
		expect.Eq(t, engine.options.timeLeft, int64(12345))
		expect.Eq(t, engine.options.moveTime, int64(0))
		expect.Eq(t, engine.options.timeInc, int64(0))
		expect.Eq(t, engine.options.movesToGo, 42)
	}
}

func TestUci020(t *testing.T) {
	mock, err := mockStdin("position startpos moves e2e4\ngo wtime 12345 btime 98765 movestogo 42\nquit\n")
	defer restoreStdin(mock)

	if err != nil {
		t.Errorf(err.Error())
	} else {
		engine := NewEngine().Uci()
		expect.Eq(t, engine.options.timeLeft, int64(98765))
		expect.Eq(t, engine.options.moveTime, int64(0))
		expect.Eq(t, engine.options.timeInc, int64(0))
		expect.Eq(t, engine.options.movesToGo, 42)
	}
}
