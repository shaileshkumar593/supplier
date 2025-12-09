package pdfvoucher

import (
	"context"
	"fmt"
	"swallow-supplier/config"
	customError "swallow-supplier/error"
	"swallow-supplier/utils"
	"swallow-supplier/utils/client"
)

// Trip required fields for accessing Customer Orchestrator
type VoucherPdf struct {
	Service *client.Request
	Ctx     context.Context
	Host    string
}

const (
	// ServiceName represents the name of the service
	ServiceName = "Voucher_Pdf_Generator"
)

// New initialize Trip
func New(ctx context.Context) (vp *VoucherPdf, err error) {
	vp = &VoucherPdf{}
	cf := config.Instance()
	vp.Service = &client.Request{}
	vp.Service.CustomRequest = client.CustomRequest{
		RequestTimeout: client.RequestTimeout,
		Retries:        5,
	}

	request := client.NewRequest(vp.Service.CustomRequest)
	request.AddHeader("Accept", client.ContentTypeJSON)
	request.AddHeader("User-Agent", cf.TripUserAgent)
	request.AddHeader("X-MM-Request-ID", utils.GenerateUUID("", true))
	fmt.Println("4")

	vp.Service = request
	vp.Ctx = ctx
	fmt.Println("5")

	if host := cf.PdfVoucherUrl; host != "" {
		vp.Host = host
	} else {
		return vp, customError.NewError(ctx, "connection_error", fmt.Sprintf(customError.ErrExternalServiceNotConfigured.Error(), ServiceName), nil)
	}

	return vp, err
}
