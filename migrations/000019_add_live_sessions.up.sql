-- Ad-hoc live classes started by a teacher (admin) for a batch taking a course.
CREATE TABLE live_sessions (
    id          UUID PRIMARY KEY,
    batch_id    UUID NOT NULL,
    course_id   UUID NOT NULL,
    host_id     UUID NOT NULL,
    status      TEXT NOT NULL DEFAULT 'live', -- 'live' | 'ended'
    started_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at    TIMESTAMPTZ,

    CONSTRAINT fk_live_sessions_batch  FOREIGN KEY (batch_id)  REFERENCES batches(id) ON DELETE CASCADE,
    CONSTRAINT fk_live_sessions_course FOREIGN KEY (course_id) REFERENCES course(id)  ON DELETE CASCADE,
    CONSTRAINT fk_live_sessions_host   FOREIGN KEY (host_id)   REFERENCES users(id)   ON DELETE CASCADE
);

-- At most one live session per batch at a time (partial unique index).
CREATE UNIQUE INDEX uq_live_sessions_one_live_per_batch
    ON live_sessions(batch_id) WHERE status = 'live';

CREATE INDEX idx_live_sessions_batch_status ON live_sessions(batch_id, status);

-- Saved whiteboard snapshots (scene JSON in R2 + optional PNG preview).
CREATE TABLE whiteboard_snapshots (
    id           UUID PRIMARY KEY,
    session_id   UUID NOT NULL,
    storage_key  TEXT NOT NULL,             -- R2 key of the .excalidraw scene JSON
    image_url    TEXT NOT NULL DEFAULT '',  -- stable URL of the PNG preview
    version      INT NOT NULL,
    created_by   UUID NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_whiteboard_snapshots_session
        FOREIGN KEY (session_id) REFERENCES live_sessions(id) ON DELETE CASCADE
);

CREATE INDEX idx_whiteboard_snapshots_session ON whiteboard_snapshots(session_id, version);
