package lape

import(
        `fmt`
)

type Position struct {
        pieces  [64]Piece
        board   [3]Bitmask // 0: white, 1: black, 2: both
        attacks *Attack
        layout  map[Piece]*Bitmask
}

func (p *Position)Initialize(g *Game) *Position {
        p.pieces = g.pieces
        p.attacks = g.attacks

        p.layout = make(map[Piece]*Bitmask)
        for piece := Piece(PAWN); piece <= Piece(KING); piece++ {
                p.layout[piece] = new(Bitmask)
                p.layout[piece | 1] = new(Bitmask)
        }

        p.setupBoard()
        return p
}

func (p *Position)Moves(color int) []*Move {
        var moves []*Move

        for side := p.board[color]; !side.IsEmpty(); {
                index := side.FirstSet()
                piece := p.pieces[index]
                moves = append(moves, p.PossibleMoves(index, piece)...)
                side.Clear(index)
        }

        fmt.Printf("%d: %v\n", len(moves), moves)
        return moves
}

func (p *Position)PossibleMoves(index int, piece Piece) []*Move {
        var moves []*Move

        targets := p.attacks.Targets(index, piece, p.board)
        for !targets.IsEmpty() {
                target := targets.FirstSet()
                moves = append(moves, new(Move).Initialize(index, target, piece, p.pieces[target]))
                targets.Clear(target)
        }

        return moves
}


func (p *Position)setupBoard() *Position {
        for i, piece := range p.pieces {
                if piece != 0 {
                        p.layout[piece].Set(i)
                        p.board[piece.Color()].Set(i)
                }
        }
        p.board[2] = p.board[0]
        p.board[2].Combine(p.board[1])

        return p
}
