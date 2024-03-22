#!/bin/bash

pip3 install -r ../requirements.txt
cp ../api/proto/disperser/disperser.proto .
python -m grpc_tools.protoc --proto_path=. ./disperser.proto --python_out=. --grpc_python_out=.
