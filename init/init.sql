-- init.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- База для Umami аналитики
SELECT 'CREATE DATABASE umami'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'umami')\gexec