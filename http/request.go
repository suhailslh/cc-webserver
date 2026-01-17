package http

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

var headerRe = regexp.MustCompile(`^([\w\-]+):\s*([\w\(\)\<\>\@\,\;\:\\\"\/\[\]\?\=\{\}\ \t\.\*\+\-]+)$`)

type Request struct {
	Method string
	URI string
	Version string
	Headers map[string]string
	Body string
}

func (r *Request) String() string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString("Method: ")
	sb.WriteString(r.Method)
	sb.WriteString("\nURI: ")
	sb.WriteString(r.URI)
	sb.WriteString("\nVersion: ")
	sb.WriteString(r.Version)
	sb.WriteString("\n---Headers---\n")
	for k, v := range r.Headers {
		sb.WriteString(k)
		sb.WriteString(": ")
		sb.WriteString(v)
		sb.WriteString("\n")
	}
	sb.WriteString("---Body---\n")
	sb.WriteString(r.Body)
	sb.WriteString("\n")
	return sb.String()
}

func (r *Request) Parse(conn net.Conn) error {
	scanner := bufio.NewScanner(conn)

	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF {
			return 0, nil, bufio.ErrFinalToken
		}
		if r.Headers == nil {
			i := strings.Index(string(data), "\r\n\r\n")
			if i == -1 {
				return 0, nil, nil
			}
			return i + 4, data[:i], nil		
		} else {
			n, err := strconv.Atoi(r.Headers[HeaderContentLength])
			if err != nil {
				return 0, nil, err
			}
			if len(data) >= n {
				return n, data[:n], bufio.ErrFinalToken
			}	
			return 0, nil, nil
		}
	})


	for scanner.Scan() {
		if r.Headers == nil {
			r.Headers = make(map[string]string)
			tokens := strings.Split(scanner.Text(), "\r\n")
			for i, token := range tokens {
				if i == 0 {
					requestline := strings.Split(token, " ")
					if len(requestline) != 3 {
						return fmt.Errorf("Invalid Request-Line: %q", token)
					}
					r.Method = requestline[0]
					r.URI = requestline[1]
					r.Version = requestline[2]
				} else {
					matches := headerRe.FindStringSubmatch(strings.TrimSpace(token))
					if len(matches) != 3 {
						return fmt.Errorf("Invalid Header: %q", token)
					}
					key := strings.ToLower(matches[1])
					value := matches[2]
					if _, ok := r.Headers[key]; ok {
						r.Headers[key] = r.Headers[key] + "," + value
					} else {
						r.Headers[key] = value
					}
				}
			}
			if _, ok := r.Headers[HeaderContentLength]; !ok {
				break
			}
		} else {
			r.Body = scanner.Text()
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	
	return nil
}
