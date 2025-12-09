package yanolja

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Category represents the MongoDB document structure for the category.
type Category struct {
	Id                 string `bson:"_id,omitempty" json:"id,omitempty"`
	SupplierGGTChannel string `bson:"supplierGGTChannel" json:"supplierGGTChannel"`                                              // MongoDB document ID
	CategoryID         int64  `bson:"categoryId" json:"categoryId" validate:"required"`                                          // Category ID
	CategoryCode       string `bson:"categoryCode" json:"categoryCode" validate:"required"`                                      // Category Code
	CategoryLevel      int32  `bson:"categoryLevel" json:"categoryLevel" validate:"required"`                                    // Category Level
	CategoryName       string `bson:"categoryName" json:"categoryName" validate:"required"`                                      // Category Name
	CategoryStatusCode string `bson:"categoryStatusCode" json:"categoryStatusCode" validate:"required,oneof=PREPARE USE UNUSED"` // Category Status Code
	ImageURL           string `bson:"imageUrl,omitempty" json:"imageUrl,omitempty"`                                              // Image URL (optional)
}

// NewCategory creates a new Category with validation.
func NewCategory(categoryID int64, categoryCode string, categoryLevel int32, categoryName string, categoryStatusCode string, imageUrl string) (*Category, error) {
	category := &Category{
		Id:                 primitive.NewObjectID().Hex(),
		SupplierGGTChannel: "YANOLJA-GGT_TRIP",
		CategoryID:         categoryID,
		CategoryCode:       categoryCode,
		CategoryLevel:      categoryLevel,
		CategoryName:       categoryName,
		CategoryStatusCode: categoryStatusCode,
		ImageURL:           imageUrl,
	}

	// Validate the struct
	validate := validator.New()
	err := validate.Struct(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}
