package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"swallow-supplier/scheduler"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	_ "github.com/jackc/pgx"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"swallow-supplier/caches"
	"swallow-supplier/config"
	svc "swallow-supplier/iface"
	repomongo "swallow-supplier/mongo/repository"

	"swallow-supplier/implementation"
	"swallow-supplier/middleware"
	"swallow-supplier/transport"
	httptransport "swallow-supplier/transport/http"
	"swallow-supplier/utils/validator"
)

func main() {
	fmt.Println("Starting main service --- all well ")
	c := config.Instance()
	var (
		httpAddr = flag.String("http.addr", ":"+c.AppPort, "HTTP listen address")
		//job      = flag.String("job", "", "Cron Job to be executed")
	)
	flag.Parse()

	// Use JSON default logger of go-kit
	var logger log.Logger
	{
		logger = log.NewJSONLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"service", c.AppServiceName,
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	// Load configs from ENV into config
	//config.Init()

	// Initialize the validator
	validator.Init()

	// Initialise Redis
	caches.Init()

	level.Info(logger).Log("redis ", "redis initialized successfully")

	//fmt.Println("Connection String:", c.DatabaseConnectionMongo)

	level.Info(logger).Log("info", "MongoDB Connection  started")
	// Create MongoDB connection
	var mongoconn = make(map[string]*mongo.Client)
	{
		var err error
		mongoconn[c.MongoDBName], err = mongo.Connect(context.TODO(), options.Client().
			ApplyURI(c.DatabaseConnectionMongo).
			SetMaxPoolSize(uint64(c.DatabaseDefaultMaxOpenConnections)).
			SetMinPoolSize(uint64(c.DatabaseDefaultMaxIdleConnections)).
			SetMaxConnIdleTime(time.Duration(c.DatabaseMaxIdleTime)*time.Second).
			SetConnectTimeout(10*time.Second)) // Explicit timeout
		if err != nil {
			level.Error(logger).Log("MongoDB connection error: %v", err)
			os.Exit(-1)
			return
		}
		//level.Info(logger).Log("info", "******** MongoDB connection established ********")
	}

	// Ensure MongoDB connection is valid
	err := mongoconn[c.MongoDBName].Ping(context.TODO(), nil)
	if err != nil {
		level.Error(logger).Log("MongoDB ping failed: %v", err)
	}

	level.Info(logger).Log("info", "******** MongoDB connection established ********")

	// Creating the service repository
	var mongorepo = make(map[string]svc.MongoRepository)
	{
		mongorepo[c.MongoDBName], err = repomongo.NewMongo(mongoconn[c.MongoDBName], c.MongoDBName, logger)
		if err != nil {
			level.Error(logger).Log("Mongo Repository creation error: %v", err)
			os.Exit(-1) // optinall  if you donot want code to exit
		}

	}

	/*
				 ctx := context.Background()

		    // Initialize MongoDB repository
		    repo := NewMongoRepository(db, logger)

		    // Ensure indexes only once
		    if err := repo.EnsureIndexes(ctx); err != nil {
		        log.Fatalf("failed to create indexes: %v", err)
			}*/

	// Create the service
	var service svc.Service
	{
		// initializing the implementation
		service = implementation.NewService(mongorepo, logger)

		// attach the middlewares here
		service = middleware.NewLoggingMiddleware(logger)(service)
	}

	// Create Go kit endpoints for your service
	// Then decorates with endpoint middlewares
	var endpoints transport.Endpoints
	{
		endpoints = transport.MakeEndpoints(service)
	}

	// Handle the HTTP endpoints required on your service
	var ctx context.Context
	var handler http.Handler
	{
		handler = httptransport.NewTransport(ctx, endpoints, mongorepo)
	}

	level.Info(logger).Log("msg", "Service Started")
	defer level.Info(logger).Log("msg", "Service Ended")

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		level.Info(logger).Log("transport", "HTTP", "addr", *httpAddr)
		server := &http.Server{
			Addr:    *httpAddr,
			Handler: handler,
		}
		errs <- server.ListenAndServe()
		level.Info(logger).Log("transport", "HTTP", "server_info", fmt.Sprintf("%+v", server))
	}()

	// Start background jobs in a goroutine
	if c.AppEnv != "LOCAL" {
		go func() {
			level.Info(logger).Log("Info", "scheduler started ")
			err := scheduler.Jobs(ctx, logger, service, mongorepo[c.MongoDBName])
			if err != nil {
				level.Error(logger).Log("scheduler", "Jobs", "error", err)
				errs <- err // Report the error back to the main error channel
			}
		}()
	}

	level.Error(logger).Log("exit", <-errs)
}
