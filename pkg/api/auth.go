package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

const tokenLife = 8 * time.Hour

type signinRequest struct {
	Password string `json:"password"`
}

type tokenPayload struct {
	Hash string `json:"hash"`
	Exp  int64  `json:"exp"`
}

func (a *App) signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if a.password == "" {
		writeError(w, http.StatusBadRequest, "пароль не задан")
		return
	}

	var req signinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Password != a.password {
		writeError(w, http.StatusUnauthorized, "неверный пароль")
		return
	}

	token, err := makeToken(a.password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}

func (a *App) auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if a.password == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil || !validateToken(cookie.Value, a.password) {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func makeToken(password string) (string, error) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))

	payload := tokenPayload{
		Hash: passwordHash(password),
		Exp:  time.Now().Add(tokenLife).Unix(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	payloadPart := base64.RawURLEncoding.EncodeToString(body)
	signingInput := header + "." + payloadPart
	signature := sign(signingInput, password)

	return signingInput + "." + signature, nil
}

func validateToken(token string, password string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}

	signingInput := parts[0] + "." + parts[1]
	expected := sign(signingInput, password)

	if !hmac.Equal([]byte(parts[2]), []byte(expected)) {
		return false
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	var payload tokenPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return false
	}

	if payload.Hash != passwordHash(password) {
		return false
	}

	if time.Now().Unix() > payload.Exp {
		return false
	}

	return true
}

func sign(data, password string) string {
	mac := hmac.New(sha256.New, []byte(password))
	_, _ = mac.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func passwordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}
