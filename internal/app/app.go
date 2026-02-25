package app

import (
	"database/sql"

	"github.com/aserawayneasera/research-tracker/internal/db"
	"github.com/aserawayneasera/research-tracker/internal/httpapi"
	"github.com/aserawayneasera/research-tracker/internal/store"
)

type App struct {
	db     *sql.DB
	router httpapi.RouterProvider
}

func New(dbPath string) (*App, error) {
	conn, err := db.Open(dbPath)
	if err != nil {
		return nil, err
	}
	if err := db.InitSchema(conn); err != nil {
		_ = conn.Close()
		return nil, err
	}

	papers := store.NewSQLitePaperStore(conn)
	r := httpapi.NewRouter(papers)

	return &App{db: conn, router: r}, nil
}

func (a *App) Router() httpapi.Handler {
	return a.router.Router()
}

func (a *App) Close() error {
	return a.db.Close()
}
