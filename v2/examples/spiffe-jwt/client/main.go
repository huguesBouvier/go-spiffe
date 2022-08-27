package main

import (
	"context"
	_ "fmt"
	_ "io/ioutil"
	"log"
	_ "net/http"
	_ "os"
	"time"

	_ "github.com/spiffe/go-spiffe/v2/spiffeid"
	_ "github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/svid/jwtsvid"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

const (
	serverURL  = "https://localhost:8443"
	socketPath = "unix:///run/spire/sockets/agent.sock"
)

func main() {
	log.Printf("Running")
	// Time out the example after 30 seconds. This prevents the example from hanging if the workloads are not properly registered with SPIRE.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create client options to setup expected socket path,
	// as default sources will use value from environment variable `SPIFFE_ENDPOINT_SOCKET`
	log.Printf("Create client")
	clientOptions := workloadapi.WithClientOptions(workloadapi.WithAddr(socketPath))

	// Create a JWTSource to fetch SVIDs
	log.Printf("Create jwt source")
	jwtSource, err := workloadapi.NewJWTSource(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Unable to create JWTSource: %v", err)
	}
	defer jwtSource.Close()

	log.Printf("get jwt")
	svid, err := jwtSource.FetchJWTSVID(ctx, jwtsvid.Params{
		Audience: "api://AzureADTokenExchange",
	})
	if err != nil {
		log.Fatalf("Unable to fetch SVID: %v", err)
	}

	test, err := jwtsvid.ParseAndValidate(svid.Marshal(), jwtSource, []string{"api://AzureADTokenExchange"})
	if err != nil {
		log.Fatalf("Validation error: %v", err)
	}

	log.Printf("%s", test)

	/*
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", svid.Marshal()))

		res, err := client.Do(req)
		if err != nil {
			log.Fatalf("Unable to connect to %q: %v", serverURL, err)
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}
		log.Printf("%s", body)*/
}
