-- name: CreateStudent :one
INSERT INTO students(
email,firstName,lastName,password
) VALUES (?,?,?,?) RETURNING *;
-- name: GetStudent :one
SELECT * FROM students WHERE email =? LIMIT 1;
-- name: GetStudentById :one
SELECT * FROM students WHERE studentId =? LIMIT 1;
-- name: GetAllStudents :many
SELECT * FROM students;
-- name: CreateCoordinator :one
INSERT INTO coordinators(
email,firstName,lastName,password
) VALUES (?,?,?,?) RETURNING *;
-- name: GetCoordinator :one
SELECT * FROM coordinators WHERE email =? LIMIT 1;
-- name: GetCoordinatorById :one
SELECT * FROM coordinators WHERE coordinatorId =? LIMIT 1;
-- name: CreateSupervisor :one
INSERT INTO supervisors(
email,firstName,lastName,password
) VALUES (?,?,?,?) RETURNING *;
-- name: GetSupervisor :one
SELECT * FROM supervisors WHERE email =? LIMIT 1;
-- name: GetSupervisorById :one
SELECT * FROM supervisors WHERE supervisorID =? LIMIT 1;
-- name: GetAllSupervisors :many
SELECT * FROM supervisors;
-- name: CreateProject :one
INSERT INTO projects(name,description
) VALUES (?,?)RETURNING *;
-- name: GetProject :one
SELECT * FROM projects WHERE projectId =? LIMIT 1;
-- name: GetAllProjects :many
SELECT * FROM projects;
-- name: AssignSupervisor :exec
UPDATE students SET supervisorId =? WHERE studentId =?;
-- name: CreateSupervisorMilestone :one
INSERT INTO supervisor_milestones(
supervisorId,name,description,due_date
) VALUES (?,?,?,?) RETURNING *;
-- name: GetSupervisorMilestone :one
SELECT * FROM supervisor_milestones WHERE milestoneId =? LIMIT 1;
-- name: GetAllSupervisorMilestones :many
SELECT * FROM supervisor_milestones;
-- name: CreateStudentMilestone :one
INSERT INTO student_milestones(
studentId,milestoneId,status,submitted_at
) VALUES (?,?,?,?) RETURNING *;
-- name: GetStudentMilestone :one
SELECT * FROM student_milestones WHERE milestoneId =? LIMIT 1;
-- name: GetAllStudentMilestones :many
SELECT * FROM student_milestones;
-- name: GetStudentMilestonesByStudentId :many
SELECT * FROM student_milestones WHERE studentId =?;
