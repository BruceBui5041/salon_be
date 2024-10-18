package customtypes

import (
	"encoding/json"
	"strings"

	"github.com/shopspring/decimal"
)

type NullDecimalString decimal.NullDecimal

func (d *NullDecimalString) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = strings.Trim(s, "\"") // Remove quotes if present

	if s == "" || s == "null" {
		*d = NullDecimalString(decimal.NullDecimal{Valid: false})
		return nil
	}

	dec, err := decimal.NewFromString(s)
	if err != nil {
		return err
	}

	*d = NullDecimalString(decimal.NullDecimal{Decimal: dec, Valid: true})
	return nil
}

func (d NullDecimalString) MarshalJSON() ([]byte, error) {
	if !decimal.NullDecimal(d).Valid {
		return []byte("null"), nil
	}
	return json.Marshal(decimal.NullDecimal(d).Decimal.String())
}

func (d NullDecimalString) GetDecimal() decimal.NullDecimal {
	return decimal.NullDecimal(d)
}

type DecimalString decimal.Decimal

func (d *DecimalString) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = strings.Trim(s, "\"") // Remove quotes if present

	dec, err := decimal.NewFromString(s)
	if err != nil {
		return err
	}

	*d = DecimalString(dec)
	return nil
}

func (d DecimalString) MarshalJSON() ([]byte, error) {
	return json.Marshal(decimal.Decimal(d).String())
}

func (d DecimalString) GetDecimal() decimal.Decimal {
	return decimal.Decimal(d)
}
