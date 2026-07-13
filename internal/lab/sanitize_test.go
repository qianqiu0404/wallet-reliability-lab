package lab

import (
	"strings"
	"testing"
)

func TestSanitizeLogRedactsSensitiveValues(t *testing.T) {
	result := SanitizeLog("token=abc api_key=123 rpc_url=https://secret private_key=0xdead request_id=safe")
	for _, secret := range []string{"abc", "123", "https://secret", "0xdead"} {
		if strings.Contains(result, secret) {
			t.Fatalf("sanitized log still contains %q: %s", secret, result)
		}
	}
	if !strings.Contains(result, "request_id=safe") {
		t.Fatalf("safe field was unexpectedly removed: %s", result)
	}
}
