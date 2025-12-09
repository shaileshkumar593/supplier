// Package validator ...
package validator

import (
	"context"
	"os"
	"strings"
	"swallow-supplier/utils/constant"

	"github.com/go-kit/log"
)

// ValidateProductVariantStatusCode ...
func ValidateProductVariantStatusCode(ctx context.Context, statusCode string) (inArry bool) {
	var (
		logger log.Logger
	)

	logger = log.NewJSONLogger(os.Stdout)
	logger = log.With(logger,
		"method", "ValidateProductVariantStatusCode",
	)
	key := strings.ToUpper(statusCode)

	for _, code := range constant.PRODUCTVARIANTSTATUSCODE {
		if strings.EqualFold(key, code) {
			inArry = true
			return
		}
	}

	return false
}
