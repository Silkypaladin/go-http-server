package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

const (
	CRLF      = "\r\n"
	OK        = "HTTP/1.1 200 OK" + CRLF + CRLF
	NOT_FOUND = "HTTP/1.1 404 Not Found" + CRLF + CRLF
)

const (
	FORWARD_SLASH = "/"
	WHITESPACE    = " "
)

type HttpHeader struct {
	Method  string
	Path    string
	Version string
}

func CreateHeader(method, path, version string) *HttpHeader {
	return &HttpHeader{
		Method:  method,
		Path:    path,
		Version: version,
	}
}

func HandleConn(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading data")
		os.Exit(1)
	}
	req := bytes.Split(buffer, []byte(CRLF))
	reqInfo := bytes.Split(req[0], []byte(WHITESPACE))
	header := CreateHeader(string(reqInfo[0]), string(reqInfo[1]), string(reqInfo[2]))
	if header.Path == FORWARD_SLASH {
		conn.Write([]byte(OK))
	} else {
		conn.Write([]byte(NOT_FOUND))
	}
}

func main() {
	fmt.Println("Starting server...")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		HandleConn(conn)
	}
}
