BEGIN;

CREATE DATABASE "cash";

CREATE SCHEMA "public";
/* migration trouble todo */


-- Пользователи. login проверяется приложением (латиница+цифры, >=8 симв.)
-- на уровне БД дополнительно защищаем уникальность логина (Unique).
CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY,
    login      TEXT        NOT NULL UNIQUE,
    password   TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

-- Таблица документов.
CREATE TABLE IF NOT EXISTS documents
(
    id         SERIAL PRIMARY KEY,
    owner_id   INTEGER     NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name       TEXT        NOT NULL,
    mime       TEXT        NOT NULL DEFAULT '',
    is_file    BOOLEAN     NOT NULL DEFAULT FALSE,
    is_public  BOOLEAN     NOT NULL DEFAULT FALSE,
    json_data  JSONB,
    file_path  TEXT        NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

-- Таблица пользователей, которым открыт доступ к документу.
CREATE TABLE IF NOT EXISTS document_grants
(
    document_id INTEGER NOT NULL REFERENCES documents (id) ON DELETE CASCADE,
    user_id     INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    PRIMARY KEY (document_id, user_id)
    );



COMMIT;