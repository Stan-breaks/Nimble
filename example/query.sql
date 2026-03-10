-- ============================================================
-- Example: Advanced Query Patterns
-- ============================================================
-- These queries demonstrate patterns for working with the
-- advanced schema. Copy into sqlc/query.sql and adapt.
-- ============================================================

-- ┌──────────────────────────────────┐
-- │  1. Role-Based CRUD              │
-- │  Same pattern, different tables  │
-- └──────────────────────────────────┘

-- name: CreateStudent :one
INSERT INTO students (email, first_name, last_name, password)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetStudentByEmail :one
SELECT * FROM students WHERE email = ? LIMIT 1;

-- name: GetStudentById :one
SELECT * FROM students WHERE id = ? LIMIT 1;

-- name: GetAllStudents :many
SELECT * FROM students;

-- name: CreateSupervisor :one
INSERT INTO supervisors (email, first_name, last_name, password)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetSupervisorByEmail :one
SELECT * FROM supervisors WHERE email = ? LIMIT 1;

-- ┌──────────────────────────────────┐
-- │  2. Assignment / Update Queries  │
-- │  Update a foreign key field     │
-- └──────────────────────────────────┘

-- Assign a supervisor to a student (nullable FK update)
-- name: AssignSupervisor :exec
UPDATE students
SET supervisor_id = ?
WHERE id = ?;

-- ┌──────────────────────────────────┐
-- │  3. Filtered Queries             │
-- │  Query by FK or status          │
-- └──────────────────────────────────┘

-- Get all students without a supervisor (unassigned)
-- name: GetUnassignedStudents :many
SELECT * FROM students
WHERE supervisor_id IS NULL;

-- Get milestones for a specific student
-- name: GetStudentMilestonesByStudentId :many
SELECT * FROM student_milestones
WHERE student_id = ?;

-- ┌──────────────────────────────────┐
-- │  4. Insert with Defaults         │
-- │  Let SQLite fill in defaults    │
-- └──────────────────────────────────┘

-- Status defaults to 'pending', submitted_at stays NULL
-- name: CreateStudentMilestone :one
INSERT INTO student_milestones (student_id, milestone_id, status, submitted_at)
VALUES (?, ?, ?, ?)
RETURNING *;

-- Milestone with sequence_order and optional fields
-- name: CreateSupervisorMilestone :one
INSERT INTO supervisor_milestones (supervisor_id, name, description, due_date)
VALUES (?, ?, ?, ?)
RETURNING *;

-- ┌──────────────────────────────────┐
-- │  5. Bulk Fetch for Aggregation   │
-- │  Get all rows, aggregate in Go  │
-- └──────────────────────────────────┘

-- Fetch all milestones — aggregate counts by status in Go
-- name: GetAllStudentMilestones :many
SELECT * FROM student_milestones;

-- Fetch all projects for dashboard stats
-- name: GetAllProjects :many
SELECT * FROM projects;
