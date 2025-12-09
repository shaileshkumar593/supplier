package yanolja

type ProductsById struct {
	ProductId string `json:"productId" binding:"required" validate:"required"`
}

type VariantById struct {
	VariantIdId string `json:"productid" binding:"required" validate:"required"`
}

type ProductInventory struct {
	ProductId          string `json:"product_id" binding:"required" validate:"required"`
	InventoryDateStart string `json:"inventoryDateStart"`
	InventoryDateEnd   string `json:"inventoryDateEnd"`
}

type VariantInventory struct {
	VariantId string `json:"product_id" binding:"required" validate:"required"`
	Date      string `json:"date"`
	Time      string `json:"time"`
}

type GetProduct struct {
	ProductId int64 `json:"productId" binding:"required" validate:"required"`
}

type AllProduct struct {
	PageNumber        int32  `json:"pageNumber" binding:"required" validate:"required"`
	PageSize          int32  `json:"pageSize" binding:"required" validate:"required"`
	ProductStatusCode string `json:"productStatusCode, omitempty" validate:"oneof=WAITING_FOR_SALE IN_SALE SOLD_OUT END_OF_SALE"`
}
