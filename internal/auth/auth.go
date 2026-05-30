package auth

import (
	"errors"
	"fmt"
	"lab/api/internal/config"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func FormAuthValidation(r *http.Request, cfg config.Config) bool {
	error := r.ParseForm()
	if error != nil {
		fmt.Println("Failed to parse form in auth")
		return false
	}
	formData := r.Form
	userName := formData.Get("username")
	passWord := formData.Get("password")

	if userName == cfg.SecretUserName && passWord == cfg.SecretPassword {
		return true
	}

	return false
}

func CreateJWT(cfg config.Config) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nbf": time.Now().Unix(),
	})
	if cfg.HmacSampleSecret == "" {
		log.Fatal("No HmacSampleSecret")
	}
	hmacSecret := []byte(cfg.HmacSampleSecret)
	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		fmt.Println("Error creating jwt token")
	}
	return tokenString
}

func ValidateJWT(tokenString string, cfg config.Config) bool {
	hmacSecret := []byte(cfg.HmacSampleSecret)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return hmacSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		log.Fatal(err)
	}
	return token.Valid
}

func SetJWTCookie(w http.ResponseWriter, cfg config.Config) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    CreateJWT(cfg),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600,
	}
	http.SetCookie(w, cookie)
}

func GetJwtFromRequest(r *http.Request) (string, error) {
	cookie, error := r.Cookie("auth_token")
	if error != nil {
		fmt.Println("Failed to get auth cookie")
		return "", errors.New("Failed to get auth cookie")
	}
	return cookie.Value, nil
}

func AuthenticateRequest(r *http.Request, w http.ResponseWriter, cfg config.Config) bool {
	jwtToken, error := GetJwtFromRequest(r)
	if error != nil {
		http.Redirect(w, r, "/logIn", http.StatusSeeOther)
		return false
	}
	isValid := ValidateJWT(jwtToken, cfg)
	if !isValid {
		http.Redirect(w, r, "/logIn", http.StatusSeeOther)
		return false
	}
	return true
}
