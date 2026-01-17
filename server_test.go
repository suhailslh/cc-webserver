package main

import (
	"io"
	"net"
	"os"
	"sync"
	"testing"
	"time"
)

func TestRunConcurrentOK(t *testing.T) {
	expected := "HTTP/1.1 200 OK\r\ncontent-length: 184\r\n\r\n<!DOCTYPE html>\n<html lang=\"en\">\n  <head>\n    <title>Simple Web Page</title>\n  </head>\n  <body>\n    <h1>Test Web Page</h1>\n    <p>My web server served this page!</p>\n  </body>\n</html>\n"

	ready := make(chan bool, 1)
	interrupt := make(chan os.Signal, 1)
	go func() {
		run("localhost:8080", ready, interrupt)
	}()
	
	<-ready

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			conn, err := net.Dial("tcp", "localhost:8080")		
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			
			time.Sleep(2 * time.Second)
			_, err = conn.Write([]byte("GET / HTTP/1.1\r\n\r\n"))
			if err != nil {
				t.Fatal(err)
			}

			buf := make([]byte, len(expected))
			_, err = conn.Read(buf)
			if err != nil && err != io.EOF {
				t.Fatal(err)
			}

			actual := string(buf)
			if (actual != expected) {
				t.Errorf("%d: expected %q; actual %q", i, expected, actual)
			}
		}()
	}
	wg.Wait()
	
	interrupt <- os.Interrupt
}

func TestRunConcurrentNotFound(t *testing.T) {
	expected := "HTTP/1.1 404 Not Found\r\n\r\n"

	ready := make(chan bool, 1)
	interrupt := make(chan os.Signal, 1)
	go func() {
		run("localhost:8080", ready, interrupt)
	}()

	<-ready
	
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			conn, err := net.Dial("tcp", "localhost:8080")		
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			
			time.Sleep(2 * time.Second)
			_, err = conn.Write([]byte("GET /xyz HTTP/1.1\r\n\r\n"))
			if err != nil {
				t.Fatal(err)
			}

			buf := make([]byte, len(expected))
			_, err = conn.Read(buf)
			if err != nil && err != io.EOF {
				t.Fatal(err)
			}

			actual := string(buf)
			if (actual != expected) {
				t.Errorf("%d: expected %q; actual %q", i, expected, actual)
			}
		}()
	}
	wg.Wait()
	
	interrupt <- os.Interrupt
}
