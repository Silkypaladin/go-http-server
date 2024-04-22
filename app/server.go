package main

import (
	"bytes"
	"fmt"
	"log"
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

type HandlerFunc func(net.Conn, *Request)

type Server struct {
	Addr     string
	Handlers map[string]HandlerFunc
}

func (s *Server) AddHandler(url string, handlerFunc HandlerFunc) {
	_, exists := s.Handlers[url]
	if exists {
		log.Fatalf("Handler for %s already exists", url)
	}
	s.Handlers[url] = handlerFunc
}

func (s *Server) HandleConn(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading data")
		return
	}
	req := bytes.Split(buffer, []byte(CRLF))
	request := NewRequest(req)
	url, found := GetHandlerFuncUrl(request.URL)
	if !found {
		conn.Write([]byte(NOT_FOUND))
	}
	s.Handlers[url](conn, request)
}

func (s *Server) Serve() {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		fmt.Printf("Failed to bind to address %s", s.Addr)
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go s.HandleConn(conn)
	}
}

func GetHandlerFuncUrl(url string) (string, bool) {
	switch {
	case url == FORWARD_SLASH:
		return "/", true
	case strings.HasPrefix(url, ECHO):
		return "/echo", true
	case strings.HasPrefix(url, USER_AGENT):
		return "/user-agent", true
	default:
		return "", false
	}
}

func CreateServer(addr string) *Server {
	return &Server{
		Addr:     addr,
		Handlers: map[string]HandlerFunc{},
	}
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
	userAgent := request.Headers.Get("User-Agent")
	if len(userAgent) == 0 {
		conn.Write([]byte(INTERNAL_SERVER_ERROR))
	}
	contentType := CONTENT_TYPE + TEXT_PLAIN + CRLF
	contentLength := CONTENT_LENGTH + strconv.Itoa(len(userAgent)) + CRLF
	response := OK + CRLF + contentType + contentLength + CRLF + userAgent
	conn.Write([]byte(response))
}

func handleRootRequest(conn net.Conn, request *Request) {
	conn.Write([]byte(OK + CRLF + CRLF))
}

func main() {
	server := CreateServer("0.0.0.0:4221")
	server.AddHandler("/echo", handleEchoRequest)
	server.AddHandler("/user-agent", handleUserAgentRequest)
	server.AddHandler("/", handleRootRequest)
	server.Serve()
}
