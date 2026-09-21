package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"tickets-system-v1/internal/auth"
	"tickets-system-v1/internal/handlers"
	"tickets-system-v1/internal/store"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".envファイルが見つかりませんでした。環境変数を使用します。")
	}

	auth.InitAuth()

	store.InitStore()

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/categories", handlers.CategoriesHandler).Methods("GET")
	api.HandleFunc("/events", handlers.SSEHandler).Methods("GET")
	api.HandleFunc("/numbers", handlers.GetNumbersHandler).Methods("GET")

	adminAPI := api.PathPrefix("").Subrouter()
	adminAPI.Use(auth.CheckAuth)
	adminAPI.HandleFunc("/numbers/current", handlers.UpdateCurrentHandler).Methods("POST")
	adminAPI.HandleFunc("/numbers/waiting", handlers.UpdateWaitingHandler).Methods("POST")

	r.HandleFunc("/", handlers.IndexHandler).Methods("GET")
	r.HandleFunc("/admin/login", handlers.AdminLoginHandler).Methods("GET")
	r.HandleFunc("/admin/login", handlers.AdminLoginPostHandler).Methods("POST")
	r.HandleFunc("/admin/logout", handlers.AdminLogoutHandler).Methods("POST")

	adminPages := r.PathPrefix("/admin").Subrouter()
	adminPages.Use(auth.CheckAuth)
	adminPages.HandleFunc("/dashboard", handlers.AdminDashboardHandler).Methods("GET")

	r.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
	}).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("サーバーが起動しました。http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
