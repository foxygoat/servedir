// servedir command start http server, serving files in current directory
//       go get foxygo.at/servedir
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func usage() {
	w := flag.CommandLine.Output()
	fmt.Fprintf(w, "usage: %s [-a] [-p <port>] [<dir>]\n\n", os.Args[0])
	fmt.Fprintf(w, "Simple HTTP server, serving files from given directory.\n\n")
	flag.PrintDefaults()
	fmt.Fprintf(w, "  <dir> defaults to current directory if not specified\n")
}

func dir(args []string) string {
	if len(args) == 0 {
		return "."
	}
	return args[0]
}

func listenAddr(port int, allInterfaces bool) string {
	if allInterfaces {
		return fmt.Sprintf(":%d", port)
	}
	return fmt.Sprintf("127.0.0.1:%d", port)
}

func listenAddrURL(addr *net.TCPAddr) string {
	if addr.IP.IsLoopback() {
		return fmt.Sprintf("http://localhost:%d", addr.Port)
	}
	if addr.IP.IsUnspecified() {
		if h, err := os.Hostname(); err == nil {
			return fmt.Sprintf("http://%s:%d", h, addr.Port)
		}
	}
	return "http://" + addr.String()
}

type config struct {
	dir        string
	listenAddr string
}

func parseFlags() config {
	port := flag.Int("p", 0, "port number (default: os chosen free port)")
	allInterfaces := flag.Bool("a", false, "listen on all interfaces not just localhost")
	flag.Usage = usage
	flag.Parse()
	return config{
		dir:        dir(flag.Args()),
		listenAddr: listenAddr(*port, *allInterfaces),
	}
}

func main() {
	cfg := parseFlags()
	http.Handle("/", http.FileServer(http.Dir(cfg.dir)))
	listener, err := net.Listen("tcp", cfg.listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	addr := listener.Addr().(*net.TCPAddr)
	fmt.Printf("Starting HTTP server on %s\n", listenAddrURL(addr))
	log.Fatal(http.Serve(listener, nil))
}
