CREATE TABLE attendance_session (
    id UUID PRIMARY KEY,
    batch_id UUID NOT NULL,
    lesson_date DATE NOT NULL,
    created_by UUID,

    CONSTRAINT fk_attendance_session_batch
        FOREIGN KEY (batch_id)
        REFERENCES batches(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_attendance_session_created_by
        FOREIGN KEY (created_by)
        REFERENCES users(id)
        ON DELETE SET NULL,

    CONSTRAINT uq_attendance_session
        UNIQUE (batch_id, lesson_date)
);

CREATE TABLE attendance (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL,
    student_id UUID NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('present', 'absent')),

    CONSTRAINT fk_attendance_session
        FOREIGN KEY (session_id)
        REFERENCES attendance_session(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_attendance_student
        FOREIGN KEY (student_id)
        REFERENCES students(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_attendance
        UNIQUE (session_id, student_id)
);

CREATE INDEX idx_attendance_session_batch
    ON attendance_session(batch_id);

CREATE INDEX idx_attendance_student
    ON attendance(student_id);
