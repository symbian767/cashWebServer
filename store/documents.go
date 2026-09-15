package store

import (
	"context"
	"document-service/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

// SaveDocument сохраняет документ в базе данных PostgreSQL.
func (s *Store) SaveDocument(ctx context.Context, doc models.Document) (models.Document, error) {
	const query = `
		INSERT INTO service_documents.files (name, file, public, mime)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING
		RETURNING id;
	`
	if err := s.pool.QueryRow(ctx, query, doc.Name, doc.File, doc.Public, doc.Mime).Scan(&doc.ID); err != nil {
		return models.Document{}, err
	}
	return doc, nil
}

func (s *Store) GetDocumentList(ctx context.Context, limit, offset uint64) (models.DocumentsList, error) {
	query := sqlSelectDocument
	if limit != 0 {
		query = fmt.Sprintf("%s limit %d", query, limit)
		if offset != 0 {
			query = fmt.Sprintf("%s offset %d", query, offset)
		}
	}

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, scanDocument)
}

func (s *Store) GetDocument(ctx context.Context, id uint64) (models.Document, error) {
	rows, err := s.pool.Query(ctx, sqlSelectDocument+` HAVING files.id = $1`, id)
	if err != nil {
		return models.Document{}, err
	}

	doc, err := pgx.CollectExactlyOneRow(rows, scanDocument)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Document{}, fmt.Errorf("документ не найден: %w", err)
		}
		return models.Document{}, err
	}
	return doc, nil
}

// Удаление документа
func (s *Store) DeleteDocument(ctx context.Context, id uint64) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM service_documents.files WHERE files.id = $1`, id)
	return err
}

// scanDocument — общий коллектор строки для GetDocument и GetDocumentList.
func scanDocument(row pgx.CollectableRow) (models.Document, error) {
	var doc models.Document
	var jsonData []byte
	if err := row.Scan(&doc.ID, &doc.Name, &doc.File, &doc.Public, &doc.Mime, &doc.Created, &jsonData); err != nil {
		return doc, err
	}
	if err := json.Unmarshal(jsonData, &doc.UsersFileLinkList); err != nil {
		log.Println("Ошибка при парсинге JSON:", err)
	}
	return doc, nil
}
