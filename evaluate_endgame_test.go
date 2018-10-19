// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

import(`github.com/michaeldv/donna/expect`; `testing`)

// King with 2+ pawns vs. king.
func TestEndgame000(t *testing.T) {
	score := NewGame(`Ke1,a4,a5`, `M,Ka8`).start().Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame010(t *testing.T) {
	score := NewGame(`Ke1,h4,h6`, `M,Kg8`).start().Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame020(t *testing.T) {
	score := NewGame(`Kh4`, `Kg8,h6,h2`).start().Evaluate()
	expect.True(t, score != 0)
}

func TestEndgame030(t *testing.T) {
	score := NewGame(`Kc4`, `Ka5,a3,a4`).start().Evaluate()
	expect.True(t, score != 0)
}

// No pawns left.
func TestEndgame100(t *testing.T) {
	score := NewGame(`Ke1,Bc1,a2,b2`, `Kd8,Bc8,Nb8`).start().Evaluate()
	expect.True(t, score != 0)
}

func TestEndgame110(t *testing.T) {
	score := NewGame(`Ke1,Bc1`, `Kd8,d5`).start().Evaluate()
	expect.True(t, score == 0)
}

func TestEndgame120(t *testing.T) {
	score := NewGame(`Ke1,Nb1`, `Kd8,a5`).start().Evaluate()
	expect.True(t, score != 0)
}

// KPK bitbase.
func TestEndgame200(t *testing.T) {
	score := NewGame(`Kf1`, `Kh1,h2`).start().Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame201(t *testing.T) {
	score := NewGame(`Kf1`, `M,Kh1,h2`).start().Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame202(t *testing.T) {
	score := NewGame(`Ka8,a7`, `Kc7`).start().Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame203(t *testing.T) {
	score := NewGame(`Ka8,a7`, `M,Kc7`).start().Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame210(t *testing.T) {
	score := NewGame(`Kf4`, `Kh5,h7`).start().Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame211(t *testing.T) {
	score := NewGame(`Kf4`, `M,Kh5,h7`).start().Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame212(t *testing.T) {
	score := NewGame(`Ka5,a2`, `Kc6`).start().Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame213(t *testing.T) {
	score := NewGame(`Ka5,a2`, `M,Kc6`).start().Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame220(t *testing.T) {
	score := NewGame(`Kf6,e6`, `Kf8`).start().Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame221(t *testing.T) {
	score := NewGame(`Kf6,e6`, `M,Kf8`).start().Evaluate()
	expect.Eq(t, score, -WhiteWinning)
}

func TestEndgame222(t *testing.T) {
	score := NewGame(`Kd1`, `Kd3,e3`).start().Evaluate()
	expect.Eq(t, score, BlackWinning)
}

func TestEndgame223(t *testing.T) {
	score := NewGame(`Kd1`, `M,Kd3,e3`).start().Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame230(t *testing.T) {
	score := NewGame(`Kf6,e6`, `Ke8`).start().Evaluate()
	expect.Eq(t, score, WhiteWinning)
}

func TestEndgame231(t *testing.T) {
	score := NewGame(`Kf6,e6`, `M,Ke8`).start().Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame232(t *testing.T) {
	score := NewGame(`Ke1`, `Kd3,e3`).start().Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame233(t *testing.T) {
	score := NewGame(`Ke1`, `M,Kd3,e3`).start().Evaluate()
	expect.Eq(t, score, -BlackWinning)
}

func TestEndgame240(t *testing.T) {
	score := NewGame(`Ke6,e4`, `Ke8`).start().Evaluate()
	expect.Eq(t, score, WhiteWinning)
}

func TestEndgame241(t *testing.T) {
	score := NewGame(`Ke6,e4`, `M,Ke8`).start().Evaluate()
	expect.Eq(t, score, -WhiteWinning)
}

func TestEndgame242(t *testing.T) {
	score := NewGame(`Kd1`, `Kd3,d5`).start().Evaluate()
	expect.Eq(t, score, BlackWinning)
}

func TestEndgame243(t *testing.T) {
	score := NewGame(`Kd1`, `M,Kd3,d5`).start().Evaluate()
	expect.Eq(t, score, -BlackWinning)
}

func TestEndgame250(t *testing.T) {
	score := NewGame(`Ka1,e4`, `Ka4`).start().Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame251(t *testing.T) {
	score := NewGame(`Ka1,e4`, `M,Ka4`).start().Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame252(t *testing.T) {
	score := NewGame(`Kh5`, `Kh8,d5`).start().Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame253(t *testing.T) {
	score := NewGame(`Kh5`, `M,Kh8,d5`).start().Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame260(t *testing.T) {
	score := NewGame(`Ka1,g4`, `Ka4`).start().Evaluate()
	expect.Eq(t, score, WhiteWinning)
}

func TestEndgame261(t *testing.T) {
	score := NewGame(`Ka1,g4`, `M,Ka4`).start().Evaluate()
	expect.Eq(t, score, -WhiteWinning)
}

func TestEndgame262(t *testing.T) {
	score := NewGame(`Kh5`, `Kh8,b5`).start().Evaluate()
	expect.Eq(t, score, BlackWinning)
}

func TestEndgame263(t *testing.T) {
	score := NewGame(`Kh5`, `M,Kh8,b5`).start().Evaluate()
	expect.Eq(t, score, -BlackWinning)
}

// kingAndPawnVsKingAndPawn
func TestEndgame300(t *testing.T) {
	score := NewGame(`Ke1,e4`, `M,Ke8,e5`).start().Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame310(t *testing.T) {
	score := NewGame(`Ke1,a4`, `Kf8,h7`).start().Evaluate()
	expect.Eq(t, score, WhiteWinning)
}

func TestEndgame320(t *testing.T) {
	score := NewGame(`Kc1,a2`, `M,Ke6,h5`).start().Evaluate()
	expect.Eq(t, score, -BlackWinning)
}

func TestEndgame330(t *testing.T) {
	score := NewGame(`Kd1,h2`, `Ka8,a4`).start().Evaluate()
	expect.Eq(t, score, WhiteWinning)
}

func TestEndgame340(t *testing.T) {
	score := NewGame(`Kd1,h3`, `M,Kb8,a4`).start().Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame350(t *testing.T) {
	score := NewGame(`Kf5,h3`, `Kd5,h4`).start().Evaluate()
	expect.Eq(t, score, 30)
}

func TestEndgame360(t *testing.T) {
	score := NewGame(`Kf5,h3`, `M,Kd5,h4`).start().Evaluate()
	expect.Eq(t, score, -10)
}

func TestEndgame370(t *testing.T) {
	score := NewGame(`Kh6,h3`, `Kf6,h4`).start().Evaluate()
	expect.Eq(t, score, 0)
}

func TestEndgame380(t *testing.T) {
	score := NewGame(`Kf1,h3`, `M,Kh1,h4`).start().Evaluate()
	expect.Eq(t, score, 0)
}

// Opposite-colored donkeys.

// Don't drop a score when a side has extra piece.
func TestEndgame400(t *testing.T) {
	score := NewGame(`Ke1,Bf1,Nc3,Nf3,f2,g3,h4`, `Ke8,Bf8,Nf6,f7,g6,h5`).start().Evaluate()
	expect.Eq(t, score, 444) // Extra Nc3.

	score = NewGame(`Ke1,Bf1,Nf3,f2,g3,h4`, `Ke8,Ra8,Bf8,Nf6,f7,g6,h5`).start().Evaluate()
	expect.Eq(t, score, -709) // Extra Ra8 for black.
}

// Don't drop a score when a side has more that 2 extra pawns.
func TestEndgame410(t *testing.T) {
	score := NewGame(`Ke1,Bf1,Nf3,a2,b2,f2,g3,h4`, `Ke8,Bf8,Nf6,f7,g6,h5`).start().Evaluate()
	expect.Eq(t, score, 95) // Extra a2,b2 pawns, drop the score (266 -> 95).

	score = NewGame(`Ke1,Bf1,Nf3,f2,g3,h4`, `Ke8,Bf8,Nf6,a7,b7,c7,f7,g6,h5`).start().Evaluate()
	expect.Eq(t, score, -373) // Extra a7,b7,c7 for black, don't drop the score.
}

// Draw if single passer and a king blocks it on safe color square.
func TestEndgame420(t *testing.T) {
	score := NewGame(`Ke1,Bf1,a6`, `Ka7,Bf8`).start().Evaluate() // King blocks on a7.
	expect.Eq(t, score, 0)

	score = NewGame(`Kf6,Be2,e7`, `Ke8,Bf2`).start().Evaluate() // King on e8 is not blocking (Bh5+).
	expect.Eq(t, score, 206)
}

// Draw if single passer and a bishop controls a square in front of it.
func TestEndgame430(t *testing.T) {
	score := NewGame(`Ke1,Bg3`, `Ke8,Bc8,h3`).start().Evaluate() // Bg3 controls h2.
	expect.Eq(t, score, 0)

	score = NewGame(`Kd6,Bb8`, `M,Ke8,Bc8,h3`).start().Evaluate() // Bb8 is blocked by Kd6 and doesn't control h2.
	expect.Eq(t, score, 295)
}
