package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/tim-harding/fatal-encounters-server/routes/cityroute"
	"github.com/tim-harding/fatal-encounters-server/routes/stateroute"
	"github.com/tim-harding/fatal-encounters-server/shared"
)

func main() {
	defer shared.Db.Close()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/city", cityroute.HandleRoute)
	r.Get("/state", stateroute.HandleRoute)
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal(err)
	}
}
