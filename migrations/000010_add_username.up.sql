-- Students are now identified by a unique username.
-- Email stays on the record as contact info but is no longer a unique key.
ALTER TABLE students ADD COLUMN username VARCHAR(255) NOT NULL UNIQUE;
ALTER TABLE students DROP CONSTRAINT students_email_key;

-- Login accounts get an optional unique username:
-- students authenticate by username, admins by email.
ALTER TABLE users ADD COLUMN username VARCHAR(255) UNIQUE;

-- Students share emails (e.g. siblings under a parent's email), so student
-- accounts store NO email and log in by username only. Email stays UNIQUE for
-- admins, but must allow NULL for students. (Postgres UNIQUE permits many NULLs.)
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;
