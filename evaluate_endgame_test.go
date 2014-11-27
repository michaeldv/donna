// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`github.com/michaeldv/donna/expect`; `testing`)

// King with 2+ pawns vs. king.
func TestEndgame000(t *testing.T) {
	p := NewGame(`Ke1,a4,a5`, `M,Ka8`).start()
	score := p.Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame010(t *testing.T) {
	p := NewGame(`Ke1,h4,h6`, `M,Kg8`).start()
	score := p.Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame020(t *testing.T) {
	p := NewGame(`Kh4`, `Kg8,h6,h2`).start()
	score := p.Evaluate()
	expect.True(t, score != 0)
}

func TestEndgame030(t *testing.T) {
	p := NewGame(`Kc4`, `Ka5,a3,a4`).start()
	score := p.Evaluate()
	expect.True(t, score != 0)
}

// No pawns left.
func TestEndgame100(t *testing.T) {
	p := NewGame(`Ke1,Bc1,a2,b2`, `Kd8,Bc8,Nb8`).start()
	score := p.Evaluate()
	expect.True(t, score != 0)
}

func TestEndgame110(t *testing.T) {
	p := NewGame(`Ke1,Bc1`, `Kd8,d5`).start()
	score := p.Evaluate()
	expect.True(t, score == 0)
}

func TestEndgame120(t *testing.T) {
	p := NewGame(`Ke1,Nb1`, `Kd8,a5`).start()
	score := p.Evaluate()
	expect.True(t, score != 0)
}

// KPK bitbase.
func TestEndgame200(t *testing.T) {
	game := NewGame(`Kf1`, `Kh1,h2`)
	white := game.start()
	expect.Eq(t, white.Evaluate(), 0)
}

func TestEndgame201(t *testing.T) {
	game := NewGame(`Kf1`, `M,Kh1,h2`)
	black := game.start()
	expect.Eq(t, black.Evaluate(), 0)
}

func TestEndgame202(t *testing.T) {
	game := NewGame(`Ka8,a7`, `Kc7`)
	white := game.start()
	expect.Eq(t, white.Evaluate(), 0)
}

func TestEndgame203(t *testing.T) {
	game := NewGame(`Ka8,a7`, `M,Kc7`)
	black := game.start()
	expect.Eq(t, black.Evaluate(), 0)
}

func TestEndgame210(t *testing.T) {
	game := NewGame(`Kf4`, `Kh5,h7`)
	white := game.start()
	expect.Eq(t, white.Evaluate(), 0)
}

func TestEndgame211(t *testing.T) {
	game := NewGame(`Kf4`, `M,Kh5,h7`)
	black := game.start()
	expect.Eq(t, black.Evaluate(), 0)
}

func TestEndgame212(t *testing.T) {
	game := NewGame(`Ka5,a2`, `Kc6`)
	white := game.start()
	expect.Eq(t, white.Evaluate(), 0)
}

func TestEndgame213(t *testing.T) {
	game := NewGame(`Ka5,a2`, `M,Kc6`)
	black := game.start()
	expect.Eq(t, black.Evaluate(), 0)
}

func TestEndgame220(t *testing.T) {
	game := NewGame(`Kf6,e6`, `Kf8`)
	white := game.start()
	expect.Eq(t, white.Evaluate(), 0)
}

func TestEndgame221(t *testing.T) {
	game := NewGame(`Kf6,e6`, `M,Kf8`)
	black := game.start()
	expect.Eq(t, black.Evaluate(), WhiteWinning)
}

func TestEndgame222(t *testing.T) {
	game := NewGame(`Kd1`, `Kd3,e3`)
	white := game.start()
	expect.Eq(t, white.Evaluate(), BlackWinning)
}

func TestEndgame223(t *testing.T) {
	game := NewGame(`Kd1`, `M,Kd3,e3`)
	black := game.start()
	expect.Eq(t, black.Evaluate(), 0)
}

func TestEndgame230(t *testing.T) {
	game := NewGame(`Kf6,e6`, `Ke8`)
	white := game.start()
	expect.Eq(t, white.Evaluate(), WhiteWinning)
}

func TestEndgame231(t *testing.T) {
	game := NewGame(`Kf6,e6`, `M,Ke8`)
	black := game.start()
	expect.Eq(t, black.Evaluate(), 0)
}

func TestEndgame232(t *testing.T) {
	game := NewGame(`Ke1`, `Kd3,e3`)
	white := game.start()
	expect.Eq(t, white.Evaluate(), 0)
}

func TestEndgame233(t *testing.T) {
	game := NewGame(`Ke1`, `M,Kd3,e3`)
	black := game.start()
	expect.Eq(t, black.Evaluate(), BlackWinning)
}

func TestEndgame240(t *testing.T) {
	game := NewGame(`Ke6,e4`, `Ke8`)
	white := game.start()
	expect.Eq(t, white.Evaluate(), WhiteWinning)
}

func TestEndgame241(t *testing.T) {
	game := NewGame(`Ke6,e4`, `M,Ke8`)
	black := game.start()
	expect.Eq(t, black.Evaluate(), WhiteWinning)
}

func TestEndgame242(t *testing.T) {
	game := NewGame(`Kd1`, `Kd3,d5`)
	white := game.start()
	expect.Eq(t, white.Evaluate(), BlackWinning)
}

func TestEndgame243(t *testing.T) {
	game := NewGame(`Kd1`, `M,Kd3,d5`)
	black := game.start()
	expect.Eq(t, black.Evaluate(), BlackWinning)
}

func TestEndgame250(t *testing.T) {
	game := NewGame(`Ka1,e4`, `Ka4`)
	white := game.start()
	expect.Eq(t, white.Evaluate(), 0)
}

func TestEndgame251(t *testing.T) {
	game := NewGame(`Ka1,e4`, `M,Ka4`)
	black := game.start()
	expect.Eq(t, black.Evaluate(), 0)
}

func TestEndgame252(t *testing.T) {
	game := NewGame(`Kh5`, `Kh8,d5`)
	white := game.start()
	expect.Eq(t, white.Evaluate(), 0)
}

func TestEndgame253(t *testing.T) {
	game := NewGame(`Kh5`, `M,Kh8,d5`)
	black := game.start()
	expect.Eq(t, black.Evaluate(), 0)
}

func TestEndgame260(t *testing.T) {
	game := NewGame(`Ka1,g4`, `Ka4`)
	white := game.start()
	expect.Eq(t, white.Evaluate(), WhiteWinning)
}

func TestEndgame261(t *testing.T) {
	game := NewGame(`Ka1,g4`, `M,Ka4`)
	black := game.start()
	expect.Eq(t, black.Evaluate(), WhiteWinning)
}

func TestEndgame262(t *testing.T) {
	game := NewGame(`Kh5`, `Kh8,b5`)
	white := game.start()
	expect.Eq(t, white.Evaluate(), BlackWinning)
}

func TestEndgame263(t *testing.T) {
	game := NewGame(`Kh5`, `M,Kh8,b5`)
	black := game.start()
	expect.Eq(t, black.Evaluate(), BlackWinning)
}
