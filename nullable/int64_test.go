package nullable

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInt64_UnmarshalJSON(t *testing.T) {
	type TestStruct struct {
		Key Int64
	}

	tcs := []struct {
		JSON  string
		Valid bool
		Error bool
	}{
		{
			JSON:  "{\"key\": null}",
			Valid: false,
			Error: false,
		},
		{
			JSON:  "{\"key\": 0}",
			Valid: true,
			Error: false,
		},
		{
			JSON:  "{\"key\": \"0\"}",
			Valid: false,
			Error: true,
		},
		{
			JSON:  "{\"key\": false}",
			Valid: false,
			Error: true,
		},
		{
			JSON:  "{\"key\": 0.5}",
			Valid: false,
			Error: true,
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("test case %v", i), func(t *testing.T) {
			var testStruct TestStruct

			err := json.Unmarshal([]byte(tc.JSON), &testStruct)

			if tc.Error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.Valid, testStruct.Key.Valid)
		})
	}
}
