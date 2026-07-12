CREATE TABLE students (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    full_name TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    dob DATE NOT NULL,
    gender TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE batches (
    id UUID PRIMARY KEY,
    batch_name TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    capacity INTEGER NOT NULL CHECK (capacity > 0),
    days TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE student_batches (
    id UUID PRIMARY KEY,
    student_id UUID NOT NULL,
    batch_id UUID NOT NULL,
    enrolled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    approved_by UUID,

    CONSTRAINT fk_student_batches_student
        FOREIGN KEY (student_id)
        REFERENCES students(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_student_batches_batch
        FOREIGN KEY (batch_id)
        REFERENCES batches(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_student_batches_approved_by
        FOREIGN KEY (approved_by)
        REFERENCES users(id)
        ON DELETE SET NULL,

    CONSTRAINT uq_student_batch
        UNIQUE (student_id, batch_id)
);

CREATE TABLE student_course_access (
    id UUID PRIMARY KEY,
    batch_id UUID NOT NULL,
    course_id UUID NOT NULL,
    granted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    granted_by UUID,

    CONSTRAINT fk_student_course_access_batch
        FOREIGN KEY (batch_id)
        REFERENCES batches(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_student_course_access_course
        FOREIGN KEY (course_id)
        REFERENCES course(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_student_course_access_granted_by
        FOREIGN KEY (granted_by)
        REFERENCES users(id)
        ON DELETE SET NULL,

    CONSTRAINT uq_batch_course
        UNIQUE (batch_id, course_id)
);

CREATE INDEX idx_students_email
    ON students(email);

CREATE INDEX idx_student_batches_student
    ON student_batches(student_id);

CREATE INDEX idx_student_batches_batch
    ON student_batches(batch_id);

CREATE INDEX idx_student_course_access_batch
    ON student_course_access(batch_id);

CREATE INDEX idx_student_course_access_course
    ON student_course_access(course_id);