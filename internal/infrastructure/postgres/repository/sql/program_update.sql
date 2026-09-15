UPDATE programs SET
    title = $2, short_desc = $3, full_desc = $4, learn = $5, age = $6,
    fee = $7, level = $8, bonus = $9, accent = $10, image_url = $11,
    published = $12, updated_at = NOW()
WHERE id = $1
RETURNING id, slug, title, short_desc, full_desc, learn, age, fee, level, bonus, accent, image_url, order_no, published, created_at, updated_at
