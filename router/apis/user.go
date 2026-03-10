package apis

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"nimblestack/database"
	"nimblestack/router/middleware"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserApi struct {
	queries *database.Queries
}

type userResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
}

func NewUserApi(queries *database.Queries) *UserApi {
	return &UserApi{
		queries: queries,
	}
}

func (h *UserApi) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value(middleware.JWTClaimsKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	email, ok := claims["email"].(string)
	if !ok {
		http.Error(w, `{"error":"invalid token claims"}`, http.StatusUnauthorized)
		return
	}

	user, err := h.queries.GetUserByEmail(ctx, email)
	if err != nil {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		log.Printf("user not found: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse{
		ID:        strconv.FormatInt(user.ID, 10),
		Email:     user.Email,
	})
}
