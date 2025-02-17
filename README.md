# 0GDA

## Overview

0GDA is a decentralized data availability (DA) service with deep consideration in security, scalability, and decentralization. It is also the first DA solution with a built-in data storage layer. Users interact with the 0G DA Client to submit and store their data into [0G DA Node](https://github.com/0glabs/0g-da-node) for later retrieval.

To dive deep into the technical details, continue reading the [0G DA spec](https://docs.0g.ai/da/0g-da-deep-dive).

## Integration

Check out [this example](https://github.com/0glabs/0g-da-example-rust) for how to integrate the 0GDA into your own applications.

For detailed public APIs, visit the [gRPC API](docs/api/) section.

## Deployment

### Installation

1. Install dependencies

*   For Linux

    ```bash
    sudo apt-get update
    sudo apt-get install cmake
    ```
*   For Mac

    ```bash
    brew install cmake
    ```

2. Install Go

* For Linux
  *   Download the Go installer

      ```bash
      wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
      ```
  *   Extract the archive

      ```bash
      rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
      ```
  *   Add `/usr/local/go/bin` to the `PATH` environment variable by adding the following line to your `~/.profile`

      ```bash
      export PATH=$PATH:/usr/local/go/bin
      ```
*   For Mac

    * Download the Go installer from [https://go.dev/dl/go1.22.0.darwin-amd64.pkg](https://go.dev/dl/go1.22.0.darwin-amd64.pkg).
    
    * Follow the installation prompts after opening the package


3.  Download the source code

    ```bash
    git clone -b v1.0.0-testnet https://github.com/0glabs/0g-da-client.git
    ```


### Configuration

| Field                                      | Description                                                        |
| ------------------------------------------ | ------------------------------------------------------------------ |
| `--chain.rpc`                              | JSON RPC node endpoint for the blockchain network.                 |
| `--chain.private-key`                      | Hex-encoded signer private key.                                    |
| `--chain.receipt-wait-rounds`              | Maximum retries to wait for transaction receipt.                   |
| `--chain.receipt-wait-interval`            | Interval between retries when waiting for transaction receipt.     |
| `--chain.gas-limit`                        | Transaction gas limit.                                             |
| `--combined-server.use-memory-db`          | Whether to use mem-db for blob storage.                            |
| `--combined-server.storage.kv-db-path`     | Path for level db.                                                 |
| `--combined-server.storage.time-to-expire` | Expiration duration for blobs in level db.                         |
| `--combined-server.log.level-file`         | File log level.                                                    |
| `--combined-server.log.level-std`          | Standard output log level.                                         |
| `--combined-server.log.path`               | Log file path.                                                     |
| `--disperser-server.grpc-port`             | Server listening port.                                             |
| `--disperser-server.retriever-address`     | GRPC host for retriever.                                           |
| `--batcher.da-entrance-contract`           | Hex-encoded da-entrance contract address.                          |
| `--batcher.da-signers-contract`            | Hex-encoded da-signers contract address.                           |
| `--batcher.finalizer-interval`             | Interval for finalizing operations.                                |
| `--batcher.finalized-block-count`          | Default number of blocks between finalized block and latest block. |
| `--batcher.confirmer-num`                  | Number of Confirmer threads.                                       |
| `--batcher.max-num-retries-for-sign`       | Number of retries before signing fails.                            |
| `--batcher.batch-size-limit`               | Maximum batch size in MiB.                                         |
| `--batcher.encoding-request-queue-size`    | Size of the encoding request queue.                                |
| `--batcher.encoding-interval`              | Interval between blob encoding requests.                           |
| `--batcher.pull-interval`                  | Interval for pulling from the encoded queue.                       |
| `--batcher.signing-interval`               | Interval between slice signing requests.                           |
| `--batcher.signed-pull-interval`           | Interval for pulling from the signed queue.                        |
| `--encoder-socket`                         | GRPC host of the encoder.                                          |
| `--encoding-timeout`                       | Total time to wait for a response from encoder.                    |
| `--signing-timeout`                        | Total time to wait for a response from signer.                     |

### Run

1.  Build combined server

    ```bash
    make build
    ```
2.  Run combined server

    Update the following command by referencing the configuration details from the links below.
    
    Reference:
    
    * [Testnet Configuration](https://docs.0g.ai/run-a-node/testnet-information)
    * [Contract Addresses](https://docs.0g.ai/run-a-node/testnet-information#contract-addresses)

    ```bash
    ./bin/combined \
        --chain.rpc L1_RPC_ENDPOINT \
        --chain.private-key YOUR_PRIVATE_KEY \
        --chain.receipt-wait-rounds 180 \
        --chain.receipt-wait-interval 1s \
        --chain.gas-limit 2000000 \
        --combined-server.use-memory-db \
        --combined-server.storage.kv-db-path ./../run/ \
        --combined-server.storage.time-to-expire 2592000 \
        --disperser-server.grpc-port 51001 \
        --batcher.da-entrance-contract ENTRANCE_CONTRACT_ADDR \
        --batcher.da-signers-contract 0x0000000000000000000000000000000000001000 \
        --batcher.finalizer-interval 20s \
        --batcher.confirmer-num 3 \
        --batcher.max-num-retries-for-sign 3 \
        --batcher.finalized-block-count 50 \
        --batcher.batch-size-limit 500 \
        --batcher.encoding-interval 3s \
        --batcher.encoding-request-queue-size 1 \
        --batcher.pull-interval 10s \
        --batcher.signing-interval 3s \
        --batcher.signed-pull-interval 20s \
        --encoder-socket DA_ENCODER_SERVER \
        --encoding-timeout 300s \
        --signing-timeout 60s \
        --chain-read-timeout 12s \
        --chain-write-timeout 13s \
        --combined-server.log.level-file trace \
        --combined-server.log.level-std trace \
        --combined-server.log.path ./../run/run.log
    ```

### Run via Docker
For Docker deployment instructions, see the [Docker Guide](run/README.md).

### (Optional) Setup DA Retriever
If you want to use the retrieval feature of the DA client, you also need to set up the DA Retriever. Please refer to this [document](https://github.com/0glabs/0g-da-retriever) for more details.


## Contributing

To make contributions to the project, please follow the [contribution guidelines](contributing.md).
