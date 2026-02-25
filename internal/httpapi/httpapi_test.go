package httpapi_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aserawayneasera/research-tracker/internal/db"
	"github.com/aserawayneasera/research-tracker/internal/httpapi"
	"github.com/aserawayneasera/research-tracker/internal/store"
)

func newTestServer(t *testing.T) (*httptest.Server, func()) {
	t.Helper()

	conn, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := conn.Ping(); err != nil {
		t.Fatalf("ping db: %v", err)
	}
	if err := db.InitSchema(conn); err != nil {
		t.Fatalf("schema: %v", err)
	}

	ps := store.NewSQLitePaperStore(conn)
	r := httpapi.NewRouter(ps).Router()

	srv := httptest.NewServer(r)
	cleanup := func() {
		srv.Close()
		_ = conn.Close()
	}
	return srv, cleanup
}

func TestCreateAndGetPaper(t *testing.T) {
	srv, cleanup := newTestServer(t)
	defer cleanup()

	body := map[string]any{
		"title":  "LGF for Small Object Detection",
		"venue":  "PRL",
		"year":   2026,
		"status": "submitted",
		"tags":   "retinanet,small",
		"notes":  "seed=42",
	}
	b, _ := json.Marshal(body)

	res, err := http.Post(srv.URL+"/papers", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status: got %d", res.StatusCode)
	}

	var created map[string]any
	_ = json.NewDecoder(res.Body).Decode(&created)
	_ = res.Body.Close()

	id := int(created["id"].(float64))
	getRes, err := http.Get(srv.URL + "/papers/" + itoa(id))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d", getRes.StatusCode)
	}
	_ = getRes.Body.Close()
}

func itoa(i int) string {
	// minimal dependency
	if i == 0 {
		return "0"
	}
	sign := ""
	if i < 0 {
		sign = "-"
		i = -i
	}
	buf := make([]byte, 0, 12)
	for i > 0 {
		d := i % 10
		buf = append([]byte{byte('0' + d)}, buf...)
		i /= 10
	}
	return sign + string(buf)
}
