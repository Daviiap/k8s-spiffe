package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

var (
	socketPath     = os.Getenv("socketPath")
	serverAddress  = os.Getenv("serverAddress")
	clientSPIFFEID = os.Getenv("clientSPIFFEID")
)

func fetchSVID() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	svid, err := workloadapi.FetchX509SVID(ctx, workloadapi.WithAddr(socketPath))

	if err != nil {
		fmt.Println("Error fetching SVID")
	} else {
		fmt.Println("Success fetching SVID")
		fmt.Println(svid.ID)
	}
}

func main() {
	fetchSVID()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clientID := spiffeid.RequireFromString(clientSPIFFEID)

	listener, err := spiffetls.ListenWithMode(ctx, "tcp", serverAddress,
		spiffetls.MTLSServerWithSourceOptions(
			tlsconfig.AuthorizeID(clientID),
			workloadapi.WithClientOptions(workloadapi.WithAddr(socketPath)),
		))
	if err != nil {
		log.Fatalf("Unable to create TLS listener: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			go handleError(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Printf("Error reading incoming data: %v", err)
		return
	}
	log.Printf("Client says: %q", req)

	if _, err = conn.Write([]byte("Hello client!\n")); err != nil {
		log.Printf("Unable to send response: %v", err)
		return
	}
}

func handleError(err error) {
	log.Printf("Unable to accept connection: %v", err)
}
