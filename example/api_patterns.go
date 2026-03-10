package example

// ============================================================
// Example: Advanced API Patterns
// ============================================================
// These patterns demonstrate techniques used in real Go APIs
// built on the Nimble template. Reference these when building
// your own handlers.
// ============================================================

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

// ┌──────────────────────────────────┐
// │  1. Role-Based Auth Switch       │
// │  Route logic by user role       │
// └──────────────────────────────────┘
//
// When you have multiple user types (student, admin, etc.)
// that share the same endpoint but hit different tables:
//
//   func (h *AuthApi) Login(w http.ResponseWriter, r *http.Request) {
//       var req loginRequest
//       json.NewDecoder(r.Body).Decode(&req)
//
//       switch req.Role {
//       case "student":
//           user, err := h.queries.GetStudentByEmail(ctx, req.Email)
//           // validate password, generate token...
//       case "supervisor":
//           user, err := h.queries.GetSupervisorByEmail(ctx, req.Email)
//           // validate password, generate token...
//       default:
//           http.Error(w, `{"error":"invalid role"}`, http.StatusBadRequest)
//       }
//   }

// ┌──────────────────────────────────┐
// │  2. Safe Response Types          │
// │  Never leak passwords/internals │
// └──────────────────────────────────┘

// Define a "safe" struct that omits sensitive fields.
// Map from the database model to this before encoding.

type safeStudent struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	// Password is intentionally excluded
}

// func mapToSafe(student database.Student) safeStudent {
//     return safeStudent{
//         ID:        student.ID,
//         FirstName: student.FirstName,
//         LastName:  student.LastName,
//         Email:     student.Email,
//     }
// }

// ┌──────────────────────────────────┐
// │  3. Dashboard Aggregation        │
// │  Fetch all, compute stats in Go │
// └──────────────────────────────────┘
//
// When SQLite doesn't do complex GROUP BY easily with SQLC,
// fetch all rows and aggregate in Go:

type dashboardStats struct {
	Total     int `json:"total"`
	OnTrack   int `json:"on_track"`
	Delayed   int `json:"delayed"`
	AtRisk    int `json:"at_risk"`
	Completed int `json:"completed"`
}

// func (h *DashApi) GetProjectsData(w http.ResponseWriter, r *http.Request) {
//     ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
//     defer cancel()
//
//     milestones, _ := h.queries.GetAllStudentMilestones(ctx)
//
//     stats := dashboardStats{}
//     seen := make(map[int64]bool)
//
//     for _, m := range milestones {
//         if seen[m.StudentID] {
//             continue  // count each student once
//         }
//         seen[m.StudentID] = true
//         stats.Total++
//
//         switch m.Status {
//         case "on track":  stats.OnTrack++
//         case "delayed":   stats.Delayed++
//         case "at risk":   stats.AtRisk++
//         case "completed": stats.Completed++
//         }
//     }
//
//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(stats)
// }

// ┌──────────────────────────────────┐
// │  4. Nullable FK Assignment       │
// │  Update a nullable foreign key  │
// └──────────────────────────────────┘
//
// Use sql.NullInt64 when updating nullable FK columns:
//
//   func (h *DashApi) AssignSupervisor(w http.ResponseWriter, r *http.Request) {
//       var req struct {
//           StudentID    int64 `json:"student_id"`
//           SupervisorID int64 `json:"supervisor_id"`
//       }
//       json.NewDecoder(r.Body).Decode(&req)
//
//       err := h.queries.AssignSupervisor(ctx, database.AssignSupervisorParams{
//           StudentID:    req.StudentID,
//           SupervisorID: sql.NullInt64{Int64: req.SupervisorID, Valid: true},
//       })
//   }

// ┌──────────────────────────────────┐
// │  5. Context Timeouts             │
// │  Always set deadlines on DB ops │
// └──────────────────────────────────┘

// ExampleHandler demonstrates proper context timeout usage.
func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	// Short timeout for simple queries
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Longer timeout for writes or multi-step operations
	ctx2, cancel2 := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel2()

	// Use ctx/ctx2 with your queries...
	_ = ctx
	_ = ctx2
}

// ┌──────────────────────────────────┐
// │  6. Duplicate Check Pattern      │
// │  Check-then-insert with ErrNoRows│
// └──────────────────────────────────┘

// ExampleDuplicateCheck shows how to check if a record exists
// before inserting, using sql.ErrNoRows.
func ExampleDuplicateCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	email := "user@example.com"

	// Try to find existing user
	_, err := queryUserByEmail(ctx, email) // your SQLC-generated function
	if err == nil {
		// User exists — conflict
		http.Error(w, `{"error":"email already registered"}`, http.StatusConflict)
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		// Actual database error
		log.Printf("Database error: %v", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// err IS sql.ErrNoRows → user doesn't exist, safe to create
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ready to create"})
}

// Stub for compilation
func queryUserByEmail(_ context.Context, _ string) (any, error) {
	return nil, sql.ErrNoRows
}

// Ensure imports are used
var _ = time.Second
var _ = log.Println
