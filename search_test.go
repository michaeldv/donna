// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`)

func TestSearch000(t *testing.T) {
        position := NewGame().InitialPosition().Start(White)
        expect(t, position.Perft(0), int64(1))
}

func TestSearch010(t *testing.T) {
        position := NewGame().InitialPosition().Start(White)
        expect(t, position.Perft(1), int64(20))
}

func TestSearch020(t *testing.T) {
        position := NewGame().InitialPosition().Start(White)
        expect(t, position.Perft(2), int64(400))
}

func TestSearch030(t *testing.T) {
        position := NewGame().InitialPosition().Start(White)
        expect(t, position.Perft(3), int64(8902))
}

func TestSearch040(t *testing.T) {
        position := NewGame().InitialPosition().Start(White)
        expect(t, position.Perft(4), int64(197281))
}

func TestSearch050(t *testing.T) {
        position := NewGame().InitialPosition().Start(White)
        expect(t, position.Perft(5), int64(4865609))
}