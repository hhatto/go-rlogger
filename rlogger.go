package rlogger

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

const CHUNK_SIZE int32 = 8388608
const HEADER_SIZE int32 = 12
const HEADER_VERSION int32 = 1
const HEADER_PSH int32 = 1

var Buffs sync.Pool

type RLogger struct {
	sock *net.UnixAddr
	conn *net.UnixConn
}

func init() {
	Buffs = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
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
	buf := Buffs.Get().(*bytes.Buffer)
	msgBuf := Buffs.Get().(*bytes.Buffer)
	defer func() {
		Buffs.Put(msgBuf)
		Buffs.Put(buf)
	}()

	offset := 0
	for {
		ret := bytes.IndexByte(msg[offset:], '\n')
		if ret == -1 {
			appendPacket(msgBuf, now, msg[offset:])
			break
		}
		appendPacket(msgBuf, now, msg[offset:offset+ret])
		offset += ret + 1
	}

	pktLen := int32(msgBuf.Len()) + headerLen

	scratch := make([]byte, HEADER_SIZE)
	scratch[0] = uint8(HEADER_VERSION)
	scratch[1] = uint8(HEADER_PSH)
	binary.BigEndian.PutUint16(scratch[2:], uint16(headerLen))
	binary.BigEndian.PutUint32(scratch[4:], 0)
	binary.BigEndian.PutUint32(scratch[8:], uint32(pktLen))
	scratch = append(scratch, tag...)

	buf.Write(scratch)
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
