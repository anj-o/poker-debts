package main

import (
	game "anj/pokercalc/internal"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
)

var sessionManager *scs.SessionManager
var nextPlayerIndex int

func main() {
	// Initialize the session.
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	mux := http.NewServeMux()

	// Handle POST and GET requests.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// POST -- Calculate
			r.ParseForm()
			Calculate(r.Form)

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

	// Add the middleware.
	//muxWithSessionMiddleware := sessionManager.LoadAndSave(mux)

	// Start the server.
	fmt.Println("listening on :3000")
	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Printf("error listening: %v", err)
	}
}

func Calculate(formValues url.Values) string {

	// Parse form values
	beginningStack, _ := strconv.Atoi(formValues.Get("beginning-stack"))
	buyInCost, _ := strconv.ParseFloat(formValues.Get("buy-in"), 32)
	numPlayers := len(formValues) / 3
	fmt.Print(beginningStack, numPlayers)

	// Create game instance
	g := game.Game{
		StartChips: beginningStack,
		BuyInCost:  float32(buyInCost),
		Players:    ParsePlayers(formValues),
	}
	if err := g.ValidateGame(); err != nil {
		return fmt.Sprintf("%v", err)
	}
	g.CalculateNet()

	for _, player := range g.Players {
		fmt.Printf("%v is owed: %v\n", player.Name, player.Owed)
	}

	return "Calculated!"
}

func ParsePlayers(m url.Values) []*game.Player {

	numPlayers := len(m) / 3
	var p []*game.Player
	for i := 1; i <= numPlayers; i++ {

		nameKey := "p" + strconv.Itoa(i) + "-name"
		buyInsKey := "p" + strconv.Itoa(i) + "-buy-ins"
		endChipsKey := "p" + strconv.Itoa(i) + "-final-stack"

		parsedBuyIns, _ := strconv.Atoi(m.Get(buyInsKey))
		parsedEndChips, _ := strconv.Atoi(m.Get(endChipsKey))

		p = append(p, &game.Player{
			Name:     m.Get(nameKey),
			BuyIns:   parsedBuyIns,
			EndChips: parsedEndChips,
			Owed:     0,
		})
	}

	return p
}
