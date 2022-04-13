package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
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
	serverSPIFFEID = os.Getenv("serverSPIFFEID")
)

func main() {
	// Setup context
	ctxtest, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	svid, err := workloadapi.FetchX509SVID(ctxtest, workloadapi.WithAddr(socketPath))

	if err != nil {
		fmt.Println("Error fetching SVID")
	} else {
		fmt.Println("Success fetching SVID")
		fmt.Println(svid.ID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Allowed SPIFFE ID
	serverID := spiffeid.RequireFromString(serverSPIFFEID)

	// Create a TLS connection.
	// The client expects the server to present an SVID with the spiffeID: 'spiffe://example.org/ns/app/server'

	// An alternative when creating Dial is using `spiffetls.Dial` that uses environment variable `SPIFFE_ENDPOINT_SOCKET`
	conn, err := spiffetls.DialWithMode(ctx, "tcp", serverAddress,
		spiffetls.MTLSClientWithSourceOptions(
			tlsconfig.AuthorizeID(serverID),
			workloadapi.WithClientOptions(workloadapi.WithAddr(socketPath)),
		))

	if err != nil {
		log.Fatalf("Unable to create TLS connection: %v", err)
	}

	defer conn.Close()

	// Send a message to the server using the TLS connection
	fmt.Fprintf(conn, "Hello server\n")

	// Read server response
	status, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil && err != io.EOF {
		log.Fatalf("Unable to read server response: %v", err)
	}

	log.Printf("Server says: %q", status)
}
