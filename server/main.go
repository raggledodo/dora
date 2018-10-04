package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/raggledodo/dora/data"
)

var cfg Config

type Config struct {
	Port  uint
	PbDir string
}

const (
	DefPorts = 8581
	DefPbDir = "/tmp/protodb"
)

func main() {
	flag.UintVar(&cfg.Port, "port", DefPorts, "server port")
	flag.StringVar(&cfg.PbDir, "pbdir", DefPbDir, "protobuf storage path")
	flag.Parse()

	b, err := json.Marshal(cfg)
	if err == nil {
		log.Print("config: ", string(b))
	}

	srv := NewDoraServer(fmt.Sprintf(":%d", cfg.Port), data.NewPbFS(cfg.PbDir))
	wg := new(sync.WaitGroup)
	wg.Add(1)
	srv.Start(wg)
}
