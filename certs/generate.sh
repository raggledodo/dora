#!/usr/bin/env bash

# taken from https://github.com/philips/grpc-gateway-example/blob/master/certs/Makefile
# run as bash script to executed from out of directory and easily from bazel
THIS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )";
openssl genrsa -out $THIS_DIR/server.key 2048
openssl req -new -x509 -key $THIS_DIR/server.key -out $THIS_DIR/server.pem -days 365
