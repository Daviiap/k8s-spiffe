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

func getSVID() {
	ctxtest, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	svid, err := workloadapi.FetchX509SVID(ctxtest, workloadapi.WithAddr(socketPath))

	if err != nil {
		fmt.Println("Error fetching SVID")
	} else {
		fmt.Println("Success fetching SVID")
		fmt.Println(svid.ID)
	}
}

func main() {
	getSVID()

	// Setup context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Allowed SPIFFE ID
	fmt.Println("Getting client spiffeid")
	clientID := spiffeid.RequireFromString(clientSPIFFEID)
	fmt.Println("Client spiffeID:", clientID.String())

	// Creates a TLS listener
	// The server expects the client to present an SVID with the spiffeID: 'spiffe://example.org/ns/app/client'
	//
	// An alternative when creating Listen is using `spiffetls.Listen` that uses environment variable `SPIFFE_ENDPOINT_SOCKET`
	fmt.Println("Creating listenner")
	listener, err := spiffetls.ListenWithMode(ctx, "tcp", serverAddress,
		spiffetls.MTLSServerWithSourceOptions(
			tlsconfig.AuthorizeID(clientID),
			workloadapi.WithClientOptions(workloadapi.WithAddr(socketPath)),
		))
	fmt.Println("Listenner created")

	if err != nil {
		log.Fatalf("Unable to create TLS listener: %v", err)
	}

	defer listener.Close()

	// Handle connections
	fmt.Println("Listening on:", listener.Addr().String())
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

	// Read incoming data into buffer
	req, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Printf("Error reading incoming data: %v", err)
		return
	}
	log.Printf("Client says: %q", req)

	// Send a response back to the client
	if _, err = conn.Write([]byte("Hello client\n")); err != nil {
		log.Printf("Unable to send response: %v", err)
		return
	}
}

func handleError(err error) {
	log.Printf("Unable to accept connection: %v", err)
}
