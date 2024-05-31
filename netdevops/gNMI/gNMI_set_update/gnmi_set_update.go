package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		if part == "" {
			continue // Skip empty parts
		}
		elem := &gnmi.PathElem{Name: part}
		elems = append(elems, elem)
	}
	return &gnmi.Path{Elem: elems}
}

func gnmiSet(address, caCert, clientCert, clientKey, username, password string, updateData []map[string]interface{}) (*gnmi.SetResponse, error) {
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

	// 创建更新请求
	var updateList []*gnmi.Update
	for _, item := range updateData {
		path := item["path"].(string)
		val := item["value"]

		jsonVal, err := json.Marshal(val)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal JSON value: %w", err)
		}
		update := &gnmi.Update{
			Path: parsePath(path),
			Val: &gnmi.TypedValue{
				Value: &gnmi.TypedValue_JsonIetfVal{
					JsonIetfVal: jsonVal,
				},
			},
		}
		updateList = append(updateList, update)
	}

	setRequest := &gnmi.SetRequest{
		Update: updateList,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	response, err := client.Set(ctx, setRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to set values: %w", err)
	}

	return response, nil
}

func main() {
	address := "10.10.1.11:9339"
	caCert := "../../cert/ca.cer"
	clientCert := "../../cert/gnmiclient.pem"
	clientKey := "../../cert/gnmiclient-key.pem"
	username := "admin"
	password := "Cisc0123"

	// 读取 JSON 文件
	jsonFile, err := ioutil.ReadFile("update_data.json")
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	var updateData []map[string]interface{}
	if err := json.Unmarshal(jsonFile, &updateData); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	response, err := gnmiSet(address, caCert, clientCert, clientKey, username, password, updateData)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal response: %v", err)
	}
	fmt.Println(string(jsonResponse))
}
