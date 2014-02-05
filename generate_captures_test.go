// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`)

// Piece captures.
func TestGenCaptures000(t *testing.T) {
        game := NewGame().Setup(`Ka1,Qd1,Rh1,Bb3,Ne5`, `Ka8,Qd8,Rh8,Be6,Ng6`)
        white := game.Start(White).StartMoveGen(1).pieceCaptures(White)
        expect(t, white.allMoves(), `[Ne5xg6 Bb3xe6 Rh1xh8 Qd1xd8]`)

        black := game.Start(Black).StartMoveGen(1).pieceCaptures(Black)
        expect(t, black.allMoves(), `[Ng6xe5 Be6xb3 Rh8xh1 Qd8xd1]`)
}
