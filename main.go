package main

import (
	game "anj/pokercalc/internal"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
)

var sessionManager *scs.SessionManager
var nextPlayerIndex int

func main() {
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	mux := http.NewServeMux()

	// Handle POST and GET requests.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// POST -- Calculate
			r.ParseForm()
			result := Calculate(r.Form)
			component := chain(result)
			component.Render(r.Context(), w)

			return
		} else {
			// GET -- Home
			component := home()
			component.Render(r.Context(), w)
			nextPlayerIndex = 3
			return
		}
	})

	mux.HandleFunc("/add-player", func(w http.ResponseWriter, r *http.Request) {
		add_player(nextPlayerIndex).Render(r.Context(), w)
		nextPlayerIndex++
	})
	mux.HandleFunc("/remove-player", func(w http.ResponseWriter, r *http.Request) {
		remove_player()
		nextPlayerIndex--
	})

	fmt.Println("listening on :3000")
	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Printf("error listening: %v", err)
	}
}

func ParsePlayers(m url.Values) map[string]*game.Player {

	// So long as there are only two non-player inputs, integer division will yield the number of players
	numPlayers := len(m) / 3
	p := make(map[string]*game.Player)

	// Cycle over the player indices and store in the Game struct
	for i := 1; i <= numPlayers; i++ {

		nameKey := "p" + strconv.Itoa(i) + "-name"
		buyInsKey := "p" + strconv.Itoa(i) + "-buy-ins"
		endChipsKey := "p" + strconv.Itoa(i) + "-final-stack"

		parsedBuyIns, _ := strconv.Atoi(m.Get(buyInsKey))
		parsedEndChips, _ := strconv.Atoi(m.Get(endChipsKey))

		p[m.Get(nameKey)] = &game.Player{
			BuyIns:   parsedBuyIns,
			EndChips: parsedEndChips,
			Owed:     0,
		}
	}

	return p
}

func Calculate(formValues url.Values) []string {

	// Parse form values
	beginningStack, _ := strconv.Atoi(formValues.Get("beginning-stack"))
	buyInCost, _ := strconv.ParseFloat(formValues.Get("buy-in"), 32)

	// Create game instance
	g := game.Game{
		StartChips: beginningStack,
		BuyInCost:  float32(buyInCost),
		Players:    ParsePlayers(formValues),
	}
	if err := g.ValidateGame(); err != nil {
		return []string{fmt.Sprintf("%v", err)}
	}
	g.CalculateNet()

	ledger := GetPayoutChain(g)

	for _, line := range ledger {
		fmt.Println(line)
	}

	return ledger
}

func GetPayoutChain(g game.Game) []string {

	ledger := make([]string, 0)

	// While players still owe money...
	for len(g.Players) > 0 && !AllSettledUp(g) {

		// Most indebted will pay the most owed player directly
		var mostIndebtedPlayer string
		var mostIndebtedPlayerAmount float32

		var mostOwedPlayer string
		var mostOwedPlayerAmount float32

		for name, stat := range g.Players {
			if stat.Owed < mostIndebtedPlayerAmount {
				mostIndebtedPlayer = name
				mostIndebtedPlayerAmount = stat.Owed
			}

			if stat.Owed > mostOwedPlayerAmount {
				mostOwedPlayer = name
				mostOwedPlayerAmount = stat.Owed
			}
		}

		if mostIndebtedPlayerAmount+mostOwedPlayerAmount > 0 {
			// If mostIndebtedPlayer *CANNOT* make mostOwedPlayer whole...

			// Update ledger
			entry := fmt.Sprintf("%s pays $%.2f to %s", mostIndebtedPlayer, -mostIndebtedPlayerAmount, mostOwedPlayer)
			ledger = append(ledger, entry)

			// mostOwedPlayer has a new balance, update it
			g.Players[mostOwedPlayer].Owed += mostIndebtedPlayerAmount

			// mostIndebtedPlayer is paid up, remove them from further consideration
			delete(g.Players, mostIndebtedPlayer)

		} else if mostIndebtedPlayerAmount+mostOwedPlayerAmount < 0 {
			// If mostIndebtedPlayer *CAN* make mostOwedPlayer whole...

			// Update ledger
			entry := fmt.Sprintf("%s pays $%.2f to %s", mostIndebtedPlayer, mostOwedPlayerAmount, mostOwedPlayer)
			ledger = append(ledger, entry)

			// mostIndebtedPlayer has a new balance, update it
			g.Players[mostIndebtedPlayer].Owed += mostOwedPlayerAmount

			// mostOwedPlayer is paid up, remove them from further consideration
			delete(g.Players, mostOwedPlayer)
		} else {
			// mostIndebtedPlayer and mostOwedPlayer have the same balance, remove them both
			entry := fmt.Sprintf("%s pays $%.2f to %s", mostIndebtedPlayer, mostOwedPlayerAmount, mostOwedPlayer)
			ledger = append(ledger, entry)

			delete(g.Players, mostOwedPlayer)
			delete(g.Players, mostIndebtedPlayer)
		}
	}

	return ledger

}

func AllSettledUp(g game.Game) bool {

	settled := true
	for _, player := range g.Players {
		if math.Abs(float64(player.Owed)) >= 0.01 {
			settled = false
		}
	}
	return settled
}
