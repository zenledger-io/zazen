package nullable

import "database/sql"

// String convenience initializer
func NewString(s string) String {
	return String{
		sql.NullString{
			String: s, Valid: true,
		},
	}
}

// String convenience initializer for invalid (nil)
func NullString() String {
	return String{
		sql.NullString{
			Valid: false,
		},
	}
}

// String is an alias for sql.NullString data type
type String struct {
	sql.NullString
}

// String Nullable conformance
func (n String) IsNull() bool {
	return !n.Valid
}

// MarshalJSON for String
func (n String) MarshalJSON() ([]byte, error) {
	return marshalJSON(n.String, n.Valid)
}

// UnmarshalJSON for String
func (n *String) UnmarshalJSON(b []byte) error {
	return unmarshalJSON(b, &n.String, &n.Valid)
}
