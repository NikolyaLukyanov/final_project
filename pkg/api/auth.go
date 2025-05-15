package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtKey     = []byte("super-secret")
	cookieName = "token"
)

func signinHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Ошибка чтения запроса", http.StatusBadRequest)
		return
	}

	expected := os.Getenv("TODO_PASSWORD") // Читается каждый раз — но после godotenv.Load()
	if expected == "" || body.Password != expected {
		http.Error(w, "Неверный пароль", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": expected,
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    tokenString,
		Path:     "/",
		Expires:  time.Now().Add(8 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		expected := os.Getenv("TODO_PASSWORD")
		if expected == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie(cookieName)
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["hash"] != expected {
			http.Error(w, "Access denied", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
