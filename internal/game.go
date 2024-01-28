package game

import "fmt"

type Player struct {
	Name     string `form:"Name" binding:"required"`
	BuyIns   int    `form:"BuyIns" binding:"required"`
	EndChips int    `form:"EndStackValue" binding:"required"`
	Owed     float32
}

type Game struct {
	StartChips int       `form:"StartStackValue" binding:"required"`
	BuyInCost  float32   `form:"BuyInCost" binding:"required"`
	Players    []*Player `form:"Players" binding:"required""`
}

//type Player struct {
//	Name     string `json:"Name" binding:"required"`
//	BuyIns   int    `json:"BuyIns" binding:"required"`
//	EndChips int    `json:"EndStackValue" binding:"required"`
//	Owed     float32
//}
//
//type Game struct {
//	StartChips int       `json:"StartStackValue" binding:"required"`
//	BuyInCost  float32   `json:"BuyInCost" binding:"required"`
//	Players    []*Player `json:"Players" binding:"required"`
//}

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
		err := fmt.Errorf("expected chip total of %v, got %v (%v buy-ins at %v chips each)", expectedTotal, countedChips, totalBuyIns, g.StartChips)
		return err
	}
	return nil
}
