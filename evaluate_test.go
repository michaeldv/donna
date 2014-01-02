// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`)

func TestEvaluate000(t *testing.T) {
        game := NewGame().InitialPosition()
        score := game.Start().Evaluate()
        expect(t, score, 0)
}

func TestEvaluate010(t *testing.T) { // After 1. e2-e4
        game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
                                `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e7,f7,g7,h7`)
        score := game.Start().Evaluate()
        expect(t, score, 5)
}

func TestEvaluate020(t *testing.T) { // After 1. e2-e4 e7-e5
        game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
                                `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`)
        score := game.Start().Evaluate()
        expect(t, score, 0)
}

func TestEvaluate030(t *testing.T) { // After 1. e2-e4 e7-e5 2. Ng1-f3
        game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Nf3,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
                                `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`)
        score := game.Start().Evaluate()
        expect(t, score, 55)
}

func TestEvaluate040(t *testing.T) { // After 1. e2-e4 e7-e5 2. Ng1-f3 Nb8-c6
        game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Nf3,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
                                `Ra8,Nc6,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`)
        score := game.Start().Evaluate()
        expect(t, score, 10)
}

func TestEvaluate050(t *testing.T) { // After 1. e2-e4 e7-e5 2. Ng1-f3 Nb8-c6 3. Nb1-c3 Ng8-f6
        game := NewGame().Setup(`Ra1,Nc3,Bc1,Qd1,Ke1,Bf1,Nf3,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
                                `Ra8,Nc6,Bc8,Qd8,Ke8,Bf8,Nf6,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`)
        score := game.Start().Evaluate()
        expect(t, score, 0)
}

// Doubled pawns.
func TestEvaluate100(t *testing.T) {
        game := NewGame().Setup(`Ke1,h2,h3`, `Ke8,a7,a6`)
        score := game.Start().Evaluate()

        expect(t, score, 0)
}

func TestEvaluate110(t *testing.T) {
        game := NewGame().Setup(`Ke1,h2,h3`, `Ke8,a7,h7`)
        score := game.Start().Evaluate()

        expect(t, score, -59)
}

func TestEvaluate120(t *testing.T) {
        game := NewGame().Setup(`Ke1,f4,f5`, `Ke8,a7,h7`)
        score := game.Start().Evaluate()

        expect(t, score, -49)
}
