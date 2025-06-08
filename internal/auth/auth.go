package auth

import (
	"api/service"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Middleware interface {
	WithCORS(next http.HandlerFunc) http.HandlerFunc
	AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc
}

type middleware struct {
	UserService service.UserService
}

func NewMiddleware(useService service.UserService) Middleware {
	return &middleware{UserService: useService}
}

func getSessionFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("better-auth.session_token")
	if err != nil {
		return "", err
	}
	sessionToken := strings.Split(cookie.Value, ".")
	if len(sessionToken) < 1 {
		return "", err
	}
	fmt.Printf("found token %v", sessionToken[0])
	return sessionToken[0], nil
}

func (m *middleware) AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := getSessionFromRequest(r)
		if err != nil {
			http.Error(w, "No credentials", http.StatusUnauthorized)
			return
		}
		user, sError := m.UserService.GetUserByCookie(token)
		if sError != nil {
			http.Error(w, "Not found", http.StatusUnauthorized)
			return
		}

		email, userId := user.Email, user.UserId
		ctx := context.WithValue(r.Context(), "email", email)
		ctx = context.WithValue(ctx, "userId", userId)
		fmt.Printf("Found user with %v %v", userId, email)
		next(w, r.WithContext(ctx))
	}
}

func (m *middleware) WithCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent) // 204: OK but no body
			return
		}
		next(w, r)
	}

}
