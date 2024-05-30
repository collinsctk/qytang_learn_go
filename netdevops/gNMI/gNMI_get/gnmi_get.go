package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"os"
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

func gnmiGet(address, caCert, clientCert, clientKey, username, password string, paths []struct {
	prefix *gnmi.Path
	path   *gnmi.Path
}) {
	creds, err := loadTLSCredentials(caCert, clientCert, clientKey)
	if err != nil {
		log.Fatalf("failed to load TLS credentials: %v", err)
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(&loginCreds{
		username: username,
		password: password,
	}))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := gnmi.NewGNMIClient(conn)

	for _, p := range paths {
		getRequest := &gnmi.GetRequest{
			Prefix:   p.prefix,
			Path:     []*gnmi.Path{p.path},
			Encoding: gnmi.Encoding_JSON_IETF,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		response, err := client.Get(ctx, getRequest)
		if err != nil {
			log.Fatalf("failed to get response for path %v: %v", p.path, err)
		}

		var results []map[string]interface{}
		for _, notification := range response.Notification {
			for _, update := range notification.Update {
				if val, ok := update.Val.GetValue().(*gnmi.TypedValue_JsonIetfVal); ok {
					var result map[string]interface{}
					if err := json.Unmarshal(val.JsonIetfVal, &result); err != nil {
						log.Fatalf("failed to unmarshal JSON for path %v: %v", p.path, err)
					}
					result["path"] = pathToMap(update.Path)
					results = append(results, result)
				}
			}
		}

		jsonResponse, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			log.Fatalf("failed to marshal response for path %v: %v", p.path, err)
		}
		fmt.Printf("Response for path %v:\n%s\n", p.path, string(jsonResponse))
	}
}

func pathToMap(path *gnmi.Path) map[string]interface{} {
	result := make(map[string]interface{})
	elems := []interface{}{}
	for _, elem := range path.Elem {
		e := map[string]interface{}{"name": elem.Name}
		if len(elem.Key) > 0 {
			e["key"] = elem.Key
		}
		elems = append(elems, e)
	}
	result["elem"] = elems
	return result
}

func main() {
	address := "10.10.1.11:9339"
	caCert := "../../cert/ca.cer"
	clientCert := "../../cert/gnmiclient.pem"
	clientKey := "../../cert/gnmiclient-key.pem"
	username := "admin"
	password := "Cisc0123"

	paths := []struct {
		prefix *gnmi.Path
		path   *gnmi.Path
	}{
		// Path 1: "openconfig:/interfaces/interface/state/counters"
		{
			prefix: &gnmi.Path{Origin: "openconfig"},
			path: &gnmi.Path{
				Elem: []*gnmi.PathElem{
					{Name: "interfaces"},
					{Name: "interface"},
					{Name: "state"},
					{Name: "counters"},
				},
			},
		},
		// Path 2: "openconfig:/interfaces/interface[name=GigabitEthernet1]/state"
		{
			prefix: &gnmi.Path{Origin: "openconfig"},
			path: &gnmi.Path{
				Elem: []*gnmi.PathElem{
					{Name: "interfaces"},
					{Name: "interface", Key: map[string]string{"name": "GigabitEthernet1"}},
					{Name: "state"},
				},
			},
		},
		// Path 3: "rfc7951:/cpu-usage/cpu-utilization"
		//{
		//	prefix: &gnmi.Path{Origin: "rfc7951"},
		//	path: &gnmi.Path{
		//		Elem: []*gnmi.PathElem{
		//			{Name: "cpu-usage"},
		//			{Name: "cpu-utilization"},
		//		},
		//	},
		//},
		// Path 4: "rfc7951:/memory-statistics/memory-statistic"
		{
			prefix: &gnmi.Path{Origin: "rfc7951"},
			path: &gnmi.Path{
				Elem: []*gnmi.PathElem{
					{Name: "memory-statistics"},
					{Name: "memory-statistic"},
				},
			},
		},
	}

	gnmiGet(address, caCert, clientCert, clientKey, username, password, paths)
}
