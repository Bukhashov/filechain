package user

import (
	"fmt"
	"context"
	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/Bukhashov/filechain/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, user *model.User) error
	EmailControl(ctx context.Context, d *Dto) error
	FindByEmail(ctx context.Context, u *user) error
	UpdateIamge(ctx context.Context, u *model.User) error
}

type repository struct {
	client *pgxpool.Pool
	logger logging.Logger
}
func (r *repository) Create(ctx context.Context, u *model.User) (err error) {
	q := `
		INSERT INTO users (
			name, email
		)
		VALUES (
			$1, $2
		)
		RETURNING id
	`
	
	if err = r.client.QueryRow(ctx, q, u.Name, u.Email).Scan(&u.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL ERROR Massage: %s Detail: %s Where: %s Code: %s SQLSelect: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Info(newErr)

			return err
		}
		return nil
	}

	return nil
}
func (r *repository) EmailControl(ctx context.Context, d *Dto) (err error){
	q := `
		SELECT id, email
		FROM users
		WHERE email=$1
	`
	if err := r.client.QueryRow(ctx, q, d.Email).Scan(&d.ID, &d.Email); err != nil {
		return err
	}
	return nil
}

func (r *repository) FindByEmail(ctx context.Context, u *user) (err error) {
	q := `
		SELECT id, name, image
		FROM users
		WHERE email=$1
	`
	err = r.client.QueryRow(ctx, q, u.Dto.Email).Scan(&u.Model.ID, &u.Model.Name, &u.Model.Image); if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateIamge(ctx context.Context, u *model.User) (err error) {
	q := `
		UPDATE users
		SET image=$1
		WHERE id=$2
	`
	_, err = r.client.Exec(ctx, q, u.Image, u.ID); if err != nil {
		return err
	}
	return  nil
}

func NewStorage(client *pgxpool.Pool, logger *logging.Logger) Repository {
	return &repository{
		client: client,
		logger: *logger,
	}
}