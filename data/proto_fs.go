package data

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	pb "github.com/golang/protobuf/proto"
	"github.com/raggledodo/dora/proto"
)

const mainProto = "dora_pb.pb"

var _ Database = &pbFS{}

type pbFS struct {
	path  string
	store *proto.TestStorage
}

func NewPbFS(dirname string) Database {
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		os.MkdirAll(dirname, 0700)
	}
	path := filepath.Join(dirname, mainProto)
	out := &pbFS{path: path}
	store, err := out.startTransaction()
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("creating new file @ %s", path)
		} else {
			log.Printf("Failed to read data from %s: %v", path, err)
		}
	}
	out.store = store
	return out
}

func (pfs *pbFS) ListTestcases(filter *Filter) (
	map[string]*proto.GeneratedTest, error) {
	if pfs.store == nil {
		store, err := pfs.startTransaction()
		if err != nil {
			log.Printf(
				"Failed to read data from %s for testcase retrieval: %v",
				pfs.path, err)
			return nil, err
		}
		pfs.store = store
	}
	schemas := make(map[string]*proto.GeneratedTest)
	store := pfs.store.GetStorage()
	limit := filter.NcaseLimit
	if filter.Names != nil {
		for _, tname := range filter.Names {
			test := store[tname]
			schemas[tname] = filterTest(test, limit)
		}
	} else {
		for tname, test := range store {
			schemas[tname] = filterTest(test, limit)
		}
	}
	return schemas, nil
}

func (pfs *pbFS) AddTestcase(tname string, gcase *proto.GeneratedCase) error {
	tx, err := pfs.startTransaction()
	if err != nil {
		return err
	}
	store := tx.GetStorage()
	if test, ok := store[tname]; ok {
		test.Cases = append(test.Cases, gcase)
	} else {
		store[tname] = &proto.GeneratedTest{
			Cases: []*proto.GeneratedCase{gcase},
		}
	}
	return pfs.completeTransaction(tx)
}

func (pfs *pbFS) RemoveTestcases(tnames []string) error {
	tx, err := pfs.startTransaction()
	if err != nil {
		return err
	}
	store := tx.GetStorage()
	for _, tname := range tnames {
		if _, ok := store[tname]; !ok {
			return fmt.Errorf(
				"Failed to delete test %s: does not exist", tname)
		}
		delete(store, tname)
	}
	return pfs.completeTransaction(tx)
}

func (pfs *pbFS) startTransaction() (*proto.TestStorage, error) {
	transaction := &proto.TestStorage{
		Storage: make(map[string]*proto.GeneratedTest),
	}
	if _, err := os.Stat(pfs.path); !os.IsExist(err) {
		// path does not exist or other err
		return nil, err
	}
	b, err := ioutil.ReadFile(pfs.path)
	if err != nil {
		return nil, err
	}
	if err = pb.Unmarshal(b, transaction); err != nil {
		return nil, err
	}
	return transaction, nil
}

func (pfs *pbFS) completeTransaction(tx *proto.TestStorage) error {
	b, err := pb.Marshal(tx)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(pfs.path, b, 0644); err != nil {
		return err
	}
	pfs.store = tx
	return nil
}

func filterTest(test *proto.GeneratedTest, limit uint32) *proto.GeneratedTest {
	if limit == 0 {
		return test
	}

	if test != nil {
		return &proto.GeneratedTest{
			Cases: test.Cases[:limit],
		}
	}
	return nil
}
