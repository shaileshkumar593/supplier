package trip_to_yanolja

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"swallow-supplier/caches/cache"
	"swallow-supplier/config"
	svc "swallow-supplier/iface"
	"swallow-supplier/request_response/trip"
	"swallow-supplier/request_response/yanolja"
	"swallow-supplier/utils"
	"swallow-supplier/utils/constant"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
)

func PreOrderMapper(ctx context.Context, mrepo svc.MongoRepository, logger log.Logger, tripReq trip.PreorderRequest) (req yanolja.WaitingForOrder, plu map[string]string, err error) {
	level.Info(logger).Log(
		"function name", "PreOrderMapper",
	)

	if tripReq.Contacts != nil {
		req.Customer.Name = strings.ToLower(tripReq.Contacts[0].Name)
		req.Customer.Tel = "010-0000-0000" //tripReq.Contacts[0].Mobile
		req.Customer.Email = strings.ToLower(tripReq.Contacts[0].Email)
	} else {
		req.Customer.Name = "GGT"
		req.Customer.Tel = "010-0000-0000"
		req.Customer.Email = "ggt@gmail.com"
	}

	itemLen := len(tripReq.Items)
	selectVariant := make([]yanolja.VariantInfo, 0)
	var channelCode string
	plus := make(map[string]string)

	// Initialize the cache layer
	cacheLayer, err := cache.New(config.Instance().CacheName)
	if err != nil {
		level.Error(logger).Log("msg", "Error initializing cache layer", "error", err)
		return req, nil, fmt.Errorf("error initializing cache layer: %w", err)
	}
	// plu as hash id
	if itemLen >= 1 {
		for _, val := range tripReq.Items {
			var variantInfo yanolja.VariantInfo

			plu, err := cacheLayer.Get(ctx, val.PLU)
			logger.Log("plu", plu)
			if err != nil && plu == "" {
				level.Error(logger).Log("error", "plu not in redis")
				return req, nil, fmt.Errorf("plu not in redis")

			} else if err == nil && plu == "" {
				level.Error(logger).Log("error", "Cache empty")
				return req, nil, fmt.Errorf("redis empty")
			}

			plus[val.PLU] = plu
			level.Info(logger).Log("plu ", plu)
			detail := strings.Split(plu, "|") // ProductID-ProductVersion-VariantID

			variantInfo.ProductID, err = strconv.ParseInt(detail[0], 10, 64)
			if err != nil {
				level.Error(logger).Log("error :", err)
				return req, nil, err
			}
			intConvrt, err := strconv.ParseInt(detail[1], 10, 64)
			if err != nil {
				level.Error(logger).Log("error :", err)
				return req, nil, err
			}

			variantInfo.ProductVersion = int32(intConvrt)

			variantInfo.VariantID, err = strconv.ParseInt(detail[2], 10, 64)
			if err != nil {
				level.Error(logger).Log("error :", err)
				return req, nil, err
			}

			if len(detail[3]) == 0 {
				variantInfo.Date = nil
			} else {
				variantInfo.Date = &detail[3]
			}

			// check for ""  as string value
			if len(detail[4]) == 0 {
				variantInfo.Time = nil
			} else {
				variantInfo.Time = &detail[4]
			}
			variantInfo.Quantity = int32(val.Quantity)
			variantInfo.Currency = "KRW" // need to be checked with subham
			variantInfo.PartnerSalePrice = float32(val.SalePrice)
			if val.DistributionChannel == "GGT_EVERLAND" {
				variantInfo.CostPrice = utils.RemoveMargineFromProductCostPrice(variantInfo.ProductID, val.Cost)

				/* for _, margine := range constant.MARGINEDETAIL {
					if margine.ProductId == variantInfo.ProductID {
						switch margine.MargineType {
						case "flat":
							variantInfo.CostPrice = float32(val.Cost + float64(margine.Value))
						case "percent":
							mergineval := float64(1 + margine.Value/100)
							variantInfo.CostPrice = float32(math.Floor(val.Cost) / mergineval)
						}

					}
				} */

			} else {
				margineVal := constant.MARKUPPERCENTAGE

				for _, margine := range constant.MARGINEDETAIL {
					if margine.ProductId == variantInfo.ProductID && margine.MargineType == constant.PERCENTAGE {
						margineVal = margine.Value
						break
					}
				}

				variantInfo.CostPrice = float32(math.Floor(val.Cost / float64(1+(margineVal/100)))) // now changes
			}

			selectVariant = append(selectVariant, variantInfo)
			req.TotalSelectedVariantsQuantity = req.TotalSelectedVariantsQuantity + int32(val.Quantity)
			channelCode = val.DistributionChannel
			if len(val.Passengers) > 0 {
				req.ActualCustomer.Name = strings.ToLower(val.Passengers[0].Name)
				req.ActualCustomer.Tel = val.Passengers[0].Mobile
				req.ActualCustomer.Email = ""
			} else {
				level.Error(logger).Log("error", "Passengers names are empty")
				return req, nil, fmt.Errorf("passengers names cannot be empty")
			}
		}
	}
	fmt.Println("39")
	req.PartnerOrderChannelCode = channelCode

	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // Create a new random generator
	//req.PartnerOrderID = strconv.Itoa(1000000 + rng.Intn(9000000))
	req.PartnerOrderID = tripReq.OtaOrderID
	req.PartnerOrderGroupID = strconv.Itoa(100000000 + rng.Intn(99999999999999)) //tripReq.SequenceID  not needed

	//req.PartnerOrderID = tripReq.OtaOrderID   // uncomment while doing QA with trip
	req.SelectVariants = selectVariant

	level.Info(logger).Log("selectVariant in mapper : ", selectVariant)

	level.Info(logger).Log("selectVariants in mapper : ", req.SelectVariants)
	level.Info(logger).Log("request from mapper : ", req)
	fmt.Println("********************* req ************************* ", req)

	return req, plus, nil
}
