package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/EvanTedesco/mtls-grpc-server-spike/greet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func LoadKeyPair() credentials.TransportCredentials {
	certificate, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		panic("Load client certification failed: " + err.Error())
	}

	ca, err := ioutil.ReadFile("certs/ca.crt")
	if err != nil {
		panic("can't read certs/ca.crt file")
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(ca) {
		panic("invalid CA file")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs:      caPool,
	}

	return credentials.NewTLS(tlsConfig)
}

func main() {
	conn, err := grpc.Dial("localhost:10200", grpc.WithTransportCredentials(LoadKeyPair()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := greet.NewGreetingClient(conn)

	var name = "world"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	resp, err := client.SayHello(context.Background(), &greet.SayHelloRequest{Name: name})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.GetGreet())
}