-- Denormalized score so a submission keeps its value even if the quiz is later
-- edited (which replaces questions/answers). Auto-computed on submit and grade.
ALTER TABLE quiz_submissions
    ADD COLUMN score INTEGER NOT NULL DEFAULT 0;
