package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type loginCreds struct {
	username string
	password string
}

func (c *loginCreds) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"username": c.username,
		"password": c.password,
	}, nil
}

func (c *loginCreds) RequireTransportSecurity() bool {
	return true
}

func loadTLSCredentials(caFile, certFile, keyFile string) (credentials.TransportCredentials, error) {
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("could not read CA certificate: %w", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	clientCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("could not load client key pair: %w", err)
	}

	config := &tls.Config{
		Certificates:       []tls.Certificate{clientCert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,
	}
	return credentials.NewTLS(config), nil
}

func parsePath(path string) *gnmi.Path {
	elems := []*gnmi.PathElem{}
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if part == "" || strings.Contains(part, ":") {
			continue // Skip empty parts and prefixes
		}
		elem := &gnmi.PathElem{}
		if strings.Contains(part, "[") {
			nameAndKey := strings.SplitN(part, "[", 2)
			elem.Name = nameAndKey[0]
			keyPart := strings.TrimSuffix(nameAndKey[1], "]")
			keyVals := strings.Split(keyPart, "=")
			elem.Key = map[string]string{keyVals[0]: keyVals[1]}
		} else {
			elem.Name = part
		}
		elems = append(elems, elem)
	}
	return &gnmi.Path{Elem: elems}
}

func gnmiGet(address, caCert, clientCert, clientKey, username, password, path string) (map[string]interface{}, error) {
	creds, err := loadTLSCredentials(caCert, clientCert, clientKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS credentials: %w", err)
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(&loginCreds{
		username: username,
		password: password,
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	client := gnmi.NewGNMIClient(conn)

	getRequest := &gnmi.GetRequest{
		Path:     []*gnmi.Path{parsePath(path)},
		Encoding: gnmi.Encoding_JSON_IETF,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	response, err := client.Get(ctx, getRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}

	for _, notification := range response.Notification {
		for _, update := range notification.Update {
			if val, ok := update.Val.GetValue().(*gnmi.TypedValue_JsonIetfVal); ok {
				var result map[string]interface{}
				if err := json.Unmarshal(val.JsonIetfVal, &result); err != nil {
					return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
				}
				return result, nil
			}
		}
	}
	return nil, fmt.Errorf("no valid response found")
}

func main() {
	address := "10.10.1.11:9339"
	caCert := "../../cert/ca.cer"
	clientCert := "../../cert/gnmiclient.pem"
	clientKey := "../../cert/gnmiclient-key.pem"
	username := "admin"
	password := "Cisc0123"
	path := "openconfig:/interfaces/interface[name=GigabitEthernet1]/state"

	response, err := gnmiGet(address, caCert, clientCert, clientKey, username, password, path)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal response: %v", err)
	}
	fmt.Println(string(jsonResponse))
}
