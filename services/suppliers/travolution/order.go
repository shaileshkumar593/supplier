package travolution

import (
	"context"
	"fmt"
	"net/http"
	"swallow-supplier/request_response/travolution"
	"swallow-supplier/utils/client"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// TravolutionOrder  will fetch all type of order
func (t *Travolution) TravolutionOrder(ctx context.Context, req travolution.OrderRequest, typeOfBooking string) (res travolution.Response, err error) {
	logger := log.With(t.Service.Logger, "method", "OrderBooking")
	level.Info(logger).Log("info", "Service OrderBooking")

	var payload interface{}

	switch typeOfBooking {
	case "TK":
		/* optionId, err := strconv.Atoi(req.Option)
		if err != nil {
			level.Info(logger).Log("error", err)
			return res, err
		} */

		option, ok := req.Option.(float64) //strconv.Atoi(req.Option)
		if !ok {
			level.Info(logger).Log("option is not int")
			return res, fmt.Errorf("option is not int, got %T", req.Option)
		}

		optionId := int(option)

		payload = travolution.TicketRequest{
			Product:               req.Product,
			Option:                optionId,
			UnitAmounts:           req.UnitAmounts,
			ReferenceNumber:       req.ReferenceNumber,
			VoucherSendType:       int(req.VoucherSendType),
			TravelerName:          req.TravelerName,
			TravelerContactEmail:  req.TravelerContactEmail,
			TravelerContactNumber: req.TravelerContactNumber,
			TravelerNationality:   req.TravelerNationality,
		}
		// Handle Ticket type
	case "BK":
		option, ok := req.Option.(float64) //strconv.Atoi(req.Option)
		if !ok {
			level.Info(logger).Log("option is not int")
			return res, fmt.Errorf("option is not int, got %T", req.Option)
		}

		optionId := int(option)

		payload = travolution.BookingRequest{
			Product:               req.Product,
			Option:                optionId,
			UnitAmounts:           req.UnitAmounts,
			BookingDate:           req.BookingDate,
			BookingTime:           req.BookingTime,
			BookingAdditionalInfo: req.BookingAdditionalInfo,
			ReferenceNumber:       req.ReferenceNumber,
			VoucherSendType:       req.VoucherSendType,
			TravelerName:          req.TravelerName,
			TravelerContactEmail:  req.TravelerContactEmail,
			TravelerContactNumber: req.TravelerContactNumber,
			TravelerNationality:   req.TravelerNationality,
		}

	case "PAS", "PKG":
		optionId, ok := req.Option.(string) //strconv.Atoi(req.Option)
		if !ok {
			level.Info(logger).Log("option is not int")
			return res, fmt.Errorf("option is not int, got %T", req.Option)
		}
		payload = travolution.PASSOrPKGRequest{
			Product:               req.Product,
			Option:                optionId,
			UnitAmounts:           req.UnitAmounts,
			ReferenceNumber:       req.ReferenceNumber,
			VoucherSendType:       int(req.VoucherSendType),
			TravelerName:          req.TravelerName,
			TravelerContactEmail:  req.TravelerContactEmail,
			TravelerContactNumber: req.TravelerContactNumber,
			TravelerNationality:   req.TravelerNationality,
		}
	default:
		// Unknown / fallback
		res.Body = fmt.Sprintf("Unsupported type: %s", typeOfBooking)
		res.Code = "400"
		return res, nil
	}

	host := t.Host + "/api/partner/v1.1/orders"
	response, err := t.Service.Send(
		ctx,
		ServiceName,
		host,
		http.MethodPost,
		client.ContentTypeJSON,
		payload,
	)
	level.Error(logger).Log("info", "url", host, "response", response, "err", err)

	return ResponseConvertor(t.Ctx, response, logger, err)
}

// TravolutionGetOrder fetch order by orderNumber
func (t *Travolution) TravolutionGetOrder(ctx context.Context, orderNumber string) (res travolution.Response, err error) {
	logger := log.With(t.Service.Logger, "method", "TravolutionGetOrder")
	level.Info(logger).Log("info", "Service TravolutionGetOrder")

	host := t.Host + fmt.Sprintf("/api/partner/v1.1/orders/%s", orderNumber)
	response, err := t.Service.Send(
		ctx,
		ServiceName,
		host,
		http.MethodGet,
		client.ContentTypeJSON,
		nil,
	)
	level.Error(logger).Log("info", "url", host, "response", response, "err", err)

	return ResponseConvertor(t.Ctx, response, logger, err)
}

// TravolutionCancelOrder  cancel order based on orderNumber
func (t *Travolution) TravolutionCancelOrder(ctx context.Context, orderNumber string) (res travolution.Response, err error) {
	logger := log.With(t.Service.Logger, "method", "OrderBooking")
	level.Info(logger).Log("info", "Service OrderBooking")

	host := t.Host + "/api/partner/v1.1/orders/"
	response, err := t.Service.Send(
		ctx,
		ServiceName,
		host,
		http.MethodDelete,
		client.ContentTypeJSON,
		orderNumber,
	)
	level.Error(logger).Log("info", "url", host, "response", response, "err", err)

	return ResponseConvertor(t.Ctx, response, logger, err)
}
