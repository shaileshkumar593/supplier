package mmbintransaction

// ReversalResponse ...
type ReversalResponse struct {
	PushToQueue    bool
	CheckStatus    int64
	IniateReversal bool
	ReversalLevel  int64
	NoAction       bool
	Type           string
}

// AuthorizeResponse ...
type AuthorizeResponse struct {
	CheckStatus     int64
	IniateAuthorize bool
	NoAction        bool
}

// ConfigRequest ...
type ConfigRequest struct {
	ID             string `json:"id"`
	ProductID      string `json:"product_id"`
	BinSponsorID   string `json:"bin_sponsor_id"`
	CustomerTypeID string `json:"customer_type_id"`
	ResourceTypeID string `json:"resource_type_id"`
	Version        string `json:"version"`
	Key            string `json:"key"`
	IsActive       string `json:"is_active"`
	ProductCode    string `json:"product_code" validate:"required"`
	ResourceType   string `json:"resource_type" validate:"required"`
	CustomerType   string `json:"customer_type" validate:"required"`
}

// Configs ...
type Configs struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

// ConfigResponse ...
type ConfigResponse struct {
	ID             string    `json:"id"`
	ProductCode    string    `json:"product_code"`
	CustomerType   string    `json:"customer_type"`
	BinSponsorCode string    `json:"bin_sponsor_code"`
	ResourceType   string    `json:"resource_type"`
	Version        string    `json:"version"`
	Records        []Configs `json:"records"`
}
