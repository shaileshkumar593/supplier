package test

/*import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"

	"swallow-supplier/config"
	service "swallow-supplier/iface"
	"swallow-supplier/implementation"
	"swallow-supplier/metrics"
	"swallow-supplier/repository"
	"swallow-supplier/transport"
	httptransport "swallow-supplier/transport/http"
	"swallow-supplier/utils/validator"
)

// GetHandlers retrieves the http handlers
func GetHandlers() http.Handler {
	c := config.Instance()
	logger := log.NewNopLogger()
	var dbconn = make(map[string]*sqlx.DB)
	{
		var err error

		// Connect to the database
		if dbconn[c.DefaultDBName], err = sqlx.Connect(config.Instance().DatabaseDriverPostgres, config.Instance().DatabaseConnection); err != nil {
			panic(err)
		}

	}

	// Initialize the metrics factory
	metrics.Init()

	// Initialize the validator
	validator.Init()

	var repo map[string]service.Repository
	{
		var err error

		repo[c.DefaultDBName], err = repository.New(dbconn[c.DefaultDBName], logger)
		if err != nil {
			panic(err)
		}

	}

	// Creating the service repository
	var mongorepo = make(map[string]ser.MongoRepository)
	{
		var err error

		// Assuming you have a NewMongoRepository function to create the repository
		repo[c.DefaultDBName], err = mongorepo.NewMongoRepository(mongoDB[c.DefaultDBName], logger)
		if err != nil {
			panic(err)
		}
	}
	// Create the service
	var svc service.Service
	{
		// initializing the implementation
		svc = implementation.NewService(repo, mongorepo, logger)
	}

	// Create Go kit endpoints for your service
	// Then decorates with endpoint middlewares
	var endpoints transport.Endpoints
	{
		endpoints = transport.MakeEndpoints(svc)
	}

	// Handle the HTTP endpoints required on your service
	var handler http.Handler
	{
		var ctx context.Context
		handler = httptransport.NewTransport(ctx, endpoints, repo)
	}

	return handler
}

// ConvertToRequestPayload converts interface payload to string
func ConvertToRequestPayload(m interface{}) string {
	data, _ := json.Marshal(m)
	return string(data)
}
*/
