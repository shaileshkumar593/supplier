package trip

import (
	"context"
	"fmt"
	"swallow-supplier/config"
	customError "swallow-supplier/error"
	"swallow-supplier/utils"
	"swallow-supplier/utils/client"
)

// Trip required fields for accessing Customer Orchestrator
type Trip struct {
	Service *client.Request
	Ctx     context.Context
	Host    string
}

const (
	// ServiceName represents the name of the service
	ServiceName = "Trip"
)

// New initialize Trip
func New(ctx context.Context) (c *Trip, err error) {
	c = &Trip{}
	cf := config.Instance()
	c.Service = &client.Request{}
	c.Service.CustomRequest = client.CustomRequest{
		RequestTimeout: client.RequestTimeout,
		Retries:        5,
	}

	request := client.NewRequest(c.Service.CustomRequest)
	request.AddHeader("Accept", client.ContentTypeJSON)
	request.AddHeader("User-Agent", cf.TripUserAgent)
	request.AddHeader("X-MM-Request-ID", utils.GenerateUUID("", true))
	fmt.Println("4")

	c.Service = request
	c.Ctx = ctx
	fmt.Println("5")

	if host := cf.Trip; host != "" {
		c.Host = host
	} else {
		return c, customError.NewError(ctx, "connection_error", fmt.Sprintf(customError.ErrExternalServiceNotConfigured.Error(), ServiceName), nil)
	}

	return c, err
}
