package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/aserawayneasera/research-tracker/internal/store"
)

type Handler = http.Handler

type RouterProvider interface {
	Router() Handler
}

type router struct {
	papers store.PaperStore
}

func NewRouter(papers store.PaperStore) RouterProvider {
	return &router{papers: papers}
}

func (r *router) Router() Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-Id"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	api := NewHandlers(r.papers)

	mux.Get("/health", api.Health)

	mux.Route("/papers", func(cr chi.Router) {
		cr.Get("/", api.ListPapers)
		cr.Post("/", api.CreatePaper)
		cr.Get("/{id}", api.GetPaper)
		cr.Put("/{id}", api.UpdatePaper)
		cr.Delete("/{id}", api.DeletePaper)
	})

	return mux
}
