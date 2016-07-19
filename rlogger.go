package rlogger

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

const CHUNK_SIZE int32 = 8388608
const HEADER_SIZE int32 = 12
const HEADER_VERSION int32 = 1
const HEADER_PSH int32 = 1

type RLogger struct {
	sock *net.UnixAddr
	conn *net.UnixConn
}

func NewRLogger(socketPath string) *RLogger {
	r := &RLogger{}
	r.sock, r.conn = createUnixDomainSocket("unix", socketPath)
	return r
}

func getRLoggerPacket(now int32, msg []byte) (buf *bytes.Buffer, err error) {
	buf = new(bytes.Buffer)
	msgLen := int32(len(msg))

	if err = binary.Write(buf, binary.BigEndian, int32(now)); err != nil {
		log.Printf("binary.write(time) error. err=%v\n", err)
		return
	}

	if err = binary.Write(buf, binary.BigEndian, int32(msgLen)); err != nil {
		log.Printf("binary.write(msgLen) error. err=%v\n", err)
		return
	}

	if err = binary.Write(buf, binary.BigEndian, msg); err != nil {
		log.Printf("binary.write(msg) error. err=%v\n", err)
		return
	}

	return buf, nil
}

func (r *RLogger) Write(tag, msg []byte) (int, error) {
	return write(r.conn, tag, msg)
}

func write(w io.Writer, tag, msg []byte) (int, error) {
	now := int32(time.Now().Unix())
	headerLen := HEADER_SIZE + int32(len(tag))
	buf := new(bytes.Buffer)
	msgBuf := new(bytes.Buffer)

	for _, line := range bytes.Split(msg, []byte("\n")) {
		pkt, _ := getRLoggerPacket(now, line)
		msgBuf.Write(pkt.Bytes())
	}

	msgLen := int32(msgBuf.Len())

	if err := binary.Write(buf, binary.BigEndian, int8(HEADER_VERSION)); err != nil {
		log.Println("binary.write(HEADER_VERSION) error")
		return 0, err
	}
	if err := binary.Write(buf, binary.BigEndian, int8(HEADER_PSH)); err != nil {
		log.Println("binary.write(HEADER_PSH) error")
		return 0, err
	}
	if err := binary.Write(buf, binary.BigEndian, int16(headerLen)); err != nil {
		log.Println("binary.write(headerLen) error")
		return 0, err
	}
	if err := binary.Write(buf, binary.BigEndian, int32(0)); err != nil {
		log.Printf("binary.write(offset) error. err=%v\n", err)
		return 0, err
	}
	if err := binary.Write(buf, binary.BigEndian, int32(headerLen+msgLen)); err != nil {
		log.Println("binary.write(headerLen+msgLen) error")
		return 0, err
	}

	if err := binary.Write(buf, binary.BigEndian, tag); err != nil {
		log.Printf("binary.write(tag) error. err=%v\n", err)
		return 0, err
	}

	buf.Write(msgBuf.Bytes())
	return w.Write(buf.Bytes())
}

func (r *RLogger) Close() {
	r.conn.Close()
}

func createUnixDomainSocket(t, addr string) (*net.UnixAddr, *net.UnixConn) {
	unixAddr, err := net.ResolveUnixAddr(t, addr)
	if err != nil {
		panic(fmt.Sprintf("create unix domain socket error. err=%v", err))
	}

	conn, err := net.DialUnix("unix", nil, unixAddr)
	if err != nil {
		panic(fmt.Sprintf("connect unix domain socket error. err=%v", err))
	}

	return unixAddr, conn
}
