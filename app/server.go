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
	USER_AGENT    = "/user-agent"
	ECHO          = "/echo/"
	FORWARD_SLASH = "/"
	WHITESPACE    = " "
)

type Request struct {
	Method  string
	URL     string
	Version string
	Headers map[string]string
}

func (r *Request) ParseHeaders(headers [][]byte) {
	r.Headers = map[string]string{}
	for _, v := range headers {
		if string(v) == "" {
			// \r\n before request body, all headers parsed
			break
		}
		v := string(v)
		h := strings.Split(v, ":")
		name, value := strings.Trim(h[0], " "), strings.Trim(h[1], " ")
		r.Headers[name] = value
	}
}

func createRequest(buffer [][]byte) *Request {
	reqInfo := bytes.Split(buffer[0], []byte(WHITESPACE))

	req := &Request{
		Method:  string(reqInfo[0]),
		URL:     string(reqInfo[1]),
		Version: string(reqInfo[2]),
	}
	req.ParseHeaders(buffer[1:])
	return req
}

func handleEchoRequest(conn net.Conn, request *Request) {
	data, found := strings.CutPrefix(request.URL, ECHO)

	if !found {
		conn.Write([]byte(INTERNAL_SERVER_ERROR))
	}
	contentType := CONTENT_TYPE + TEXT_PLAIN + CRLF
	contentLength := CONTENT_LENGTH + strconv.Itoa(len(data)) + CRLF
	response := OK + CRLF + contentType + contentLength + CRLF + data
	conn.Write([]byte(response))
}

func handleUserAgentRequest(conn net.Conn, request *Request) {
	userAgent, ok := request.Headers["User-Agent"]
	if !ok {
		conn.Write([]byte(INTERNAL_SERVER_ERROR))
	}
	contentType := CONTENT_TYPE + TEXT_PLAIN + CRLF
	contentLength := CONTENT_LENGTH + strconv.Itoa(len(userAgent)) + CRLF
	response := OK + CRLF + contentType + contentLength + CRLF + userAgent
	conn.Write([]byte(response))
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading data")
		return
	}
	req := bytes.Split(buffer, []byte(CRLF))
	request := createRequest(req)
	switch {
	case request.URL == FORWARD_SLASH:
		conn.Write([]byte(OK + CRLF + CRLF))
	case strings.HasPrefix(request.URL, ECHO):
		handleEchoRequest(conn, request)
	case strings.HasPrefix(request.URL, USER_AGENT):
		handleUserAgentRequest(conn, request)
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
		handleConn(conn)
	}
}
