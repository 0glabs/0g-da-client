# Glossary

## Data Packaging

**Blob:** Blobs are the fundamental unit of data posted to ZGDA by users.

**Batch:** Batch is an aggregated data of multiple blobs together with the KZG commitments to each blob.

## System Components

**Disperser**. The Disperser is an off-chain service which is responsible for uploading the data into s3 buckets and sending the blob to the Batcher for further process. The disperser is an untrusted system component.

**Batcher**. The Batcher is an off-chain component which accepts the blob requests from the Disperser and batches multiple blobs together with their KZG commitments into one data packet and send out to the ZeroG Storage nodes for data storage.

**Retriever**. The Retriever is an off-chain service which implements a protocol for receiving data blobs from ZeroG Storage nodes. It is also responsible for verifying the authentication and correctness of the data.

