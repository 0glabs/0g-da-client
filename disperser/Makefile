clean:
	rm -rf ./bin

build: build_server build_batcher build_encoder # build_dataapi

build_batcher:
	go build -o ./bin/batcher ./cmd/batcher

build_server:
	go build -o ./bin/server ./cmd/apiserver

build_encoder:
	go build -o ./bin/encoder ./cmd/encoder

build_dataapi:
	go build -o ./bin/dataapi ./cmd/dataapi

run_batcher: build_batcher
	./bin/batcher \
	--batcher.pull-interval 10s \
	--chain.rpc ETH_RPC_ENDPOINT \
	--chain.private-key YOUR_PRIVATE_KEY \
	--batcher.aws.region us-east-1 \
	--batcher.aws.access-key-id localstack \
	--batcher.aws.secret-access-key localstack \
	--batcher.aws.endpoint-url http://0.0.0.0:4566 \
	--batcher.s3-bucket-name test-eigenda-blobstore \
	--batcher.dynamodb-table-name test-BlobMetadata \
	--encoder-socket 0.0.0.0:34000 \
	--batcher.batch-size-limit 10000 \
	--batcher.srs-order 3000 \
	--encoding-timeout 10s \
	--chain-read-timeout 12s \
	--chain-write-timeout 13s \
	--batcher.storage.node-urls http://0.0.0.0:5678,http://0.0.0.0:6789 \
	--batcher.storage.flow-contract FLOW_CONTRACT_ADDR

run_server: build_server
	./bin/server \
	--disperser-server.grpc-port 51001 \
	--disperser-server.s3-bucket-name test-eigenda-blobstore \
	--disperser-server.dynamodb-table-name test-BlobMetadata \
	--disperser-server.aws.region us-east-1 \
	--disperser-server.aws.access-key-id localstack \
	--disperser-server.aws.secret-access-key localstack \
	--disperser-server.aws.endpoint-url http://0.0.0.0:4566 \

run_encoder: build_encoder
	./bin/encoder \
	--disperser-encoder.grpc-port 34000 \
  --disperser-encoder.metrics-http-port 9109 \
  --kzg.g1-path ../inabox/resources/kzg/g1.point \
  --kzg.g2-path ../inabox/resources/kzg/g2.point \
  --kzg.cache-path ../inabox/resources/kzg/SRSTables \
  --kzg.srs-order 3000 \
  --kzg.num-workers 12 \
  --disperser-encoder.log.level-std trace \
  --disperser-encoder.log.level-file trace