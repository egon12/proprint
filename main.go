package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	listenAddr := "0.0.0.0:9000"
	if len(os.Args) > 2 {
		listenAddr = os.Args[2]
	}

	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		panic(err)
	}

	log.Printf("listening at %s\n", listenAddr)

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
		}
		go Dial(c)
	}
}

func Dial(source net.Conn) {
	target, err := net.Dial("tcp", os.Args[1])
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
		log.Println(from + ":\n" + b.String())
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

func printUsage() {
	fmt.Println(`
proprint

usage: proprint <dest-addr> <list-addr>

dest-addr: service or destination address that you want to use
list-addr: new port or addr that you want to listen to (optional)
           by default it will listen at 0.0.0.0:9000

example:

# proxy to postgres service
proprint 127.0.0.1:5432 :9000
`)
}
