SELECT id, slug, title, short_desc, full_desc, learn, age, fee, level, bonus, accent, image_url, order_no, published, created_at, updated_at FROM programs WHERE slug = $1 AND published = TRUE
