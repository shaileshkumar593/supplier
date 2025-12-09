package yanolja

import (
	"context"
	"fmt"
	"swallow-supplier/config"
	customError "swallow-supplier/error"
	"swallow-supplier/utils/client"
)

// Yanolja required fields for accessing services of yanolja
type Yanolja struct {
	Service *client.Request
	Ctx     context.Context
	Host    string
	Route   string
}

const (
	// ServiceName represents the name of the service
	ServiceName = "Yanolja"
)

// New initialize Yanolja
func New(ctx context.Context) (c *Yanolja, err error) {
	c = &Yanolja{}
	cf := config.Instance()

	c.Service = &client.Request{}
	c.Service.CustomRequest = client.CustomRequest{
		RequestTimeout: client.RequestTimeout,
		Retries:        5,
	}

	request := client.NewRequest(c.Service.CustomRequest)
	request.AddHeader("accept", `*/*`)
	request.AddHeader("channel-code", cf.ChannelCode)
	request.AddHeader("x-api-key", cf.YanoljaApiKey)

	c.Service = request
	c.Ctx = ctx

	if host := cf.YanoljaDomain; host != "" {
		c.Host = host
	} else {
		return c, customError.NewError(ctx, "connection_error", fmt.Sprintf(customError.ErrExternalServiceNotConfigured.Error(), ServiceName), nil)
	}

	return c, err
}
