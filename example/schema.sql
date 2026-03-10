-- ============================================================
-- Example: Advanced Schema Patterns
-- ============================================================
-- This file demonstrates patterns you can use when extending
-- your Nimble project with more complex table relationships.
-- Copy what you need into sqlc/schema.sql and adapt it.
-- ============================================================

-- ┌──────────────────────────────────┐
-- │  1. Role-Based Tables            │
-- │  Separate tables per user role   │
-- │  with shared fields              │
-- └──────────────────────────────────┘

CREATE TABLE IF NOT EXISTS coordinators (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  email TEXT NOT NULL UNIQUE,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS supervisors (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  email TEXT NOT NULL UNIQUE,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS students (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  email TEXT NOT NULL UNIQUE,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  password TEXT NOT NULL,
  supervisor_id INTEGER,
  project_id INTEGER,
  FOREIGN KEY (supervisor_id) REFERENCES supervisors(id),
  FOREIGN KEY (project_id) REFERENCES projects(id)
);

-- ┌──────────────────────────────────┐
-- │  2. Foreign Key Relationships    │
-- │  One-to-many and many-to-one    │
-- └──────────────────────────────────┘

CREATE TABLE IF NOT EXISTS projects (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  description TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- A supervisor defines milestone templates.
-- Students track their progress against these templates.

CREATE TABLE IF NOT EXISTS supervisor_milestones (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  supervisor_id INTEGER NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  submission_file_name TEXT,
  due_date DATE,
  sequence_order INTEGER,           -- ordering without gaps
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (supervisor_id) REFERENCES supervisors(id)
);

-- ┌──────────────────────────────────┐
-- │  3. Status Tracking Table        │
-- │  Track progress with defaults   │
-- └──────────────────────────────────┘

CREATE TABLE IF NOT EXISTS student_milestones (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  student_id INTEGER NOT NULL,
  milestone_id INTEGER NOT NULL,
  status TEXT DEFAULT 'pending',     -- 'pending' | 'in-progress' | 'completed'
  submitted_at DATETIME,
  FOREIGN KEY (student_id) REFERENCES students(id),
  FOREIGN KEY (milestone_id) REFERENCES supervisor_milestones(id)
);
