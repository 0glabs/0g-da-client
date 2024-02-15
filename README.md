# ZeroGDA

## Overview

ZeroGDA is a decentralized data availability (DA) service with deep consideration in security, scalability and decentralization. It is also the first DA solution with a built-in data storage layer. Users interact with ZeroGDA to submit and store their data into [ZeroG Storage](https://github.com/zero-gravity-labs/zerog-storage-client) for later retrieval.

To dive deep into the technical details, continue reading [ZeroGDA protocol spec](docs/overview.md).&#x20;

## API Documentation

The ZeroGDA public API is documented [here](docs/docs/).

## Prerequisites

* For local test environment, [aws-cli](https://aws.amazon.com/cn/cli/) is required.

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

Set the cli arguments of run\_batcher in Makefile to proper values, then run batcher and server:

```
make run_batcher

make run_server
```
