package route

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Post struct {
	ID        int       `json:"id" example:"1"`
	Title     string    `json:"title" example:"My Blog Post"`
	Content   string    `json:"content" example:"This is the content of my blog post"`
	CreatedAt time.Time `json:"created_at" example:"2025-01-15T10:30:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-01-15T10:30:00Z"`
}

type storePostRequest struct {
	Title   string `json:"title" example:"My New Post" binding:"required"`
	Content string `json:"content" example:"Content of the new post" binding:"required"`
}

type updatePostRequest struct {
	Title   string `json:"title" example:"Updated Post Title" binding:"required"`
	Content string `json:"content" example:"Updated post content" binding:"required"`
}

var postStore = map[int]*Post{}
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
			postId := len(postStore) + 1
			post := &Post{
				ID:        postId,
				Title:     req.Title,
				Content:   req.Content,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			postStore[postId] = post
			mu.Unlock()
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(post); err != nil {
				log.Println(err.Error())
			}
		})

		r.Put("/{postId}", func(w http.ResponseWriter, r *http.Request) {
			postIdParam := chi.URLParam(r, "postId")
			w.Header().Set("Content-Type", "application/json")
			postId, err := strconv.Atoi(postIdParam)
			if err != nil {
				http.Error(w, "Post id is invalid", http.StatusBadRequest)
				return
			}
			var req *updatePostRequest
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
			defer mu.Unlock()
			post, ok := postStore[postId]
			if !ok {
				http.Error(w, "No post found", http.StatusNotFound)
				return
			}
			post.Title = req.Title
			post.Content = req.Content
			post.UpdatedAt = time.Now()
			postStore[postId] = post
			if err := json.NewEncoder(w).Encode(post); err != nil {
				log.Println(err.Error())
			}
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
			_, ok := postStore[postId]
			if !ok {
				http.Error(w, "No post found", http.StatusNotFound)
				return
			}
			delete(postStore, postId)
			w.WriteHeader(http.StatusNoContent)
		})
	})
}
