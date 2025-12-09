package trip

// PreorderRequest represents the pre-order creation request body
type PreorderRequest struct {
	SequenceID string    `json:"sequenceId" binding:"required"`
	OtaOrderID string    `json:"otaOrderId" binding:"required"`
	Contacts   []Contact `json:"contacts"`
	Items      []Item    `json:"items" binding:"required"`
}
type Contact struct {
	Name             string `json:"name"`
	Mobile           string `json:"mobile"`
	IntlCode         string `json:"intlCode"`
	OptionalMobile   string `json:"optionalMobile"`
	OptionalIntlCode string `json:"optionalIntlCode"`
	Email            string `json:"email"`
}

type Item struct {
	PLU                    string          `json:"PLU" binding:"required"`
	Locale                 string          `json:"locale" binding:"required"`
	DistributionChannel    string          `json:"distributionChannel"`
	UseStartDate           string          `json:"useStartDate" binding:"required"`
	UseEndDate             string          `json:"useEndDate" binding:"required"`
	Remark                 string          `json:"remark"`
	Price                  float64         `json:"price" binding:"required"`
	PriceCurrency          string          `json:"priceCurrency" binding:"required"`
	SalePrice              float64         `json:"salePrice" binding:"required"`
	SalePriceCurrency      string          `json:"salePriceCurrency" binding:"required"`
	Cost                   float64         `json:"cost"`
	CostCurrency           string          `json:"costCurrency"`
	SuggestedPrice         float64         `json:"suggestedPrice"`
	SuggestedPriceCurrency string          `json:"suggestedPriceCurrency"`
	Quantity               int             `json:"quantity" binding:"required"`
	Passengers             []Passenger     `json:"passengers"`
	Adjunctions            []Adjunction    `json:"adjunctions"`
	Deposit                Deposit         `json:"deposit"`
	ExpressDelivery        ExpressDelivery `json:"expressDelivery"`
}

type Passenger struct {
	PassengerID      string  `json:"passengerId" binding:"required"`
	Name             string  `json:"name"`
	FirstName        string  `json:"firstName"`
	LastName         string  `json:"lastName"`
	Mobile           string  `json:"mobile"`
	IntlCode         string  `json:"intlCode"`
	CardType         string  `json:"cardType"`
	CardNo           string  `json:"cardNo"`
	BirthDate        string  `json:"birthDate"`
	AgeType          string  `json:"ageType"`
	Gender           string  `json:"gender"`
	NationalityCode  string  `json:"nationalityCode"`
	NationalityName  string  `json:"nationalityName"`
	CardIssueCountry string  `json:"cardIssueCountry"`
	CardIssuePlace   string  `json:"cardIssuePlace"`
	CardIssueDate    string  `json:"cardIssueDate"`
	CardValidDate    string  `json:"cardValidDate"`
	CardIssueNumber  string  `json:"cardIssueNumber"`
	BirthPlace       string  `json:"birthPlace"`
	Height           float64 `json:"height"`
	Weight           float64 `json:"weight"`
	MyopiaDegreeL    float64 `json:"myopiaDegreeL"`
	MyopiaDegreeR    float64 `json:"myopiaDegreeR"`
	ShoeSize         float64 `json:"shoeSize"`
}

type Adjunction struct {
	Name        string `json:"name" binding:"required"`
	NameCode    string `json:"nameCode"`
	Content     string `json:"content" binding:"required"`
	ContentCode string `json:"contentCode"`
}

type Deposit struct {
	Type           int     `json:"type" binding:"required"`
	Amount         float64 `json:"amount" binding:"required"`
	AmountCurrency string  `json:"amountCurrency" binding:"required"`
}

type ExpressDelivery struct {
	Type     string `json:"type" binding:"required"`
	Name     string `json:"name"`
	Mobile   string `json:"mobile"`
	IntlCode string `json:"intlCode"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	District string `json:"district"`
	Address  string `json:"address"`
}

// PreOrderResponse represents the pre-order creation response body
type PreOrderResponse struct {
	Header ResponseHeader `json:"header" binding:"required"`
	Body   ResponseBody   `json:"body"`
}

// ResponseHeader represents the response header
type ResponseHeader struct {
	ResultCode    string `json:"resultCode"`
	ResultMessage string `json:"resultMessage"`
}

// ResponseBody represents the response body
type ResponseBody struct {
	OtaOrderID      string         `json:"otaOrderId"`
	SupplierOrderID string         `json:"supplierOrderId"`
	Items           []ResponseItem `json:"items"`
}

// ResponseItem represents a response item object
type ResponseItem struct {
	PLU        string      `json:"PLU"`
	Inventorys []Inventory `json:"inventorys"`
}

// Inventory represents an inventory object
type Inventory struct {
	UseDate  string `json:"useDate"`
	Quantity int    `json:"quantity"`
}
