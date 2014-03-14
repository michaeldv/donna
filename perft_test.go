// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`)

func TestPerft000(t *testing.T) {
        position := NewGame().InitialPosition().Start(White)
        expect(t, position.Perft(0), int64(1))
}

func TestPerft010(t *testing.T) {
        position := NewGame().InitialPosition().Start(White)
        expect(t, position.Perft(1), int64(20))
}

func TestPerft020(t *testing.T) {
        position := NewGame().InitialPosition().Start(White)
        expect(t, position.Perft(2), int64(400))
}

func TestPerft030(t *testing.T) {
        position := NewGame().InitialPosition().Start(White)
        expect(t, position.Perft(3), int64(8902))
}

func TestPerft040(t *testing.T) {
        position := NewGame().InitialPosition().Start(White)
        expect(t, position.Perft(4), int64(197281))
}

func TestPerft050(t *testing.T) {
        position := NewGame().InitialPosition().Start(White)
        expect(t, position.Perft(5), int64(4865609))
}