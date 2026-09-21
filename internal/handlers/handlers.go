package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"tickets-system-v1/internal/auth"
	"tickets-system-v1/internal/store"
)

var eventNotifier chan struct{}

func init() {
	eventNotifier = make(chan struct{}, 1)
}

func broadcastUpdate() {
	select {
	case eventNotifier <- struct{}{}:
	default:
	}
}

func SSEHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	sendState(w, flusher)

	notify := r.Context().Done()

	for {
		select {
		case <-notify:
			return
		case <-eventNotifier:
			sendState(w, flusher)
		}
	}
}

func sendState(w http.ResponseWriter, flusher http.Flusher) {
	store.Data.Mutex.RLock()
	data, err := json.Marshal(store.Data)
	store.Data.Mutex.RUnlock()

	if err != nil {
		log.Printf("Error marshaling data for SSE: %v", err)
		return
	}
	
	fmt.Fprintf(w, "data: %s\n\n", string(data))
	flusher.Flush()
}

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store.ConfigData.Categories)
}

func GetNumbersHandler(w http.ResponseWriter, r *http.Request) {
	store.Data.Mutex.RLock()
	defer store.Data.Mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store.Data)
}

func UpdateCurrentHandler(w http.ResponseWriter, r *http.Request) {
	var numbers []int
	if err := json.NewDecoder(r.Body).Decode(&numbers); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	store.Data.Mutex.Lock()
	store.Data.CurrentNumbers = numbers
	store.Data.Mutex.Unlock()

	store.SaveData()
	broadcastUpdate()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func UpdateWaitingHandler(w http.ResponseWriter, r *http.Request) {
	var numbers []int
	if err := json.NewDecoder(r.Body).Decode(&numbers); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	store.Data.Mutex.Lock()
	store.Data.WaitingNumbers = numbers
	store.Data.Mutex.Unlock()

	store.SaveData()
	broadcastUpdate()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func serveHTML(w http.ResponseWriter, r *http.Request, filePath string) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	http.ServeFile(w, r, filePath)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	serveHTML(w, r, "templates/index.html")
}

func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	serveHTML(w, r, "templates/login.html")
}

func AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	serveHTML(w, r, "templates/admin.html")
}

func AdminLoginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	password := r.FormValue("password")

	if password == auth.AdminPass {
		token := auth.GenerateSessionToken()
		auth.CreateSession(token)

		http.SetCookie(w, &http.Cookie{
			Name:     auth.CookieName,
			Value:    token,
			Expires:  time.Now().Add(auth.SessionTimeout),
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/admin/login?error=invalid", http.StatusSeeOther)
	}
}

func AdminLogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(auth.CookieName)
	if err == nil {
		auth.DestroySession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.CookieName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
	})
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}
