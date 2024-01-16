clean:
	rm -rf ./bin

build: clean
	# cd .. && make protoc
	go mod tidy
	go build -o ./bin/server ./cmd

run: build
	./bin/server \
	--retriever.hostname localhost \
	--retriever.grpc-port 32011 \
	--retriever.storage.node-url http://0.0.0.0:5678 \
	--retriever.storage.node-url http://0.0.0.0:6789 \
	--retriever.storage.kv-url http://0.0.0.0:7890 \
	--retriever.storage.kv-stream-id 000000000000000000000000000000000000000000000000000000000000f2bd \
	--retriever.storage.flow-contract 0xD46b1B603377085E082b48B5e7AB9b7417633401 \
	--retriever.log.level-std debug \
	--kzg.g1-path ../inabox/resources/kzg/g1.point.300000 \
	--kzg.g2-path ../inabox/resources/kzg/g2.point.300000 \
	--kzg.cache-path ../inabox/resources/kzg/SRSTables \
	--kzg.srs-order 300000