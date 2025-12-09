package gorm

import (
	"errors"
	"strings"

	"swallow-supplier/common/go-tools/array"
)

// operator
const (
	OperatorBetween          = "BETWEEN"
	OperatorLike             = "LIKE"
	OperatorEqual            = "="
	OperatorNotEqual         = "!="
	OperatorGreaterThan      = ">"
	OperatorGreaterThanEqual = ">="
	OperatorLessThanEqual    = "<="
	OperatorLessThan         = "<"
	OperatorIn               = "IN"
)

// Condition ...
type Condition struct {
	Field    string
	Operator string
	Value    interface{}
}

// Conditions list of Document
type Conditions []Condition

// Validate Operator
func (me *Condition) Validate() error {
	if me.Operator == "" {
		return errors.New(SQLInvalidOperator)
	}

	if ok, _ := array.InArray(strings.ToUpper(me.Operator), []string{
		OperatorBetween,
		OperatorLike,
		OperatorIn,
		OperatorEqual,
		OperatorNotEqual,
		OperatorGreaterThan,
		OperatorGreaterThanEqual,
		OperatorLessThanEqual,
		OperatorLessThan}); !ok {
		return errors.New(SQLInvalidOperator)
	}

	return nil
}
