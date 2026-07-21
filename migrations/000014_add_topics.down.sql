DROP TABLE IF EXISTS topics;

ALTER TABLE lesson
    ADD COLUMN description TEXT,
    ADD COLUMN content TEXT,
    ADD COLUMN youtube_url TEXT;
