package repository

import (
	"context"
	"fmt"
	"time"

	"swallow-supplier/request_response/heartbeat"

	"github.com/go-kit/kit/log/level"
)

const (
	// StatusAlive ...
	StatusAlive = "Connection alive"

	// StatusUnavailable ...
	StatusUnavailable = "Connection unavailable"
)

// GetHeartBeat
// Pings database availability
func (r *mongoRepository) GetHeartBeatFromMongo(ctx context.Context) (res heartbeat.MongoResponse, err error) {
	level.Info(r.logger).Log("repo-method", "GetHeartBeatFromMongo")

	// Trace the database segments to the enabled metric
	/*m, _ := metric.New(config.Instance().MetricMethod)
	m.TraceDatabase(
		ctx,
		"heartbeat",
		"PING",
		"postgres_ping",
		nil,
	)*/
	res.Timestamp = fmt.Sprintf("%s", time.Now().UTC().Format(time.RFC3339))
	res.AppMongoDb = StatusAlive
	if err := r.db.Client().Ping(ctx, nil); err != nil {
		res.AppMongoDb = StatusUnavailable + ": " + err.Error()
		return res, err
	}
	return res, nil
}
