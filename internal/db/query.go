package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	GetUserIDsQuery      = "SELECT user_id FROM users;"
	InsertAPIResultQuery = "INSERT INTO debank_users_assets (body, user_id, created_at) VALUES ($1, $2, now());"

	GetTokenPositionsQuery = `SELECT value->'net_usd_value' AS usd, user_id, created_at
    FROM debank_users_assets, jsonb_array_elements(debank_users_assets.body);`

	GetTokenPositionsBeforeQuery = `SELECT value->'net_usd_value' AS usd, user_id, created_at
    FROM debank_users_assets, jsonb_array_elements(debank_users_assets.body)
	WHERE created_at < $1;`
)

func GetUserList(ctx context.Context, db *pgxpool.Pool) ([]string, error) {
	rows, err := db.Query(ctx, GetUserIDsQuery)
	if err != nil {
		log.Printf("get user id's: %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	var address string

	userList := make([]string, 0, 0)

	for rows.Next() {
		err = rows.Scan(&address)
		if err != nil {
			log.Printf("get user id's: %s", err.Error())
			return nil, err
		}
		userList = append(userList, address)
	}
	return userList, nil
}

func GetTokenPositions(ctx context.Context, db *pgxpool.Pool) (pgx.Rows, error) {
	rows, err := db.Query(ctx, GetTokenPositionsQuery)
	if err != nil {
		return nil, fmt.Errorf("get token positions: %s", err.Error())
	}

	return rows, err
}

func GetTokenPositionsBefore(ctx context.Context, db *pgxpool.Pool, before string) (pgx.Rows, error) {
	rows, err := db.Query(ctx, GetTokenPositionsBeforeQuery, before)
	if err != nil {
		return nil, fmt.Errorf("get token positions before: %s", err.Error())
	}

	return rows, nil
}

func InsertAPIResult(ctx context.Context, db *pgxpool.Pool, data []byte, userID string) error {
	_, err := db.Exec(ctx, InsertAPIResultQuery, data, userID)
	if err != nil {
		return fmt.Errorf("insert api result: %s", err.Error())
	}

	return nil
}
