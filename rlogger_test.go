package rlogger

import (
	"bytes"
	"testing"
)

func TestInnerWriteOneLine(t *testing.T) {
	buf := new(bytes.Buffer)
	write(buf, []byte("tag"), []byte("msg"))

	b := buf.Bytes()

	if buf.Len() != 26 {
		t.Errorf("invalid packet length. len=%v, pkt=%v", buf.Len(), b)
	}

	// check version
	if b[0] != 1 {
		t.Errorf("invalid version: %v", b[0])
	}

	// check header size
	if b[3] != 15 {
		t.Errorf("invalid header size: %v", b[3])
	}

	// check tag string
	if string(b[12:15]) != "tag" {
		t.Errorf("invalid tag string: %s", b[12:15])
	}

	// check msg size
	if b[22] != 3 {
		t.Errorf("invalid message size: %v", b[22])
	}

	// check msg string
	if string(b[23:]) != "msg" {
		t.Errorf("invalid message string: %s", b[23:])
	}
}
