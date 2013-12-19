// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`fmt`; `testing`)

func TestEvaluate000(t *testing.T) {
        game := NewGame().InitialPosition()
        position := NewPosition(game, game.pieces, WHITE, 0)
        score := position.Evaluate()
        expect(t, fmt.Sprintf(`%d`, score), `0`)
}

func TestEvaluate010(t *testing.T) { // After 1. e2-e4
        game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
                                `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e7,f7,g7,h7`)
        position := NewPosition(game, game.pieces, WHITE, 0)
        score := position.Evaluate()
        expect(t, fmt.Sprintf(`%d`, score), `43`)
}

func TestEvaluate020(t *testing.T) { // After 1. e2-e4 e7-e5
        game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
                                `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`)
        position := NewPosition(game, game.pieces, WHITE, 0)
        score := position.Evaluate()
        expect(t, fmt.Sprintf(`%d`, score), `0`)
}

func TestEvaluate030(t *testing.T) { // After 1. e2-e4 e7-e5 2. Ng1-f3
        game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Nf3,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
                                `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`)
        position := NewPosition(game, game.pieces, WHITE, 0)
        score := position.Evaluate()
        expect(t, fmt.Sprintf(`%d`, score), `40`)
}

func TestEvaluate040(t *testing.T) { // After 1. e2-e4 e7-e5 2. Ng1-f3 Nb8-c6
        game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Nf3,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
                                `Ra8,Nc6,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`)
        position := NewPosition(game, game.pieces, WHITE, 0)
        score := position.Evaluate()
        expect(t, fmt.Sprintf(`%d`, score), `10`)
}
