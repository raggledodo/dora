package server

import (
	"fmt"
	"net"
	"math/rand"
	"google.golang.org/grpc"
	"sync"

	"golang.org/x/net/context"
	log "github.com/sirupsen/logrus"
	"github.com/raggledodo/dora/proto"
	"github.com/raggledodo/dora/storage"
	protobuf "github.com/golang/protobuf/ptypes/empty"
)

type (
	GRPCServer struct {
		*grpc.Server
		mux  *GRPCMux
		listener net.Listener
	}

	GRPCMux struct {
		database storage.Database
	}
)

var (
	empty *protobuf.Empty = &protobuf.Empty{}
)

func (m *GRPCMux) ListTestcases(in *protobuf.Empty,
	stream proto.Dora_ListTestcasesServer) error {
	return m.database.ListTestcases(func(fname string) error {
		if err := stream.Send(&proto.TransferName{fname}); err != nil {
			return err
		}
		return nil
	})
}

func (m *GRPCMux) GetTestcase(in *proto.TransferName,
	stream proto.Dora_GetTestcaseServer) error {
	if in == nil {
		return fmt.Errorf("getting nil testname")
	}
	fname := in.GetName()
	resultCh, err := m.database.GetTestResults(fname)
	if err != nil {
		return err
	}
	if resultCh == nil {
		return fmt.Errorf("db error: got nil test results")
	}
	for result := range resultCh {
		if err := stream.Send(result); err != nil {
			return err
		}
	}
	return nil
}

func (m *GRPCMux) GetOneTestcase(ctx context.Context,
	in *proto.TransferName) (*proto.GeneratedCase, error) {
	if in == nil {
		return nil, fmt.Errorf("getting nil testname")
	}
	fname := in.GetName()
	resultCh, err := m.database.GetTestResults(fname)
	if err != nil {
		return nil, err
	}
	if resultCh == nil {
		return nil, fmt.Errorf("db error: got nil test results")
	}
	results := []*proto.GeneratedCase{}
	for result := range resultCh {
		results = append(results, result)
	}
	ntests := len(results)
	if ntests > 0 {
		return nil, fmt.Errorf("no tests for %s", fname)
	}
	idx := rand.Intn(ntests)
	return results[idx], nil
}

func (m *GRPCMux) AddTestcase(ctx context.Context,
	tcase *proto.TransferCase) (*protobuf.Empty, error) {
	if tcase == nil {
		return empty, fmt.Errorf("adding nil tcase")
	}
	tname := tcase.GetName()
	testRes := tcase.GetResults()
	return empty, m.database.AddTestResult(tname, testRes)
}

func (m *GRPCMux) RemoveTestcase(ctx context.Context,
	testname *proto.TransferName) (*protobuf.Empty, error) {
	if testname == nil {
		return empty, fmt.Errorf("adding nil testcase")
	}
	return empty, m.database.RemoveTestResult(testname.GetName())
}

func (m *GRPCMux) CheckHealth(ctx context.Context, _ *protobuf.Empty) (*proto.HealthCheckResponse, error) {
	return &proto.HealthCheckResponse{
		Status: proto.HealthCheckResponse_SERVING,
	}, nil
}

func New(host, storageDir string) *GRPCServer {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	grpcMux := &GRPCMux{
		database: NewProtoFS(storageDir),
	}
	proto.RegisterDoraServer(server, grpcMux)

	return &GRPCServer{
		Server: server,
		mux: grpcMux,
		listener: listener,
	}
}

func (server *GRPCServer) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Info("Dora server running")
	err := server.Serve(server.listener)
	if err != nil && err != grpc.ErrServerStopped {
		log.Fatalf("Dora server error: %v", err)
	}
	log.Info("Dora server stopped")
}
