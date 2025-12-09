package error

import (
	"errors"
	"net/http"
)

// StructuredErrorDetails struct
type StructuredErrorDetails struct {
	code    string
	status  int
	message string
}

// errorDetails ...
type errorDetails []StructuredErrorDetails

// GetErrorByCode get error code
func GetErrorByCode(code string) (StructuredErrorDetails, error) {
	for _, e := range AllErrors {
		if e.code == code {
			return e, nil
		}
	}

	return StructuredErrorDetails{}, errors.New("error `" + code + "` ")
}

// List of custom errors
// Error message should be on lowercase and should not end on a punctuation marks
// If more than one sentence, separate into period then space but the last sentence should not end on punctuation marks
// This will be converted properly on NewError function
var (
	ErrCmdRepository                = errors.New("unable to command repository")
	ErrResourceNotFound             = errors.New("resource not found")
	ErrApiKeyRequired               = errors.New("missing api key on header")
	ErrApiKeyInvalid                = errors.New("access denied due to invalid apikey on header")
	ErrChannelCodeRequired          = errors.New("missing authentication channel code on header")
	ErrChannelInvalid               = errors.New("access denied due to invalid channel code on header")
	ErrUndefinedRoute               = errors.New("missing auth definition on the given route")
	ErrAccessDenied                 = errors.New("access is denied for this resource")
	ErrDatabase                     = errors.New("repository Error: %s")
	ErrValidation                   = errors.New("validation errors")
	ErrInvalidBody                  = errors.New("request body is invalid")
	ErrInvalidModel                 = errors.New("invalid model %s")
	ErrMissingHeaders               = errors.New("access denied due to missing required headers")
	ErrNoRecordsFound               = errors.New("no records found")
	ErrExternalService              = errors.New("%s service error")
	ErrExternalServiceNotConfigured = errors.New("invalid %s service configuration")
	ErrForbiddenClient              = errors.New("client is not connected with vpn")
	ErrGatewayTimeout               = errors.New("client gateway timeout")
	ErrBadRequest                   = errors.New("Invalid request parameters")
	ErrTokenNotProvided             = errors.New("Token is omitted")
	ErrInvalidtoken                 = errors.New("Invalid Token")
)

// AllErrors list of all errors
var AllErrors = errorDetails{
	{
		// When the post params does not follow the model definition
		"leisure-api-0001",
		http.StatusBadRequest,
		"Invalid request",
	},
	{
		// Dealing with service level authentication when api key is invalid
		"leisure-api-0002",
		http.StatusUnauthorized,
		"Authentication failed.	",
	},
	{
		// When access to resource is not guaranteed
		"leisure-api-0003",
		http.StatusForbidden,
		"You don't have access rights.",
	},
	{
		// When service is not available
		"leisure-api-0004",
		http.StatusInternalServerError,
		"Some services are not available.",
	},
	{
		/// Dealing with service level authentication when tenant code is invalid
		"leisure-api-0005",
		http.StatusInternalServerError,
		"Server Internal Error",
	},
	{
		"leisure-api-0006",
		http.StatusBadRequest,
		"Invalid request",
	},
	{
		"leisure-api-0007",
		http.StatusInternalServerError,
		"An internal service call failed with an invalid request",
	},
	{
		"leisure-api-0008",
		http.StatusInternalServerError,
		"REST service call failed",
	},
	{
		"leisure-api-0009",
		http.StatusInternalServerError,
		"REST service errors",
	},
	{
		"leisure-api-0010",
		http.StatusInternalServerError,
		"Order Service Error",
	},
	{
		"leisure-api-0011",
		http.StatusInternalServerError,
		"Out of stock",
	},
	{
		"leisure-api-0012",
		http.StatusBadRequest,
		"Invalid request",
	},
	{
		"leisure-api-0013",
		http.StatusInternalServerError,
		"You have exceeded the maximum number of purchases at one time",
	},
	{
		"leisure-api-0015",
		http.StatusInternalServerError,
		"This is a product version that cannot be ordered",
	},
	{
		"leisure-api-0016",
		http.StatusInternalServerError,
		"I can't find the product",
	},
	{
		"leisure-api-0017",
		http.StatusInternalServerError,
		"Variant not found.",
	},
	{
		"leisure-api-0018",
		http.StatusInternalServerError,
		"This item has been discontinued",
	},
	{
		"leisure-api-0019",
		http.StatusInternalServerError,
		"This item is out of stock.",
	},
	{
		"leisure-api-0020",
		http.StatusBadRequest,
		"The date of the order is invalid.",
	},
	{
		"leisure-api-0021",
		http.StatusInternalServerError,
		"There is no order date that matches the requested date.",
	},
	{
		"leisure-api-0022",
		http.StatusBadRequest,
		"This is an invalid order.",
	},
	{
		"leisure-api-0023",
		http.StatusInternalServerError,
		"There are no order rounds that match the requested round.",
	},
	{
		"leisure-api-0024",
		http.StatusInternalServerError,
		"We can't find any personal information about your order.",
	},
	{
		"leisure-api-1001",
		http.StatusNotFound,
		"I can't find the product.",
	},
	{
		"leisure-api-1003",
		http.StatusNotFound,
		"Variant not found.",
	},
	{
		"leisure-api-1012",
		http.StatusNotFound,
		"object not found.",
	},
	{
		"leisure-api-1013",
		http.StatusNotFound,
		"Unable to find AuthorityId for channelCode.",
	},
	{
		"leisure-api-1014",
		http.StatusNotFound,
		"I can't find my partner information.",
	},
	{
		"leisure-api-1015",
		http.StatusInternalServerError,
		"mongo repository error",
	},
	{
		"leisure-api-1016",
		http.StatusInternalServerError,
		"validation error",
	},
	{
		"leisure-api-1017",
		http.StatusInternalServerError,
		"any data type validation error",
	},
	{
		"leisure-api-1018",
		http.StatusBadRequest,
		"service name is empty",
	},
	{
		"leisure-api-1019",
		http.StatusInternalServerError,
		"marshaling/unmarshalling error",
	},
	{
		"leisure-api-1020",
		http.StatusInternalServerError,
		"redis keys value empty",
	}, {
		"leisure-api-1021",
		http.StatusInternalServerError,
		"redis access failed",
	}, {
		"leisure-api-1022",
		http.StatusInternalServerError,
		"trip request error",
	}, {
		"leisure-api-00023",
		http.StatusBadRequest,
		"Invalid quantity",
	}, {
		"leisure-api-00024",
		http.StatusInternalServerError,
		"Empty inventory for Product",
	}, {
		"leisure-api-00025",
		http.StatusGatewayTimeout,
		"Client Gateway Timeout",
	}, {
		// When product version is not greater than existing one
		"leisure-api-0026",
		http.StatusBadRequest,
		"Product Version is Invalid",
	},
	// Get  information
	{
		"PRODUCT_NOT_FOUND",
		http.StatusNotFound,
		"Product information not found",
	}, {
		"UNAUTHORIZED_PRODUCT_ACCESS",
		http.StatusBadRequest,
		"Unauthorized to access the product",
	}, {
		"UNAUTHORIZED_OPTION_ACCESS",
		http.StatusBadRequest,
		"Unauthorized to access the option",
	}, {
		"OPTION_NOT_FOUND",
		http.StatusNotFound,
		"Option information not found",
	}, {
		"OPTION_NOT_ORDERABLE",
		http.StatusNotFound,
		"Option is not orderable",
	}, {
		"UNIT_NOT_FOUND",
		http.StatusBadRequest,
		"Unit information not found",
	}, {
		"UNAVAILABLE_PRODUCT_CONTENT",
		http.StatusBadRequest,
		"Product content in specified language is unavailable",
	}, {
		"PRODUCT_TYPE_ERROR",
		http.StatusBadRequest,
		"The product is not of BK type",
	},
	{
		"SCHEDULE_NOT_FOUND",
		http.StatusNotFound,
		"There are no available schedules",
	},
	{
		"INVALID_DATE_FORMAT",
		http.StatusBadRequest,
		"Date format is invalid",
	}, {
		"INVALID_TIME_FORMAT",
		http.StatusBadRequest,
		"Time format is invalid",
	}, {
		"INVALID_ADDTIONALINFO_UID",
		http.StatusBadRequest,
		"Invalid AdditionalInfoUid",
	},

	// Create Booking
	{
		"INVALID_UNITAMOUNTS_FORMAT",
		http.StatusBadRequest,
		"Invalid UnitAmounts format",
	}, {
		"MISSING_UNITAMOUNTS",
		http.StatusBadRequest,
		"UnitAmounts is missing",
	}, {
		"MISSING_PRODUCT_UID",
		http.StatusBadRequest,
		"ProductUid is missing",
	}, {
		"MISSING_OPTION_UID",
		http.StatusBadRequest,
		"OptionUid is missing",
	}, {
		"MISSING_UNIT_UID",
		http.StatusBadRequest,
		"Unit is missing",
	}, {
		"MISSING_ORDER_AMOUNT",
		http.StatusBadRequest,
		"Order amount is missing",
	}, {
		"UNAUTHORIZED_PRODUCT_ORDER",
		http.StatusBadRequest,
		"Unauthorized to order product",
	}, {
		"UNAUTHORIZED_OPTION_ORDER",
		http.StatusBadRequest,
		"Unauthorized to order option",
	}, {
		"UNIT_NOT_ORDERABLE",
		http.StatusBadRequest,
		"Unit is not orderable",
	}, {
		"MINIMUM_AMOUNT_REQUIRED",
		http.StatusBadRequest,
		"Order amounts are less than the minimal requirement",
	}, {
		"MAXIMUM_AMOUNT_REACHED",
		http.StatusBadRequest,
		"Order amounts are higher than the maximum requirement",
	}, {
		"INVALID_AMOUNT_UNIT",
		http.StatusBadRequest,
		"Invalid amounts. Must order {minAmount} units plus multiples of {unitAmount}",
	}, {
		"INVALID_SCHEDULE_DATETIME",
		http.StatusBadRequest,
		"Invalid schedule datetime",
	}, {
		"INVALID_BOOKINGDATE_FORMAT",
		http.StatusBadRequest,
		"Invalid bookingDate format",
	}, {
		"INVALID_BOOKINGTIME_FORMAT",
		http.StatusBadRequest,
		"Invalid bookingTime format",
	}, {
		"SCHEDULE_CLOSED",
		http.StatusBadRequest,
		"Schedule is closed",
	}, {
		"LOW_STOCK_FOR_SCHEDULE",
		http.StatusBadRequest,
		"Insufficient stock for the schedule",
	}, {
		"LOW_STOCK",
		http.StatusBadRequest,
		"Insufficient stock for the order",
	}, {
		"INVALID_ADDITIONALINFO_UID",
		http.StatusBadRequest,
		"Invalid AdditionalInfoUid",
	}, {
		"MISSING_ADDITIONALINFO_UID",
		http.StatusBadRequest,
		"AdditionalInfo is missing,",
	}, {
		"MISSING_ADDITIONALINFO_VALUE",
		http.StatusBadRequest,
		"Value of AdditionalInfo is missing",
	}, {
		"INVALID_ADDITIONALINFO",
		http.StatusBadRequest,
		"The bookingAdditionalInfo field is required when product is {productUid}",
	}, {
		"INVALID_ADDITIONALINFO_OBJECT",
		http.StatusBadRequest,
		"bookingAdditionalInfo must be a Object",
	}, {
		"INVALID_ADDITIONALINFO_TRAVELER_REQUIRED",
		http.StatusBadRequest,
		"The bookingAdditionalInfo.traveler field is required when product is {productUid}",
	}, {
		"INVALID_ADDITIONALINFO_TRAVELER_TYPE",
		http.StatusBadRequest,
		"bookingAdditionalInfo.traveler must be a JSON array",
	}, {
		"INVALID_ADDITIONALINFO_TRAVELER_IDX_TYPE",
		http.StatusBadRequest,
		"bookingAdditionalInfo.traveler.{idx} must be a JSON array",
	}, {
		"INVALID_ADDITIONALINFO_TRAVELER_IDX_REQUIRED",
		http.StatusBadRequest,
		"The bookingAdditionalInfo.traveler.{idx} field is required when product is {productUid}",
	}, {
		"KORAILPASS_ERROR",
		http.StatusBadRequest,
		"Varies depending on the cause of the error",
	},

	// Search Booking
	{
		"ORDER_NOT_FOUND",
		http.StatusNotFound,
		"Order information not found",
	},

	// Cancel Booking
	{
		"ORDER_EXPIRED",
		http.StatusBadRequest,
		"Order cannot be canceled as it has expired",
	}, {
		"ORDER_ALREADY_USED",
		http.StatusBadRequest,
		"Order cannot be canceled as it is already used",
	}, {
		"ORDER_ALREADY_CANCELLED",
		http.StatusBadRequest,
		"Order has already been canceled",
	}, {
		"ORDER_CANNOT_BE_CANCELED",
		http.StatusBadRequest,
		"This order cannot be canceled",
	},
}
