package store

import (
	"context"
	"encoding/base64"
	"fmt"

	"dnevnik-rg.ru/internal/models/response"
	"dnevnik-rg.ru/pkg/utils"
)

func (s *RgStore) Authorize(key int, checksum string, ip string) (*response.Auth, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	var query = `select * from auth.select_user_private($1, $2, 'default authorization')`

	var (
		auth     *response.Auth
		password string
	)

	err := s.s.QueryRow(ctx, query, key, ip).Scan(
		&password,
		&auth.Token,
	)

	if err != nil {
		return nil, err
	}

	byteArr := []byte(utils.HashSumGen(key, checksum))

	if password != base64.StdEncoding.EncodeToString(byteArr) {
		return nil, fmt.Errorf("wrong password")
	}

	query = `select role from auth.get_user_role($1)`

	err = s.s.QueryRow(ctx, query, key).Scan(
		&auth.Role,
	)

	if err != nil {
		return nil, err
	}

	return auth, nil
}
