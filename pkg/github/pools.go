// Package github provides performance optimization pools for frequently allocated objects.
// This helps reduce garbage collection pressure in high-throughput scenarios.
package github

import (
	"bytes"
	"encoding/json"
	"sync"
)

// JSONBufferPool provides a pool of bytes.Buffer objects for JSON operations
// to reduce garbage collection pressure.
var JSONBufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

// GetJSONBuffer gets a buffer from the pool for JSON operations.
func GetJSONBuffer() *bytes.Buffer {
	return JSONBufferPool.Get().(*bytes.Buffer)
}

// PutJSONBuffer returns a buffer to the pool after use.
// The buffer is reset before being returned.
func PutJSONBuffer(buf *bytes.Buffer) {
	buf.Reset()
	JSONBufferPool.Put(buf)
}

// MarshalJSONWithPool marshals an object to JSON using a pooled buffer
// to reduce allocations.
func MarshalJSONWithPool(v interface{}) ([]byte, error) {
	buf := GetJSONBuffer()
	defer PutJSONBuffer(buf)

	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}

	// Remove the trailing newline that Encode adds
	data := buf.Bytes()
	if len(data) > 0 && data[len(data)-1] == '\n' {
		data = data[:len(data)-1]
	}

	// Make a copy since we're returning the buffer to the pool
	result := make([]byte, len(data))
	copy(result, data)

	return result, nil
}

// StringSlicePool provides a pool of string slices for frequent allocations.
var StringSlicePool = sync.Pool{
	New: func() interface{} {
		slice := make([]string, 0, 10) // Pre-allocate capacity of 10
		return &slice
	},
}

// GetStringSlice gets a string slice from the pool.
func GetStringSlice() []string {
	return *StringSlicePool.Get().(*[]string)
}

// PutStringSlice returns a string slice to the pool after use.
// The slice is reset before being returned.
func PutStringSlice(slice []string) {
	// Clear the slice but keep the underlying array
	slice = slice[:0]
	StringSlicePool.Put(&slice)
}
