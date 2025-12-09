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

// Message represte the message content struct
type Message struct {
	ProgramCode          string   `json:"program_code"`
	BinSponsorCode       string   `json:"bin_sponsor_code"`
	HashID               string   `json:"hash_id"`
	FundingType          string   `json:"funding_type"`
	TransactionReference string   `json:"transaction_reference"`
	Action               string   `json:"action"`
	EventDestination     string   `json:"event_destination"`
	EventStatus          string   `json:"event_status"`
	EventRemarks         string   `json:"event_remarks"`
	Request              Request  `json:"request"`
	Response             Response `json:"response"`
}

type Request struct {
	IsMock   bool `json:"is_mock,omitempty"`
	Metadata struct {
		SourceAccount struct {
			Number string `json:"number" validate:"required"`
			Name   string `json:"name"`
			Bank   string `json:"bank,omitempty"`
		} `json:"source_account"`
		Config  interface{} `json:"config" validate:"required"`
		Network string      `json:"network"`
		Action  string      `json:"action"`
	} `json:"metadata"`
	Transactions []Transactions `json:"transactions" validate:"required"`
}

// Response ...
type Response struct {
	ID       string
	Metadata struct {
		SourceAccount struct {
			Number string `json:"number" validate:"required"`
			Name   string `json:"name"`
			Bank   string `json:"bank,omitempty"`
		} `json:"source_account"`
		Config      interface{} `json:"config" validate:"required"`
		Network     string      `json:"network"`
		FundingType string      `json:"funding_type"`
		Action      string      `json:"action"`
	} `json:"metadata"`
	Transactions ResponseTransactions `json:"transaction"`
}

// ResponseTransactions ...
type ResponseTransactions struct {
	Amount             string            `json:"amount"`
	TransactionType    string            `json:"transaction_type"`
	TransactionRef     string            `json:"transaction_ref"`
	Fee                []Fee             `json:"fee"`
	TransactionStatus  string            `json:"transaction_status"`
	TransactionRemarks string            `json:"transaction_remarks"`
	AdditionalDetails  []AdditionDetails `json:"additional_details"`
}

// AdditionDetails ...
type AdditionDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Transactions ...
type Transactions struct {
	Amount                  string `json:"amount" validate:"required"`
	TransactionType         string `json:"transaction_type"  validate:"required"`
	OriginalTransactionType string `json:"original_transaction_type"`
	TransactionRef          string `json:"transaction_ref"  validate:"required"`
	Fee                     []Fee  `json:"fee"`
}

// Fee ...
type Fee struct {
	Amount          string `json:"amount"`
	TransactionType string `json:"transaction_type"`
}

// ReversalRequest ...
type ReversalRequest struct {
	ProgramCode, FundingType, Action, HashID, TransactionType, ProcessType string
	RetryAttempt                                                           int64
	BinSponsorEnabled                                                      bool
}
