package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/raggledodo/dora/data"
)

var cfg Config

type Config struct {
	Host     string
	Port     uint
	PbDir    string
	Keyfile  string
	Certfile string
}

type Certificate struct {
	Key  *tls.Certificate
	Pool *x509.CertPool
}

const (
	DefPorts       = 10000
	DefPbDir       = "/tmp/protodb"
	DefKeyfile     = "certs/server.key"
	DefCertificate = "certs/server.crt"
	DefHost        = "localhost"
)

func init() {
	flag.StringVar(&cfg.Host, "host", DefHost,
		"host name used as common name in certificate")
	flag.UintVar(&cfg.Port, "port", DefPorts, "server port")
	flag.StringVar(&cfg.PbDir, "pbdir", DefPbDir, "protobuf storage path")
	flag.StringVar(&cfg.Keyfile, "key", DefKeyfile, "rsa private key file")
	flag.StringVar(&cfg.Certfile, "cert", DefCertificate, "certificate file")
	flag.Parse()

	b, err := json.Marshal(cfg)
	if err == nil {
		log.Print("config: ", string(b))
	}
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	kfile, err := ioutil.ReadFile(cfg.Keyfile)
	check(err)
	cert, err := ioutil.ReadFile(cfg.Certfile)
	check(err)
	key, err := tls.X509KeyPair(cert, kfile)
	check(err)
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(cert)
	if !ok {
		log.Panic("bad certs")
	}
	certificate := Certificate{
		Key:  &key,
		Pool: pool,
	}

	Serve(cfg.Host, cfg.Port, certificate, data.NewPbFS(cfg.PbDir))
}
