package main

import (
	"io"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
		}

		go Dial(c)
	}
}

func Dial(source net.Conn) {
	target, err := net.Dial("tcp", "127.0.0.1:5432")
	if err != nil {
		log.Printf("dial err: %v", err)
		return
	}
	log.Println("dial succeed")

	go Stream(source, target, "user")
	go Stream(target, source, "service")
}

func Stream(source, target net.Conn, from string) {
	var err error
	b := Buffer{}
	for {
		b.Reset()
		_, err = b.ReadFrom(source)
		if err != nil {
			log.Printf("read from %s err: %v", from, err)
			return
		}
		log.Println(from + ":" + b.String())
		_, err = b.WriteTo(target)
		if err != nil {
			log.Printf("write from %s err: %v", from, err)
			return
		}
	}
}

// Buffer is some a copy from bytes.Buffer with different
// behaviour
type Buffer struct {
	buf      []byte
	lastRead int
}

func (b *Buffer) Reset() {
	b.buf = make([]byte, 1024)
	b.lastRead = 0
}

func (b *Buffer) ReadFrom(r io.Reader) (n int, err error) {
	m, e := r.Read(b.buf)
	if m < 0 {
		panic("negative from read")
	}
	b.lastRead = m
	return m, e

}

func (b *Buffer) WriteTo(w io.Writer) (n int, err error) {
	return w.Write(b.buf[0:b.lastRead])
}

func (b *Buffer) String() string {
	return string(b.buf[0:b.lastRead])
}
