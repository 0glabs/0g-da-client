# Deployment

## Prerequisites

- For local test environment, [aws-cli](https://aws.amazon.com/cn/cli/) is required.

## LocalStack
Create LocalStack(local aws simulation) docker image and start a docker instance:
```bash
cd inabox

make deploy-localstack
```

## Disperser
Build binaries:
```
cd disperser

make build
```
Run encoder:
```
make run_encoder
```
Set the cli arguments of run_batcher in Makefile to proper values, then run batcher and server:
```
make run_batcher

make run_server
```