package store

import (
	"context"
	"document-service/models"
)

func (s *Store) CreateUser(ctx context.Context, user models.User) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO service_documents.users(login, password) VALUES ($1, crypt($2, gen_salt('bf')))`,
		user.Login, user.Password)
	return err
}

func (s *Store) AuthenticateUser(ctx context.Context, login, password string) (bool, error) {
	var userExists bool
	err := s.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT * FROM service_documents.users WHERE login = $1 AND password = crypt($2, password))`,
		login, password).Scan(&userExists)
	return userExists, err
}

func (s *Store) GetUserIDByLogin(ctx context.Context, login string) (uint64, error) {
	var id uint64
	err := s.pool.QueryRow(ctx, `SELECT id FROM service_documents.users WHERE login = $1`, login).Scan(&id)
	return id, err
}

func (s *Store) SaveUserFileLink(ctx context.Context, userFileLink models.UsersFileLink) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO service_documents.users_file_link(id_user, id_document) VALUES ($1, $2)`,
		userFileLink.IdUser, userFileLink.IdDocument)
	return err
}
