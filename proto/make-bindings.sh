#!/usr/bin/env bash
# Use the protoc image to run protoc.sh and generate the bindings.
docker run --user "$(id -u):$(id -g)" -e PROTO=pollaris.proto --mount type=bind,source="$PWD",target=/home/proto/ -it saichler/protoc:latest
docker run --user "$(id -u):$(id -g)" -e PROTO=devices.proto --mount type=bind,source="$PWD",target=/home/proto/ -it saichler/protoc:latest
docker run --user "$(id -u):$(id -g)" -e PROTO=collect.proto --mount type=bind,source="$PWD",target=/home/proto/ -it saichler/protoc:latest

# Now move the generated bindings to the models directory and clean up
mkdir -p ../go/types
mv ./types/*.pb.go ../go/types/.
rm -rf ./types
