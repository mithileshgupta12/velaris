package route

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type storePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

var postStore = []*Post{}

var mu = &sync.RWMutex{}

func PostRoutes(r *chi.Mux) {
	r.Route("/posts", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			mu.RLock()
			defer mu.RUnlock()

			w.Header().Set("Content-Type", "application/json")

			if err := json.NewEncoder(w).Encode(postStore); err != nil {
				log.Println(err.Error())
			}
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			var req *storePostRequest

			w.Header().Set("Content-Type", "application/json")

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				log.Println(err.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if req.Title == "" {
				http.Error(w, "title is a required field", http.StatusBadRequest)
				return
			}

			if req.Content == "" {
				http.Error(w, "content is a required field", http.StatusBadRequest)
				return
			}

			mu.Lock()
			post := &Post{
				ID:        len(postStore) + 1,
				Title:     req.Title,
				Content:   req.Content,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			postStore = append(postStore, post)
			mu.Unlock()

			w.WriteHeader(http.StatusCreated)

			if err := json.NewEncoder(w).Encode(post); err != nil {
				log.Println(err.Error())
			}
		})
	})

	r.Delete("/{postId}", func(w http.ResponseWriter, r *http.Request) {
		postIdParam := chi.URLParam(r, "postId")

		w.Header().Set("Content-Type", "application/json")

		postId, err := strconv.Atoi(postIdParam)
		if err != nil {
			http.Error(w, "Post id is invalid", http.StatusBadRequest)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		postFound := false
		for i, v := range postStore {
			if v.ID == postId {
				postStore = append(postStore[:i], postStore[i+1:]...)
				postFound = true
				break
			}
		}

		if !postFound {
			http.Error(w, "No post found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
