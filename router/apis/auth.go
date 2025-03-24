package apis

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"nimblestack/database"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthApi struct {
	Queries   *database.Queries
	jwtSecret []byte
	tokenExp  time.Duration
}

func NewAuthApi(queries *database.Queries, jwtSecret []byte) *AuthApi {
	return &AuthApi{
		Queries:   queries,
		jwtSecret: jwtSecret,
		tokenExp:  time.Hour,
	}
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type tokenClaim struct {
	Email string
	Role  string
}

func (h *AuthApi) generateToken(tokenClaim tokenClaim) (string, error) {
	claims := jwt.MapClaims{
		"email": tokenClaim.Email,
		"role":  tokenClaim.Role,
		"exp":   time.Now().Add(h.tokenExp).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// login
func (h *AuthApi) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid http method", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("form error: %v\n", err)
		return
	}
	req := loginRequest{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	role := r.FormValue("role")
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	switch role {
	case "student":
		user, err := h.Queries.GetStudent(ctx, req.Email)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				log.Printf("Database error: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		tokenString, err := h.generateToken(tokenClaim{Role: role, Email: req.Email})
		if err != nil {
			log.Printf("Token signing error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    tokenString,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(h.tokenExp),
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
		w.Header().Set("HX-Redirect", "/student/dashboard")
	case "supervisor":
		user, err := h.Queries.GetSupervisor(ctx, req.Email)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				log.Printf("Database error: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		tokenString, err := h.generateToken(tokenClaim{Role: role, Email: req.Email})
		if err != nil {
			log.Printf("Token signing error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    tokenString,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(h.tokenExp),
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
		w.Header().Set("HX-Redirect", "/supervisor/dashboard")
	case "coordinator":
		user, err := h.Queries.GetCoordinator(ctx, req.Email)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				log.Printf("Database error: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		tokenString, err := h.generateToken(tokenClaim{Role: role, Email: req.Email})
		if err != nil {
			log.Printf("Token signing error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    tokenString,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(h.tokenExp),
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
		w.Header().Set("HX-Redirect", "/coordinator/dashboard")
	default:
		http.Error(w, "Invalid user role", http.StatusBadRequest)
		log.Printf("invalid user role: %v\n", role)
		return
	}
}

// register
func (h *AuthApi) Reqister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid http method", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("form error: %v\n", err)
		return
	}
	role := r.FormValue("role")
	switch role {
	case "student":
		req := database.CreateStudentParams{
			Email:     r.FormValue("email"),
			Firstname: r.FormValue("firstName"),
			Lastname:  r.FormValue("lastName"),
			Password:  r.FormValue("password"),
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Password hashing error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		req.Password = string(hashedPassword)
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		_, err = h.Queries.GetStudent(ctx, req.Email)
		if err == nil {
			http.Error(w, "Email already registered", http.StatusConflict)
			return
		} else if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("Database error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if _, err := h.Queries.CreateStudent(ctx, req); err != nil {
			log.Printf("User creation error: %v", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		tokenString, err := h.generateToken(tokenClaim{Role: role, Email: req.Email})
		if err != nil {
			log.Printf("Token signing error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    tokenString,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
		w.Header().Set("HX-Redirect", "/student/dashboard")
	case "coordinator":
		req := database.CreateCoordinatorParams{
			Email:     r.FormValue("email"),
			Firstname: r.FormValue("firstName"),
			Lastname:  r.FormValue("lastName"),
			Password:  r.FormValue("password"),
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Password hashing error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		req.Password = string(hashedPassword)
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		_, err = h.Queries.GetStudent(ctx, req.Email)
		if err == nil {
			http.Error(w, "Email already registered", http.StatusConflict)
			return
		} else if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("Database error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if _, err := h.Queries.CreateCoordinator(ctx, req); err != nil {
			log.Printf("User creation error: %v", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		tokenString, err := h.generateToken(tokenClaim{Role: role, Email: req.Email})
		if err != nil {
			log.Printf("Token signing error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    tokenString,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
		w.Header().Set("HX-Redirect", "/coordinator/dashboard")
	case "supervisor":
		req := database.CreateSupervisorParams{
			Email:     r.FormValue("email"),
			Firstname: r.FormValue("firstName"),
			Lastname:  r.FormValue("lastName"),
			Password:  r.FormValue("password"),
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Password hashing error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		req.Password = string(hashedPassword)
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		_, err = h.Queries.GetStudent(ctx, req.Email)
		if err == nil {
			http.Error(w, "Email already registered", http.StatusConflict)
			return
		} else if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("Database error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if _, err := h.Queries.CreateSupervisor(ctx, req); err != nil {
			log.Printf("User creation error: %v", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		tokenString, err := h.generateToken(tokenClaim{Role: role, Email: req.Email})
		if err != nil {
			log.Printf("Token signing error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    tokenString,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
		w.Header().Set("HX-Redirect", "/supervisor/dashboard")
	default:
		http.Error(w, "Invalid user role", http.StatusBadRequest)
		log.Printf("invalid user role: %v\n", role)
	}
}
