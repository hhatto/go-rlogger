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

func appendPacket(buf *bytes.Buffer, now uint32, msg []byte) {
	scratch := make([]byte, 8)
	msgLen := uint32(len(msg))

	binary.BigEndian.PutUint32(scratch, now)
	binary.BigEndian.PutUint32(scratch[4:], msgLen)
	buf.Write(scratch) // never retruns nil
	buf.Write(msg)
}

func (r *RLogger) Write(tag, msg []byte) (int, error) {
	return write(r.conn, tag, msg)
}

func write(w io.Writer, tag, msg []byte) (int, error) {
	now := uint32(time.Now().Unix())
	headerLen := HEADER_SIZE + int32(len(tag))
	buf := new(bytes.Buffer)
	msgBuf := new(bytes.Buffer)

	offset := 0
	for len(msg[offset:]) > 0 {
		ret := bytes.IndexByte(msg[offset:], '\n')
		if ret == -1 {
			appendPacket(msgBuf, now, msg[offset:])
			break
		}
		appendPacket(msgBuf, now, msg[offset:offset+ret])
		offset += ret + 1
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

	msgBuf.WriteTo(buf)
	nw, err := buf.WriteTo(w)
	return int(nw), err
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
