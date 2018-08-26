package server

import (
	"path/filepath"
	"os"
	pb "github.com/golang/protobuf/proto"
	"github.com/raggledodo/dora/proto"
	"io/ioutil"
	"github.com/sirupsen/logrus"
	"github.com/raggledodo/dora/storage"
	"fmt"
)

const mainProto = "dora_db.pb"

var _ storage.Database = &ProtoFS{}

type ProtoFS struct {
	path string
	store *proto.CaseStorage
}

func NewProtoFS(dirname string) *ProtoFS {
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		os.MkdirAll(dirname, 0700)
	}
	path := filepath.Join(dirname, mainProto)
	store := &proto.CaseStorage{make(map[string]*proto.Cases)}
	if _, err := os.Stat(path); os.IsExist(err) {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			logrus.Warnf("Failed to read data from %s", path)
		} else if err = pb.Unmarshal(b, store); err != nil {
			logrus.Warnf("Failed to unmarshal protobuf from %s", path)
		}
	}
	return &ProtoFS{
		path: path,
		store: store,
	}
}

func (pfs *ProtoFS) ListTestcases(read func(string)error) error {
	store := pfs.store.GetStorage()
	for key := range store {
		if err := read(key); err != nil {
			return err
		}
	}
	return nil
}

func (pfs *ProtoFS) GetTestResults(testcase string) (chan *proto.GeneratedCase, error) {
	out := make(chan *proto.GeneratedCase)
	store := pfs.store.GetStorage()
	tcases, ok := store[testcase]
	if !ok {
		return nil, fmt.Errorf("testcase %s not found", testcase)
	}
	touts := tcases.GetCases()
	go func() {
		for _, tout := range touts {
			out <- tout
		}
		close(out)
	}()
	return out, nil
}

func (pfs *ProtoFS) AddTestResult(key string, data *proto.GeneratedCase) error {
	store := pfs.store.GetStorage()
	if outputs, ok := store[key]; ok {
		outputs.Cases = append(outputs.Cases, data)
	} else {
		store[key] = &proto.Cases{
			Cases: []*proto.GeneratedCase{data},
		}
	}
	return pfs.update()
}

func (pfs *ProtoFS) RemoveTestResult(key string) error {
	store := pfs.store.GetStorage()
	delete(store, key)
	return pfs.update()
}

func (pfs *ProtoFS) update() error {
	b, err := pb.Marshal(pfs.store)
	if err != nil {
		return err
	}
	ioutil.WriteFile(pfs.path, b, 0644)
	return nil
}
