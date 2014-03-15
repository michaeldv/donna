// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`)

// The very first chess puzzle I had solved as a kid.
func TestSearch000(t *testing.T) {
        move := NewGame().Setup(`Kf8,Rh1,g6`, `Kh8,Bg8,g7,h7`).Start(White).xSearch()
        expect(t, move, `Rh1-h6`)
}

// func TestSearch020(t *testing.T) {
//         move := NewGame().Setup(`Kf4,Qc2,Nc5`, `Kd4`).Start(White).xSearch()
//         expect(t, move, `Nc5-b7`)
// }
