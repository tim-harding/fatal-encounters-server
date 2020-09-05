package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/tim-harding/fatal-encounters-server/queries/city"
	"github.com/tim-harding/fatal-encounters-server/shared"
)

func main() {
	defer shared.Db.Close()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/city", city.HandleRoute)
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal(err)
	}
}
