#!/bin/bash

DATABASE_URL="postgres://yarkhan:yarkhanworkshop@localhost:5432/mahirlearning?sslmode=disable"

echo "Resetting database..."

psql "$DATABASE_URL" <<'EOF'
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO public;
EOF

echo "Database reset complete."