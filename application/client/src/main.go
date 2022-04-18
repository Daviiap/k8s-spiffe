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

	serverID := spiffeid.RequireFromString(serverSPIFFEID)

	for {
		conn, err := spiffetls.DialWithMode(ctx, "tcp", serverAddress,
			spiffetls.MTLSClientWithSourceOptions(
				tlsconfig.AuthorizeID(serverID),
				workloadapi.WithClientOptions(workloadapi.WithAddr(socketPath)),
			))

		if err != nil {
			log.Fatalf("Unable to create TLS connection: %v", err)
		}

		fmt.Fprintf(conn, "Hello server!\n")

		status, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil && err != io.EOF {
			log.Fatalf("Unable to read server response: %v", err)
		}
		log.Printf("Server says: %q", status)

		time.Sleep(5 * time.Second)
	}
}
