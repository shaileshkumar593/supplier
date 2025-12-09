package yanolja

type CategoryType struct {
	CategoryId         int64  `json:"categoryId"`
	CategoryCode       string `json:"categoryCode"`
	CategoryLevel      int32  `json:"categoryLevel"`
	CategoryName       string `json:"categoryName"`
	CategoryStatusCode string `json:"categoryStatusCode"`
	ImageUrl           string `json:"imageurl"`
}
