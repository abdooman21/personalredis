package resp

import (
	"testing"
	"time"
)

func TestRespBulkString(t *testing.T) {
	start := time.Now()

	buf := []byte{'$', '5', '\r', '\n', 'h', 'e', 'l', 'l', 'o', '\r', '\n'}

	str, _, err := readBulkString(buf)

	if err {
		t.Fatal("failed")
	}

	if str != "hello" {
		t.Fatalf("expected hello, got %q", str)
	}

	duration := time.Since(start)

	t.Logf("received : %v with len(%d)", str, len(str))
	t.Logf("test duration: %v", duration)
}
