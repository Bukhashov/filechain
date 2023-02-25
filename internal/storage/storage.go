package storage

import (
	"context"
	"fmt"

	"github.com/Bukhashov/filechain/internal/model"
	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/Bukhashov/filechain/pkg/validator"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User)(err error)
	// FindUserById(ctx context.Context, u *model.User) error
	FindUserByEmail(ctx context.Context, u *model.User)(err error)
	UpdateIamge(ctx context.Context, u *model.User)(err error)
}

type userRepository struct {
	client *pgxpool.Pool
	logger logging.Logger
}

func (r *userRepository) Create(ctx context.Context, u *model.User) (err error) {
	if err = validator.UserVaidator(u); err != nil {
		return err;
	}
	
	q := `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
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
func (r *userRepository) FindUserByEmail(ctx context.Context, u *model.User)(err error){
	q := `
		SELECT id, email, image
		FROM users
		WHERE email=$1
	`
	if err := r.client.QueryRow(ctx, q, u.Email).Scan(&u.ID, &u.Email, &u.Image); err != nil {
		return err
	}
	return nil
}
func (r *userRepository) UpdateIamge(ctx context.Context, u *model.User)(err error){
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

func NewUserStorage(client *pgxpool.Pool, logger *logging.Logger) UserRepository {
	return &userRepository{
		client: client,
		logger: *logger,
	}
}