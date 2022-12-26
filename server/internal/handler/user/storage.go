package user

import(
	"fmt"
	"context"
	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, user *UserModel) error
	FindAll(ctx context.Context) (u []UserModel, err error)
	FindOne(ctx context.Context, id string) error
	Update(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	client *pgxpool.Pool
	logger logging.Logger
}
func (r *repository) Create(ctx context.Context, u *UserModel) (error) {
	q := `
		INSERT INTO users (
			name, email, confirm, confirm_code
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
		RETURNING id
	`
	
	if err := r.client.QueryRow(ctx, q, u.Name, u.Email, ).Scan(&u.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL ERROR Massage: %s Detail: %s Where: %s Code: %s SQLSelect: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Info(newErr)

			return nil
		}
		return nil
	}

	return nil
}

func (r *repository) FindAll(ctx context.Context) (u []UserModel, err error){
	return nil, nil
}
func (r *repository) FindOne(ctx context.Context, email string) (err error) {
	var i string
	q := `
		SELECT id 
		FROM users
		WHERE email=$1 
	`

	if err := r.client.QueryRow(ctx, q, email).Scan(&i); err != nil {
		return err
	}
	return nil
}

func (r *repository) Update(ctx context.Context, id string) (err error) {
	return  nil
}
func (r *repository) Delete(ctx context.Context, id string) (err error) {
	return nil
}

func NewStorage(client *pgxpool.Pool, logger *logging.Logger) Repository {
	return &repository{
		client: client,
		logger: *logger,
	}
}