CREATE TABLE assignments (
    id UUID PRIMARY KEY,
    lesson_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    starter_code TEXT,
    due_date DATE,
    total_marks INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_assignments_lesson
        FOREIGN KEY (lesson_id)
        REFERENCES lesson(id)
        ON DELETE CASCADE
);

CREATE TABLE assignment_submissions (
    id UUID PRIMARY KEY,
    student_id UUID NOT NULL,
    assignment_id UUID NOT NULL,
    code TEXT NOT NULL,
    remarks TEXT,
    marks INTEGER,
    status TEXT NOT NULL DEFAULT 'submitted',
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_assignment_submissions_student
        FOREIGN KEY (student_id)
        REFERENCES students(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_assignment_submissions_assignment
        FOREIGN KEY (assignment_id)
        REFERENCES assignments(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_assignment_submission
        UNIQUE (student_id, assignment_id)
);

CREATE INDEX idx_assignments_lesson
    ON assignments(lesson_id);

CREATE INDEX idx_assignment_submissions_assignment
    ON assignment_submissions(assignment_id);
