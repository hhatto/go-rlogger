package rlogger

import (
	"testing"
)

const SOCKET_PATH string = "/usr/local/var/run/rlogd/rloggerd.sock"

func TestSimpleWrite(t *testing.T) {
	r := NewRLogger([]byte("example.acc.go"), SOCKET_PATH)
	defer r.Close()
	ret, err := r.Write([]byte("hogehoge1"))
	if err != nil {
		t.Errorf("write error: err=%v", err)
	}

	if ret != 43 {
		t.Errorf("invalid msg len. len=%v", ret)
	}
}

func TestTwoLineWrite(t *testing.T) {
	r := NewRLogger([]byte("example.acc.go"), SOCKET_PATH)
	defer r.Close()
	ret, err := r.Write([]byte("hogehoge2\nhogehoge3"))
	if err != nil {
		t.Errorf("write error: err=%v", err)
	}

	if ret != 60 {
		t.Errorf("invalid msg len. len=%v", ret)
	}
}
