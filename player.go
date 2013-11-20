package lape

type Player struct {
	Color  int  // 0: white, 1: black
        Can00  bool // Can castle king's side?
        Can000 bool // Can castle queen's side?
}

func (p *Player)Initialize(color int) *Player {
        p.Color = color
        p.Can00 = true
        p.Can000 = true

        return p
}

func (p *Player) IsWhite() bool {
	return p.Color == 0
}

func (p *Player) IsBlack() bool {
	return p.Color != 0
}
