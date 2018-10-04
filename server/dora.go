package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/raggledodo/dora/data"
	"github.com/raggledodo/dora/proto"
	"golang.org/x/net/context"
)

type (
	DoraServer struct {
		*grpc.Server
		service  *DoraService
		listener net.Listener
	}

	DoraService struct {
		db data.Database
	}
)

var _ proto.DoraServer = &DoraService{}

func (m *DoraService) ListTestcases(ctx context.Context, req *proto.ListRequest) (
	*proto.ListResponse, error) {
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

func (m *DoraService) AddTestcase(ctx context.Context, req *proto.AddRequest) (
	*empty.Empty, error) {
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

func NewDoraServer(host string, db data.Database) *DoraServer {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	service := &DoraService{db: db}
	proto.RegisterDoraServer(server, service)

	return &DoraServer{
		Server:   server,
		service:  service,
		listener: listener,
	}
}

func (server *DoraServer) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Print("Dora server running")
	err := server.Serve(server.listener)
	if err != nil && err != grpc.ErrServerStopped {
		log.Fatalf("Dora server error: %v", err)
	}
	log.Print("Dora server stopped")
}
