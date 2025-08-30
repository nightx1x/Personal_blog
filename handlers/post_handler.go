package api

import (
	"encoding/json"
	"fmt"
	"myblog/middleware"
	"myblog/model"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
	"time"
)

func CreatePostAPI(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/newPost.html")
	if err != nil {
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var a model.Posts

	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if a.Title == "" || len(a.Title) > 200 {
		http.Error(w, `{"error":"Title is required and must be less than 200 characters"}`, http.StatusBadRequest)
		return
	}
	if a.Content == "" {
		http.Error(w, `{"error":"Content is required"}`, http.StatusBadRequest)
		return
	}
	if _, err := time.Parse("2006-01-02", a.Date); err != nil {
		http.Error(w, `{"error":"Date must be in YYYY-MM-DD format"}`, http.StatusBadRequest)
		return
	}

	a.Author = user.Username

	files, _ := os.ReadDir("posts")
	a.ID = len(files) + 1

	filePath := fmt.Sprintf("posts/posts%d.json", a.ID)
	file, _ := os.Create(filePath)
	defer file.Close()
	json.NewEncoder(file).Encode(a)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(a)
	tmpl.Execute(w, a)
}

func GetPostsApi(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	files, _ := os.ReadDir("posts")
	var posts []model.Posts

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".json" {
			continue
		}

		data, _ := os.ReadFile(filepath.Join("posts", f.Name()))
		var pos model.Posts
		json.Unmarshal(data, &pos)
		posts = append(posts, pos)
	}

	json.NewEncoder(w).Encode(posts)
	tmpl.Execute(w, posts)
}

func GetPostApi(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/post.html")
	if err != nil {
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	idStr := r.URL.Path[len("/posts/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid post ID"}`, http.StatusBadRequest)
		return
	}
	filepath := fmt.Sprintf("posts/posts%d.json", id)
	data, err := os.ReadFile(filepath)
	if err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	var a model.Posts
	json.Unmarshal(data, &a)
	json.NewEncoder(w).Encode(a)
	tmpl.Execute(w, a)

}
func DeletePostAPI(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/posts/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid post ID"}`, http.StatusBadRequest)
		return
	}
	filepath := fmt.Sprintf("posts/posts%d.json", id)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		http.Error(w, `{"error":"post not found"}`, http.StatusNotFound)
		return
	}
}

func UpdatePostApi(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/updatePost.html")
	if err != nil {
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	idStr := r.URL.Path[len("/posts/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid post ID"}`, http.StatusBadRequest)
		return
	}

	var a model.Posts
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}
	a.ID = id

	filePath := fmt.Sprintf("posts/posts%d.json", a.ID)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, `{"error":"post not found"}`, http.StatusNotFound)
		return
	}
	file, _ := os.Create(filePath)
	defer file.Close()
	json.NewEncoder(file).Encode(a)
	json.NewEncoder(w).Encode(a)
	tmpl.Execute(w, a)
}

func CreatePostWithAuthI() http.HandlerFunc {
	return middleware.JWTMiddleware(CreatePostAPI)
}
func DeletePostWithAuthI() http.HandlerFunc {
	return middleware.JWTMiddleware(DeletePostAPI)
}
func UpdatePostWithAuthI() http.HandlerFunc {
	return middleware.JWTMiddleware(UpdatePostApi)
}

// Доробити - валідація
// Доробити GetPostApi
// Доробити DeletePostAPI
// Доробити UpdatePostApi
