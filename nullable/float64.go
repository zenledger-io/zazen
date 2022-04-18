package nullable

import "database/sql"

// Float64 convenience initializer
func NewFloat64(f float64) Float64 {
	return Float64{
		sql.NullFloat64{
			Float64: f, Valid: true,
		},
	}
}

// Float64 convenience initializer for invalid (nil)
func NullFloat64() Float64 {
	return Float64{
		sql.NullFloat64{
			Valid: false,
		},
	}
}

// Float64 is an alias for sql.NullFloat64 data type
type Float64 struct {
	sql.NullFloat64
}

// Float64 Nullable conformance
func (n Float64) IsNull() bool {
	return !n.Valid
}

// MarshalJSON for Float64
func (n Float64) MarshalJSON() ([]byte, error) {
	return marshalJSON(n.Float64, n.Valid)
}

// UnmarshalJSON for Float64
func (n *Float64) UnmarshalJSON(b []byte) error {
	return unmarshalJSON(b, &n.Float64, &n.Valid)
}
