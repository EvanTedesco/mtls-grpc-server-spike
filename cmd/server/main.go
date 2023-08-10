package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/EvanTedesco/mtls-grpc-server-spike/greet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

func main() {
	server := grpc.NewServer(
		grpc.Creds(LoadKeyPair()),
		grpc.UnaryInterceptor(tlsMiddleware),
	)

	greet.RegisterGreetingServer(server, new(GreetServer))

	go func() {
		l, err := net.Listen("tcp", ":10200")
		if err != nil {
			panic(err)
		}
		log.Println("Listening and serving...")
		if err := server.Serve(l); err != nil {
			panic(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	server.GracefulStop()
	log.Println("Shutting down")
}

func tlsMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// get client tls info
	if p, ok := peer.FromContext(ctx); ok {
		if mtls, ok := p.AuthInfo.(credentials.TLSInfo); ok {
			for _, item := range mtls.State.PeerCertificates {
				log.Println("request certificate subject:", item.Subject)
			}
		}
	}
	return handler(ctx, req)
}

type GreetServer struct {
	greet.UnimplementedGreetingServer
}

func (g *GreetServer) SayHello(ctx context.Context, req *greet.SayHelloRequest) (*greet.SayHelloResponse, error) {
	response := "Hello," + req.GetName()
	return &greet.SayHelloResponse{Greet: response}, nil
}

func LoadKeyPair() credentials.TransportCredentials {
	certificate, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
	if err != nil {
		panic("Failed to load server certification: " + err.Error())
	}

	data, err := ioutil.ReadFile("certs/ca.crt")
	if err != nil {
		panic("Failed to load CA file: " + err.Error())
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(data) {
		panic("Can't add CA cert")
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    caPool,
	}
	return credentials.NewTLS(tlsConfig)
}