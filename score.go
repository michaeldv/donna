// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

type Score struct {
	midgame int
	endgame int
}

func (s *Score) clear() *Score {
	s.midgame, s.endgame = 0, 0

	return s
}

func (s *Score) add(score Score) *Score {
	s.midgame += score.midgame
	s.endgame += score.endgame

	return s
}

func (s *Score) subtract(score Score) *Score {
	s.midgame -= score.midgame
	s.endgame -= score.endgame

	return s
}

func (s *Score) adjust(n int) *Score {
	s.midgame += n
	s.endgame += n

	return s
}

func (s Score) times(n int) Score {
	s.midgame *= n
	s.endgame *= n

	return s
}

// Calculates normalized score based on the game phase.
func (s Score) blended(phase int) int {
	return (s.midgame * phase + s.endgame * (256-phase)) / 256
}
