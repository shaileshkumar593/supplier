package heartbeat

// Response struct for heartbeat implementation
type Response struct {
	Timestamp      string `json:"timestamp"`
	AppPostgresSQL string `json:"app_postgres_sql"`
}

// Response struct for heartbeat implementation
type MongoResponse struct {
	Timestamp  string `json:"timestamp"`
	AppMongoDb string `json:"app_postgres_sql"`
}
