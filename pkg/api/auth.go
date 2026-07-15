package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var tokenSigningKey = []byte("scheduler_secret_signature_2026")

func getPassChecksum(pass string) string {
	sum := sha256.Sum256([]byte(pass))
	return hex.EncodeToString(sum[:])
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json payload"})
		return
	}

	systemPass := os.Getenv("TODO_PASSWORD")
	if input.Password != systemPass {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверный пароль"})
		return
	}

	tokenClaims := jwt.MapClaims{
		"pass_hash": getPassChecksum(systemPass),
		"exp":       time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	signedStr, err := token.SignedString(tokenSigningKey)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to sign token"})
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{"token": signedStr})
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requiredPass := os.Getenv("TODO_PASSWORD")
		if requiredPass == "" {
			next(w, r)
			return
		}

		authCookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		parsedToken, err := jwt.Parse(authCookie.Value, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return tokenSigningKey, nil
		})

		if err != nil || !parsedToken.Valid {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		storedHash, _ := claims["pass_hash"].(string)
		if storedHash != getPassChecksum(requiredPass) {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}
