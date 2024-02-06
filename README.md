# ZeroGDA

## Overview

ZeroGDA is a decentralized data availability (DA) service with deep consideration in security and high throughput. Users interact with ZeroGDA to submit their data and store into [ZeroG Storage](https://github.com/zero-gravity-labs/zerog-storage-client) for later retrieval.

To dive deep into the technical details, continue reading [ZeroGDA protocol spec](docs/spec/overview.md).

## API Documentation

The ZeroGDA public API is documented [here](api/docs/).

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
