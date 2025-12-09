package yanolja

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SubRegion represents a sub-region object within a region.
type SubRegion struct {
	RegionID          int64  `bson:"regionId" json:"regionId" validate:"required"`
	RegionCode        string `bson:"regionCode" json:"regionCode" validate:"required"`
	ParentRegionID    int64  `bson:"parentRegionId,omitempty" json:"parentRegionId,omitempty"`
	RegionName        string `bson:"regionName" json:"regionName" validate:"required"`
	RegionDescription string `bson:"regionDescription,omitempty" json:"regionDescription,omitempty"`
	RegionLevel       int32  `bson:"regionLevel" json:"regionLevel" validate:"required,gt=0,lte=5"` // Ensure region level is from 1 to 5
	IsUsed            bool   `bson:"isUsed" json:"isUsed" validate:"required"`
}

// Region represents a regional object in the database.
type Region struct {
	Id                string      `bson:"_id,omitempty" json:"id,omitempty"`
	RegionID          int64       `bson:"regionId" json:"regionId" validate:"required"`
	RegionCode        string      `bson:"regionCode" json:"regionCode" validate:"required"`
	ParentRegionID    int64       `bson:"parentRegionId,omitempty" json:"parentRegionId,omitempty"`
	RegionName        string      `bson:"regionName" json:"regionName" validate:"required"`
	RegionDescription string      `bson:"regionDescription,omitempty" json:"regionDescription,omitempty"`
	RegionLevel       int32       `bson:"regionLevel" json:"regionLevel" validate:"required,gt=0,lte=5"`
	IsUsed            bool        `bson:"isUsed" json:"isUsed" validate:"required"`
	SubRegions        []SubRegion `bson:"subRegions,omitempty" json:"subRegions" validate:"dive"`
}

// NewRegion creates a new Region with validation.
func NewRegion(regionID int64, regionCode string, parentRegionID int64, regionName string, regionDescription string, regionLevel int32, isUsed bool, subRegions []SubRegion) (*Region, error) {
	region := &Region{
		Id:                primitive.NewObjectID().Hex(),
		RegionID:          regionID,
		RegionCode:        regionCode,
		ParentRegionID:    parentRegionID,
		RegionName:        regionName,
		RegionDescription: regionDescription,
		RegionLevel:       regionLevel,
		IsUsed:            isUsed,
		SubRegions:        subRegions,
	}

	// Validate the struct
	validate := validator.New()
	err := validate.Struct(region)
	if err != nil {
		return nil, err
	}

	return region, nil
}
