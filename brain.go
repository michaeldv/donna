package donna

import ()

const CENTER = 0x0000001818000000 // 4 central squares
const EXTENDED_CENTER = 0x00003C3C3C3C0000 // 12 central squares

type Brain struct {
        player *Player
        color  int
}

func NewBrain(player *Player) *Brain {
        brain := new(Brain)

        brain.player = player
        brain.color = player.Color

        return brain
}

func (b *Brain) Evaluate(p *Position) (score int) {
        material := b.materialBalance(p)
        mobility := b.mobilityBalance(p)
        aggression := b.aggressionBalance(p)
        center := b.centerBoost(p)
        score = material + mobility + aggression + center
        Log("Score for %s is %.2f (mat: %.2f, mob: %.2f, agg: %.2f, ctr: %.2f)\n", C(b.color), score, material, mobility, aggression, center)
        return
}

func (b *Brain) materialBalance(p *Position) int {
        opposite := b.color^1

        score := 1000 * (p.count[King(b.color)] - p.count[King(opposite)]) +
                  900 * (p.count[Queen(b.color)] - p.count[Queen(opposite)]) +
                  500 * (p.count[Rook(b.color)] - p.count[Rook(opposite)]) +
                  305 * (p.count[Bishop(b.color)] - p.count[Bishop(opposite)]) +
                  300 * (p.count[Knight(b.color)] - p.count[Knight(opposite)]) +
                  100 * (p.count[Pawn(b.color)] - p.count[Pawn(opposite)])

        return score
}

func (b *Brain) mobilityBalance(p *Position) (score int) {
        score = b.movesAvailable(p, b.color) - b.movesAvailable(p, b.color^1)
        if score != 0 {
                score += 25 * (score / Abs(score))
        }
        return
}

func (b *Brain) aggressionBalance(p *Position) (score int) {
        score = b.attacksAvailable(p, b.color) - b.attacksAvailable(p, b.color^1)
        if score != 0 {
                score += 20 * (score / Abs(score))
        }
        return
}

// How many attacks for the central squares?
func (b *Brain) centerBoost(p *Position) (center int) {
        for i, piece := range p.pieces {
                if piece != 0 && piece.Color() == b.color {
                        targets := p.targets[i]
                        sq12 := targets.Intersect(EXTENDED_CENTER).Count()
                        sq04 := targets.Intersect(CENTER).Count()
                        center += 5 * (sq12 - sq04) + 3 * sq04
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
