// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import `regexp`

func NewPositionFromFEN(game *Game, fen string) *Position {
	const ( 
		row1   = `([rnbqkRNBQK1-8]+/`
		rows26 = `([rnbqkpRNBQKP1-8]+/){6}`
		row8   = `[rnbqkRNBQK1-8]+)`
		color  = `([bw])`
		castle = `(-|K?Q?k?q?)`
		enpass = `(-|[a-h][36])`
		number = `(\d+)`
		spx    = `\s*`
		sp     = `\s`
	)

	tree[node] = Position{game: game}
	p := &tree[node]

	// Expected matches of interest are as follows:
	// [1] - Pieces (entire board).
	// [3] - Color of side to move.
	// [4] - Castle rights.
	// [5] - En-passant square.
	// [6] - Number of half-moves.
	// [7] - Number of full moves.
	re := regexp.MustCompile(spx + row1 + rows26 + row8 + sp + color + sp + castle + sp + enpass + sp + number + sp + number + spx);
	matches := re.FindStringSubmatch(fen)
	if matches == nil {
		return nil
	}

	p.reversible = true
	p.board = p.outposts[White] | p.outposts[Black]
	p.hash, p.hashPawns, p.hashMaterial = p.polyglot()
	p.tally = p.valuation()

	return p
}

func (p *Position) fen() string {
	return ":-)"
}
