package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func main() {
	address := "10.10.1.11:9339"
	caCertPath := "../../cert/ca.cer"
	clientCertPath := "../../cert/gnmiclient.pem"
	clientKeyPath := "../../cert/gnmiclient-key.pem"
	username := "admin"
	password := "Cisc0123"

	// Load CA certificate
	caCert, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		log.Fatalf("could not read CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		log.Fatalf("failed to append CA certificate to pool")
	}

	// Load client certificate and key
	clientCert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		log.Fatalf("could not load client certificate and key: %v", err)
	}

	// Create TLS credentials
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
	})

	// Create a gNMI client
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("could not create gNMI client: %v", err)
	}
	defer conn.Close()

	c := gnmi.NewGNMIClient(conn)

	// Create gNMI subscription request
	subscriptionList := &gnmi.SubscriptionList{
		Prefix: &gnmi.Path{},
		Subscription: []*gnmi.Subscription{
			{
				Path: &gnmi.Path{
					Elem: []*gnmi.PathElem{
						{Name: "interfaces"},
						{Name: "interface", Key: map[string]string{"name": "GigabitEthernet1"}},
					},
				},
				Mode:           gnmi.SubscriptionMode_SAMPLE,
				SampleInterval: 1000000000, // 1 second
			},
		},
		Mode:     gnmi.SubscriptionList_STREAM,
		Encoding: gnmi.Encoding_PROTO,
	}

	// Create context with metadata for authentication
	ctx := metadata.AppendToOutgoingContext(context.Background(), "username", username, "password", password)

	// Start the subscription
	stream, err := c.Subscribe(ctx)
	if err != nil {
		log.Fatalf("could not subscribe: %v", err)
	}

	// Send subscription request
	if err := stream.Send(&gnmi.SubscribeRequest{
		Request: &gnmi.SubscribeRequest_Subscribe{
			Subscribe: subscriptionList,
		},
	}); err != nil {
		log.Fatalf("could not send subscription request: %v", err)
	}

	// Receive and handle telemetry updates
	for {
		response, err := stream.Recv()
		if err != nil {
			log.Fatalf("error receiving response: %v", err)
		}
		fmt.Printf("Received telemetry update: %v\n", response)
	}
}
