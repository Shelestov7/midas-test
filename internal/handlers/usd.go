package handlers

import (
	"context"
	"log"
	"net/http"
	"text/template"
	"time"

	dbq "midasivestment/internal/db"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type storedData struct {
	TokensPosition []TokenAssets
}

type TokenAssets struct {
	Usd         float64
	UserAddress string
	CreatedAt   time.Time
}

func Usd(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		searchBefore := r.URL.Query().Get("before")
		templateData := storedData{}
		err := templateData.collectDataTemplate(ctx, db, searchBefore)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		err = renderTemplate(w, "templates/index.html", templateData)
		if err != nil {
			log.Printf("render template: %s", err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func (t *TokenAssets) storeTokenPosition(rows pgx.Rows, data *storedData) error {
	for rows.Next() {
		err := rows.Scan(&t.Usd, &t.UserAddress, &t.CreatedAt)
		if err != nil {
			log.Printf("get user id error: %s", err.Error())
			return err
		}
		data.TokensPosition = append(data.TokensPosition, *t)
	}
	return nil
}

func renderTemplate(w http.ResponseWriter, templatePath string, data storedData) error {
	temp, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Printf("create templates: %s", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return err
	}
	err = temp.Execute(w, data)
	if err != nil {
		log.Printf("execute templates: %s", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return err
	}
	return nil
}

func (sd *storedData) collectDataTemplate(ctx context.Context, db *pgxpool.Pool, params string) error {
	position := TokenAssets{}
	if params != "" {
		rows, err := dbq.GetTokenPositionsBefore(ctx, db, params)
		if err != nil {
			log.Printf("get user id's from db: %s", err.Error())
			return err
		}
		defer rows.Close()

		err = position.storeTokenPosition(rows, sd)
		if err != nil {
			return err
		}
	} else {
		rows, err := dbq.GetTokenPositions(ctx, db)
		if err != nil {
			log.Printf("get user id's from db: %s", err.Error())
			return err
		}
		defer rows.Close()

		err = position.storeTokenPosition(rows, sd)
		if err != nil {
			return err
		}
	}
	return nil
}
