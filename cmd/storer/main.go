package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
)

type CLI struct {
	Listen *url.URL `help:"the uri to listen on" default:"tcp://0.0.0.0:9000" required:""`
	File   string   `type:"path" help:"where to forward to" required:""`
}

func main() {
	cli := &CLI{}
	kong.Parse(cli,
		kong.Name("forwarder"),
		kong.Description("forward syslog messages to another server"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)

	server, err := createUDPServer(cli.Listen)
	if err != nil {
		log.Fatalf("could not create UDP server: %s", err)
	}

	file, err := os.OpenFile(cli.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("could not open file for write: %s", err)
	}

	log.Printf("writing to file: %s", cli.File)
	listenForUDP(server, file)
}

func listenForUDP(server *net.UDPConn, writer io.Writer) {
	p := make([]byte, 2048)
	for {
		n, _, err := server.ReadFromUDP(p[0:])
		if err != nil {
			log.Printf("could not read from UDP: %s", err)
			continue
		}

		_, err = writer.Write(p[0:n])
		if err != nil {
			log.Printf("could not write to file: %s", err)
		}
	}
}

func createUDPServer(listen *url.URL) (*net.UDPConn, error) {
	port, err := strconv.Atoi(listen.Port())
	if err != nil {
		return nil, fmt.Errorf("cannot parse port: %w", err)
	}

	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(listen.Hostname()),
	}

	server, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return nil, fmt.Errorf("cannot start server: %w", err)
	}

	log.Printf("starting UDP server on %d", port)
	return server, nil
}
