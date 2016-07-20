package rlogger

import (
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"
)

var SOCKET_PATH string

func TestSimpleWrite(t *testing.T) {
	tag := []byte("this.is.tag")
	msg := []byte("hogehoge1")
	t.Logf("path=%v\n", SOCKET_PATH)
	r := NewRLogger(SOCKET_PATH)
	defer r.Close()
	ret, err := r.Write(tag, msg)
	if err != nil {
		t.Errorf("write error: err=%v", err)
	}

	if ret != 40 {
		t.Errorf("invalid msg len. len=%v", ret)
	}
}

func TestTwoLineWrite(t *testing.T) {
	tag := []byte("this.is.tag")
	msg := []byte("hogehoge2\nhogehoge3")
	r := NewRLogger(SOCKET_PATH)
	defer r.Close()
	ret, err := r.Write(tag, msg)
	if err != nil {
		t.Errorf("write error: err=%v", err)
	}

	if ret != 57 {
		t.Errorf("invalid msg len. len=%v", ret)
	}
}

func TestMain(m *testing.M) {
	dir, err := ioutil.TempDir("", "gorlogger")
	if err != nil {
		panic(err)
	}

	SOCKET_PATH = filepath.Join(dir, "rloggerd.dummy.sock")
	unixListener, err := net.ListenUnix("unix", &net.UnixAddr{Name: SOCKET_PATH, Net: "unix"})
	if err != nil {
		os.RemoveAll(dir)
		panic(err)
	}

	go func() {
		for {
			unixListener.Accept()
		}
	}()

	ret := m.Run()
	os.RemoveAll(dir)

	os.Exit(ret)
}
