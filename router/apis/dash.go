package apis

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"nimblestack/database"
	"strconv"
	"time"
)

type DashApi struct {
	Queries *database.Queries
}

func NewDashApi(queries *database.Queries) *DashApi {
	return &DashApi{
		Queries: queries,
	}
}

type safeStudent struct {
	Id        int64
	Firstname string
	Lastname  string
	Email     string
}

type safeSupervisor struct {
	Id        int64
	Firstname string

	Lastname string
	Email    string
}

type projectsData struct {
	Total     int
	OnTrack   int
	Delayed   int
	AtRisk    int
	Completed int
}

func (h *DashApi) GetAllStudents(w http.ResponseWriter, r *http.Request) {
	var safeStudents []safeStudent
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	students, err := h.Queries.GetAllStudents(ctx)
	if err != nil {
		http.Error(w, "Error with getting students", http.StatusInternalServerError)
		log.Println("Error with getting students: ", err)
		return
	}
	for _, student := range students {
		if student.Supervisorid.Int64 == 0 {
			safeStudent := safeStudent{
				Id:        student.Studentid,
				Firstname: student.Firstname,
				Lastname:  student.Lastname,
				Email:     student.Email,
			}
			safeStudents = append(safeStudents, safeStudent)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(safeStudents); err != nil {
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
		log.Printf("Failed to encode response: %v", err)
	}
}

func (h *DashApi) GetAllSupervisors(w http.ResponseWriter, r *http.Request) {
	var safeSupervisors []safeSupervisor
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	supervisors, err := h.Queries.GetAllSupervisors(ctx)
	if err != nil {
		http.Error(w, "Error with getting supervisors", http.StatusInternalServerError)
		log.Println("Error with getting supervisors: ", err)
		return
	}
	for _, supervisor := range supervisors {
		safeSupervisor := safeSupervisor{
			Id:        supervisor.Supervisorid,
			Email:     supervisor.Email,
			Firstname: supervisor.Firstname,
			Lastname:  supervisor.Lastname,
		}
		safeSupervisors = append(safeSupervisors, safeSupervisor)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(safeSupervisors); err != nil {
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
		log.Printf("Failed to encode response: %v", err)
	}
}

func (h *DashApi) AssignSupervisor(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		StudentId    int64 `json:"studentId"`
		SupervisorId int64 `json:"supervisorId"`
	}
	body := requestBody{
		StudentId:    0,
		SupervisorId: 0,
	}
	var err error
	body.StudentId, err = strconv.ParseInt(r.FormValue("studentId"), 10, 64)
	if err != nil {
		http.Error(w, "Error with getting studentId", http.StatusBadRequest)
		log.Println("Error with getting studentId: ", err)
		return
	}
	body.SupervisorId, err = strconv.ParseInt(r.FormValue("supervisorId"), 10, 64)
	if err != nil {
		http.Error(w, "Error with getting supervisorId", http.StatusBadRequest)
		log.Println("Error with getting supervisorId: ", err)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err = h.Queries.AssignSupervisor(ctx, database.AssignSupervisorParams{
		Studentid:    body.StudentId,
		Supervisorid: sql.NullInt64{Int64: body.SupervisorId, Valid: true},
	})
	if err != nil {
		http.Error(w, "Error with assigning supervisor", http.StatusInternalServerError)
		log.Println("Error with assigning supervisor: ", err)
		return
	}
	w.Header().Set("HX-Redirect", "/coordinator/dashboard")
}

func (h *DashApi) GetProjectsData(w http.ResponseWriter, r *http.Request) {
	var projectsData projectsData
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	projects, err := h.Queries.GetAllProjects(ctx)
	if err != nil {
		http.Error(w, "Error with getting projects", http.StatusInternalServerError)
		log.Println("Error with getting projects: ", err)
		return
	}
	projectsData.Total = len(projects)
	studentMilestones, err := h.Queries.GetAllStudentMilestones(ctx)
	students := make(map[int64]bool)
	for _, studentMilestone := range studentMilestones {
		if !students[studentMilestone.Studentid] {
			switch studentMilestone.Status {
			case sql.NullString{String: "on track", Valid: true}:
				projectsData.OnTrack++
			case sql.NullString{String: "delayed", Valid: true}:
				projectsData.Delayed++
			case sql.NullString{String: "at risk", Valid: true}:
				projectsData.AtRisk++
			case sql.NullString{String: "completed", Valid: true}:
				projectsData.Completed++
			}
			students[studentMilestone.Studentid] = true
		} else {

		}
	}

	if err != nil {
		http.Error(w, "Error with getting student milestones", http.StatusInternalServerError)
		log.Println("Error with getting student milestones: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(projectsData); err != nil {
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
		log.Printf("Failed to encode response: %v", err)
	}
}
