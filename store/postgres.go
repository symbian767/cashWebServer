package store

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

const sqlSelectDocument = `SELECT documents.id, name, file, public, mime, documents.created, json_agg(json_build_object('login', users.login)) FROM service_documents.documents
    INNER JOIN service_documents.users_file_link ufl on documents.id = ufl.id_document
    INNER JOIN service_documents.users ON users.id = ufl.id_user GROUP BY documents.id`

// Store — репозиторий поверх пула соединений pgx. Пул создаётся один раз
// в NewStore и переиспользуется всеми методами.
type Store struct {
	pool *pgxpool.Pool
}

// NewStore открывает пул соединений с PostgreSQL и проверяет его пингом.
func NewStore(ctx context.Context, dsn string) (*Store, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("создание пула соединений: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("проверка соединения с базой данных: %w", err)
	}
	log.Println("Connected to PostgreSQL database")
	return &Store{pool: pool}, nil
}

// Close закрывает пул соединений. Вызывается один раз при остановке сервиса.
func (s *Store) Close() {
	s.pool.Close()
}
