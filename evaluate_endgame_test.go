// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import `testing`

// King with 2+ pawns vs. king.
func TestEndgame000(t *testing.T) {
	p := NewGame(`Ke1,a4,a5`, `Ka8`).Start(Black)
	score, _ := p.EvaluateWithTrace()
	expect(t, score, 0)
}

func TestEndgame010(t *testing.T) {
	p := NewGame(`Ke1,h4,h6`, `Kg8`).Start(Black)
	score, _ := p.EvaluateWithTrace()
	expect(t, score, 0)
}

func TestEndgame020(t *testing.T) {
	p := NewGame(`Kh4`, `Kg8,h6,h2`).Start(White)
	score, _ := p.EvaluateWithTrace()
	expect(t, score != 0, true)
}

func TestEndgame030(t *testing.T) {
	p := NewGame(`Kc4`, `Ka5,a3,a4`).Start(White)
	score, _ := p.EvaluateWithTrace()
	expect(t, score != 0, true)
}
