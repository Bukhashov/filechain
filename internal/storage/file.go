package storage

import (
	"context"
	"fmt"
	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/Bukhashov/filechain/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FileStorage interface {
	New(ctx context.Context, f model.File) (err error)
}

type fileStorage struct {
	client *pgxpool.Pool
	logger logging.Logger
}

func (s *fileStorage) New(ctx context.Context, f model.File) (err error) {
	q := `
		INSERT INTO file (
			timeStamp, 
			hash, 
			prevHash,
			access,
			type,
			title,
			file
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		RETURNING id
	
	`
	s.client.QueryRow(ctx, q, f.TimeStamp, f.Hash, f.PrevHash, f.Access, f.Type, f.Title, f.File).Scan(f.Id); if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL ERROR Massage: %s Detail: %s Where: %s Code: %s SQLSelect: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			s.logger.Info(newErr)
			return nil
		}
	}

	return nil
}

func NewFileStorage(client *pgxpool.Pool, logger *logging.Logger) FileStorage {
	return &fileStorage{
		client: client,
		logger: *logger,
	}
}