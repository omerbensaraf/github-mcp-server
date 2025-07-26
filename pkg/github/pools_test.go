package github

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONBufferPool(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "simple string",
			input:    "hello",
			expected: `"hello"`,
		},
		{
			name:     "simple map",
			input:    map[string]int{"count": 42},
			expected: `{"count":42}`,
		},
		{
			name: "complex struct",
			input: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{
				Name: "test",
				Age:  25,
			},
			expected: `{"name":"test","age":25}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Test MarshalJSONWithPool
			result, err := MarshalJSONWithPool(tc.input)
			require.NoError(t, err)
			assert.JSONEq(t, tc.expected, string(result))

			// Verify result is the same as standard json.Marshal
			standardResult, err := json.Marshal(tc.input)
			require.NoError(t, err)
			assert.JSONEq(t, string(standardResult), string(result))
		})
	}
}

func TestStringSlicePool(t *testing.T) {
	// Get a slice from the pool
	slice1 := GetStringSlice()
	assert.NotNil(t, slice1)
	assert.Equal(t, 0, len(slice1))

	// Add some items
	slice1 = append(slice1, "item1", "item2", "item3")
	assert.Equal(t, 3, len(slice1))

	// Return to pool
	PutStringSlice(slice1)

	// Get another slice
	slice2 := GetStringSlice()
	assert.NotNil(t, slice2)
	assert.Equal(t, 0, len(slice2)) // Should be reset
}

func BenchmarkJSONMarshal(b *testing.B) {
	testData := map[string]interface{}{
		"name":  "test",
		"value": 123,
		"items": []string{"a", "b", "c"},
	}

	b.Run("standard", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := json.Marshal(testData)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("pooled", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := MarshalJSONWithPool(testData)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
