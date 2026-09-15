CREATE SCHEMA IF NOT EXISTS service_documents;

SET search_path TO service_documents;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Создадим таблицу пользователей
CREATE TABLE IF NOT EXISTS service_documents.users
(
    "id"          SERIAL    PRIMARY KEY,
    "login"       TEXT      NOT NULL,
    "password"    TEXT      NOT NULL,
    "created"     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Создадим таблицу документов
CREATE TABLE IF NOT EXISTS service_documents.documents
(
    "id"            TEXT      PRIMARY KEY,                -- Уникальный идентификатор файла(UUID)
    "head_user_id"  INTEGER   NOT NULL REFERENCES users (id) ON DELETE CASCADE, -- Идентификатор владельца документа
    "name"          TEXT      NOT NULL,                   -- Имя файла
    "file"          BOOLEAN   NOT NULL,                   -- Поле, указывающее, является ли запись файлом
    "public"        BOOLEAN   NOT NULL,                   -- Поле, указывающее, доступен ли файл публично
    "mime"          TEXT      NOT NULL,                   -- MIME-тип файла
    "created"       TIMESTAMP DEFAULT CURRENT_TIMESTAMP   -- Время создания записи
);

CREATE TABLE IF NOT EXISTS service_documents.users_file_link
(
    "id_user"       INTEGER REFERENCES service_documents.users on delete cascade,
    "uuid_document" INTEGER REFERENCES service_documents.documents on delete cascade ,
    primary key (id_user, uuid_document)
)