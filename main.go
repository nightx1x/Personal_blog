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

	http.HandleFunc("/", api.HomeHandler)
	http.HandleFunc("/post/", api.PostHandler)

	http.HandleFunc("/dashboard", api.DashboardWithAuth())
	http.HandleFunc("/new", api.CreatePostWithAuthI())
	http.HandleFunc("/edit/", api.UpdatePostWithAuthI())
	http.HandleFunc("/delete/", api.DeletePostWithAuthI())

	http.HandleFunc("/logout", api.LogoutHandler)

	http.HandleFunc("/login", api.LoginHandler)
	fmt.Println("Server started at :8083")
	http.ListenAndServe(":8083", nil)

}
