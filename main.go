package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

type GlobalState struct {
	Count int
}

var global GlobalState
var sessionManager *scs.SessionManager
var nextPlayerIndex int

//func getHandler(w http.ResponseWriter, r *http.Request) {
//	//userCount := sessionManager.GetInt(r.Context(), "count")
//	component := home()
//	component.Render(r.Context(), w)
//	nextPlayerIndex = 3
//
//}

//func postHandler(w http.ResponseWriter, r *http.Request) {
//	// Update state.
//	r.ParseForm()
//
//	if r.Form != nil {
//
//	}
//
//	//// Check to see if the global button was pressed.
//	//if r.Form.Has("global") {
//	//	global.Count++
//	//}
//	//if r.Form.Has("user") {
//	//	currentCount := sessionManager.GetInt(r.Context(), "count")
//	//	sessionManager.Put(r.Context(), "count", currentCount+1)
//	//}
//
//	// Display the form.
//	getHandler(w, r)
//}

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
			fmt.Print(r.Form)

			// Assume form is complete
			game := game.
			return
		} else {
			// GET
			component := home()
			component.Render(r.Context(), w)
			nextPlayerIndex = 3
			return
		}
		//getHandler(w, r)
	})

	mux.HandleFunc("/add-player", func(w http.ResponseWriter, r *http.Request) {
		add_player(nextPlayerIndex).Render(r.Context(), w)
		nextPlayerIndex++
	})
	mux.HandleFunc("/remove-player", func(w http.ResponseWriter, r *http.Request) {
		remove_player()
	})

	// Include the static content.
	//mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Add the middleware.
	muxWithSessionMiddleware := sessionManager.LoadAndSave(mux)

	// Start the server.
	fmt.Println("listening on :3000")
	if err := http.ListenAndServe(":3000", muxWithSessionMiddleware); err != nil {
		log.Printf("error listening: %v", err)
	}
}
