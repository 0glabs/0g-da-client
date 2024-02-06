# API

This part contains the API documentation for the gRPC services included in the ZeroGDA platform. Each markdown file contains the protobuf definitions for each respective service including:

* Disperser: the hosted service for users to interact with to [ZeroG Storage](https://github.com/zero-gravity-labs/zerog-storage-client).
* Retriever: a service that users can run on their own infrastructure, which exposes a gRPC endpoint for retrieval and verification of blobs from [ZeroG Storage nodes](https://github.com/zero-gravity-labs/zerog-storage-rust).
