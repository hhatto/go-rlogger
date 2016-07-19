package rlogger

import (
	"bytes"
	"testing"
)

func BenchmarkWrite(b *testing.B) {
	tag := []byte("hello.world")
	msg := []byte(`hello
hello
world
world
sun
mercury
venus
earth
mars
jupiter
saturn
`)
	buf := new(bytes.Buffer)
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		write(buf, tag, msg)
		buf.Reset()
	}
}
