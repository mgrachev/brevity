package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/mgrachev/brevity/api"
	"github.com/mgrachev/brevity/config"
	"github.com/mgrachev/brevity/web"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	db, err := sql.Open("postgres", config.PGConnectionString())
	if err != nil {
		log.Fatal(err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err.Error())
	}

	router := mux.NewRouter()
	apiHandler := api.NewHandler(db)
	router.HandleFunc("/api/v1/shortlink", apiHandler.FindOrCreateShortLink).Methods("POST")

	webHandler := web.NewHandler(db)
	router.HandleFunc("/{token}", webHandler.RedirectToOrigin).Methods("GET")

	log.Info(fmt.Sprintf("HTTP Server started at port: %s", config.AppPort()))

	port := fmt.Sprintf(":%s", config.AppPort())
	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal(err.Error())
	}
}
