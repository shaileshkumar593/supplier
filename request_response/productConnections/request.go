package productConnections

// Request struct for config on boarding implementation
type Request struct {
	ProductCode    string `json:"product_code" validate:"required"`
	DSN            string `json:"dsn" validate:"required"`
	SlaveDSN       string `json:"slave_dsn" validate:"required"`
	ConnectionType string `json:"connection_type" validate:"required"`
	IsActive       string `json:"is_active" validate:"required"`
	CreatedBy      string `json:"created_by"`
}
