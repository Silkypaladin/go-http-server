package main

const toLower = 'a' - 'A'

type Header map[string]string

func (h Header) Add(key, value string) {
	key = CanonicalHeaderKey(key)
	h[key] = value
}

func (h Header) Get(key string) string {
	if h == nil {
		return ""
	}
	v := h[CanonicalHeaderKey(key)]

	if len(v) == 0 {
		return ""
	}

	return v
}

// Following https://datatracker.ietf.org/doc/html/rfc7230#section-8.1
// Simplified methods inspired by https://github.com/golang/go/blob/master/src/net/textproto/reader.go#L643
func CanonicalHeaderKey(s string) string {
	upper := true
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !validHeaderFieldByte(c) {
			return s
		}
		if upper && 'a' <= c && c <= 'z' {
			return canonicalMIMEHeaderKey([]byte(s))
		}
		if !upper && 'A' <= c && c <= 'Z' {
			return canonicalMIMEHeaderKey([]byte(s))
		}
		upper = c == '-'
	}
	return s
}

func canonicalMIMEHeaderKey(h []byte) string {
	if len(h) == 0 {
		return ""
	}
	upper := true
	for i, c := range h {
		if upper && 'a' <= c && c <= 'z' {
			c -= toLower
		} else if !upper && 'A' <= c && c <= 'Z' {
			c += toLower
		}
		h[i] = c
		upper = c == '-'
	}
	return string(h)
}

func validHeaderFieldByte(c byte) bool {
	if c >= 65 && c <= 90 {
		// A - Z
		return true
	}
	if c >= 94 && c <= 122 {
		// ^_` a - z
		return true
	}
	if c >= 48 && c <= 57 {
		// 0 - 9
		return true
	}
	if c == 124 || c == 126 {
		// |~
		return true
	}
	if c >= 33 && c <= 46 {
		if c == 34 || c == 40 || c == 41 || c == 44 {
			return false
		}
		// !#$%&'*+-.
		return true
	}
	return false
}
