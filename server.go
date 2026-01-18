package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/suhailslh/cc-webserver/http"
)

func main() {
	ready := make(chan bool, 1)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	
	os.Exit(run(ready, interrupt))
}

func run(ready chan<- bool, interrupt <-chan os.Signal) int {
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Println(err)
		return 1
	}

	go func() {
		<-interrupt
		fmt.Println()
		log.Println("Closing...")
		listener.Close()
	}()
	
	log.Printf("Listening on %s\n", listener.Addr())

	ready <- true
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			break
		}

		go func() {
			defer conn.Close()

			var request http.Request
			err = request.Parse(conn)
			if err != nil {
				log.Println(err)
				return
			}

			log.Println(request.String())

			serve(request, conn)
		}()
	}

	return 0
}

func serve(request http.Request, conn net.Conn) {
	response := http.Response{
		Version: "HTTP/1.1",
		Headers: make(map[string]string),
	}
	
	switch request.Method {
		case http.MethodGet:
			err := response.WriteFile("www" + request.URI)
			if err != nil {
				log.Println(err)
				return
			}
	}
	
	_, err := conn.Write([]byte(response.String()))
	if err != nil {
		log.Println(err)
		return
	}
}
