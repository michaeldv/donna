// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`; `fmt`)

func TestTargets000(t *testing.T) {
        game := NewGame().InitialPosition()
        moves := game.Start().Moves()

        // Possible moves includes illegal move that causes a check.
        expect(t, fmt.Sprintf("%v", moves), `[Nb1-c3 Ng1-f3 e2-e4 d2-d4 c2-c4 e2-e3 d2-d3 f2-f4 f2-f3 c2-c3 b2-b3 g2-g4 g2-g3 b2-b4 h2-h4 a2-a3 a2-a4 h2-h3 Nb1-a3 Ng1-h3]`)
}

func TestTargets010(t *testing.T) {
        Settings.Log = false
        Settings.Fancy = false
        game := NewGame().Setup(`a2,b3,c4,d2,e6,f5,g4,h3`, `Kg8`)
        moves := game.Start().Moves()

        // Moves should be sorted by relative strength.
        expect(t, fmt.Sprintf("%v", moves), `[h3-h4 a2-a4 e6-e7 a2-a3 g4-g5 f5-f6 b3-b4 c4-c5 d2-d4 d2-d3]`)
}

func TestTargets020(t *testing.T) {
        game := NewGame().Setup(`a2,b3,c4,d2,e6,f5,g4,h3`, `a3,b4,c5,e7,f6,g5,h4,Kg8`)
        moves := game.Start().Moves()

        expect(t, fmt.Sprintf("%v", moves), `[d2-d4 d2-d3]`)
}

func TestTargets030(t *testing.T) {
        game := NewGame().Setup(`a2,e4,g2`, `b3,f5,f3,h3,Kg8`)
        moves := game.Start().Moves()

        expect(t, fmt.Sprintf("%v", moves), `[a2xb3 g2xf3 g2xh3 e4xf5 a2-a4 a2-a3 g2-g4 g2-g3 e4-e5]`)
}
