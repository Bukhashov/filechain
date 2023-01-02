package history

import (
	"context"

	"github.com/Bukhashov/filechain/internal/config"
	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)
type History interface{
	ReadAll(ctx context.Context)(err error)
	ReadOne(ctx context.Context, model *HistoryModel)(err error)
	Write(ctx context.Context)(err error)
	Genesis(ctx context.Context)(err error)
}
type history struct{
	client	*pgxpool.Pool
	logger	*logging.Logger
	config	*config.Config
	table	*string
	// [FILE] саздать етілген уақыт
	t 		int64

	address	*[]byte
	
	user	*string
	target	*string
	
}



func NewHistory(client *pgxpool.Pool, config *config.Config, logger *logging.Logger, table, user, target *string, t int64, address *[]byte) History {
	return &history{
		client: client,
		config: config,
		logger: logger,
		table: table,
		t: t,
		user: user,
		target: target,
		address: address,
	}
}