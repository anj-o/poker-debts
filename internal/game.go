package game

import "fmt"

type Player struct {
	BuyIns   int
	EndChips int
	Owed     float32
}

type Game struct {
	StartChips int
	BuyInCost  float32
	Players    map[string]*Player
}

func (g Game) CalculateNet() {
	for _, p := range g.Players {
		net := float32(p.EndChips-p.BuyIns*g.StartChips) / float32(g.StartChips)
		p.Owed = net * g.BuyInCost
	}
}

func (g Game) ValidateGame() error {
	var totalBuyIns int
	var countedChips int
	for _, p := range g.Players {
		totalBuyIns += p.BuyIns
		countedChips += p.EndChips
	}

	expectedTotal := totalBuyIns * g.StartChips
	countingError := expectedTotal - countedChips
	if countingError != 0 {
		err := fmt.Errorf("Expected chip total of %v (%v buy-ins @ %v chips each), but got %v.", expectedTotal, totalBuyIns, g.StartChips, countedChips)
		return err
	}
	return nil
}
