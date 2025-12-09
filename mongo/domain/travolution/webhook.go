package travolution

// WebhookEventData represents the order data in the webhook
type WebhookEventData struct {
	OrderNumber     string `bson:"orderNumber" json:"orderNumber" validate:"required"`
	ReferenceNumber string `bson:"referenceNumber" json:"referenceNumber" validate:"required"`
	DateAt          string `bson:"dateAt" json:"dateAt" validate:"required"` // ISO8601 string
}

// Webhook represents the webhook stored in MongoDB
type Webhook struct {
	ID        string           `bson:"_id,omitempty" json:"id"`
	EventType string           `bson:"eventType" json:"eventType" validate:"required"`
	Data      WebhookEventData `bson:"data" json:"data" validate:"required,dive"`
	CreatedAt string           `bson:"createdAt" json:"createdAt" validate:"required"`
	UpdatedAt string           `bson:"updatedAt" json:"updatedAt"`
}
