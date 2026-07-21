-- Quizzes attached to a lesson. A quiz has multiple questions; each question is
-- either multiple-choice (auto-graded) or typed (teacher-graded).
CREATE TABLE quizzes (
    id          UUID PRIMARY KEY,
    lesson_id   UUID NOT NULL,
    title       VARCHAR(200) NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_quizzes_lesson
        FOREIGN KEY (lesson_id) REFERENCES lesson(id) ON DELETE CASCADE
);
CREATE INDEX idx_quizzes_lesson ON quizzes(lesson_id);

CREATE TABLE quiz_questions (
    id             UUID PRIMARY KEY,
    quiz_id        UUID NOT NULL,
    prompt         TEXT NOT NULL,
    type           TEXT NOT NULL CHECK (type IN ('mcq', 'typed')),
    marks          INTEGER NOT NULL DEFAULT 0,
    -- mcq only: whether more than one option may be correct/selected.
    allow_multiple BOOLEAN NOT NULL DEFAULT FALSE,
    order_no       INTEGER NOT NULL,

    CONSTRAINT fk_quiz_questions_quiz
        FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE
);
CREATE INDEX idx_quiz_questions_quiz ON quiz_questions(quiz_id);

CREATE TABLE quiz_options (
    id          UUID PRIMARY KEY,
    question_id UUID NOT NULL,
    text        TEXT NOT NULL,
    is_correct  BOOLEAN NOT NULL DEFAULT FALSE,
    order_no    INTEGER NOT NULL,

    CONSTRAINT fk_quiz_options_question
        FOREIGN KEY (question_id) REFERENCES quiz_questions(id) ON DELETE CASCADE
);
CREATE INDEX idx_quiz_options_question ON quiz_options(question_id);

-- One submission per student per quiz. status = 'submitted' while any typed
-- answer still needs the teacher; 'graded' once finalized (MCQ-only quizzes are
-- graded immediately on submit).
CREATE TABLE quiz_submissions (
    id           UUID PRIMARY KEY,
    quiz_id      UUID NOT NULL,
    student_id   UUID NOT NULL,
    status       TEXT NOT NULL DEFAULT 'submitted',
    remarks      TEXT,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    graded_at    TIMESTAMPTZ,

    CONSTRAINT fk_quiz_submissions_quiz
        FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE,
    CONSTRAINT fk_quiz_submissions_student
        FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    CONSTRAINT uq_quiz_submission UNIQUE (quiz_id, student_id)
);
CREATE INDEX idx_quiz_submissions_quiz ON quiz_submissions(quiz_id);

CREATE TABLE quiz_answers (
    id            UUID PRIMARY KEY,
    submission_id UUID NOT NULL,
    question_id   UUID NOT NULL,
    answer_text   TEXT,          -- typed answers
    awarded_marks INTEGER,       -- set on submit for mcq; on grading for typed

    CONSTRAINT fk_quiz_answers_submission
        FOREIGN KEY (submission_id) REFERENCES quiz_submissions(id) ON DELETE CASCADE,
    CONSTRAINT fk_quiz_answers_question
        FOREIGN KEY (question_id) REFERENCES quiz_questions(id) ON DELETE CASCADE,
    CONSTRAINT uq_quiz_answer UNIQUE (submission_id, question_id)
);
CREATE INDEX idx_quiz_answers_submission ON quiz_answers(submission_id);

CREATE TABLE quiz_answer_options (
    answer_id UUID NOT NULL,
    option_id UUID NOT NULL,

    PRIMARY KEY (answer_id, option_id),
    CONSTRAINT fk_qao_answer
        FOREIGN KEY (answer_id) REFERENCES quiz_answers(id) ON DELETE CASCADE,
    CONSTRAINT fk_qao_option
        FOREIGN KEY (option_id) REFERENCES quiz_options(id) ON DELETE CASCADE
);
