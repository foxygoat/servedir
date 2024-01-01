package main

import (
	"bytes"
	"flag"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUsage(t *testing.T) {
	w := bytes.NewBuffer(nil)
	fs := flag.NewFlagSet("servedir", flag.ExitOnError)
	fs.SetOutput(w)
	usage(fs)
	opts := `[-a] [-p <port>] [<dir>]`
	desc := `Simple HTTP server, serving files from given directory.`
	require.Contains(t, w.String(), opts)
	require.Contains(t, w.String(), desc)
}

func TestDir(t *testing.T) {
	require.Equal(t, ".", dir(nil))
	require.Equal(t, "abc", dir([]string{"abc"}))
	require.Equal(t, "abc", dir([]string{"abc", "efg"}))
}

func TestListenAddr(t *testing.T) {
	require.Equal(t, ":0", listenAddr(0, true))
	require.Equal(t, "127.0.0.1:0", listenAddr(0, false))
}

func TestListenAddrURL(t *testing.T) {
	addr := &net.TCPAddr{
		Port: 80,
		IP:   net.IPv4(127, 0, 0, 1),
	}
	require.Equal(t, "http://localhost:80", listenAddrURL(addr))

	addr = &net.TCPAddr{
		Port: 80,
		IP:   net.IPv4(0, 0, 0, 0),
	}
	got := listenAddrURL(addr)
	require.True(t, strings.HasPrefix(got, "http://"))
	require.True(t, strings.HasSuffix(got, ":80"))

	addr = &net.TCPAddr{
		Port: 80,
	}
	require.Equal(t, "http://:80", listenAddrURL(addr))
}

func TestParseFlags(t *testing.T) {
	want := config{
		listenAddr: "127.0.0.1:0",
		dir:        ".",
	}
	got := parseFlags()
	require.Equal(t, want, got)
}

func TestParseFlagsFromEnvironment(t *testing.T) {
	want := config{
		listenAddr: ":1",
		dir:        ".",
	}
	t.Setenv("SERVEDIR_ALL_INTERFACES", "true")
	t.Setenv("SERVEDIR_PORT", "1")
	got := parseFlags()
	require.Equal(t, want, got)
}
