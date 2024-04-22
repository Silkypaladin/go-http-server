package main

import (
	"bytes"
	"strings"
)

type Request struct {
	Method  string
	URL     string
	Version string
	Headers Header
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
		r.Headers.Add(name, value)
	}
}

func NewRequest(buffer [][]byte) *Request {
	reqInfo := bytes.Split(buffer[0], []byte(WHITESPACE))

	req := &Request{
		Method:  string(reqInfo[0]),
		URL:     string(reqInfo[1]),
		Version: string(reqInfo[2]),
	}
	req.ParseHeaders(buffer[1:])
	return req
}
