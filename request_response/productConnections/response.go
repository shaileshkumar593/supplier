package productConnections

// Request struct for config on boarding implementation
type Response struct {
	ProductCode    string `json:"product_code"`
	DSN            string `json:"dsn"`
	SlaveDSN       string `json:"slave_dsn"`
	ConnectionType string `json:"connection_type"`
	IsActive       string `json:"is_active"`
}
