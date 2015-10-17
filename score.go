// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

type Score struct {
	midgame int
	endgame int
}

// Reference methods that change the score receiver in place and return a
// pointer to the updated score.
func (s *Score) clear() *Score {
	s.midgame, s.endgame = 0, 0

	return s
}

func (s *Score) add(score Score) *Score {
	s.midgame += score.midgame
	s.endgame += score.endgame

	return s
}

func (s *Score) sub(score Score) *Score {
	s.midgame -= score.midgame
	s.endgame -= score.endgame

	return s
}

func (s *Score) apply(weight Score) *Score {
	s.midgame = s.midgame * weight.midgame / 100
	s.endgame = s.endgame * weight.endgame / 100

	return s
}


func (s *Score) adjust(n int) *Score {
	s.midgame += n
	s.endgame += n

	return s
}

func (s *Score) scale(n int) *Score {
	s.midgame = s.midgame * n / 100
	s.endgame = s.endgame * n / 100

	return s
}

// Value methods that return newly updated score value.
func (s Score) plus(score Score) Score {
	s.midgame += score.midgame
	s.endgame += score.endgame

	return s
}

func (s Score) minus(score Score) Score {
	s.midgame -= score.midgame
	s.endgame -= score.endgame

	return s
}

func (s Score) times(n int) Score {
	s.midgame *= n
	s.endgame *= n

	return s
}

// Calculates normalized score based on the game phase.
func (s Score) blended(phase int) int {
	return (s.midgame * phase + s.endgame * (256 - phase)) / 256
}
