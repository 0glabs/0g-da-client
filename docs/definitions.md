# Definitions

## Data Packaging

**Blob:** Blobs are the fundamental unit of data posted to ZGDA by users.

**Batch:** Batch is an aggregated data of multiple blobs

## System Components

**Disperser**. The Disperser is an off-chain component which is responsible for packaging data blobs in a specific way, distributing their data among the DA nodes, aggregating certifications from the nodes, and then pushing the aggregated certificate to the chain. The disperser is an untrusted system component.

**Retriever**. The Retriever is an off-chain component which implements a protocol for receiving data blobs from the set of DA nodes.

