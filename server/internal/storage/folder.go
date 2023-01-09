package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/Bukhashov/filechain/internal/handler/user"
	"github.com/Bukhashov/filechain/internal/model"
	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FolderStorage interface {
	Create(ctx context.Context, f *model.Folder)(err error)
	Update(ctx context.Context, f *model.Folder)(err error)
	GetAllAddress(ctx context.Context, u *user.Dto) ([]*model.Folder, error)
	GetFolder(ctx context.Context, uid *int64, address []byte)(folder *model.Folder, err error)
}

type folderStorage struct {
	client *pgxpool.Pool
	logger logging.Logger
}

func (s *folderStorage) Create(ctx context.Context, f *model.Folder)(err error) {
	q := `
		INSERT INTO folder (
			address, name, userId, file, access
		)
		VALUES (
			$1, $2, $3, $4, $5
		)
		RETURNING id
	`
	
	if err = s.client.QueryRow(ctx, q, f.Addres, f.Name, f.UserId, f.File, f.Access).Scan(&f.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL ERROR Massage: %s Detail: %s Where: %s Code: %s SQLSelect: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			s.logger.Info(newErr)
			return nil
		}
		return nil
	}
	return nil
}
func (s *folderStorage) Update(ctx context.Context, f *model.Folder)(err error) {
	q := `
		UPDATE folder
		SET file=$1
		WHERE address=$2 AND userid=$3
	`
	_, err = s.client.Exec(ctx, q, f.File, f.Addres, f.UserId); if err != nil {
		s.logger.Info(err)
		return err
	}
	return nil
}
func (s *folderStorage) GetAllAddress(ctx context.Context, u *user.Dto) (folderModel []*model.Folder, err error) {
	q := `
		SELECT address
		FROM folder
		WHERE userId=$1
	`
	rows, err := s.client.Query(ctx, q, u.ID); if err != nil {
		return nil, err
	}

	for rows.Next() {
		var f model.Folder
		if err = rows.Scan(&f.Addres); err != nil {
			return nil, err
		}
		folderModel = append(folderModel, &f)
	}
	return folderModel, nil
}
func (s *folderStorage) GetFolder(ctx context.Context, uid *int64, address []byte)(folder *model.Folder, err error) { 
	q := `
		SELECT id, name, file, access, address
		FROM folder
		WHERE userid=$1
	`
	rows, err := s.client.Query(ctx, q, uid); if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL ERROR Massage: %s Detail: %s Where: %s Code: %s SQLSelect: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			s.logger.Info(newErr)

			return nil, err
		}
		fmt.Printf("\nerr get folder %v", err)
	}

	var f model.Folder
	for rows.Next() {
		if err = rows.Scan(&f.ID, &f.Name, &f.File, &f.Access, &f.Addres); err != nil {
			return nil, err
		}
		if bytes.Equal(f.Addres, address) {
			return &f, nil
		}
	}

	err = errors.New("not fund")
	return nil, err
	
}

func NewFolderStorage(client *pgxpool.Pool, logger *logging.Logger) FolderStorage {
	return &folderStorage{
		client: client,
		logger: *logger,
	}
}