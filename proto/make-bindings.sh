#!/usr/bin/env bash
set -e
wget https://raw.githubusercontent.com/saichler/l8types/refs/heads/main/proto/services.proto
# Use the protoc image to run protoc.sh and generate the bindings.
docker run --user "$(id -u):$(id -g)" -e PROTO=pollaris.proto --mount type=bind,source="$PWD",target=/home/proto/ -it saichler/protoc:latest
docker run --user "$(id -u):$(id -g)" -e PROTO=targets.proto --mount type=bind,source="$PWD",target=/home/proto/ -it saichler/protoc:latest
docker run --user "$(id -u):$(id -g)" -e PROTO=jobs.proto --mount type=bind,source="$PWD",target=/home/proto/ -it saichler/protoc:latest
rm -rf service.proto

# Now move the generated bindings to the models directory and clean up
rm -rf ../go/types
mkdir -p ../go/types
mv ./types/* ../go/types/.
rm -rf ./types
