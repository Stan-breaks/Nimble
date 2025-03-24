-- Coordinators table: Manages assignment of supervisors to students.
CREATE TABLE coordinators (
  coordinatorId INTEGER PRIMARY KEY AUTOINCREMENT,
  firstName TEXT NOT NULL,
  lastName TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL
);
-- Supervisors table: holds supervisor info
CREATE TABLE supervisors (
  supervisorId INTEGER PRIMARY KEY AUTOINCREMENT,
  firstName TEXT NOT NULL,
  lastName TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL
);

-- Projects table: holds project details
CREATE TABLE projects (
  projectId INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  description TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Students table: each student is linked to a supervisor and (optionally) a project
CREATE TABLE students (
  studentId INTEGER PRIMARY KEY AUTOINCREMENT,
  email TEXT NOT NULL UNIQUE,
  firstName TEXT NOT NULL,
  lastName TEXT NOT NULL,
  password TEXT NOT NULL,
  supervisorId INTEGER,
  projectId INTEGER,
  FOREIGN KEY (supervisorId) REFERENCES supervisors(supervisor_id),
  FOREIGN KEY (projectId) REFERENCES projects(project_id)
);

-- Supervisor Milestones table: milestone templates defined by a supervisor.
CREATE TABLE supervisor_milestones (
  milestoneId INTEGER PRIMARY KEY AUTOINCREMENT,
  supervisorId INTEGER NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  submission_fileName TEXT,
  due_date DATE,
  sequence_order INTEGER,  -- optional: order in which milestones should be completed
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (supervisorId) REFERENCES supervisors(supervisor_id)
);

-- Student Milestones table: tracks each student’s progress on a given milestone template.
CREATE TABLE student_milestones (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  studentId INTEGER NOT NULL,
  milestoneId INTEGER NOT NULL,
  status TEXT DEFAULT 'pending',  -- e.g., 'pending', 'in-progress', 'completed'
  submitted_at DATETIME,
  FOREIGN KEY (studentId) REFERENCES students(student_id),
  FOREIGN KEY (milestoneId) REFERENCES supervisor_milestones(milestone_id)
);

