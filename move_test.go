// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`)

// PxQ, NxQ, BxQ, RxQ, QxQ, KxQ
func TestMove000(t *testing.T) {
        game := NewGame().Setup(`Kd6,Qd1,Ra5,Nc3,Bc4,e4`, `Kh8,Qd5`)
        p := game.Start(White)
        expect(t, p.NewMove(E4, D5).value(), 1182) // PxQ
        expect(t, p.NewMove(C3, D5).value(), 1180) // NxQ
        expect(t, p.NewMove(C4, D5).value(), 1178) // BxQ
        expect(t, p.NewMove(A5, D5).value(), 1176) // RxQ
        expect(t, p.NewMove(D1, D5).value(), 1174) // QxQ
        expect(t, p.NewMove(D6, D5).value(), 1172) // KxQ
}

// PxR, NxR, BxR, RxR, QxR, KxR
func TestMove010(t *testing.T) {
        game := NewGame().Setup(`Kd6,Qd1,Ra5,Nc3,Bc4,e4`, `Kh8,Rd5`)
        p := game.Start(White)
        expect(t, p.NewMove(E4, D5).value(), 1150) // PxR
        expect(t, p.NewMove(C3, D5).value(), 1148) // NxR
        expect(t, p.NewMove(C4, D5).value(), 1146) // BxR
        expect(t, p.NewMove(A5, D5).value(), 1144) // RxR
        expect(t, p.NewMove(D1, D5).value(), 1142) // QxR
        expect(t, p.NewMove(D6, D5).value(), 1140) // KxR
}

// PxB, NxB, BxB, RxB, QxB, KxB
func TestMove020(t *testing.T) {
        game := NewGame().Setup(`Kd6,Qd1,Ra5,Nc3,Bc4,e4`, `Kh8,Bd5`)
        p := game.Start(White)
        expect(t, p.NewMove(E4, D5).value(), 1118) // PxB
        expect(t, p.NewMove(C3, D5).value(), 1116) // NxB
        expect(t, p.NewMove(C4, D5).value(), 1114) // BxB
        expect(t, p.NewMove(A5, D5).value(), 1112) // RxB
        expect(t, p.NewMove(D1, D5).value(), 1110) // QxB
        expect(t, p.NewMove(D6, D5).value(), 1108) // KxB
}

// PxN, NxN, BxN, RxN, QxN, KxN
func TestMove030(t *testing.T) {
        game := NewGame().Setup(`Kd6,Qd1,Ra5,Nc3,Bc4,e4`, `Kh8,Nd5`)
        p := game.Start(White)
        expect(t, p.NewMove(E4, D5).value(), 1086) // PxN
        expect(t, p.NewMove(C3, D5).value(), 1084) // NxN
        expect(t, p.NewMove(C4, D5).value(), 1082) // BxN
        expect(t, p.NewMove(A5, D5).value(), 1080) // RxN
        expect(t, p.NewMove(D1, D5).value(), 1078) // QxN
        expect(t, p.NewMove(D6, D5).value(), 1076) // KxN
}

// PxP, NxP, BxP, RxP, QxP, KxP
func TestMove040(t *testing.T) {
        game := NewGame().Setup(`Kd6,Qd1,Ra5,Nc3,Bc4,e4`, `Kh8,d5`)
        p := game.Start(White)
        expect(t, p.NewMove(E4, D5).value(), 1054) // PxP
        expect(t, p.NewMove(C3, D5).value(), 1052) // NxP
        expect(t, p.NewMove(C4, D5).value(), 1050) // BxP
        expect(t, p.NewMove(A5, D5).value(), 1048) // RxP
        expect(t, p.NewMove(D1, D5).value(), 1046) // QxP
        expect(t, p.NewMove(D6, D5).value(), 1044) // KxP
}
