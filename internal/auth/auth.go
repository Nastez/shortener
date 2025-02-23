package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var secretKey = []byte("supersecretkey")

func WithAuth(w http.ResponseWriter, r *http.Request) (string, error) {
	var userID string
	cookie, err := r.Cookie("user_id")
	if err != nil {
		// Создаем новую куку, если кука не найдена
		userIDCookie := uuid.NewString()
		signature := signValue(userIDCookie)
		cookieValue := fmt.Sprintf("%s|%s", userIDCookie, signature)

		// устанавливаем куку в запрос
		http.SetCookie(w, &http.Cookie{
			Name:    "user_id",
			Value:   cookieValue,
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
		})

		return userIDCookie, nil
	}

	// Проверяем подпись куки
	parts := splitCookie(cookie.Value)
	if len(parts) != 2 || !validateCookie(parts[0], parts[1]) {
		http.Error(w, "Invalid cookie", http.StatusUnauthorized)
		return "", err
	}
	userID = parts[0]

	return userID, nil
}

// signValue подписывает значение с помощью HMAC-SHA256
func signValue(value string) string {
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}

// validateCookie проверяет подпись куки
func validateCookie(value, signature string) bool {
	expected := signValue(value)
	return hmac.Equal([]byte(expected), []byte(signature))
}

// splitCookie разбивает значение куки на ID и подпись
func splitCookie(value string) []string {
	parts := make([]string, 2)
	sep := "|"
	for i, v := range []byte(value) {
		if string(v) == sep {
			parts[0] = value[:i]
			parts[1] = value[i+1:]
			break
		}
	}
	return parts
}
