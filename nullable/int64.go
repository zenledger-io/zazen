package nullable

import "database/sql"

// Int64 convenience initializer
func NewInt64(i int64) Int64 {
	return Int64{
		sql.NullInt64{
			Int64: i, Valid: true,
		},
	}
}

// Int64 convenience initializer for invalid (nil)
func NullInt64() Int64 {
	return Int64{
		sql.NullInt64{
			Valid: false,
		},
	}
}

// Int64 is an alias for sql.NullInt64 data type
type Int64 struct {
	sql.NullInt64
}

// Int64 Nullable conformance
func (n Int64) IsNull() bool {
	return !n.Valid
}

// MarshalJSON for Int64
func (n Int64) MarshalJSON() ([]byte, error) {
	return marshalJSON(n.Int64, n.Valid)
}

// UnmarshalJSON for Int64
func (n *Int64) UnmarshalJSON(b []byte) error {
	return unmarshalJSON(b, &n.Int64, &n.Valid)
}
