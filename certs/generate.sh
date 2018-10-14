#!/usr/bin/env bash

# taken from https://github.com/philips/grpc-gateway-example/blob/master/certs/Makefile
# run as bash script to executed from out of directory and easily from bazel
THIS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )";

openssl req -x509 -newkey rsa:4096 -days 365 -nodes -subj '/CN=localhost' \
    -keyout $THIS_DIR/server.key -out $THIS_DIR/server.crt;
