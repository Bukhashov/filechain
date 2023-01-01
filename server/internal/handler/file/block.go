package block

import(
	"net/http"
	"github.com/Bukhashov/filechain/internal/config"
	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Filechain interface {
	New(w http.ResponseWriter, req *http.Request)
	Get(w http.ResponseWriter, req *http.Request)
	All(w http.ResponseWriter, req *http.Request)
	Update(w http.ResponseWriter, req *http.Request)
}

type filechain struct {
	client 	*pgxpool.Pool
	logger	*logging.Logger
	config	config.Config
	Blocks	[]*Block
}

func NewBlock(client *pgxpool.Pool, config config.Config, logger *logging.Logger) Filechain{
	return &filechain {
		config: config,
		logger: logger,
		client: client,
	}
}





