# document-service

REST-сервис на Go для хранения документов с разграничением доступа по пользователям: регистрация, аутентификация по логину/паролю, загрузка/получение/удаление файлов, выдача прав на документ конкретным пользователям.

## Стек

- Go 1.26.5
- [gorilla/mux](https://github.com/gorilla/mux) — роутинг
- [pgx/v5](https://github.com/jackc/pgx) (`pgxpool`) — драйвер PostgreSQL
- [godotenv](https://github.com/joho/godotenv) + [envconfig](https://github.com/kelseyhightower/envconfig) — конфигурация из `.env`
- PostgreSQL (расширение `pgcrypto` для хэширования паролей)

## Запуск

# поднять окружение
```
`make run-service`
```

```
Сервис слушает `http://127.0.0.1:8080`.

## Аутентификация

1. Регистрация пользователя требует токен администратора (`DOCUMENT_ADMIN_TOKEN`) в теле запроса — доступна только тому, у кого он есть.
2. `/api/auth` по логину/паролю выдаёт сессионный токен.
3. Все остальные ручки (`/api/docs*`) требуют этот токен в заголовке `Authorization`.

**Требования к логину:** минимум 8 символов, только латиница и цифры.

**Требования к паролю:** минимум 8 символов, минимум одна заглавная и одна строчная буква, минимум одна цифра, минимум один спецсимвол.

## API

### `POST /api/register`
Регистрация нового пользователя (только с админ-токеном).

```json
{
  "token": "token-admin",
  "login": "testuser1",
  "password": "Passw0rd!"
}
```

### `POST /api/auth`
Аутентификация, возвращает сессионный токен.

```json
{
  "login": "testuser1",
  "password": "Passw0rd!"
}
```

Ответ:

```json
{ "data": { "token": "..." } }
```

### `DELETE /api/auth/{token}`
Завершение сессии (logout).

### `POST /api/docs`
Загрузка документа. `multipart/form-data`, заголовок `Authorization: <session-token>`.

Поля формы: `file` (сам файл), `name`, `mime`.

### `GET /api/docs?limit=&offset=`
Список документов. Заголовок `Authorization` обязателен.

### `GET /api/docs/{id}`
Получение одного документа по id.

### `DELETE /api/docs/{id}`
Удаление документа по id.

## Формат ответа

Все ответы — JSON вида:

```json
{ "data": { ... } }
```

или при ошибке:

```json
{ "error": { "code": 500, "text": "..." } }
```

## Схема БД

Таблицы находятся в схеме `service_documents`: `users`, `files`, `users_file_link` (связь многие-ко-многим между пользователями и документами, на которые им выдан доступ). Подробности — в `migration.sql`.
