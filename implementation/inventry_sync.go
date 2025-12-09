package implementation

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"swallow-supplier/config"
	customError "swallow-supplier/error"
	"swallow-supplier/request_response/yanolja"
	tripservice "swallow-supplier/services/distributors/trip"
	yanoljasvc "swallow-supplier/services/suppliers/yanolja"
	"swallow-supplier/utils"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// InventrySync for syncing inventory to trip
func (s *service) InventorySync(ctx context.Context) (resp yanolja.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "InventrySync",
		"Request ID", requestID,
	)

	// Defer when panic
	defer func(context.Context) {
		r := recover()
		if r == nil {
			return
		}

		level.Error(logger).Log("error", "processing request went into panic mode", "panic", r)
		resp.Code = "500"
		err = fmt.Errorf("panic occurred: %v", r)

	}(ctx)

	products, err := s.mongoRepository[config.Instance().MongoDBName].FetchAllProductsWithinDateRange(ctx)
	if err != nil {
		resp.Code = "500"
		resp.Body = err.Error()
		return resp, err
	}

	var req yanolja.ProductInventory
	allProductInventory := make([]yanolja.Inventories, 0)

	// Start date: tomorrow
	startDate := time.Now().UTC().AddDate(0, 0, 1)
	req.InventoryDateStart = startDate.Format("2006-01-02")

	// End date: 90 days from tomorrow
	endDate := startDate.AddDate(0, 0, 90)
	req.InventoryDateEnd = endDate.Format("2006-01-02")

	var inventrysvc, _ = yanoljasvc.New(ctx)

	for _, product := range products {
		fmt.Println("Processing Product ID:", product.ProductID)
		req.ProductId = strconv.FormatInt(product.ProductID, 10)

		inventryresp, err := inventrysvc.GetProductsInventories(ctx, req)
		if err != nil {
			level.Error(logger).Log("error", "request to yanolja client raised error", err)
			return resp, err
		}
		level.Info(logger).Log("Info for inventryresp body ", inventryresp.Body)

		// Assuming inventryresp.Body contains a JSON object with an array `variantInventories`
		var inventories yanolja.Inventories

		// Try to unmarshal the response body into the expected structure
		bodyBytes, err := json.Marshal(inventryresp.Body)
		if err != nil {
			level.Error(logger).Log("error", "failed to marshal inventory response body", "productID", product.ProductID, "error", err)
			resp.Code = "500"
			resp.Body = fmt.Sprintln("Failed to marshal inventory response for Product ID:", product.ProductID)
			return resp, err
		}

		err = json.Unmarshal(bodyBytes, &inventories)
		if err != nil {
			level.Error(logger).Log("error", "failed to unmarshal inventory response body", "productID", product.ProductID, "error", err)
			resp.Code = "500"
			resp.Body = fmt.Sprintln("Failed to unmarshal inventory response for Product ID:", product.ProductID)
			return resp, err
		}

		// Validate the inventory
		if len(inventories.VariantInventories) == 0 {
			level.Error(logger).Log("error ", "inventory response is empty for product", "productID", product.ProductID)
			resp.Code = "500"
			resp.Body = fmt.Sprintln("Empty inventory for Product ID:", product.ProductID)
			return resp, customError.NewError(ctx, "leisure-api-00024", fmt.Sprintf("Empty inventory for Product ID %v :", product.ProductID), "InventorySync")
		}

		// Process the inventory and append to the final result
		for _, productoption := range product.ProductOptionGroups {
			for i, variant := range productoption.Variants {
				if i >= len(inventories.VariantInventories) {
					level.Error(logger).Log("error", "inventory length mismatch", "productID", product.ProductID)
					fmt.Println("Inventory mismatch for Product ID:", product.ProductID)
					break
				}

				// Update inventory fields based on product options
				inventories.VariantInventories[i].Price = variant.Price
				inventories.VariantInventories[i].IsSchedule = productoption.IsSchedule
				inventories.VariantInventories[i].IsRound = productoption.IsRound
				inventories.VariantInventories[i].ProductVersion = product.ProductVersion
			}
		}

		/* fmt.Println("$$$$$$$$$$$$$$$$$ Appending inventory for Product ID:", product.ProductID, "Inventory Length:", len(inventories.VariantInventories))
		fmt.Println("$$$$$$$$$$$$$$$$ Index of product is ------->", i) */
		allProductInventory = append(allProductInventory, inventories)

	}

	//fmt.Println("Final allProductInventory Length:", len(allProductInventory))

	if len(allProductInventory) == 0 {
		level.Error(logger).Log("error ", "no data to be sent to trip")
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("empty request for trip, %v", err), "InventorySync")
	}

	reqToTrip := yanolja.InventoryToTrip{
		Message: "Inventory",
		Data:    allProductInventory,
	}
	fmt.Println("***************** request to trip ****************** :", reqToTrip)
	var tripsvc, _ = tripservice.New(ctx)
	resp, err = tripsvc.InventorySyncToTrip(ctx, reqToTrip)
	if err != nil {
		resp.Code = "500"
		resp.Body = err.Error()
		return resp, err
	}

	return resp, nil
}

// removeGarbageTimestamps cleans up invalid timestamps from raw JSON string
func removeGarbageTimestamps(rawData string) string {
	// Regex to match garbage timestamp patterns
	garbageTimestampPattern := regexp.MustCompile(`,?\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+Z `)
	// Replace garbage timestamps with empty string
	cleanedData := garbageTimestampPattern.ReplaceAllString(rawData, "")
	return cleanedData
}
