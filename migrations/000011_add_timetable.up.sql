-- The weekly schedule moves out of batches.days into its own table.
ALTER TABLE batches DROP COLUMN days;

-- A timetable row is a recurring rule: "this batch takes this course on these
-- weekdays, at this time". Actual class dates are generated on the fly between
-- the batch's start_date and end_date, so changing those keeps the schedule in
-- sync without touching this table.
CREATE TABLE timetable (
    id UUID PRIMARY KEY,
    batch_id UUID NOT NULL,
    course_id UUID NOT NULL,
    -- Bitmask of ISO weekdays: bit (weekday-1), Mon=1, Tue=2, Wed=4 … Sun=64.
    weekdays INTEGER NOT NULL CHECK (weekdays > 0),
    start_time TEXT NOT NULL,
    end_time TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_timetable_batch
        FOREIGN KEY (batch_id)
        REFERENCES batches(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_timetable_course
        FOREIGN KEY (course_id)
        REFERENCES course(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_timetable_batch ON timetable(batch_id);
CREATE INDEX idx_timetable_course ON timetable(course_id);
