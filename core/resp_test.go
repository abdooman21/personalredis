package resp

import (
	"reflect"
	"testing"
	"time"
)

func TestRespBulkString(t *testing.T) {
	start := time.Now()

	buf := []byte{'$', '5', '\r', '\n', 'h', 'e', 'l', 'l', 'o', '\r', '\n'}

	//str, l, err := readBulkString(buf)
	str, l, err := DecodeOne(buf)
	if err {
		t.Fatal("failed")
	}

	if str != "hello" {
		t.Fatalf("expected hello, got %q", str)
	}
	if l != 11 {
		t.Fatal("len wrong")
	}
	duration := time.Since(start)

	t.Logf("received : %v with len(%d)", str, l)
	t.Logf("test duration: %v", duration)
}

func TestDecodeOne(t *testing.T) {
	tests := []struct {
		name    string
		buf     []byte
		want    any
		wantLen int
		wantErr bool
	}{
		// simple string
		{
			name:    "simple string",
			buf:     []byte("+OK\r\n"),
			want:    "OK",
			wantLen: 5,
			wantErr: false,
		},

		// error
		{
			name:    "error string",
			buf:     []byte("-ERR unknown command\r\n"),
			want:    "ERR unknown command",
			wantLen: 22,
			wantErr: false,
		},
		// integer
		{
			name:    "integer",
			buf:     []byte(":1000\r\n"),
			want:    int64(1000),
			wantLen: 7,
			wantErr: false,
		},
		{
			name:    "posIntegerWOptional+",
			buf:     []byte(":+1000\r\n"),
			want:    int64(1000),
			wantLen: 8,
			wantErr: false,
		},
		{
			name:    "negInteger",
			buf:     []byte(":-1000\r\n"),
			want:    int64(-1000),
			wantLen: 8,
			wantErr: false,
		},

		// bulk string
		{
			name:    "bulk string",
			buf:     []byte("$5\r\nhello\r\n"),
			want:    "hello",
			wantLen: 11,
			wantErr: false,
		},
		{
			name:    "empty bulk string",
			buf:     []byte("$0\r\n\r\n"),
			want:    "",
			wantLen: 6,
			wantErr: false,
		},

		// array
		{
			name: "array of bulk strings",
			buf:  []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"),
			want: []any{"hello", "world"},
			// *2\r\n = 4
			// each bulk string = 11
			wantLen: 26,
			wantErr: false,
		},

		// invald
		{
			name:    "unknown type",
			buf:     []byte("?abc\r\n"),
			want:    nil,
			wantLen: 0,
			wantErr: true,
		},
		{
			name:    "empty buffer",
			buf:     []byte{},
			want:    nil,
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotLen, err := DecodeOne(tt.buf)

			if err != tt.wantErr {
				t.Fatalf("expected error=%v, got=%v", tt.wantErr, err)
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("expected %#v, got %#v", tt.want, got)
			}

			if gotLen != tt.wantLen {
				t.Fatalf("expected len=%d, got=%d", tt.wantLen, gotLen)
			}
		})
	}
}
