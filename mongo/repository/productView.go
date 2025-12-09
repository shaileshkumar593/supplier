package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"swallow-supplier/caches/cache"
	"swallow-supplier/config"
	"swallow-supplier/mongo/domain/yanolja"
	"swallow-supplier/utils"
	"swallow-supplier/utils/constant"

	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-kit/log"

	"github.com/go-kit/log/level"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var c = config.Instance()

// UpdateOrInsertProductView updates or inserts a product document into the "productview" and "PLUDetail" collection.
func (r *mongoRepository) UpdateOrInsertProductView(ctx context.Context, products []yanolja.Product) (err error) {
	level.Info(r.logger).Log(
		"method", "UpdateOrInsertProductView",
	)

	productSuccess := make(map[int64]string)

	plucollection := r.db.Collection("plus")
	collection := r.db.Collection("productview")
	productCollection := r.db.Collection("products")
	var npluhash = make(map[string]string)

	for _, product := range products {
		value, _ := productSuccess[product.ProductID]
		if value == "done" {
			continue
		}
		var transformedVariants []yanolja.Variant
		var schedulesSet = make(map[string]struct{})
		var roundsSet = make(map[string]struct{})
		var pluhash = make(map[string]string)

		// Generate PLUs for syncing product
		plus := GeneratePLUs(product)

		for _, val := range plus {
			for key, val := range val.PluHash {
				pluhash[key] = val
				npluhash[key] = val
			}
		}

		for _, optionGroup := range product.ProductOptionGroups {
			// Add variants to transformed variants
			transformedVariants = append(transformedVariants, optionGroup.Variants...)

			// Collect unique schedules and rounds
			for _, option := range optionGroup.ProductOptions {
				for _, item := range option.ProductOptionItems {
					for _, schedule := range item.Schedules {
						schedulesSet[schedule] = struct{}{}
					}
					for _, round := range item.Rounds {
						roundsSet[round] = struct{}{}
					}
				}
			}
		}

		// Convert unique sets to slices
		var schedules []string
		var rounds []string
		for schedule := range schedulesSet {
			schedules = append(schedules, schedule)
		}
		for round := range roundsSet {
			rounds = append(rounds, round)
		}

		// Prepare the `ProductView` document for upsert
		productView := yanolja.ProductView{
			ProductID:                   product.ProductID,
			SupplierName:                constant.SUPPLIERYANOLJA,
			ProductName:                 product.ProductName,
			ProductVersion:              product.ProductVersion,
			ProductInfo:                 product.ProductInfo,
			Price:                       product.Price,
			ProductStatusCode:           product.ProductStatusCode,
			ProductTypeCode:             product.ProductTypeCode,
			SalePeriod:                  product.SalePeriod,
			ProductBriefIntroduction:    product.ProductBriefIntroduction,
			Variants:                    transformedVariants,
			Categories:                  product.Categories,
			Regions:                     product.Regions,
			Images:                      product.Images,
			IsIntegratedVoucher:         product.IsIntegratedVoucher,
			IsRefundableAfterExpiration: product.IsRefundableAfterExpiration,
			IsUsed:                      product.IsUsed,
			PluDetails:                  plus,
			CreatedAt:                   product.CreatedAt,
			UpdatedAt:                   product.UpdatedAt,
		}

		err = PlusCreation(ctx, r.logger, pluhash, plucollection)
		if err != nil {
			level.Error(r.logger).Log("error ", "plu update error ")
			return err
		}
		// Define the upsert filter based on `ProductID`
		filter := bson.M{"productId": product.ProductID}
		update := bson.M{"$set": productView}

		opts := options.Update().SetUpsert(true)
		_, err := collection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return fmt.Errorf("failed to upsert product view: %w", err)
		}

		productSuccess[product.ProductID] = "done"

		filter = bson.M{"productId": product.ProductID, "viewScheduleStatus": false}
		update = bson.M{"$set": bson.M{"productId": product.ProductID, "viewScheduleStatus": true}}
		_, err = productCollection.UpdateOne(
			ctx,
			filter,
			update,
		)
		if err != nil {
			level.Error(r.logger).Log("error ", fmt.Sprintf("failed to update viewScheduleStatus for productId %d", product.ProductID))
			return fmt.Errorf("failed to update viewScheduleStatus for product: %w", err)
		}

		level.Info(r.logger).Log("info", fmt.Sprintf("Product %d successfully updated/inserted into productview to collections.", product.ProductID))
	}

	return nil
}

// PlusCreation create plu from productView
func PlusCreation(ctx context.Context, logger log.Logger, pluhash map[string]string, plucollection *mongo.Collection) error {
	level.Info(logger).Log("Info", "Updating PLU detail in plu_detail collection")

	documentId := "6510b8b8c9e77a3f4d6aab09"
	objectID, err := primitive.ObjectIDFromHex(documentId)
	if err != nil {
		return fmt.Errorf("invalid ObjectID: %w", err)
	}

	filter := bson.M{"_id": objectID}

	update := bson.M{
		"$set": bson.M{
			"supplierName":  constant.SUPPLIERYANOLJA,
			"lastUpdatedAt": time.Now().UTC().Format(time.RFC3339), // Always update
		},
		"$push": bson.M{
			"pluhash": bson.M{
				"$each": convertMapToBSONArray(pluhash), // Convert and append new key-value pairs
			},
		},
	}

	opts := options.Update().SetUpsert(true) // Ensure the document exists

	_, err = plucollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		level.Error(logger).Log("error", "Failed to upsert PLU detail", "err", err)
		return fmt.Errorf("failed to upsert document: %w", err)
	}

	level.Info(logger).Log("Info", "PLU detail successfully upserted")
	return nil
}

// Converts map[string]string to an array of BSON key-value pairs
func convertMapToBSONArray(data map[string]string) []bson.M {
	result := make([]bson.M, 0, len(data))
	for key, value := range data {
		result = append(result, bson.M{key: value})
	}
	return result
}

// Helper function to transform Variants to VariantPLU
func transformedVariantsToVariantPLU(variants []yanolja.Variant) []yanolja.VariantPLU {
	var variantPLUs []yanolja.VariantPLU
	for _, variant := range variants {
		variantPLUs = append(variantPLUs, yanolja.VariantPLU{
			VariantId:   variant.VariantID,
			VariantName: variant.VariantName,
		})
	}
	return variantPLUs
}

// UpdateOrInsertRedisRecord saves or updates a struct in Redis as a serialized JSON string.
func UpdateOrInsertRedisRecord(ctx context.Context, logger log.Logger, plus map[string]string) error {
	if ctx == nil {
		err := errors.New("context is nil")
		level.Error(logger).Log("msg", "Invalid context provided", "error", err)
		return err
	}

	// Initialize the cache layer
	cacheLayer, err := cache.New(config.Instance().CacheName)
	if err != nil {
		level.Error(logger).Log("msg", "Error initializing cache layer", "error", err)
		return fmt.Errorf("error initializing cache layer: %w", err)
	}

	for key, plu := range plus {
		// Insert or update the record in Redis
		if _, err := cacheLayer.SetNX(ctx, key, plu, 1*time.Minute); err != nil {
			level.Error(logger).Log(
				"msg", "Failed to set key in Redis",
				"key", key,
				"value", plu, // Log the serialized JSON for debugging
				"error", err,
			)
			return fmt.Errorf("failed to set key %s in Redis: %w", key, err)
		}
	}

	level.Info(logger).Log("msg", "Successfully updated/inserted key in Redis")
	return nil
}

// GetRecentProducts retrieves products updated/inserted in the last 12 hours in UTC.
func (r *mongoRepository) GetRecentProducts(ctx context.Context) ([]yanolja.ProductView, error) {
	// Calculate 12 hours ago in UTC.
	utcNow := time.Now().UTC()
	twelveHoursAgo := utcNow.Add(-1200 * time.Hour)

	// Define the filter.
	filter := bson.M{
		"updatedAt": bson.M{"$gte": twelveHoursAgo.Format(time.RFC3339)},
	}

	// Define the options (e.g., sorting by updatedAt descending).
	options := options.Find()
	options.SetSort(bson.D{{"updatedAt", -1}})

	collection := r.db.Collection("productview")
	// Execute the query.
	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode the results.
	var products []yanolja.ProductView
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

// generate plu format data "10012383|80|10268832|2025-02-06|10:00"
func GeneratePLUs(product yanolja.Product) []yanolja.PluDetail {
	pluDetails := make([]yanolja.PluDetail, 0)
	existingPLUs := make(map[string]struct{}) // Track unique PLUs

	for _, optionGroup := range product.ProductOptionGroups {
		for _, option := range optionGroup.ProductOptions {
			var plus []string
			var pluHash = make(map[string]string)

			for _, item := range option.ProductOptionItems {
				for _, variant := range optionGroup.Variants {
					var generatedPLUs []string

					switch option.ProductOptionTypeCode {
					case "ROUND":
						generatedPLUs = generatePLUsRecursive(product.ProductID, int(product.ProductVersion), variant.VariantID, item.Schedules, item.Rounds)
					case "SCHEDULE":
						generatedPLUs = generatePLUsRecursive(product.ProductID, int(product.ProductVersion), variant.VariantID, item.Schedules, nil)
					case "LIST":
						generatedPLUs = []string{fmt.Sprintf("%d|%d|%d|%s|%s", product.ProductID, product.ProductVersion, variant.VariantID, "", "")}
					}

					// Deduplicate PLUs before storing
					for _, plu := range generatedPLUs {
						if _, exists := existingPLUs[plu]; !exists {
							existingPLUs[plu] = struct{}{}
							hash := utils.GenerateDeterministicCode(plu, 15)
							plus = append(plus, plu)
							pluHash[hash] = plu
						}
					}
				}
			}

			if len(plus) > 0 {
				pluDetails = append(pluDetails, yanolja.PluDetail{
					ProductOptionTypeCode: option.ProductOptionTypeCode,
					PLU:                   plus,
					PluHash:               pluHash,
				})
			}
		}
	}
	return pluDetails
}

// Recursive helper function to generate PLUs
func generatePLUsRecursive(productID int64, productVersion int, variantID int64, schedules []string, rounds []string) []string {
	var plus []string

	if len(schedules) > 0 {
		for _, schedule := range schedules {
			if len(rounds) > 0 {
				for _, round := range rounds {
					plu := fmt.Sprintf("%d|%d|%d|%s|%s", productID, productVersion, variantID, schedule, round)
					plus = append(plus, plu)
				}
			} else {
				plu := fmt.Sprintf("%d|%d|%d|%s|%s", productID, productVersion, variantID, schedule, "")
				plus = append(plus, plu)
			}
		}
	} else {
		plu := fmt.Sprintf("%d|%d|%d|%s|%s", productID, productVersion, variantID, "", "")
		plus = append(plus, plu)
	}
	return plus
}

// GetPLUDetails fetches the PLU details for a given productId
func (r *mongoRepository) GetPLUDetails(ctx context.Context, productID int64) (productview yanolja.ProductView, err error) {
	level.Info(r.logger).Log(
		"method", "GetPLUDetails",
	)

	collection := r.db.Collection("productview")

	filter := bson.M{
		"productId": productID,
	}

	err = collection.FindOne(ctx, filter).Decode(&productview)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return productview, mongo.ErrNoDocuments
		}
		return productview, fmt.Errorf("failed to fetch productview for productId %d: %w", productID, err)
	}

	return productview, nil
}

// GetAllProducts  get all products
func (r *mongoRepository) GetAllProductViews(ctx context.Context) (products []yanolja.ProductView, err error) {
	level.Info(r.logger).Log("methodName ", "GetAllProducts")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filter := bson.M{}
	// Read environment (e.g., "dev", "local", "prod")
	if c.AppEnv == "dev" || c.AppEnv == "LOCAL" {
		filter = bson.M{
			"productId": bson.M{"$in": constant.ProductForGlobaltix},
		}
	}

	collection := r.db.Collection("productview")

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		level.Error(r.logger).Log("msg", "Failed to fetch products", "err", err)
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &products); err != nil {
		level.Error(r.logger).Log("msg", "Failed to decode products", "err", err)
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	// Apply costPrice modification
	for i := range products {
		for j := range products[i].Variants {
			cp := products[i].Variants[j].Price.CostPrice
			products[i].Variants[j].Price.CostPrice = math.Ceil(cp + cp*0.03)
		}
	}

	return products, nil
}

// GetProductByProductId fetches a single product by productId
func (r *mongoRepository) GetProductViewByProductId(ctx context.Context, productId int64) (product yanolja.ProductView, err error) {
	level.Info(r.logger).Log("methodName", "GetProductByProductId", "productID", productId)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if c.AppEnv == "dev" || c.AppEnv == "LOCAL" {
		var flag = false

		for _, pid := range constant.ProductForGlobaltix {
			if pid == productId {
				flag = true
				break
			}
		}

		if flag == false {
			level.Error(r.logger).Log("msg", "No product found", "productID", productId)
			return product, mongo.ErrNoDocuments
		}
	}

	filter := bson.M{"productId": productId}

	collection := r.db.Collection("productview")

	// Use FindOne since we're fetching by a unique productId
	err = collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			level.Error(r.logger).Log("msg", "No product found", "productID", productId)
			return product, mongo.ErrNoDocuments
		}
		level.Error(r.logger).Log("msg", "Failed to fetch product", "err", err)
		return product, fmt.Errorf("failed to fetch product: %w", err)
	}

	if len(product.Variants) > 0 {

		for j, variant := range product.Variants {

			cp := variant.Price.CostPrice
			newCost := math.Ceil(cp + (cp * 0.03))
			level.Info(r.logger).Log(
				"productId", product.ProductID,
				"variantId", product.Variants[j].VariantID,
				"oldCostPrice", cp,
				"newCostPrice", newCost,
			)
			product.Variants[j].Price.CostPrice = newCost
		}
	}
	return product, nil
}

// below code create issue that each time create new plu
/* // Generate SHA-256 hash of the PLU
func generateHash(plu string) string {
	hash := sha256.Sum256([]byte(plu))
	return hex.EncodeToString(hash[:])
}

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func generateUniqueAlphanumeric(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[randomIndex.Int64()]
	}
	return string(result), nil
}
*/
