package lape

import(
        `fmt`
)

type Position struct {
        pieces  [64]Piece
        board   Bitmask
        attacks *Bitboard
        sides   [2]Bitmask
        kings   [2]Bitmask
        queens  [2]Bitmask
        rooks   [2]Bitmask
        bishops [2]Bitmask
        knights [2]Bitmask
        pawns   [2]Bitmask
}

func (p *Position)Initialize(g *Game) *Position {
        p.pieces = g.pieces
        p.attacks = g.attacks

        p.setupBoard()
        return p
}

func (p *Position)Moves(color int) []*Move {
        var moves []*Move

        for side := p.sides[color]; !side.IsEmpty(); {
                index := side.FirstSet()
                piece := p.pieces[index]
                moves = append(moves, p.PossibleMoves(piece, index)...)
                side.Clear(index)
        }

        fmt.Printf("%d: %v\n", len(moves), moves)
        return moves
}

func (p *Position)PossibleMoves(piece Piece, index int) []*Move {
        var moves []*Move

        if piece.IsKnight() {
                attacks := p.attacks.Knight[index]
                attacks.Exclude(p.sides[piece.Color()])
                for !attacks.IsEmpty() {
                        target := attacks.FirstSet()
                        moves = append(moves, new(Move).Initialize(index, target, piece, p.pieces[target]))
                        attacks.Clear(target)
                }
        }

        return moves
}


func (p *Position)setupBoard() *Position {
        for i, piece := range p.pieces {
                if piece == 0 {
                        continue
                }
                color := piece.Color()
                p.sides[color].Set(i)

                switch piece.Type() {
                case KING:
                        p.kings[color].Set(i)
                case QUEEN:
                        p.queens[color].Set(i)
                case ROOK:
                        p.rooks[color].Set(i)
                case BISHOP:
                        p.bishops[color].Set(i)
                case KNIGHT:
                        p.knights[color].Set(i)
                case PAWN:
                        p.pawns[color].Set(i)
                }
        }
        p.board = p.sides[0]
        p.board.Combine(p.sides[1])

        return p
}
