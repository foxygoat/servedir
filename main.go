// The servedir command starts an HTTP server, serving files from either
// the current directory or a specified directory, on the next free
// ephemeral port or a specified port.
//
//	go run foxygo.at/servedir@latest --help
//	usage: servedir [-a] [-p <port>] [<dir>]
//
//	Simple HTTP server, serving files from given directory.
//
//	  -a	listen on all interfaces not just localhost (env: SERVEDIR_ALL_INTERFACES)
//	  -p int
//	        port number (default: os chosen free port) (env: SERVEDIR_PORT)
//	  <dir> defaults to current directory if not specified
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func usage(fs *flag.FlagSet) {
	w := fs.Output()
	fmt.Fprintf(w, "usage: %s [-a] [-p <port>] [<dir>]\n\n", os.Args[0])
	fmt.Fprintf(w, "Simple HTTP server, serving files from given directory.\n\n")
	fs.PrintDefaults()
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

func listenAddrURL(address net.Addr) string {
	addr, ok := address.(*net.TCPAddr)
	if !ok {
		return "<unknown address>"
	}
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

func parseFlags(args ...string) config {
	var portDefault int = 0
	if s, ok := os.LookupEnv("SERVEDIR_PORT"); ok {
		if v, err := strconv.ParseInt(s, 0, strconv.IntSize); err == nil {
			portDefault = int(v)
		}
	}

	var allInterfacesDefault bool = false
	if s, ok := os.LookupEnv("SERVEDIR_ALL_INTERFACES"); ok {
		if v, err := strconv.ParseBool(s); err == nil {
			allInterfacesDefault = v
		}
	}

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	port := fs.Int("p", portDefault, "port number (default: os chosen free port)")
	allInterfaces := fs.Bool("a", allInterfacesDefault, "listen on all interfaces not just localhost")
	fs.Usage = func() { usage(fs) }
	fs.Parse(args) //nolint:errcheck // ExitOnError means this does not return an error
	return config{
		dir:        dir(fs.Args()),
		listenAddr: listenAddr(*port, *allInterfaces),
	}
}

func main() {
	cfg := parseFlags(os.Args[1:]...)
	http.Handle("/", http.FileServer(http.Dir(cfg.dir)))
	listener, err := net.Listen("tcp", cfg.listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	server := newServer(cfg.dir)
	fmt.Printf("Starting HTTP server on %s\n", listenAddrURL(listener.Addr()))
	log.Fatal(server.Serve(listener))
}

func newServer(dir string) *http.Server {
	fileServer := http.FileServer(http.Dir(dir))
	h := func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add("Cache-Control", "no-cache")
		if strings.HasSuffix(req.URL.Path, ".wasm") {
			resp.Header().Set("content-type", "application/wasm")
		}
		fileServer.ServeHTTP(resp, req)
	}
	return &http.Server{Handler: http.HandlerFunc(h)}
}
