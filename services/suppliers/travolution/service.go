package travolution

import (
	"context"
	"fmt"
	"swallow-supplier/config"
	customError "swallow-supplier/error"
	"swallow-supplier/utils/client"
)

// Yanolja required fields for accessing services of yanolja
type Travolution struct {
	Service *client.Request
	Ctx     context.Context
	Host    string
	Route   string
}

const (
	// ServiceName represents the name of the service
	ServiceName = "Travolution"
)

// New initialize Yanolja
func New(ctx context.Context) (t *Travolution, err error) {
	t = &Travolution{}
	cf := config.Instance()

	t.Service = &client.Request{}
	t.Service.CustomRequest = client.CustomRequest{
		RequestTimeout: client.RequestTimeout,
		Retries:        5,
	}

	request := client.NewRequest(t.Service.CustomRequest)
	request.AddHeader("accept", `*/*`)
	request.AddHeader("channel-code", cf.ChannelCode)
	request.AddHeader("Authorization", cf.TravolutionAuthorizationKey)

	t.Service = request
	t.Ctx = ctx

	if host := cf.TravolutionDomain; host != "" {
		t.Host = host
	} else {
		return t, customError.NewError(ctx, "connection_error", fmt.Sprintf(customError.ErrExternalServiceNotConfigured.Error(), ServiceName), nil)
	}

	return t, err
}
