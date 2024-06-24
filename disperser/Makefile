clean:
	rm -rf ./bin

build: build_server build_batcher build_combined 

build_batcher:
	go build -o ./bin/batcher ./cmd/batcher

build_server:
	go build -o ./bin/server ./cmd/apiserver

build_combined: build_server build_batcher
	go build -o ./bin/combined ./cmd/combined_server

run_batcher: build_batcher
	./bin/batcher \
	--batcher.pull-interval 5s \
	--chain.rpc ETH_RPC_ENDPOINT \
	--chain.private-key YOUR_PRIVATE_KEY \
	--chain.receipt-wait-rounds 180 \
	--chain.receipt-wait-interval 1s \
	--chain.gas-limit 2000000 \
	--batcher.finalizer-interval 300s \
	--batcher.confirmer-num 3 \
	--batcher.aws.region us-east-1 \
	--batcher.aws.access-key-id localstack \
	--batcher.aws.secret-access-key localstack \
	--batcher.aws.endpoint-url http://0.0.0.0:4566 \
	--batcher.s3-bucket-name test-zgda-blobstore \
	--batcher.dynamodb-table-name test-BlobMetadata \
	--encoder-socket 0.0.0.0:34000 \
	--batcher.batch-size-limit 50 \
	--encoding-timeout 15s \
	--chain-read-timeout 12s \
	--chain-write-timeout 13s \
	--batcher.storage.node-url http://0.0.0.0:5678 \
	--batcher.storage.node-url http://0.0.0.0:6789 \
	--batcher.storage.kv-url http://0.0.0.0:7890 \
	--batcher.storage.kv-stream-id 000000000000000000000000000000000000000000000000000000000000f2bd \
	--batcher.storage.flow-contract	FLOW_CONTRACT_ADDR

run_server: build_server
	./bin/server \
	--chain.rpc ETH_RPC_ENDPOINT \
	--disperser-server.grpc-port 51001 \
	--disperser-server.s3-bucket-name test-zgda-blobstore \
	--disperser-server.dynamodb-table-name test-BlobMetadata \
	--disperser-server.aws.region us-east-1 \
	--disperser-server.aws.access-key-id localstack \
	--disperser-server.aws.secret-access-key localstack \
	--disperser-server.aws.endpoint-url http://0.0.0.0:4566

run_combined: build_combined
	./bin/combined \
	--chain.rpc ETH_RPC_ENDPOINT \
	--chain.private-key YOUR_PRIVATE_KEY \
	--chain.receipt-wait-rounds 180 \
	--chain.receipt-wait-interval 1s \
	--chain.gas-limit 2000000 \
	--combined-server.use-memory-db \
	--combined-server.storage.da-entrance-contract ENTRANCE_CONTRACT_ADDR \
	--combined-server.storage.da-signers-contract SIGNERS_CONTRACT_ADDR \
	--combined-server.storage.kv-db-path /data/db \
	--disperser-server.s3-bucket-name test-zgda-blobstore \
	--disperser-server.dynamodb-table-name test-BlobMetadata \
	--disperser-server.grpc-port 51001 \
	--batcher.s3-bucket-name test-zgda-blobstore \
	--batcher.dynamodb-table-name test-BlobMetadata \
	--batcher.pull-interval 5s \
	--batcher.finalizer-interval 20s \
	--batcher.confirmer-num 3 \
	--batcher.max-num-retries-for-sign 2 \
	--batcher.finalized-block-count 50 \
	--encoder-socket 0.0.0.0:34000 \
	--batcher.batch-size-limit 50 \
	--encoding-timeout 15s \
	--chain-read-timeout 12s \
	--chain-write-timeout 13s