package cronjob

import (
	"context"
	"fmt"

	svc "swallow-supplier/iface"
	"swallow-supplier/mongo/domain/yanolja"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"

	"go.mongodb.org/mongo-driver/bson"
)

// MonitorProductUpdates polls the product collection for changes and updates/inserts into ProductView.
func MonitorProductUpdates(ctx context.Context, mrepo svc.MongoRepository, logger log.Logger) (err error) {
	db := mrepo.GetMongoDb(ctx)
	productCollection := db.Collection("products")

	// --- Product_View
	cursor, err := productCollection.Find(ctx, bson.M{"viewScheduleStatus": false})
	if err != nil {
		level.Error(logger).Log("msg", "Error querying product for viewSchedule", "err", err)
		return fmt.Errorf("failed to execute find query for viewSchedule: %w", err)
	}
	defer cursor.Close(ctx)

	var products []yanolja.Product
	err = cursor.All(ctx, &products)
	if err != nil {
		level.Error(logger).Log("msg", "Error decoding product documents", "err", err)
		// defer will close cursor
	}
	// Close current cursor explicitly (optional; defer will handle it)
	cursor.Close(ctx)

	if len(products) != 0 {
		level.Info(logger).Log("info", "UpdateOrInsertProductView proceeds")
		if err = mrepo.UpdateOrInsertProductView(ctx, products); err != nil {
			level.Error(logger).Log("msg", "Error updating ProductView", "err", err)
		} else {
			level.Info(logger).Log("msg", "Successfully processed productview")
		}
	}

	// --- productImage
	cursor, err = productCollection.Find(ctx, bson.M{"imageScheduleStatus": false})
	if err != nil {
		level.Error(logger).Log("msg", "Error querying product for imageSchedule", "err", err)
		return fmt.Errorf("failed to execute find query for imageSchedule: %w", err)
	}
	defer cursor.Close(ctx)

	products = nil
	err = cursor.All(ctx, &products)
	if err != nil {
		level.Error(logger).Log("msg", "Error decoding product documents", "err", err)
		cursor.Close(ctx)
	}
	cursor.Close(ctx)

	if len(products) != 0 {
		level.Info(logger).Log("info", "ProductImagesForProcessing proceeds")
		if err = ProductImagesForProcessing(ctx, logger, mrepo, products); err != nil {
			level.Error(logger).Log("error", "Failed to schedule ProductImagesForProcessing job", "err", err)
			return err
		}

		// --- contentSync to trip (requires image + productview updated first)
		cursor, err = productCollection.Find(ctx, bson.M{"contentScheduleStatus": false})
		if err != nil {
			level.Error(logger).Log("msg", "Error querying product for contentSchedule", "err", err)
			return fmt.Errorf("failed to execute find query for contentSchedule: %w", err)
		}
		defer cursor.Close(ctx)

		products = nil
		err = cursor.All(ctx, &products)
		if err != nil {
			level.Error(logger).Log("msg", "Error decoding product documents", "err", err)
			cursor.Close(ctx)
			return err // <--- FIX: explicit return of the outer err
		}
		cursor.Close(ctx)

		if len(products) != 0 {
			if err = ProductSyncToTrip(ctx, logger, mrepo); err != nil {
				level.Error(logger).Log("error", "Failed to schedule job ProductSyncToTrip", "err", err)
				return err
			}
		}
	}

	return nil
}
