package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

const (
	socketPath = "unix:///run/spire/sockets/agent.sock"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for {
		svid, err := workloadapi.FetchX509SVID(ctx, workloadapi.WithAddr(socketPath))

		if err != nil {
			fmt.Println("Error fetching SVID")
		} else {
			fmt.Println("Success fetching SVID")
			fmt.Println(svid.ID)
		}

		time.Sleep(5 * time.Second)
	}
}
