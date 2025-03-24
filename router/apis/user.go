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
	Queries *database.Queries
}

type safeUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Role      string `json:"role"`
}

func NewUserApi(queries *database.Queries) *UserApi {
	return &UserApi{
		Queries: queries,
	}
}

func (h *UserApi) GetCurrentUSer(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.JWTClaimsKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	switch claims["role"].(string) {
	case "student":
		user, err := h.Queries.GetStudent(ctx, claims["email"].(string))
		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			log.Printf("user not found: %v\n", err.Error())
			return
		}
		safeUser := safeUser{
			ID:        strconv.FormatInt(user.Studentid, 10),
			Email:     user.Email,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Role:      claims["role"].(string),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(safeUser)
	case "supervisor":
		user, err := h.Queries.GetSupervisor(ctx, claims["email"].(string))
		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			log.Printf("user not found: %v\n", err.Error())
			return
		}

		safeUser := safeUser{
			ID:        strconv.FormatInt(user.Supervisorid, 10),
			Email:     user.Email,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Role:      claims["role"].(string),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(safeUser)

	case "coordinator":
		user, err := h.Queries.GetCoordinator(ctx, claims["email"].(string))
		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			log.Printf("user not found: %v\n", err.Error())
			return
		}
		safeUser := safeUser{
			ID:        strconv.FormatInt(user.Coordinatorid, 10),
			Email:     user.Email,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Role:      claims["role"].(string),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(safeUser)

	default:
		http.Error(w, "invalid user role", http.StatusBadRequest)
		return
	}
}
