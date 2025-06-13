package auth

import (
	"api/service"
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Middleware interface {
	WithCORS(next http.HandlerFunc) http.HandlerFunc
	AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc
}

type middleware struct {
	UserService service.UserService
	PublicKey   ed25519.PublicKey
}

func NewMiddleware(useService service.UserService) Middleware {
	return &middleware{UserService: useService, PublicKey: getPublicKeyFromJWK(os.Getenv("JWKS_PUBLIC_KEY"))}
}

func (m *middleware) AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			fmt.Println("No credentials")
			http.Error(w, "No credentials", http.StatusUnauthorized)
			return
		}

		parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
				return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
			}
			return m.PublicKey, nil
		})

		if err != nil || !parsed.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := parsed.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		if exp, isOk := claims["exp"].(float64); isOk {
			if int64(exp) < time.Now().Unix() {
				http.Error(w, "Token Expired", http.StatusUnauthorized)
				return
			}
		}

		issClaim, ok := claims["iss"].(string)
		if !ok {
			http.Error(w, "Missing or invalid issuer claim", http.StatusUnauthorized)
			return
		}
		if issClaim != os.Getenv("FRONTEND_URL") {
			http.Error(w, "Invalid issuer", http.StatusUnauthorized)
			return
		}

		id, found := claims["id"]

		if !found {
			fmt.Printf("Token is missing id claim")
			http.Error(w, "Token is missing id claim", http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "user-id", id)
		next(w, r.WithContext(ctx))
	}
}

func (m *middleware) WithCORS(next http.HandlerFunc) http.HandlerFunc {
	frontendUrl := os.Getenv("FRONTEND_URL")
	fmt.Printf("furl %v \n", frontendUrl)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", frontendUrl)
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

type JWK struct {
	Crv string `json:"crv"`
	X   string `json:"x"`
	Kty string `json:"kty"`
}

func getPublicKeyFromJWK(jwkJson string) ed25519.PublicKey {
	var jwk JWK
	if err := json.Unmarshal([]byte(jwkJson), &jwk); err != nil {
		fmt.Errorf("failed to parse JWK: %w", err)
		log.Fatal("Failed to parse JWK")
	}

	if jwk.Kty != "OKP" || jwk.Crv != "Ed25519" {
		fmt.Errorf("unsupported key type or curve: %s/%s", jwk.Kty, jwk.Crv)
		log.Fatal("unsupported key type or curve")
	}

	pubBytes, err := base64.RawURLEncoding.DecodeString(jwk.X)
	if err != nil {
		fmt.Errorf("failed to decode x: %w", err)
		log.Fatal("failed to decode")
	}

	if len(pubBytes) != ed25519.PublicKeySize {
		fmt.Errorf("invalid public key size: got %d", len(pubBytes))
		log.Fatal("invalid public key size got :%d", len(pubBytes))
	}

	return pubBytes
}
