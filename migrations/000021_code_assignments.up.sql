-- Coding assignments: a language flag, hidden test cases, and autograde results.

ALTER TABLE assignments
    ADD COLUMN language TEXT NOT NULL DEFAULT '';

ALTER TABLE assignment_submissions
    ADD COLUMN auto_score   INTEGER,   -- scaled to total_marks; NULL until graded
    ADD COLUMN tests_total  INTEGER,
    ADD COLUMN tests_passed INTEGER;

-- Hidden test cases (expected_stdout is never exposed to students).
CREATE TABLE assignment_test_cases (
    id              UUID PRIMARY KEY,
    assignment_id   UUID NOT NULL REFERENCES assignments(id) ON DELETE CASCADE,
    stdin           TEXT NOT NULL DEFAULT '',
    expected_stdout TEXT NOT NULL DEFAULT '',
    weight          INTEGER NOT NULL DEFAULT 1,
    ordinal         INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_test_cases_assignment ON assignment_test_cases(assignment_id);

-- Per-test outcome for a submission (what the student/teacher can review).
CREATE TABLE submission_test_results (
    id            UUID PRIMARY KEY,
    submission_id UUID NOT NULL REFERENCES assignment_submissions(id) ON DELETE CASCADE,
    test_case_id  UUID NOT NULL REFERENCES assignment_test_cases(id) ON DELETE CASCADE,
    passed        BOOLEAN NOT NULL,
    actual_stdout TEXT NOT NULL DEFAULT '',
    stderr        TEXT NOT NULL DEFAULT '',
    timed_out     BOOLEAN NOT NULL DEFAULT FALSE,
    duration_ms   INTEGER NOT NULL DEFAULT 0,
    ordinal       INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_sub_results_submission ON submission_test_results(submission_id);
