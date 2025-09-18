// @title Posts API
// @version 1.0
// @description A simple blog posts API
// @host localhost:8000
// @BasePath /

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

// Post represents a blog post
type Post struct {
	ID        int       `json:"id" example:"1"`
	Title     string    `json:"title" example:"My Blog Post"`
	Content   string    `json:"content" example:"This is the content of my blog post"`
	CreatedAt time.Time `json:"created_at" example:"2025-01-15T10:30:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-01-15T10:30:00Z"`
}

// storePostRequest represents the request body for creating a post
type storePostRequest struct {
	Title   string `json:"title" example:"My New Post" binding:"required"`
	Content string `json:"content" example:"Content of the new post" binding:"required"`
}

// updatePostRequest represents the request body for updating a post
type updatePostRequest struct {
	Title   string `json:"title" example:"Updated Post Title" binding:"required"`
	Content string `json:"content" example:"Updated post content" binding:"required"`
}

var postStore = map[int]*Post{}
var mu = &sync.RWMutex{}

func PostRoutes(r *chi.Mux) {
	r.Route("/posts", func(r chi.Router) {
		// @Summary Get all posts
		// @Description Retrieve all posts from the store
		// @Tags posts
		// @Produce json
		// @Success 200 {object} map[string]Post "Posts retrieved successfully"
		// @Router /posts [get]
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			mu.RLock()
			defer mu.RUnlock()
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(postStore); err != nil {
				log.Println(err.Error())
			}
		})

		// @Summary Create a new post
		// @Description Create a new blog post
		// @Tags posts
		// @Accept json
		// @Produce json
		// @Param post body storePostRequest true "Post data"
		// @Success 201 {object} Post "Post created successfully"
		// @Failure 400 {string} string "Bad Request - missing required fields"
		// @Failure 500 {string} string "Internal Server Error"
		// @Router /posts [post]
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

		// @Summary Update a post
		// @Description Update an existing post by ID
		// @Tags posts
		// @Accept json
		// @Produce json
		// @Param postId path int true "Post ID"
		// @Param post body updatePostRequest true "Updated post data"
		// @Success 200 {object} Post "Post updated successfully"
		// @Failure 400 {string} string "Bad Request - invalid ID or missing fields"
		// @Failure 404 {string} string "Post not found"
		// @Failure 500 {string} string "Internal Server Error"
		// @Router /posts/{postId} [put]
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

		// @Summary Delete a post
		// @Description Delete a post by ID
		// @Tags posts
		// @Param postId path int true "Post ID"
		// @Success 204 "Post deleted successfully"
		// @Failure 400 {string} string "Bad Request - invalid ID"
		// @Failure 404 {string} string "Post not found"
		// @Router /posts/{postId} [delete]
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
