package history

import (
	"context"
	"fmt"

	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	Write(ctx context.Context, h *HistoryModel) (err error)
	ReadOne(ctx context.Context, h *HistoryModel) (err error)
	ReadAll(ctx context.Context, h *HistoryDto)([]HistoryModel, error)
}

type storage struct {
	client 		*pgxpool.Pool
	logger 		*logging.Logger
	// Table name
	// Деректер қорындағы таблицанын атауы
	table	*string
	// folder_history
	// file_history
}


func (s *storage) Write(ctx context.Context, h *HistoryModel)(err error){
	q := `
		INSERT INTO %s (
			hash,
			addres,
			timeStamp,
			proveHash,
			user,
			target
		)
		VALUES (
			$1, $2, $3, $4, &5, $6
		)
		RETURNING id
	`

	q = fmt.Sprintf(q, s.table)

	if err = s.client.QueryRow(ctx, q, h.Hash, h.Addres, h.TimeStamp, h.ProveHash, h.User, h.Target).Scan(&h.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL ERROR Massage: %s Detail: %s Where: %s Code: %s SQLSelect: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			s.logger.Info(newErr)

			return nil
		}
		return nil
	}

	return nil
}
func (s *storage) ReadOne(ctx context.Context, h *HistoryModel)(err error) {
	q := `
		SELECT 
			id,
			hash,
			timeStamp,
			proveHash,
			user,
			target
		FROM
			%s
		WHERE
			addres=$1
	`
	q = fmt.Sprintf(q, s.table)
	
	if err = s.client.QueryRow(ctx, q, h.Addres).Scan(&h.ID, h.Hash, h.TimeStamp, h.ProveHash, h.User, h.Target); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL ERROR Massage: %s Detail: %s Where: %s Code: %s SQLSelect: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			s.logger.Info(newErr)

			return nil
		}
		return nil
	}

	return nil
}
func (s *storage) ReadAll(ctx context.Context, h *HistoryDto)(historyes []HistoryModel, err error){
	q := `
		SELECT 
			id,
			hash,
			timeStamp,
			proveHash,
			user,
			target
		FROM
			%s
		WHERE
			addres=$1
	`
	q = fmt.Sprintf(q, s.table)

	rows, err := s.client.Query(ctx, q, h.Addres); if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	historyes = []HistoryModel{}
	
	for rows.Next(){
		var h HistoryModel
		if err = rows.Scan(&h.ID, &h.Hash, &h.TimeStamp, &h.ProveHash, &h.User, &h.Target); err != nil {
			return nil, err
		}
		historyes = append(historyes, h)
	}

	return historyes, nil
}


func NewStorage(client *pgxpool.Pool, logger *logging.Logger, table *string) Storage{
	return &storage{
		client: client,
		logger: logger,
		table: table,
	}
}