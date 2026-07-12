CREATE TABLE lesson_progress (
    id UUID PRIMARY KEY,
    student_id UUID NOT NULL,
    lesson_id UUID NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT TRUE,
    completed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_lesson_progress_student
        FOREIGN KEY (student_id)
        REFERENCES students(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_lesson_progress_lesson
        FOREIGN KEY (lesson_id)
        REFERENCES lesson(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_lesson_progress
        UNIQUE (student_id, lesson_id)
);

CREATE INDEX idx_lesson_progress_student
    ON lesson_progress(student_id);
