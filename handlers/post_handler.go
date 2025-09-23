package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"myblog/middleware"
	"myblog/model"
	"myblog/validator"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func parseTemplate(templateName string) *template.Template {
	tmpl, err := template.ParseFiles(fmt.Sprintf("templates/%s", templateName))
	if err != nil {
		panic(fmt.Sprintf("Error parsing template %s, %v", templateName, err))
	}
	return tmpl
}

func getPost() []model.Posts {
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

	return posts
	// json.NewEncoder(w).Encode(posts)
}

func getPostbyID(id int) *model.Posts {

	filepath := fmt.Sprintf("posts/posts%d.json", id)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil
	}

	var a model.Posts
	json.Unmarshal(data, &a)

	return &a

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	posts := getPost()
	tmpl := parseTemplate("home.html")
	tmpl.Execute(w, posts)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/posts/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid post ID"}`, http.StatusBadRequest)
		return
	}

	post := getPostbyID(id)
	if post == nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	tmpl := parseTemplate("postpage.html")
	tmpl.Execute(w, post)

}

func DashBoardHandler(w http.ResponseWriter, r *http.Request) {
	posts := getPost()
	tmpl := parseTemplate("dashboard.html")
	tmpl.Execute(w, posts)
	return
}

func NewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := parseTemplate("newPost.html")
		tmpl.Execute(w, nil)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	date := r.FormValue("date")

	if err := validator.ValidatePost(title, content, date); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Create
	files, _ := os.ReadDir("posts")

	var a model.Posts
	a.Title = title
	a.Content = content
	a.Date = date
	a.Author = user.Username
	a.ID = len(files) + 1

	filePath := fmt.Sprintf("posts/posts%d.json", a.ID)
	file, _ := os.Create(filePath)
	defer file.Close()
	json.NewEncoder(file).Encode(a)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/edit/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid post ID"}`, http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		post := getPostbyID(id)
		if post != nil {
			http.Error(w, `{"error":"post not found"}`, http.StatusNotFound)
			return
		}
		tmpl := parseTemplate("updatePost.html")
		tmpl.Execute(w, post)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	date := r.FormValue("date")

	if err := validator.ValidatePost(title, content, date); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var a model.Posts
	a.Title = title
	a.Content = content
	a.Date = date
	a.Author = user.Username
	a.ID = id

	filePath := fmt.Sprintf("posts/posts%d.json", a.ID)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, `{"error":"post not found"}`, http.StatusNotFound)
		return
	}
	file, _ := os.Create(filePath)
	defer file.Close()
	json.NewEncoder(file).Encode(a)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/delete/")
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
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

}

func CreatePostWithAuthI() http.HandlerFunc {
	return middleware.CookieAuthMidleware(NewHandler)
}
func DeletePostWithAuthI() http.HandlerFunc {
	return middleware.CookieAuthMidleware(DeleteHandler)
}
func UpdatePostWithAuthI() http.HandlerFunc {
	return middleware.CookieAuthMidleware(EditHandler)
}
func DashboardWithAuth() http.HandlerFunc {
	return middleware.CookieAuthMidleware(DashBoardHandler)
}
