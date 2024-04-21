package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	CRLF                  = "\r\n"
	OK                    = "HTTP/1.1 200 OK"
	NOT_FOUND             = "HTTP/1.1 404 Not Found" + CRLF + CRLF
	BAD_REQUEST           = "HTTP/1.1 400 Bad Request" + CRLF + CRLF
	INTERNAL_SERVER_ERROR = "HTTP/1.1 500 Internal Server Error" + CRLF + CRLF
	CONTENT_TYPE          = "Content-Type: "
	CONTENT_LENGTH        = "Content-Length: "
)

const (
	TEXT_PLAIN = "text/plain"
)

const (
	ECHO          = "/echo/"
	FORWARD_SLASH = "/"
	WHITESPACE    = " "
)

type RequestHeader struct {
	Method  string
	Path    string
	Version string
}

func CreateHeader(method, path, version string) *RequestHeader {
	return &RequestHeader{
		Method:  method,
		Path:    path,
		Version: version,
	}
}

func HandleEchoRequest(conn net.Conn, header *RequestHeader) {
	data, found := strings.CutPrefix(header.Path, ECHO)

	if !found {
		conn.Write([]byte(INTERNAL_SERVER_ERROR))
	}
	contentType := CONTENT_TYPE + TEXT_PLAIN + CRLF
	contentLength := CONTENT_LENGTH + strconv.Itoa(len(data)) + CRLF
	response := OK + CRLF + contentType + contentLength + CRLF + data
	fmt.Println(response)
	conn.Write([]byte(response))
}

func HandleServerError(conn net.Conn) {
	conn.Write([]byte(INTERNAL_SERVER_ERROR))
}

func HandleConn(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading data")
		conn.Close()
		return
	}
	req := bytes.Split(buffer, []byte(CRLF))
	reqInfo := bytes.Split(req[0], []byte(WHITESPACE))
	header := CreateHeader(string(reqInfo[0]), string(reqInfo[1]), string(reqInfo[2]))
	switch {
	case header.Path == FORWARD_SLASH:
		conn.Write([]byte(OK + CRLF + CRLF))
	case strings.HasPrefix(header.Path, ECHO):
		HandleEchoRequest(conn, header)
	default:
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
