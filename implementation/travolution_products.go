package implementation

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime/debug"
	"strconv"
	"swallow-supplier/config"
	customError "swallow-supplier/error"
	domain "swallow-supplier/mongo/domain/travolution"
	"swallow-supplier/request_response/travolution"
	travolutionSvc "swallow-supplier/services/suppliers/travolution"

	"swallow-supplier/utils"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// GetAllproducts
func (s *service) GetAllproducts(ctx context.Context, req travolution.ProductReq) (resp travolution.Response, err error) {
	var requestID string

	requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "GetAllproducts",
		"Request ID", requestID,
	)
	// Defer panic recovery
	defer func(ctx context.Context) {
		if r := recover(); r != nil {
			level.Info(logger).Log("info", "processing request went into panic mode", "panic", r)
			resp.Code = "500"
		}
	}(ctx)

	level.Info(logger).Log(" info ", "travolution service call")

	var tsvc, _ = travolutionSvc.New(ctx)
	resp, err = tsvc.GetProducts(ctx, req)
	if err != nil {
		level.Error(logger).Log("treavolution error ", err)

		if resp.Code == "403" {
			err = customError.NewError(ctx, "leisure-api-0003", fmt.Sprint(customError.ErrForbiddenClient.Error(), "GetProductByUid"), nil)
		} else {
			err = fmt.Errorf("request to travolution client raised error")
			resp.Code = "500"
			resp.Body = err
		}
		return resp, err
	}

	level.Info(logger).Log("response ", resp)
	return resp, nil
}

// GetProductByUid
func (s *service) GetProductByUid(ctx context.Context, req travolution.ProductReq) (resp travolution.Response, err error) {
	var requestID string

	requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "GetProductByUid",
		"Request ID", requestID,
	)
	// Defer panic recovery
	defer func(ctx context.Context) {
		if r := recover(); r != nil {
			level.Info(logger).Log("info", "processing request went into panic mode", "panic", r)
			resp.Code = "500"
		}
	}(ctx)

	level.Info(logger).Log(" info ", "travolution service call")

	var tsvc, _ = travolutionSvc.New(ctx)
	resp, err = tsvc.GetProducts(ctx, req)
	if err != nil {
		level.Error(logger).Log("treavolution error ", err)

		if resp.Code == "403" {
			err = customError.NewError(ctx, "leisure-api-0003", fmt.Sprint(customError.ErrForbiddenClient.Error(), "GetProductByUid"), nil)
		} else {
			err = fmt.Errorf("request to travolution client raised error")
			resp.Code = "500"
			resp.Body = err
		}
		return resp, err
	}

	level.Info(logger).Log("response ", resp)
	return resp, nil
}

// PostCreateAllProduct
func (s *service) PostCreateAllProduct(ctx context.Context) (resp travolution.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "PostCreateAllProduct",
		"request_id", requestID,
	)

	defer func() {
		if r := recover(); r != nil {
			level.Error(logger).Log("event", "panic_recovered", "panic", r, "stacktrace", string(debug.Stack()))
			resp.Code = "500"
			resp.Body = "internal server error"
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()

	level.Info(logger).Log("msg", "calling Travolution service")

	// Create Travolution service client
	tsvc, err := travolutionSvc.New(ctx)
	if err != nil {
		level.Error(logger).Log("event", "travolution_service_init_failed", "error", err)
		resp.Code = "500"
		resp.Body = "failed to initialize travolution service"
		return resp, err
	}
	fmt.Println("========1========")

	// Initial fetch for all products
	req := travolution.ProductReq{}
	allProductResp, err := tsvc.GetProducts(ctx, req)
	if err != nil {
		return s.wrapTravolutionError(ctx, logger, resp, err)
	}

	fmt.Println("========2========")
	var genericList []map[string]any
	var rawProducts []travolution.RawProduct

	if genericSlice, ok := allProductResp.Body.([]interface{}); ok {
		genericList = make([]map[string]interface{}, 0, len(genericSlice))

		for _, item := range genericSlice {
			if m, ok := item.(map[string]interface{}); ok {
				genericList = append(genericList, m)
			} else {
				level.Error(logger).Log(
					"msg", "skipping non-map item in products slice",
					"type", fmt.Sprintf("%T", item),
				)
			}
		}

		level.Info(logger).Log(
			"msg", "unwrapped products from []interface{} -> []map[string]interface{}",
			"type", fmt.Sprintf("%T", resp.Body),
			"count", len(genericList),
		)

	} else {
		level.Error(logger).Log(
			"msg", "resp.Body is not []interface{}",
			"actual_type", fmt.Sprintf("%T", resp.Body),
		)
		return resp, fmt.Errorf("unexpected resp.Body type: %T", resp.Body)
	}

	//  Only proceed if we actually have data
	if len(genericList) == 0 {
		level.Warn(logger).Log("msg", "no products found in response")
		return resp, nil
	}

	// Marshal -> Unmarshal into RawProduct
	bytesJson, err := json.Marshal(genericList)
	if err != nil {
		level.Error(logger).Log(
			"msg", "failed to marshal generic products list",
			"error", err,
		)
		return resp, fmt.Errorf("failed to marshal generic products list: %w", err)
	}

	if err := json.Unmarshal(bytesJson, &rawProducts); err != nil {
		level.Error(logger).Log(
			"msg", "failed to unmarshal into RawProduct",
			"error", err,
		)
		return resp, fmt.Errorf("failed to unmarshal into RawProduct: %w", err)
	}

	var insertedIDs []string

	for _, product := range rawProducts {
		var optionList []travolution.RawProductOption
		req.ProductUid = product.UID
		req.Lang = "en"

		fmt.Println("========4========")

		// Fetch product details by productId
		productResp, err := tsvc.GetProducts(ctx, req)
		if err != nil {
			return s.wrapTravolutionError(ctx, logger, productResp, err)
		}
		fmt.Println("******************* productResp.Body ****************** ", productResp.Body)
		fmt.Printf("############## type of productResp.Body %T", productResp.Body)
		fmt.Println("========5========")

		rawProductDetailMap, ok := productResp.Body.(map[string]any)
		if !ok {
			resp.Code = "500"
			resp.Body = "invalid product detail type"
			return resp, fmt.Errorf("unexpected product detail type")
		}

		jsonBytes, err := json.Marshal(rawProductDetailMap)
		if err != nil {
			level.Error(logger).Log("msg", "failed to marshal map to JSON", "err", err)
			return resp, fmt.Errorf("marshal error: %w", err)
		}

		var rawProductDetail travolution.RawProduct
		// Unmarshal JSON to struct
		if err := json.Unmarshal(jsonBytes, &rawProductDetail); err != nil {
			level.Error(logger).Log("msg", "failed to unmarshal JSON to struct", "err", err)
			return resp, fmt.Errorf("unmarshal error: %w", err)
		}

		fmt.Println("========6========")

		// Fetch options
		optionReq := travolution.OptionRequest{
			ProductUid: rawProductDetail.UID,
			Lang:       "en",
		}

		optionsResp, err := tsvc.GetOptions(ctx, optionReq)
		if err != nil {
			return s.wrapTravolutionError(ctx, logger, optionsResp, err)
		}

		fmt.Println("******************* optionsResp.Body ****************** ", optionsResp.Body)
		fmt.Printf("############## type of optionsResp.Body %T", optionsResp.Body)

		fmt.Println("========7========")

		// Marshal back to JSON
		optionsJson, err := json.Marshal(optionsResp.Body)
		if err != nil {
			return resp, fmt.Errorf("marshal error: %w", err)
		}

		// Step 1: Unmarshal into slice of maps for normalization
		var tmpOptions []map[string]interface{}
		if err := json.Unmarshal(optionsJson, &tmpOptions); err != nil {
			level.Error(logger).Log(
				"msg", "failed to unmarshal options JSON into temp list",
				"error", err,
			)
			return resp, fmt.Errorf("failed to unmarshal options JSON: %w", err)
		}

		// Step 2: Normalize UID for each element
		for i, tmp := range tmpOptions {
			if v, ok := tmp["uid"]; ok {
				switch val := v.(type) {
				case string:
					tmpOptions[i]["uid"] = val
				case float64:
					tmpOptions[i]["uid"] = strconv.FormatInt(int64(val), 10)
				case json.Number:
					tmpOptions[i]["uid"] = val.String()
				default:
					return resp, fmt.Errorf("unsupported uid type in options: %T", v)
				}
			}
		}

		// Step 3: Marshal normalized list back to JSON
		normalizedOptions, err := json.Marshal(tmpOptions)
		if err != nil {
			level.Error(logger).Log(
				"msg", "failed to re-marshal normalized options JSON",
				"error", err,
			)
			return resp, fmt.Errorf("failed to re-marshal normalized options JSON: %w", err)
		}

		// Step 4: Unmarshal normalized JSON into struct slice
		if err := json.Unmarshal(normalizedOptions, &optionList); err != nil {
			level.Error(logger).Log(
				"msg", "failed to unmarshal normalized options JSON into RawOption",
				"error", err,
			)
			return resp, fmt.Errorf("failed to unmarshal normalized options JSON into RawOption: %w", err)
		}

		fmt.Println("\n========8====optionList====", optionList)

		// iterate  each option
		for idx, opt := range optionList {
			fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&& ", opt.UID)
			optionReq.OptionUid = opt.UID

			var optionRecord travolution.RawProductOption
			var unitAndPriceList []travolution.RawOptionUnitPrice
			var ScheduleList []domain.BookingSchedule
			fmt.Println("========9========")

			// Fetch option details (e.g., Notice)
			optionDetailResp, err := tsvc.GetOptions(ctx, optionReq)
			if err != nil {
				//return s.wrapTravolutionError(ctx, logger, optionDetailResp, err)
				return resp, err
			}

			fmt.Println("========10========")
			optionJson, err := DataConversionTOjson(optionDetailResp.Body, logger, "option")
			if err != nil {
				return resp, err
			}

			fmt.Println("=================  PPPPPPPPP  ================================")
			// Step 1: Unmarshal into a temp map
			var tmp map[string]interface{}
			if err := json.Unmarshal(optionJson, &tmp); err != nil {
				level.Error(logger).Log(
					"msg", "failed to unmarshal option JSON into temp map",
					"error", err,
				)
				return resp, fmt.Errorf("failed to unmarshal option JSON: %w", err)
			}

			// Step 2: Normalize UID (force to string)
			if v, ok := tmp["uid"]; ok {
				switch val := v.(type) {
				case string:
					tmp["uid"] = val
				case float64: // default JSON unmarshal for numbers
					tmp["uid"] = strconv.FormatInt(int64(val), 10)
				case json.Number:
					tmp["uid"] = val.String()
				default:
					return resp, fmt.Errorf("unsupported uid type: %T", v)
				}
			}

			// Step 3: Marshal normalized map back to JSON
			normalized, err := json.Marshal(tmp)
			if err != nil {
				level.Error(logger).Log(
					"msg", "failed to re-marshal normalized option JSON",
					"error", err,
				)
				return resp, fmt.Errorf("failed to re-marshal normalized option JSON: %w", err)
			}

			// Step 4: Unmarshal normalized JSON into struct
			//var optionRecord travolution.RawProductOption
			if err := json.Unmarshal(normalized, &optionRecord); err != nil {
				level.Error(logger).Log(
					"msg", "failed to unmarshal normalized JSON into RawProductOption",
					"error", err,
				)
				return resp, fmt.Errorf("failed to unmarshal into RawProductOption: %w", err)
			}

			// Step 5: Assign notice if optionRecord is not empty
			if !reflect.DeepEqual(optionRecord, travolution.RawProductOption{}) {
				optionList[idx].Notice = optionRecord.Notice
			}

			fmt.Println("========11========", "options")

			// ******************  Fetch units and prices  ****************************************//
			unitAndPriceReq := travolution.UnitPriceRequest{
				ProductUid: rawProductDetail.UID,
				OptionUid:  opt.UID,
			}
			fmt.Println("========12========")

			unitPriceResp, err := tsvc.GetUnits(ctx, unitAndPriceReq)
			if err != nil {
				return s.wrapTravolutionError(ctx, logger, unitPriceResp, err)
			}

			fmt.Println("========13========  ", unitPriceResp.Body)

			fmt.Println("############## type of unitPriceResp.Body #######  ", reflect.TypeOf(unitPriceResp.Body))

			// Marshal back to JSON
			unitAndPriceJson, err := json.Marshal(unitPriceResp.Body)
			if err != nil {
				return resp, fmt.Errorf("marshal error: %w", err)
			}

			// Step 1: Unmarshal into slice of maps for normalization
			var tmpList []map[string]interface{}
			if err := json.Unmarshal(unitAndPriceJson, &tmpList); err != nil {
				level.Error(logger).Log(
					"msg", "failed to unmarshal unitPrice JSON into temp list",
					"error", err,
				)
				return resp, fmt.Errorf("failed to unmarshal unitPrice JSON: %w", err)
			}

			// Step 2: Normalize UID for each element of unit and price
			for i, tmp := range tmpList {
				if v, ok := tmp["uid"]; ok {
					switch val := v.(type) {
					case string:
						tmpList[i]["uid"] = val
					case float64:
						tmpList[i]["uid"] = strconv.FormatInt(int64(val), 10)
					case json.Number:
						tmpList[i]["uid"] = val.String()
					default:
						return resp, fmt.Errorf("unsupported uid type: %T", v)
					}
				}
			}

			// Step 3: Marshal normalized list back to JSON
			normalizedunit, err := json.Marshal(tmpList)
			if err != nil {
				level.Error(logger).Log(
					"msg", "failed to re-marshal normalized unitPrice JSON",
					"error", err,
				)
				return resp, fmt.Errorf("failed to re-marshal normalized unitPrice JSON: %w", err)
			}

			// Step 4: Unmarshal normalized JSON into struct slice
			if err := json.Unmarshal(normalizedunit, &unitAndPriceList); err != nil {
				level.Error(logger).Log(
					"msg", "failed to unmarshal normalized unitPrice JSON into RawOptionUnitPrice",
					"error", err,
				)
				return resp, fmt.Errorf("failed to unmarshal normalized unitPrice JSON into RawOptionUnitPrice: %w", err)
			}

			optionList[idx].UnitsPrice = unitAndPriceList

			// ***************************  Fetch booking schedules  ****************************//
			if rawProductDetail.Type == "BK" {
				bookingScheduleReq := travolution.BookingScheduleReq{
					ProductUid: rawProductDetail.UID,
					OptionUid:  opt.UID,
				}
				fmt.Println("========14========")

				scheduleResp, err := tsvc.GetBookingSchedules(ctx, bookingScheduleReq)
				if err != nil {
					return s.wrapTravolutionError(ctx, logger, scheduleResp, err)
				}

				fmt.Println("========15========")

				bookingScheduleJson, err := DataConversionToJsonForListType(scheduleResp.Body, logger, "schedule")
				if err != nil {
					return resp, err
				}

				if err := json.Unmarshal(bookingScheduleJson, &ScheduleList); err != nil {
					level.Error(logger).Log(
						"msg", "failed to unmarshal into rawOptionUnitAndPrice",
						"error", err,
					)
					return resp, fmt.Errorf("failed to unmarshal into rawOptionUnitAndPrice: %w", err)
				}

				fmt.Println("========16========")
			}

			optionList[idx].BookingSchedules = ScheduleList

			// Fetch additional booking info if applicable
			if rawProductDetail.HasBookingAdditionalInfo == "Y" {
				additionalBookingInfoReq := travolution.BookingAdditionalInfoRequest{
					ProductUID: rawProductDetail.UID,
					OptionUID:  opt.UID,
				}

				fmt.Println("========17========")

				additionalInfoResp, err := tsvc.GetBookingAdditionalInfo(ctx, additionalBookingInfoReq)
				if err != nil {

					return s.wrapTravolutionError(ctx, logger, additionalInfoResp, err)
				}

				fmt.Println("========18========")

				if addInfo, ok := additionalInfoResp.Body.([]travolution.RawBookingAdditionalInfo); ok {
					optionList[idx].AdditionalBookingInfo = addInfo
				}
			}
		}

		fmt.Println("========19========")

		// Save enriched product
		rawProductDetail.Options = optionList
		insertedID, err := s.mongoRepository[config.Instance().MongoDBName].UpsertTravolutionProduct(ctx, rawProductDetail)
		if err != nil {
			level.Error(logger).Log("event", "insert_failed", "error", err)
			resp.Code = "500"
			resp.Body = err.Error()
			return resp, customError.NewError(ctx, "repository error", fmt.Sprint(customError.ErrDatabase.Error(), "InsertTravolutionProduct"), nil)
		}

		insertedIDs = append(insertedIDs, insertedID)
	}

	fmt.Println("========20========")

	return travolution.Response{Code: "200", Body: insertedIDs}, nil
}

// wrapTravolutionError standardizes error wrapping for Travolution calls
func (s *service) wrapTravolutionError(ctx context.Context, logger log.Logger, resp travolution.Response, err error) (travolution.Response, error) {
	level.Error(logger).Log("error", err)
	if resp.Code == "403" {
		return resp, customError.NewError(ctx, "leisure-api-0003", fmt.Sprint(customError.ErrForbiddenClient.Error(), "PostCreateAllProduct"), nil)
	}
	return travolution.Response{Code: "500", Body: "request to travolution client raised error"}, err
}

func DataConversionToJsonForListType(data interface{}, logger log.Logger, types string) (dataBytes []byte, err error) {
	var genericList []map[string]any

	fmt.Printf("############## type of data is ***************** %T", data)

	if genericSlice, ok := data.([]any); ok {
		genericList = make([]map[string]any, 0, len(genericSlice))

		for _, item := range genericSlice {
			if m, ok := item.(map[string]any); ok {
				genericList = append(genericList, m)
			} else {
				level.Error(logger).Log(
					"msg", fmt.Sprintf("skipping non-map item in %s slice", types),
					"type", fmt.Sprintf("%T", item),
				)
			}
		}

		/* level.Info(logger).Log(
			"msg", "unwrapped products from []interface{} -> []map[string]interface{}",
			"type", fmt.Sprintf("%T", data),
			"count", len(genericList),
		) */

	} else {
		level.Error(logger).Log(
			"msg", "resp.Body is not []interface{}",
			"actual_type", fmt.Sprintf("%T", data),
		)
		return nil, fmt.Errorf("unexpected resp.Body type: %T", data)
	}

	//  Only proceed if we actually have data
	if len(genericList) == 0 {
		level.Warn(logger).Log("msg", fmt.Sprintf(" no %s available in %s api response", types, types))
		return nil, fmt.Errorf("%s", fmt.Sprintf(" no %s available in %s api response", types, types))
	}

	// Marshal -> Unmarshal into RawProduct
	jsonbytes, err := json.Marshal(genericList)
	if err != nil {
		level.Error(logger).Log(
			"msg", "failed to marshal generic products list",
			"error", err,
		)
		return nil, fmt.Errorf("failed to marshal generic products list: %w", err)
	}

	return jsonbytes, nil

}

func DataConversionTOjson(respData interface{}, logger log.Logger, types string) (databyte []byte, err error) {
	rawProductDetailMap, ok := respData.(map[string]any)
	if !ok {
		level.Error(logger).Log("error ", "type conversion error ")
		return nil, fmt.Errorf("%s", fmt.Sprintf("invalid %s detail", types))
	}

	jsonBytes, err := json.Marshal(rawProductDetailMap)
	if err != nil {
		level.Error(logger).Log("msg", "failed to marshal map to JSON", "err", err)
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	return jsonBytes, nil
}
