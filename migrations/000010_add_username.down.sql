ALTER TABLE users ALTER COLUMN email SET NOT NULL;
ALTER TABLE users DROP COLUMN username;

ALTER TABLE students ADD CONSTRAINT students_email_key UNIQUE (email);
ALTER TABLE students DROP COLUMN username;
