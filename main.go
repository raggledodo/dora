package main

import (
	"github.com/raggledodo/dora/server"
	"fmt"
	"flag"
	"sync"
	"github.com/sirupsen/logrus"
	"encoding/json"
)

var cfg Config

type Config struct {
	Port uint
	PbDir string
}

const (
	DefPorts = 8581
	DefPbDir = "/tmp/protodb"
)

func main() {
	flag.UintVar(&cfg.Port, "port", DefPorts, "server port")
	flag.StringVar(&cfg.PbDir,"pbdir", DefPbDir, "protobuf storage path")
	flag.Parse()

	b, err := json.Marshal(cfg)
	if err == nil {
		logrus.Info("config: ", string(b))
	}

	srv := server.New(fmt.Sprintf(":%d", cfg.Port), cfg.PbDir)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	srv.Start(wg)
}
