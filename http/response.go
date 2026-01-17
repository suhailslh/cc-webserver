package http

import (
	"os"
	"strconv"
	"strings"
)

type Response struct {
	Version string
	StatusCode string
	ReasonPhrase string
	Headers map[string]string
	Body string
}

func (r *Response) String() string {
	var sb strings.Builder
	sb.WriteString(r.Version)
	sb.WriteString(" ")
	sb.WriteString(r.StatusCode)
	sb.WriteString(" ")
	sb.WriteString(r.ReasonPhrase)
	sb.WriteString("\r\n")
	for k, v := range r.Headers {
		sb.WriteString(k)
		sb.WriteString(": ")
		sb.WriteString(v)
		sb.WriteString("\r\n")
	}
	sb.WriteString("\r\n")
	sb.WriteString(r.Body)
	return sb.String()
}

func (r *Response) WriteFile(path string) error {
	if path == "www/" {
		path = "www/index.html"
	}
	
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			r.StatusCode = "404"
			r.ReasonPhrase = "Not Found"
			return nil
		}
		return err
	}
	
	r.StatusCode = "200"
	r.ReasonPhrase = "OK"
	r.Headers[HeaderContentLength] = strconv.Itoa(len(data))
	r.Body = string(data)
	
	return nil
}
