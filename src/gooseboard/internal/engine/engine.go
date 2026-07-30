package engine

import (
	"encoding/json"
	"net/http"

	"gooseboard/internal/schema"
	"gooseboard/internal/store"
)

type Engine struct {
	Panel schema.Panel
	Store store.Store
}

func New(p schema.Panel, s store.Store) *Engine {
	return &Engine{Panel: p, Store: s}
}

func (e *Engine) Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/nav", e.handleNav)
	mux.HandleFunc("/api/pages/", e.handlePage)

	for name := range e.Panel.Entities {
		entity := name
		mux.HandleFunc("/api/entities/"+entity, e.handleEntityCollection(entity))
	}

	return mux
}

func (e *Engine) handleNav(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(e.Panel.Nav)
}

func (e *Engine) handlePage(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/api/pages/"):]
	page, ok := e.Panel.Pages[title]
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(page)
}

func (e *Engine) handleEntityCollection(entity string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rows, err := e.Store.List(entity, nil)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			json.NewEncoder(w).Encode(rows)
		case http.MethodPost:
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			id, err := e.Store.Create(entity, body)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"id": id})
		default:
			http.Error(w, "method not allowed", 405)
		}
	}
}
