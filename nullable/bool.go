package nullable

import "database/sql"

// Bool convenience initializer
func NewBool(b bool) Bool {
	return Bool{
		sql.NullBool{
			Bool: b, Valid: true,
		},
	}
}

// Bool convenience initializer for invalid (nil)
func NullBool() Bool {
	return Bool{
		sql.NullBool{
			Valid: false,
		},
	}
}

// Bool is an alias for sql.NullBool data type
type Bool struct {
	sql.NullBool
}

// Bool Nullable conformance
func (n Bool) IsNull() bool {
	return !n.Valid
}

// MarshalJSON for Bool
func (n Bool) MarshalJSON() ([]byte, error) {
	return marshalJSON(n.Bool, n.Valid)
}

// UnmarshalJSON for Bool
func (n *Bool) UnmarshalJSON(b []byte) error {
	return unmarshalJSON(b, &n.Bool, &n.Valid)
}
