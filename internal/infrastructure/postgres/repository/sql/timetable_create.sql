INSERT INTO timetable (id, batch_id, course_id, weekdays, start_time, end_time, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW());
