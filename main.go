package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/tim-harding/fatal-encounters-server/routes/cityroute"
	"github.com/tim-harding/fatal-encounters-server/routes/enumroute"
	"github.com/tim-harding/fatal-encounters-server/routes/incidentroute"
	"github.com/tim-harding/fatal-encounters-server/routes/stateroute"
	"github.com/tim-harding/fatal-encounters-server/shared"
)

var enumTables = []string{
	"agency",
	"cause",
	"county",
	"race",
	"use_of_force",
}

func main() {
	defer shared.Db.Close()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/city", func(r chi.Router) {
		r.Get("/", cityroute.HandleBaseRoute)
		r.Get("/{id}", cityroute.HandleIDRoute)
	})
	r.Route("/state", func(r chi.Router) {
		r.Get("/", stateroute.HandleBaseRoute)
		r.Get("/{id}", stateroute.HandleIDRoute)
	})
	for _, table := range enumTables {
		route := fmt.Sprintf("/%s", table)
		r.Route(route, func(r chi.Router) {
			r.Get("/", enumroute.HandleBaseRouteFactory(table))
			r.Get("/{id}", enumroute.HandleIDRouteFactory(table))
		})
	}
	r.Route("/incident", func(r chi.Router) {
		r.Get("/filter", incidentroute.HandleIncidentFilterRoute)
		r.Get("/position", incidentroute.HandleIncidentPositionRoute)
		r.Get("/detail/{id:[0-9,]+}", incidentroute.HandleIncidentDetailRoute)
		r.Get("/count", incidentroute.HandleCountRoute)
	})
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal(err)
	}
}
