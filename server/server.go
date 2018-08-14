package server

import (
	"fmt"
	"net"
	"google.golang.org/grpc"
	"sync"

	"golang.org/x/net/context"
	log "github.com/sirupsen/logrus"
	"github.com/raggledodo/dora/proto"
	"github.com/raggledodo/dora/storage"
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
	nothing *proto.Nothing = &proto.Nothing{}
)

func (m *GRPCMux) AddTestcase(ctx context.Context,
	testcase *proto.Testcase) (*proto.Nothing, error) {
	if testcase == nil {
		return nothing, fmt.Errorf("adding nil testcase")
	}
	fname := testcase.GetName()
	testRes := testcase.GetResults()
	return nothing, m.database.AddTestResult(fname, testRes)
}

func (m *GRPCMux) ListTestcases(in *proto.Nothing,
	stream proto.Dora_ListTestcasesServer) error {
	return m.database.ListTestcases(func(fname string) error {
		if err := stream.Send(&proto.Testname{fname}); err != nil {
			return err
		}
		return nil
	})
}

func (m *GRPCMux) GetTestcase(in *proto.Testname,
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
