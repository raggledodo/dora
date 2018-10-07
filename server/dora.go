package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/raggledodo/dora/data"
	"github.com/raggledodo/dora/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type (
	DoraService struct {
		db data.Database
	}
	Credentials struct {
		Key  *tls.Certificate
		Pool *x509.CertPool
	}
)

var _ proto.DoraServer = &DoraService{}

func (m *DoraService) ListTestcases(ctx context.Context,
	req *proto.ListRequest) (*proto.ListResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("nil ListRequest")
	}
	filter := data.ListReqToFilter(req)
	tests, err := m.db.ListTestcases(filter)
	if err != nil {
		return nil, err
	}
	return &proto.ListResponse{
		Tests: tests,
	}, nil
}

func (m *DoraService) AddTestcase(ctx context.Context,
	req *proto.AddRequest) (*empty.Empty, error) {
	if req == nil {
		return nil, fmt.Errorf("nil AddRequest")
	}
	if err := m.db.AddTestcase(req.Name, req.Payload); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (m *DoraService) RemoveTestcase(ctx context.Context,
	req *proto.RemoveRequest) (*empty.Empty, error) {
	if req == nil {
		return nil, fmt.Errorf("nil RemoveRequest")
	}
	if err := m.db.RemoveTestcases(req.Names); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (m *DoraService) CheckHealth(ctx context.Context, _ *empty.Empty) (
	*proto.HealthCheckResponse, error) {
	return &proto.HealthCheckResponse{
		Status: proto.HealthCheckResponse_SERVING,
	}, nil
}

func Serve(host string, certificate Certificate, db data.Database) {
	log.Printf("Serving on %s", host)

	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewClientTLSFromCert(certificate.Pool, host))}

	grpcServer := grpc.NewServer(opts...)
	proto.RegisterDoraServer(grpcServer, &DoraService{db: db})
	ctx := context.Background()

	dcreds := credentials.NewTLS(&tls.Config{
		ServerName: host,
		RootCAs:    certificate.Pool,
	})
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}

	mux := http.NewServeMux()

	gwmux := runtime.NewServeMux()
	err := proto.RegisterDoraHandlerFromEndpoint(ctx, gwmux, host, dopts)
	if err != nil {
		log.Fatalf("Failed to register http: %v", err)
	}
	mux.Handle("/", gwmux)

	conn, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal(err)
	}

	httpServer := &http.Server{
		Addr: host,
		Handler: http.HandlerFunc(
			func(writer http.ResponseWriter, req *http.Request) {
				if req.ProtoMajor == 2 && strings.Contains(
					req.Header.Get("Content-Type"), "application/grpc") {
					grpcServer.ServeHTTP(writer, req)
				} else {
					mux.ServeHTTP(writer, req)
				}
			}),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*certificate.Key},
			NextProtos:   []string{"h2"},
		},
	}

	err = httpServer.Serve(tls.NewListener(conn, httpServer.TLSConfig))
	if err != nil {
		log.Fatalf("Failed to listen and serve: %v", err)
	}
}
