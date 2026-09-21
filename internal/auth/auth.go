package auth

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	AdminPass      string
	sessionStore   = make(map[string]time.Time)
	sessionMutex   sync.Mutex
	CookieName     = "admin_session"
	SessionTimeout = 24 * time.Hour
)

func InitAuth() {
	AdminPass = os.Getenv("ADMIN_PASSWORD")
	if AdminPass == "" {
		AdminPass = "password"
	}
}

func GenerateSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(CookieName)
		if err != nil {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}

		sessionMutex.Lock()
		expiry, exists := sessionStore[cookie.Value]
		if !exists || time.Now().After(expiry) {
			delete(sessionStore, cookie.Value)
			sessionMutex.Unlock()
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		sessionStore[cookie.Value] = time.Now().Add(SessionTimeout)
		sessionMutex.Unlock()

		next.ServeHTTP(w, r)
	})
}

func CreateSession(token string) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	sessionStore[token] = time.Now().Add(SessionTimeout)
}

func DestroySession(token string) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	delete(sessionStore, token)
}
