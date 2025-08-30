package main

import (
	"fmt"
	api "myblog/handlers"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found")
	}

	if _, err := os.Stat("posts"); os.IsNotExist(err) {
		os.Mkdir("posts", 0755)
	}
}

func main() {

	http.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			api.GetPostsApi(w, r)
		case http.MethodPost:
			api.CreatePostWithAuthI()(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			api.GetPostApi(w, r)
		case http.MethodPut:
			api.UpdatePostWithAuthI()(w, r)
		case http.MethodDelete:
			api.DeletePostWithAuthI()(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/login", api.LoginHandler)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)

}
