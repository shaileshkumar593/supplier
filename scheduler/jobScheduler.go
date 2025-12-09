package scheduler

import (
	"context"
	"fmt"

	//"fmt"
	svc "swallow-supplier/iface"
	"swallow-supplier/scheduler/cronjob"

	//"swallow-supplier/scheduler/cronjob"

	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"

	"github.com/robfig/cron/v3"
)

// Jobs will handle background jobs
func Jobs(ctx context.Context, logger log.Logger, s svc.Service, mrepo svc.MongoRepository) (err error) {
	level.Info(logger).Log(
		"method", "Jobs ",
		"starttime", time.Now(),
	)

	if ctx == nil {
		ctx = context.Background()
	}

	// Setup cron scheduler
	job := cron.New()

	// Schedule the first execution after 1 minute
	/* time.AfterFunc(2*time.Minute, func() {

		resp, err := s.InsertAllProduct(ctx)
		level.Info(logger).Log("resp ", resp)
		if err == nil {
			level.Info(logger).Log("msg", "Initial execution of MonitorProductUpdates started")

			err := cronjob.MonitorProductUpdates(ctx, mrepo, logger)
			if err != nil {
				level.Error(logger).Log("error", "Initial execution of MonitorProductUpdates failed", "err", err)
			} else {
				level.Info(logger).Log("msg", "Initial execution of MonitorProductUpdates completed successfully")
			}
		}

	}) */

	// Setup cron scheduler for every 5 second for product update
	_, err = job.AddFunc("@every 5s", func() {
		//level.Info(logger).Log("msg", "Scheduled execution of MonitorProductUpdates every 5 second")

		cronjob.MonitorProductUpdates(ctx, mrepo, logger)
	})
	if err != nil {
		level.Error(logger).Log("error", "Failed to schedule MonitorProductUpdates job", "err", err)
		return err
	}

	// Add jobs to the scheduler to run every 60 seconds  for redis update for plu
	_, err = job.AddFunc("@every 60s", func() {
		cronjob.SchedulePluUpsertToRedis(ctx, mrepo, logger)

	})
	if err != nil {
		level.Error(logger).Log("error", "Failed to schedule SyncItemIdDetail job", "err", err)
		return err
	}

	// Add jobs to the scheduler to run every 2 seconds  always need in order creation
	_, err = job.AddFunc("@every 2s", func() {
		cronjob.SyncItemIdDetail(ctx, logger, mrepo)

	})
	if err != nil {
		level.Error(logger).Log("error", "Failed to schedule SyncItemIdDetail job", "err", err)
		return err
	}

	// Add jobs to the scheduler to run every 2 seconds  always need in order creation
	// below code needs to be removed
	_, err = job.AddFunc("@every 2s", func() {
		cronjob.SyncTripRequestToredis(ctx, mrepo, logger)

	})
	if err != nil {
		level.Error(logger).Log("error", fmt.Sprintf("Failed to schedule SyncTripRequestToredis job %w", err))
		return err
	}

	// Start the cron scheduler
	level.Info(logger).Log("msg", "Starting the cron scheduler...")
	job.Start()

	// Keep the scheduler running by waiting for context cancellation
	<-ctx.Done()
	job.Stop() // Stop the scheduler gracefully on context cancellation
	return nil
}
