-- Batch-scoped announcements: only students enrolled in the batch see them.
CREATE TABLE announcements (
    id UUID PRIMARY KEY,
    batch_id UUID NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_announcements_batch
        FOREIGN KEY (batch_id)
        REFERENCES batches(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_announcements_batch ON announcements(batch_id);
