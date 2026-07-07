#!/bin/bash
set -euo pipefail

CONNECTION_URL="postgresql://yarkhan:yarkhanworkshop@localhost:5432/mahirlearning"

if [ -z "$CONNECTION_URL" ]; then
  echo "Usage: $0 <database_url> or set DATABASE_URL" >&2
  exit 1
fi

psql "$CONNECTION_URL" -v ON_ERROR_STOP=1 <<'SQL'
INSERT INTO role (name) VALUES
  ('admin'),
  ('parent'),
  ('student')
ON CONFLICT DO NOTHING;
SQL