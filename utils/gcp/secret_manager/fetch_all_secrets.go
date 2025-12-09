package secret_manager

import (
	"context"
	"fmt"
	"hash/crc32"
	"os"
	"strings"
	"sync"

	secretspb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/api/cloudresourcemanager/v3"
)

var (
	projectNumber string
	mutex         sync.Mutex
	secretCache   = make(map[string]string)
)

// getProjectNumber retrieves the project number from GCP
func getProjectNumber(ctx context.Context, logger log.Logger) string {
	mutex.Lock()
	defer mutex.Unlock()

	if projectNumber != "" {
		return projectNumber // Return cached project number
	}

	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		level.Error(logger).Log("msg", "GCP_PROJECT_ID environment variable is required")
		os.Exit(1)
	}

	cloudResourceManagerService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		level.Error(logger).Log("msg", "Failed to initialize Cloud Resource Manager client", "err", err)
		os.Exit(1)
	}

	projectResourceName := "projects/" + projectID
	projectResp, err := cloudResourceManagerService.Projects.Get(projectResourceName).Do()
	if err != nil {
		level.Error(logger).Log("msg", "Failed to retrieve project details", "projectID", projectID, "err", err)
		os.Exit(1)
	}

	projectNumber = extractProjectNumber(projectResp.Name)
	if projectNumber == "" {
		level.Error(logger).Log("msg", "Failed to extract Project Number", "projectResponse", projectResp)
		os.Exit(1)
	}

	level.Info(logger).Log("msg", "✅ Retrieved Project Number", "projectNumber", projectNumber)
	return projectNumber
}

// extractProjectNumber parses "projects/409931392284" and returns "409931392284"
func extractProjectNumber(resourceName string) string {
	parts := strings.Split(resourceName, "/")
	if len(parts) > 1 {
		return parts[1] // Extract project number
	}
	return ""
}

// fetchSecret retrieves a specific secret from Google Secret Manager.
func (s *SecretManagerService) FetchSecret(ctx context.Context, secretID string) (string, error) {
	// Get project number if not already set
	projectNum := getProjectNumber(ctx, s.logger)
	secretName := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectNum, secretID)

	req := &secretspb.AccessSecretVersionRequest{Name: secretName}

	result, err := s.client.AccessSecretVersion(ctx, req)
	if err != nil {
		level.Error(s.logger).Log("msg", "Failed to access secret", "secret", secretID, "err", err)
		if strings.Contains(err.Error(), "PermissionDenied") {
			return "", fmt.Errorf("insufficient permissions to access secret %s", secretID)
		}
		return "", fmt.Errorf("failed to access secret %s: %w", secretID, err)
	}

	// Validate checksum for data integrity
	crc32c := crc32.MakeTable(crc32.Castagnoli)
	checksum := int64(crc32.Checksum(result.Payload.Data, crc32c))
	if checksum != *result.Payload.DataCrc32C {
		level.Error(s.logger).Log("msg", "Data corruption detected", "secret", secretID)
		return "", fmt.Errorf("data corruption detected in secret %s", secretID)
	}

	level.Info(s.logger).Log("msg", "Successfully fetched secret", "secret", secretID)
	return string(result.Payload.Data), nil
}

// GetSecrets retrieves secrets using their full secret names.
/*func (s *SecretManagerService) GetAllSecrets(secretIDs []string) map[string]string {
	ctx := context.Background()
	secrets := make(map[string]string)

	for _, secretID := range secretIDs {
		secretValue, err := s.fetchSecret(ctx, secretID)
		if err != nil {
			level.Error(s.logger).Log("msg", "Failed to fetch secret", "secret", secretID, "err", err)
			continue // Log error and continue
		}
		secrets[secretID] = secretValue
	}

	level.Info(s.logger).Log("msg", "Successfully retrieved secrets", "total_secrets", len(secrets))
	return secrets
}*/
