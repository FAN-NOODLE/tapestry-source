package user

import "testing"

func TestNormalizeAppCode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		valid    bool
	}{
		{name: "legacy caller", expected: AppCodeWebAdmin, valid: true},
		{name: "web admin", input: AppCodeWebAdmin, expected: AppCodeWebAdmin, valid: true},
		{name: "store operation admin", input: "  " + AppCodeStoreOperationAdmin + " ", expected: AppCodeStoreOperationAdmin, valid: true},
		{name: "unknown app", input: "unknown", valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, valid := NormalizeAppCode(tt.input)
			if actual != tt.expected || valid != tt.valid {
				t.Fatalf("NormalizeAppCode(%q) = (%q, %v), want (%q, %v)", tt.input, actual, valid, tt.expected, tt.valid)
			}
		})
	}
}

func TestClientForAppCode(t *testing.T) {
	tests := []struct {
		name     string
		appCode  string
		expected string
	}{
		{name: "old admin", appCode: AppCodeWebAdmin, expected: "console"},
		{name: "store operation admin", appCode: AppCodeStoreOperationAdmin, expected: "store_operation_admin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual := ClientForAppCode(tt.appCode); actual != tt.expected {
				t.Fatalf("ClientForAppCode(%q) = %q, want %q", tt.appCode, actual, tt.expected)
			}
		})
	}
}
