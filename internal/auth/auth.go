package auth

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/microcosm-cc/bluemonday"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Middleware interface {
	WithCORS(next http.HandlerFunc) http.HandlerFunc
	AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc
	AuthenticateStoreMiddleware(next http.HandlerFunc) http.HandlerFunc
	CleanXSS(next http.HandlerFunc) http.HandlerFunc
	WithSecurityHeaders(next http.HandlerFunc) http.HandlerFunc
}

type middleware struct {
	PublicKey ed25519.PublicKey
	Sanitizer *bluemonday.Policy
}

func NewMiddleware() Middleware {
	return &middleware{
		PublicKey: getPublicKeyFromJWK(os.Getenv("JWKS_PUBLIC_KEY")),
		Sanitizer: bluemonday.StrictPolicy(),
	}
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
		env := os.Getenv("ENVIRONMENT")
		if env == "local" || env == "dev" {
			w.Header().Set("Access-Control-Allow-Origin", frontendUrl)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent) // 204: OK but no body
			return
		}
		next(w, r)
	}
}

// TODO
func (m *middleware) AuthenticateStoreMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "store-id", 1)
		next(w, r.WithContext(ctx))
	}
}

func (m *middleware) CleanXSS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			m.WithCORS(next)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		contentType := r.Header.Get("Content-Type")
		if contentType == "" {
			log.Printf("Rejected %s request to %s: missing Content-Type", r.Method, r.URL.Path)
			http.Error(w, "Missing Content-Type header", http.StatusUnsupportedMediaType)
			return
		}

		if r.Method != http.MethodOptions && !strings.HasPrefix(contentType, "application/json") {
			log.Printf("Rejected %s request to %s: unsupported Content-Type %q", r.Method, r.URL.Path, contentType)
			http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Could not read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if len(body) == 0 {
			// Allow empty JSON body (e.g. POST with no payload)
			r.Body = io.NopCloser(bytes.NewBuffer([]byte{}))
			next.ServeHTTP(w, r)
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err = sanitizeJSON(data, m.Sanitizer, 0); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		cleanedBody, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "Could not re-encode JSON", http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(cleanedBody))
		r.ContentLength = int64(len(cleanedBody))

		next.ServeHTTP(w, r)
		next(w, r)
	}
}

func (m *middleware) WithSecurityHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
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
		log.Fatal("Failed to parse JWK")
	}

	if jwk.Kty != "OKP" || jwk.Crv != "Ed25519" {
		log.Fatal("unsupported key type or curve")
	}

	pubBytes, err := base64.RawURLEncoding.DecodeString(jwk.X)
	if err != nil {
		log.Fatal("failed to decode")
	}

	if len(pubBytes) != ed25519.PublicKeySize {
		log.Fatal("invalid public key size got :%d", len(pubBytes))
	}

	return pubBytes
}

func sanitizeJSON(data map[string]interface{}, sanitizer *bluemonday.Policy, depth int) error {
	var err error
	if depth >= 30 {
		return errors.New("object nesting past acceptable limit")
	}
	for k, v := range data {
		switch val := v.(type) {
		case string:
			data[k] = sanitizer.Sanitize(val)
		case map[string]interface{}:
			err = sanitizeJSON(val, sanitizer, depth+1)
			if err != nil {
				return err
			}
		case []interface{}:
			data[k], err = sanitizeJSONArray(val, sanitizer, depth+1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func sanitizeJSONArray(arr []interface{}, sanitizer *bluemonday.Policy, depth int) ([]interface{}, error) {
	var err error
	//dont allow input depth over 30
	if depth >= 30 {
		return nil, errors.New("object nesting past acceptable limit")
	}
	for i, v := range arr {
		switch val := v.(type) {
		case string:
			arr[i] = sanitizer.Sanitize(val)
		case map[string]interface{}:
			err = sanitizeJSON(val, sanitizer, depth+1)
			if err != nil {
				return nil, err
			}
		case []interface{}:
			arr[i], err = sanitizeJSONArray(val, sanitizer, depth+1)
			if err != nil {
				return nil, err
			}
		}
	}
	return arr, nil
}
