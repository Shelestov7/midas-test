package handlers

import (
	"context"
	"log"
	"net/http"
	"sync"

	"midasivestment/internal/api"
	dbq "midasivestment/internal/db"

	"github.com/jackc/pgx/v4/pgxpool"
)

type UserList interface {
	GetUserList() ([]string, error)
}

func FetchProtocolList(ctx context.Context, db *pgxpool.Pool, httpClient *api.HTTPClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := GetProtocolList(ctx, db, httpClient)
		if err != nil {
			log.Printf("Get protocol list error: %s", err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func GetProtocolList(ctx context.Context, db *pgxpool.Pool, httpClient *api.HTTPClient) error {
	userList, err := dbq.GetUserList(ctx, db)
	if err != nil {
		log.Println("Get user list error")
		return err
	}

	wg := sync.WaitGroup{}
	for _, userID := range userList {
		wg.Add(1)
		go func(ctx context.Context, userID string) {
			resp, err := httpClient.UserSimpleProtocolList(userID)
			if err != nil {
				log.Printf("api request for user id %q error: %s", userID, err.Error())
				return
			}
			err = dbq.InsertAPIResult(ctx, db, resp, userID)
			if err != nil {
				log.Printf("api request for user id %q error: %s", userID, err.Error())
				return
			}
			wg.Done()
		}(ctx, userID)
	}
	wg.Wait()
	return nil
}
