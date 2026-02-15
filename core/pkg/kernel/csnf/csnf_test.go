package csnf

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/unicode/norm"
)

func TestCanonicalize(t *testing.T) {
	tests := []struct {
		name      string
		inputJSON string
		wantErr   bool
		check     func(*testing.T, interface{})
	}{
		{
			name:      "Float Rejection",
			inputJSON: `{"value": 3.14}`,
			wantErr:   true,
		},
		{
			name:      "Integer Accepted",
			inputJSON: `{"value": 42}`,
			wantErr:   false,
			check: func(t *testing.T, v interface{}) {
				m := v.(map[string]interface{})
				assert.Equal(t, int64(42), m["value"])
			},
		},
		{
			name: "NFC Normalization",
			// "café" in NFD form: 'c', 'a', 'f', 'e', 0x301 (combining acute)
			inputJSON: `{"name": "cafe\u0301"}`,
			wantErr:   false,
			check: func(t *testing.T, v interface{}) {
				m := v.(map[string]interface{})
				name := m["name"].(string)
				// Expect NFC: 'c', 'a', 'f', 0xE9 (é)
				expected := norm.NFC.String("cafe\u0301")
				assert.Equal(t, expected, name)
			},
		},
		{
			name:      "Null Stripping",
			inputJSON: `{"a": 1, "b": null}`,
			wantErr:   false,
			check: func(t *testing.T, v interface{}) {
				m := v.(map[string]interface{})
				_, hasB := m["b"]
				assert.False(t, hasB, "null field 'b' should be stripped")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var input interface{}
			err := json.Unmarshal([]byte(tt.inputJSON), &input)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			got, err := Canonicalize(input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Canonicalize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}
