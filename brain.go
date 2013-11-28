package lape

import (
        `fmt`
        `math`
)

const CENTER = 0x0000001818000000 // 4 central squares
const EXTENDED_CENTER = 0x00003C3C3C3C0000 // 12 central squares

type Brain struct {
        player *Player
        color  int
}

func (b *Brain) Initialize(player *Player) *Brain {
        b.player = player
        b.color = player.Color

        return b
}

func (b *Brain) Evaluate(p *Position) float64 {
        material := b.materialBalance(p)
        mobility := b.mobilityBalance(p)
        aggression := b.aggressionBalance(p)
        center := b.centerBoost(p)
        fmt.Printf("Score for %s is %.2f (mat: %.2f, mob: %.2f, agg: %.2f, ctr: %.2f)\n", C(b.color), math.Abs(material + mobility + aggression + center), material, mobility, aggression, center)
        return math.Abs(material + mobility + aggression + center)
}

func (b *Brain) materialBalance(p *Position) float64 {
        opposite := b.color^1

        score := 1000 * (p.count[King(b.color)] - p.count[King(opposite)]) +
                    9 * (p.count[Queen(b.color)] - p.count[Queen(opposite)]) +
                    5 * (p.count[Rook(b.color)] - p.count[Rook(opposite)]) +
                    3 * (p.count[Bishop(b.color)] - p.count[Bishop(opposite)]) +
                    3 * (p.count[Knight(b.color)] - p.count[Knight(opposite)]) +
                    1 * (p.count[Pawn(b.color)] - p.count[Pawn(opposite)])

        return float64(score) + 0.1 * float64(p.count[Bishop(b.color)] - p.count[Bishop(opposite)])
}

func (b *Brain) mobilityBalance(p *Position) float64 {
        return 0.25 * float64(b.movesAvailable(p, b.color) - b.movesAvailable(p, b.color^1))
}

func (b *Brain) aggressionBalance(p *Position) float64 {
        return 0.20 * float64(b.attacksAvailable(p, b.color) - b.attacksAvailable(p, b.color^1))
}

// How many attacks for the central squares?
func (b *Brain) centerBoost(p *Position) (center float64) {
        for i, piece := range p.pieces {
                if piece != 0 && piece.Color() == b.color {
                        targets := p.targets[i]
                        sq12 := targets.Intersect(EXTENDED_CENTER).Count()
                        sq04 := targets.Intersect(CENTER).Count()
                        center += 0.05 * float64(sq12 - sq04) + 0.3 * float64(sq04)
                }
        }
        return
}

// Number of moves available for all pieces of certain color.
func (b *Brain) movesAvailable(p *Position, color int) (moves int) {
        for i, piece := range p.pieces {
                if piece != 0 && piece.Color() == color {
                        moves += p.targets[i].Count()
                }
        }
        return
}

// How many times pieces of opposite color are being attacked?
func (b *Brain) attacksAvailable(p *Position, color int) (attacks int) {
        for i, piece := range p.pieces {
                if piece != 0 && piece.Color() == color {
                        targets := p.targets[i]
                        attacks += targets.Intersect(p.board[color^1]).Count()
                }
        }
        return
}
