package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func FormAuthValidation(r *http.Request, config Config) bool {
	error := r.ParseForm()
	if error != nil {
		fmt.Println("Failed to parse form in auth")
		return false
	}
	formData := r.Form
	userName := formData.Get("username")
	passWord := formData.Get("password")

	if userName == config.SecretUserName && passWord == config.SecretPassword {
		return true
	}

	return false
}

func createJWT(config Config) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nbf": time.Now().Unix(),
	})
	if config.HmacSampleSecret == "" {
		log.Fatal("No HmacSampleSecret")
	}
	hmacSampleSecret := []byte(config.HmacSampleSecret)
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		fmt.Println("Error creating jwt token")
	}
	return tokenString
}

func validateJWT(tokenString string, config Config) bool {
	hmacSampleSecret := []byte(config.HmacSampleSecret)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return hmacSampleSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		log.Fatal(err)
	}
	return token.Valid
}

func setJWTCookie(w http.ResponseWriter, config Config) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    createJWT(config),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600,
	}
	http.SetCookie(w, cookie)
}

func getJwtFromRequest(r *http.Request) (string, error) {
	cookie, error := r.Cookie("auth_token")
	if error != nil {
		fmt.Println("Failed to get auth cookie")
		return "", errors.New("Failed to get auth cookie")
	}
	return cookie.Value, nil
}

func AuthenticateRequest(r *http.Request, w http.ResponseWriter, config Config) bool {
	jwtToken, error := getJwtFromRequest(r)
	if error != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return false
	}
	isValid := validateJWT(jwtToken, config)
	if !isValid {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return false
	}
	return true
}
