package babel

import (
	"context"
	"fmt"
	"swallow-supplier/config"
	customError "swallow-supplier/error"
	"swallow-supplier/utils/client"
)

// Babel required fields for accessing Customer Orchestrator
type Babel struct {
	Service *client.Request
	Ctx     context.Context
	Host    string
	Route   string
	APIKey  string
}

const (
	// ServiceName represents the name of the service
	ServiceName = "Babel"
)

// New initialize Yanolja
func New(ctx context.Context) (c *Babel, err error) {
	c = &Babel{}
	cf := config.Instance()

	c.Service = &client.Request{}
	c.Service.CustomRequest = client.CustomRequest{
		RequestTimeout: client.RequestTimeout,
		Retries:        5,
	}

	request := client.NewRequest(c.Service.CustomRequest)
	request.AddHeader("Content-Type", `application/json`)
	request.AddHeader("channel-code", cf.ChannelCode)
	request.AddHeader("Authorization", cf.BabelApiKey)

	c.Service = request
	c.Ctx = ctx

	if host := cf.BabelDomain; host != "" {
		c.Host = host
	} else {
		return c, customError.NewError(ctx, "connection_error", fmt.Sprintf(customError.ErrExternalServiceNotConfigured.Error(), ServiceName), nil)
	}

	return c, err
}
