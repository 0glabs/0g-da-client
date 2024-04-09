# 0GDA

## Overview

0GDA is a decentralized data availability (DA) service with deep consideration in security, scalability and decentralization. It is also the first DA solution with a built-in data storage layer. Users interact with ZeroGDA to submit and store their data into [0G Storage](https://github.com/0glabs/0g-storage-client) for later retrieval.

To dive deep into the technical details, continue reading [0GDA protocol spec](docs/).

## Integration

Check out [this example](https://github.com/0glabs/0g-da-example-rust) for how to integrate the 0GDA into your own applications.

For detailed public APIs, visit [gRPC API](docs/api/) section.

## Deployment

* For local test environment, [aws-cli](https://aws.amazon.com/cli/) is required.
* [Local Stack setup](./#localstack)
* [Disperser](./#disperser)
* [Retriever](./#retriever)

### LocalStack

Create LocalStack(local aws simulation) docker image and start a docker instance:

```bash
cd inabox

make deploy-localstack
```

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

To make contributions to the project, please follow the guidelines [here](contributing.md).
