package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/raggledodo/dora/data"
	"github.com/raggledodo/dora/insecure"
)

var cfg Config

type Config struct {
	Port  uint
	PbDir string
}

const (
	DefPorts = 58581
	DefPbDir = "/tmp/protodb"
)

func init() {
	flag.UintVar(&cfg.Port, "port", DefPorts, "server port")
	flag.StringVar(&cfg.PbDir, "pbdir", DefPbDir, "protobuf storage path")
	flag.Parse()

	b, err := json.Marshal(cfg)
	if err == nil {
		log.Print("config: ", string(b))
	}
}

var (
	demoKeyPair  *tls.Certificate
	demoCertPool *x509.CertPool
)

func main() {
	var err error
	pair, err := tls.X509KeyPair([]byte(insecure.Cert), []byte(insecure.Key))
	if err != nil {
		panic(err)
	}
	demoKeyPair = &pair
	demoCertPool = x509.NewCertPool()
	ok := demoCertPool.AppendCertsFromPEM([]byte(insecure.Cert))
	if !ok {
		panic("bad certs")
	}

	Serve(fmt.Sprintf("localhost:%d", cfg.Port), data.NewPbFS(cfg.PbDir))
}
