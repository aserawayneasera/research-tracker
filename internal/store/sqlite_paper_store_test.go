package store_test

import (
	"database/sql"
	"testing"

	"github.com/aserawayneasera/research-tracker/internal/db"
	"github.com/aserawayneasera/research-tracker/internal/models"
	"github.com/aserawayneasera/research-tracker/internal/store"
)

func TestListFilters(t *testing.T) {
	conn, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer conn.Close()

	if err := db.InitSchema(conn); err != nil {
		t.Fatalf("schema: %v", err)
	}

	ps := store.NewSQLitePaperStore(conn)

	_, _ = ps.Create(&models.Paper{Title: "A", Venue: "X", Year: 2025, Status: "draft", Tags: "cv", Notes: ""})
	_, _ = ps.Create(&models.Paper{Title: "B", Venue: "Y", Year: 2026, Status: "submitted", Tags: "vision", Notes: ""})

	items, total, err := ps.List(store.ListParams{Status: "submitted", Limit: 10, Offset: 0})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("expected 1, got total=%d len=%d", total, len(items))
	}
	if items[0].Title != "B" {
		t.Fatalf("expected B, got %s", items[0].Title)
	}
}
