package trip

import (
	"swallow-supplier/mongo/domain/trip"
)

type ProductContentSync struct {
	Message string               `json:"message"`
	Data    []trip.ProuctContent `json:"data"`
}

type PackageContentSync struct {
	Message string                `json:"message"`
	Data    []trip.PackageContent `json:"data"`
}

type TripMessageForSync struct {
	Message string              `json:"message" validate:"required,oneof='product' 'package'"`
	Status  []ContentSyncStatus `json:"status"`
}

type ContentSyncStatus struct {
	SupplierProductId string `json:"supplierProductId"`
	SyncStatus        string `json:"syncStatus" validate:"required,oneof='Sync' 'NotSync'"`
}
