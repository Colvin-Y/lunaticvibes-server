#!/bin/bash
RUN_NAME="lunaticvibes-server"
git submodule update --remote --merge && sh gen_proto.sh
rm -rf  go.sum && go mod tidy && go mod vendor
rm -rf output
mkdir output
GOOS=linux GOARCH=amd64 go build -gcflags="all=-N -l" -o output/${RUN_NAME}
chmod +x output/${RUN_NAME}

