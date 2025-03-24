package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
)

type contextKey string

const JWTClaimsKey contextKey = "jwtClaims"

func CheckAuth(next http.HandlerFunc, jwtSecret []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			http.Redirect(w, r, "/auth/login", http.StatusUnauthorized)
			return
		}
		tokenString := cookie.Value
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.SetCookie(w, &http.Cookie{
				Name:     "auth_token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			})
			http.Redirect(w, r, "auth/login", http.StatusSeeOther)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Redirect(w, r, "/auth/login", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), JWTClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
