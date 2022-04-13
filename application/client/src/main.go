package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

var (
	socketPath     = os.Getenv("socketPath")
	serverAddress  = os.Getenv("serverAddress")
	serverSPIFFEID = os.Getenv("serverSPIFFEID")
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a `workloadapi.X509Source`, it will connect to Workload API using provided socket path
	// If socket path is not defined using `workloadapi.SourceOption`, value from environment variable `SPIFFE_ENDPOINT_SOCKET` is used.
	source, err := workloadapi.NewX509Source(ctx, workloadapi.WithClientOptions(workloadapi.WithAddr(socketPath)))
	if err != nil {
		log.Fatalf("Unable to create X509Source %v", err)
	}

	defer source.Close()

	// Allowed SPIFFE ID
	serverID := spiffeid.RequireFromString(serverSPIFFEID)

	// Create a `tls.Config` to allow mTLS connections, and verify that presented certificate has SPIFFE ID `spiffe://example.org/server`
	tlsConfig := tlsconfig.MTLSClientConfig(source, source, tlsconfig.AuthorizeID(serverID))
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	for {
		r, err := client.Get(serverAddress)

		if err != nil {
			log.Fatalf("Error connecting to %q: %v", serverAddress, err)
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Fatalf("Unable to read body: %v", err)
		}

		log.Printf("%s", body)

		time.Sleep(10 * time.Second)
	}
}
