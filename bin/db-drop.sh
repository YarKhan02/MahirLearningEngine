#!/bin/bash

CYAN='\033[1;36m'
NO_COLOR='\033[0m'
LABEL="db-drop"
printf "${CYAN}== ${LABEL}${NO_COLOR}\n"

CONNECTION_URL="postgresql://yarkhan:yarkhanworkshop@localhost:5432/postgres"

# Terminate active connections first (important!)
psql $CONNECTION_URL -c "
    SELECT pg_terminate_backend(pg_stat_activity.pid)
    FROM pg_stat_activity
    WHERE datname = 'mahirlearning' AND pid <> pg_backend_pid();
"

# Drop the DB only if it exists
psql $CONNECTION_URL -tc "SELECT 1 FROM pg_database WHERE datname = 'mahirlearning'" | grep -q 1 &&
  psql $CONNECTION_URL -c "DROP DATABASE mahirlearning" ||
  echo "mahirlearning DB does not exist — skipping drop."