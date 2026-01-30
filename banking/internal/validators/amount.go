package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

const maxAmountValue = 1000000000

// RegisterCustomValidators registers custom validator tags.
func RegisterCustomValidators(v *validator.Validate) error {
	if err := v.RegisterValidation("amount", ValidateAmount); err != nil {
		return err
	}
	return nil
}

// ValidateAmount validates decimal amount format, precision, and range.
func ValidateAmount(fl validator.FieldLevel) bool {
	amountStr := fl.Field().String()
	if amountStr == "" {
		return false
	}

	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		return false
	}

	if amount.Exponent() < -2 {
		return false
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return false
	}

	maxAmount := decimal.NewFromInt(maxAmountValue)
	if amount.GreaterThan(maxAmount) {
		return false
	}

	return true
}
