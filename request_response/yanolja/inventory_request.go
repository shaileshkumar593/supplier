package yanolja

import "swallow-supplier/mongo/domain/yanolja"

type Inventories struct {
	VariantInventories []InventoryDetailOfVariant `json:"variantInventories"`
}
type InventoryDetailOfVariant struct {
	ProductID                  int64                      `json:"productId"`
	VariantID                  int64                      `json:"variantId"`
	ProductVersion             int32                      `json:"productVersion"`
	InventoryTypeCode          string                     `json:"inventoryTypeCode"`
	QuantityPerPerson          int                        `json:"quantityPerPerson"`
	QuantityPerPurchase        int                        `json:"quantityPerPurchase"`
	Quantity                   int                        `json:"quantity"`
	IsSchedule                 bool                       `json:"isSchedule,omitempty"` // Is Schedule
	IsRound                    bool                       `json:"isRound,omitempty"`
	Price                      yanolja.VariantPrice       `json:"price"`
	VariantScheduleInventories []VariantScheduleInventory `json:"variantScheduleInventories,omitempty"`
}

type VariantScheduleInventory struct {
	Date                    string                  `json:"date"`
	Quantity                int                     `json:"quantity"`
	VariantRoundInventories []VariantRoundInventory `json:"variantRoundInventories,omitempty"`
}
type VariantRoundInventory struct {
	Time     string `json:"time"`
	Quantity int    `json:"quantity"`
}

type InventoryToTrip struct {
	Message string        `json:"message"`
	Data    []Inventories `json:"data"`
}
