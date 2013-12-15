// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`; `fmt`)

func TestTargets010(t *testing.T) {
        Settings.Log = false
        Settings.Fancy = false
        game := NewGame().Setup(`a2,b3,c4,d2,e6,f5,g4,h3`, `Kg8`)
        moves := NewPosition(game, game.pieces, WHITE, Bitmask(0)).Moves()

        expect(t, fmt.Sprintf("%v", moves), `[a2-a3 a2-a4 d2-d3 d2-d4 b3-b4 h3-h4 c4-c5 g4-g5 f5-f6 e6-e7]`)
}

func TestTargets020(t *testing.T) {
        game := NewGame().Setup(`a2,b3,c4,d2,e6,f5,g4,h3`, `a3,b4,c5,e7,f6,g5,h4,Kg8`)
        moves := NewPosition(game, game.pieces, WHITE, Bitmask(0)).Moves()

        expect(t, fmt.Sprintf("%v", moves), `[d2-d3 d2-d4]`)
}

func TestTargets030(t *testing.T) {
        game := NewGame().Setup(`a2,e4,g2`, `b3,f5,f3,h3,Kg8`)
        moves := NewPosition(game, game.pieces, WHITE, Bitmask(0)).Moves()

        expect(t, fmt.Sprintf("%v", moves), `[a2xb3 g2xf3 g2xh3 e4xf5 a2-a3 a2-a4 g2-g3 g2-g4 e4-e5]`)
}

func TestTargets040(t *testing.T) {
        game := NewGame().Setup(`Kf1,Nd2`, `Kc1,c2,c3`)
        moves := NewPosition(game, game.pieces, BLACK, Bitmask(0)).Moves()

        expect(t, fmt.Sprintf("%v", moves), `[kc1xd2 c3xd2 kc1-d1 kc1-b2]`)
}
