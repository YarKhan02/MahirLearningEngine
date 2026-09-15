UPDATE programs SET published = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, slug, title, short_desc, full_desc, learn, age, fee, level, bonus, accent, image_url, order_no, published, created_at, updated_at
