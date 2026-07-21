ALTER TABLE lesson
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS content,
    DROP COLUMN IF EXISTS youtube_url;

CREATE TABLE topics (
    id           UUID PRIMARY KEY,
    lesson_id    UUID NOT NULL,
    title        VARCHAR(150) NOT NULL,
    description  TEXT,
    content      TEXT,
    youtube_url  TEXT,
    order_no     INT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_topics_lesson
        FOREIGN KEY (lesson_id)
        REFERENCES lesson(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_topics_lesson_order
        UNIQUE (lesson_id, order_no)
        DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX idx_topics_lesson_id ON topics(lesson_id);
