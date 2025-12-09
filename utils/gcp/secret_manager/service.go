package secret_manager

import (
	"context"
	"fmt"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// SecretManager interface defines methods for managing secrets.
type SecretManager interface {
	FetchSecret(ctx context.Context, secretID string) (string, error)
}

// SecretManagerService implements SecretManager.
type SecretManagerService struct {
	client    *secretmanager.Client
	projectID string
	logger    log.Logger
}

// NewSecretManager initializes a new SecretManagerService.
func NewSecretManager(logger log.Logger) (SecretManager, error) {
	projectID := os.Getenv("GCP_PROJECT_ID")
	level.Info(logger).Log("projectId ", projectID)
	if projectID == "" {
		level.Error(logger).Log("msg", "GCP_PROJECT_ID environment variable is missing")
		return nil, fmt.Errorf("GCP_PROJECT_ID environment variable is missing")
	}

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		level.Error(logger).Log("msg", "Failed to create Secret Manager client", "err", err)
		return nil, fmt.Errorf("failed to create Secret Manager client: %w", err)
	}

	return &SecretManagerService{
		client:    client,
		projectID: projectID,
		logger:    logger,
	}, nil
}
