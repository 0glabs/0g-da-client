# Introduction

## Overview

0G DA is a decentralized data availability (DA) service with deep consideration in security, scalability and decentralization. It is also the first DA solution with a built-in data storage layer. Users interact with 0G DA to submit and store their data into[ 0G Storage](https://github.com/0glabs/0g-storage-node) for later retrieval.

## Integration

Check out [this example](https://github.com/0glabs/0g-da-example-rust) for how to integrate the 0G DA into your own applications.

For detailed public APIs, visit [gRPC API](../0G%20DA/broken-reference/) section.

## Deployment

* For local test environment, [aws-cli](https://aws.amazon.com/cli/) is required.
* [Local Stack setup](./#localstack)
* [Disperser](./#disperser)
* [Retriever](./#retriever)

### Disperser

1. Build binaries:

```
cd disperser
make build
```

2. Run encoder:

```
make run_encoder
```

3. Set the cli arguments of run\_batcher in Makefile to proper values. Full list of available configuration parameters are showing below.

```
# default
--batcher.pull-interval 5s
--chain.receipt-wait-rounds 180
--chain.receipt-wait-interval 1s
--chain.gas-limit 2000000
--batcher.finalizer-interval 300s
--batcher.confirmer-num 3
--encoder-socket 0.0.0.0:34000
--batcher.batch-size-limit 50
--batcher.srs-order 300000
--encoding-timeout 10s
--chain-read-timeout 12s
--chain-write-timeout 13s
--batcher.storage.node-url http://0.0.0.0:5678
--batcher.storage.node-url http://0.0.0.0:6789
--batcher.storage.kv-url http://0.0.0.0:7890
--batcher.storage.kv-stream-id 000000000000000000000000000000000000000000000000000000000000f2bd
--batcher.aws.region us-east-2

# custom
# aws
--batcher.aws.access-key-id localstack
--batcher.aws.secret-access-key localstack
--batcher.s3-bucket-name test-zgda-blobstore
--batcher.dynamodb-table-name test-BlobMetadata
# chain
--chain.rpc ETH_RPC_ENDPOINT
--chain.private-key YOUR_PRIVATE_KEY
--batcher.storage.flow-contract FLOW_CONTRACT_ADDR
```

4. Then run batcher and the main disperser server:

```
make run_batcher

make run_server
```

### Retriever

1. Build binaries:

```
cd retriever
make build
```

2. Run the main retriever server:

```
make run
```

## Contributing

To make contributions to the project, please follow the guidelines [here](../../contributing.md).
