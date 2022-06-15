package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"midasivestment/internal/api"
	"midasivestment/internal/db"
	"midasivestment/internal/handlers"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	apiURL   = "https://openapi.debank.com"
	Timeout  = time.Second * 10
	grabTime = time.Hour * 4
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx := context.Background()

	dbConf := db.ConfigureDBConnection(ctx)
	database, err := db.InitDBCConnection(ctx, dbConf)

	if err != nil {
		panic(err)
	}

	client := api.NewHTTPClient(ctx, apiURL, http.Client{Timeout: Timeout})

	router := mux.NewRouter()

	server := api.InitHTTPServer(router)

	router.HandleFunc("/fetch", handlers.FetchProtocolList(ctx, database, client))
	router.HandleFunc("/usd", handlers.Usd(ctx, database)).Methods("GET")
	errChan := make(chan error)

	go func(errChan chan error) {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}(errChan)
	log.Printf("Server Started on port %s", server.Addr)

	go func(ctx context.Context, db *pgxpool.Pool) {
		for {
			err = handlers.GetProtocolList(ctx, db, client)
			if err != nil {
				log.Printf("get protocol list: %s", err.Error())
			}

			<-time.After(grabTime)
		}
	}(ctx, database)

	select {
	case err = <-errChan:
		log.Fatalln(err)

	case <-done:
		log.Print("Server Stopped")
		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer func() {
			database.Close()
			cancel()
		}()

		if err = server.Shutdown(ctx); err != nil {
			log.Fatalf("Server Shutdown Failed:%+v", err)
		}
		log.Print("Server Successfully Stopped")
	}
}
