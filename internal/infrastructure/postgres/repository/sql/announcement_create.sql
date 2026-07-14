INSERT INTO announcements (id, batch_id, title, description, created_at)
VALUES ($1, $2, $3, $4, NOW());
