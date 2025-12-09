package config

import (
	"context"
	"fmt"
	"os"

	"swallow-supplier/utils/gcp/secret_manager"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/joho/godotenv"

	"github.com/kelseyhightower/envconfig"
)

// AppConfig ...
type AppConfig struct {
	AppEnv              string `envconfig:"APP_ENV"`
	AppPort             string `envconfig:"APP_PORT"`
	AppDomain           string `envconfig:"APP_DOMAIN"`
	AppServiceName      string `envconfig:"APP_SERVICE_NAME"`
	AuthScopeConfigPath string `envconfig:"AUTH_SCOPE_DEF_PATH"`
	AuthorizationKey    string `envconfig:"AUTHORIZATION_KEY"`

	// postgres

	// MongoDb
	DatabaseConnectionMongo string `envconfig:"DB_CONNECT_STRING_MONGO"`
	MongoDBName             string `envconfig:"MONGO_DB_NAME"`

	//Database config
	DatabaseDefaultMaxPoolSize        int `envconfig:"DB_MAX_POOL_SIZE"`
	DatabaseMaxConnIdleTime           int `envconfig:"DB_MAX_CONN_IDLE_TIME"`
	DatabaseMaxLifeTime               int `envconfig:"DB_MAX_LIFE_TIME_IN_MINUTES"`
	DatabaseMaxIdleTime               int `envconfig:"DB_MAX_IDLE_TIME_IN_SECONDS"`
	DatabaseDefaultMaxOpenConnections int `envconfig:"DB_DEFAULT_MAX_OPEN_CONNECTIONS"`
	DatabaseDefaultMaxIdleConnections int `envconfig:"DB_DEFAULT_MAX_IDLE_CONNECTIONS"`

	//Redis
	RedisURL  string `envconfig:"REDIS_URL"`
	CacheName string `envconfig:"CACHE_NAME"`

	//GCP
	GooglePlaceIdKey     string `envconfig:"GCP_GOOGLE_PLACEID_KEY_IP_RESTRICTED"`
	GooglePlaceIdBaseUrl string `envconfig:"GOOGLE_PLACE_ID_BASEURL"`
	GcpProjectId         string `envconfig:"GCP_PROJECT_ID"`

	//Circuit Breaker
	CircuitBreakerEnable       string `envconfig:"CIRCUIT_BREAKER_ENABLE"`
	CircuitBreakerRequests     string `envconfig:"CIRCUIT_BREAKER_REQUESTS"`
	CircuitBreakerFailureRatio string `envconfig:"CIRCUIT_BREAKER_FAILURE_RATIO"`

	//GGT
	ChannelCode string `envconfig:"CHANNEL_CODE"`

	// Supplier
	// Yanolja config
	YanoljaApiKey string `envconfig:"YANOLJA_API_KEY"`
	YanoljaDomain string `envconfig:"YANOLJA_DOMAIN"`

	//Babel API
	BabelApiKey string `envconfig:"BABEL_API_KEY"`
	BabelDomain string `envconfig:"BABEL_DOMAIN"`

	// Travolution
	TravolutionAuthorizationKey string `envconfig:"TRAVOLUTION_AUTHORIZATION_KEY"`
	TravolutionDomain           string `envconfig:"TRAVOLUTION_DOMAIN"`

	//distributor Config
	//Trip
	Trip           string `envconfig:"TRIP"`
	TripMockHTTP   string `envconfig:"TRIP_MOCK_HTTP"`
	TripUserAgent  string `envconfig:"TRIP_USER_AGENT"`
	TripAdminToken string `envconfig:"TRIP_ADMIN_TOKEN"`
	TripSyncUrl    string `envconfig:"TRIP_SYNC_URL"`
	PdfVoucherUrl  string `envconfig:"PDF_VOUCHER"`
	//Excel
	YGTFilePath     string `envconfig:"YGT_FILE_PATH"`
	CleanedDataPath string `envconfig:"CLEANED_DATA_PATH"`

	//SCHEDULE
	Schedule int `envconfig:"SCHEDULE"`
}

var (
	instance *AppConfig
	logger   log.Logger
)

// Init initializes the configuration based on the environment
func Init() {
	// Initialize structured logging
	logger = log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	var cfg AppConfig
	appEnv := os.Getenv("APP_ENV")

	level.Info(logger).Log("************AppEnv ************** ", appEnv)
	if appEnv == "LOCAL" || appEnv == "" {
		// Load .env file
		if err := godotenv.Load(".env"); err != nil {
			level.Error(logger).Log("msg", "Failed to load .env file", "err", err)
			os.Exit(1)
		}
		level.Info(logger).Log("msg", "Loaded environment variables from .env")
	} else {
		ctx := context.Background()
		secretManager, err := secret_manager.NewSecretManager(logger)
		if err != nil {
			level.Error(logger).Log("error ", "failed to connect secret manager")
			os.Exit(1)
		}

		fmt.Println("************** All Secrets ********************* ")
		cfg.BabelApiKey, _ = secretManager.FetchSecret(ctx, "BABEL_API_KEY")
		cfg.GooglePlaceIdKey, _ = secretManager.FetchSecret(ctx, "GCP_GOOGLE_PLACEID_KEY_IP_RESTRICTED")

		if appEnv == "PRODUCTION" {

			cfg.YanoljaApiKey, _ = secretManager.FetchSecret(ctx, "YANOLJA_API_KEY_PROD")
			cfg.AuthorizationKey, _ = secretManager.FetchSecret(ctx, "YANOLJA_AUTHORIZATION_KEY_PROD")
			cfg.DatabaseConnectionMongo, _ = secretManager.FetchSecret(ctx, "MONGO_DB_URI_PROD")
			level.Info(logger).Log("msg", "Running in PRODUCTION mode, using system environment variables and secrets")

		} else if appEnv == "DEVELOPMENT" {
			cfg.YanoljaApiKey, _ = secretManager.FetchSecret(ctx, "YANOLJA_API_KEY_DEV")
			cfg.AuthorizationKey, _ = secretManager.FetchSecret(ctx, "YANOLJA_AUTHORIZATION_KEY_DEV")
			cfg.DatabaseConnectionMongo, _ = secretManager.FetchSecret(ctx, "MONGO_DB_URI_DEV")

			level.Info(logger).Log("msg", "Running in DEVELOPMENT mode, using system environment variables and secrets")
		}
	}

	// Load configurations from environment variables
	if err := envconfig.Process("", &cfg); err != nil {
		level.Error(logger).Log("msg", "Failed to load configs from environment variables", "err", err)
		os.Exit(1)
	}

	instance = &cfg
	level.Info(logger).Log("msg", "Configuration initialized successfully")
}

// Instance returns the singleton configuration instance
func Instance() *AppConfig {
	if instance == nil {
		Init()
	}
	return instance
}
