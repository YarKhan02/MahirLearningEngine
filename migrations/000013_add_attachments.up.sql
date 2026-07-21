CREATE TABLE attachments (
    id            UUID PRIMARY KEY,
    r2_key        TEXT NOT NULL UNIQUE,
    file_name     TEXT NOT NULL,
    content_type  TEXT NOT NULL,
    size_bytes    BIGINT,
    resource_type TEXT NOT NULL,          -- polymorphic owner, e.g. 'course'
    resource_id   UUID NOT NULL,          -- e.g. the course id
    uploaded_by   UUID,                   -- admin user; kept if the user is removed
    status        TEXT NOT NULL DEFAULT 'pending',  -- 'pending' | 'confirmed'
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    confirmed_at  TIMESTAMPTZ,
    deleted_at    TIMESTAMPTZ,

    CONSTRAINT fk_attachments_uploaded_by
        FOREIGN KEY (uploaded_by)
        REFERENCES users(id)
        ON DELETE SET NULL
);

-- Fast lookup of a resource's live materials.
CREATE INDEX idx_attachments_resource
    ON attachments (resource_type, resource_id)
    WHERE deleted_at IS NULL;
