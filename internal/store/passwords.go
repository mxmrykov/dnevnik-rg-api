package store

import (
	"context"
	"time"

	"dnevnik-rg.ru/internal/models"
)

func (s *RgStore) NewPassword(password models.Password) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select from passwords.new_password_row($1, $2, $3, $4)`

	_, err := s.s.Exec(ctx, query, password.Key, password.CheckSum, password.Token, password.LastUpdate)

	return err
}

func (s *RgStore) DeletePassword(key int) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select * from passwords.delete_password_row($1)`

	_, err := s.s.Exec(ctx, query, key)

	return err
}

func (s *RgStore) AutoUpdateToken(token string, key int) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	const query = `select * from auth.auto_upd_token($1, $2, $3)`

	_, err := s.s.Exec(ctx, query, key, token, time.Now().Format(time.RFC3339))

	return err
}
