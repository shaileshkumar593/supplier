package yanolja

import (
	"swallow-supplier/mongo/domain/yanolja"
)

type Upsert_Product struct {
	TraceID     string          `json:"traceID"`
	Body        yanolja.Product `json:"body"`
	Page        NumberOfPage    `json:"page"`
	Collection  bool            `json:"collection"`
	ContentType string          `json:"contentType"`
}
