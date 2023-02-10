package storage

import (
	"context"

	"ghostorange/internal/app/storage/psqldb"
)

type (
	Storage interface {
		GetPasswordHash(cxt context.Context, userName string) (string, string, error)
		AddUser(ctx context.Context, username, hash string) (string, error)
		UserExists(ctx context.Context, username string) (bool, error)

		Count(ctx context.Context, dataType int, user string) (int, error)
		GetData(ctx context.Context, dataType int) (any, error)
		AddData(ctx context.Context, dataType int, userID string, data any) error
		UpdateData(ctx context.Context, dataType int, userID string, data any) error
	}
	config interface {
		DBDSN() string
	}
)

func New(cfg config) (Storage, error) {
	return psqldb.New(cfg.DBDSN())
}
