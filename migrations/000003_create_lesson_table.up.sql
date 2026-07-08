CREATE TABLE lesson (
    id UUID PRIMARY KEY,
    course_id UUID NOT NULL,
    title VARCHAR(50) NOT NULL,
    description TEXT,
    order_no INT NOT NULL,
    youtube_url TEXT,
    content TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_lesson_course
        FOREIGN KEY (course_id)
        REFERENCES course(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_lesson_course_order
        UNIQUE (course_id, order_no)
);

CREATE INDEX idx_lesson_course_id ON lesson(course_id);