package handler

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

var postStore = map[int]*Post{}
var mu = &sync.RWMutex{}

type storePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type updatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PostHandler struct{}

func NewPostHandler() *PostHandler {
	return &PostHandler{}
}

func (ph *PostHandler) Index(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(postStore); err != nil {
		log.Println(err.Error())
	}
}

func (ph *PostHandler) Store(w http.ResponseWriter, r *http.Request) {
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
}

func (ph *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
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
}

func (ph *PostHandler) Destroy(w http.ResponseWriter, r *http.Request) {
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
}
