package folder

import (
	"fmt"
	"context"
	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	Create(ctx context.Context, f *FolderModel)(err error)

	
}

type storage struct {
	client *pgxpool.Pool
	logger logging.Logger
}

func (s *storage) Create(ctx context.Context, f *FolderModel)(err error) {
	q := `
		INSERT INTO folder (
			name, addres, access
		)
		VALUES (
			$1, $2, $3
		)
		RETURNING id
	`
	if err = s.client.QueryRow(ctx, q, f.Name, f.Addres, f.Access).Scan(&f.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL ERROR Massage: %s Detail: %s Where: %s Code: %s SQLSelect: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			s.logger.Info(newErr)

			return nil
		}
		return nil
	}
	return nil
}

func NewStorage(client *pgxpool.Pool, logger *logging.Logger) Storage {
	return &storage{
		client: client,
		logger: *logger,
	}
}