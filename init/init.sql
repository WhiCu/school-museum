-- init.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Visitors table (auto-created by bun ORM, DDL here for reference)
CREATE TABLE IF NOT EXISTS visitors (
    id            BIGSERIAL PRIMARY KEY,
    ip            TEXT NOT NULL UNIQUE,
    user_agent    TEXT DEFAULT '',
    page          TEXT DEFAULT '',
    referrer      TEXT DEFAULT '',
    screen_width  INT DEFAULT 0,
    screen_height INT DEFAULT 0,
    language      TEXT DEFAULT '',
    visit_count   INT NOT NULL DEFAULT 1,
    first_visit_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_visit_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS visitors_ip_idx ON visitors (ip);