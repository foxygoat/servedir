// servedir command start http server, serving files in current directory
//       go run foxygo.at/servedir
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Starting HTTP server on http://localhost:%d\n", listener.Addr().(*net.TCPAddr).Port)
	log.Fatal(http.Serve(listener, nil))
}
