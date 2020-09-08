package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/tim-harding/fatal-encounters-server/routes/cityroute"
	"github.com/tim-harding/fatal-encounters-server/routes/enumroute"
	"github.com/tim-harding/fatal-encounters-server/shared"
)

var (
	enumTables = []string{
		"agency",
		"cause",
		"county",
		"race",
		"use_of_force",
	}
)

func main() {
	defer shared.Db.Close()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/city", cityroute.HandleRoute)
	for _, table := range enumTables {
		route := fmt.Sprintf("/%s", table)
		r.Get(route, enumroute.HandleRouteFactory(table))
	}
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal(err)
	}
}
