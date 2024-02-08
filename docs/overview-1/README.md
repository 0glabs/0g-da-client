# Security Guarantee

The overall security guarantee provided by ZGDA is actually a composite of two guarantees, one at blob dispersal phase and another one at blob retrieval phase.

## Dispersal

The main guarantee at the dispersal phase is implemented in the ZeroG Storage module. In particular, the storage module is responsible for upholding the following guarantee:

* ZeroG Storage Contract: receive merkle root from the batcher and emit on-chain events
* ZeroG Storage Node: receive full batch data and verify the data with the submitted on-chain merkle root (by listening on corresponding events).

The merkle root is constructed from the multiple blobs in the batch by the batcher. Its purpose is to verify that certain blob is in the batch. The storage node is responsible for verifying the correctness of the full batch data using the root.

## Retrieval

The security of ZGDA is guaranteed at retrieval phase. When an end user posts a blob of data to ZGDA, the [disperser](../architecture/disperser.md) determines which place to store the data and does two things:

1. Directly store the data into the pre-configured s3 bucket.
2. Send the blob request into a queue for the [batcher](../architecture/batcher.md) to batch multiple blobs together and send out to ZeroG Storage Node for DA. The batcher will also append the [KZG commitment](../pkg/kzg.md) to the batch for later verification use.

When an end user wants to retrieve a blob of data from ZGDA, he can directly call the disperser to download the blob from the s3 provided that he trusts the disperser. Otherwise, he can start his own [retriever](../architecture/retriever.md) service, and retrieve the blob from ZeroG Storage Node. Here comes with the security guarantee. The retriever service will verify the KZG commitment to check the authentication of the data and send the authenticated data back to the user.

In this way, ZGDA provides not only security guarantee to the data but also efficiency for the user to quickly retrieve the data from the disperser.

The ZGDA retrievers expect for blobs to correspond to evaluations of a polynomial of a certain degree. The blob payload delivered to the receiver contains a KZG polynomial commitment identifying the polynomial, as well as a separate commitment allowing the retriever to verify its degree.

The receiver will perform the following checks for each retrieval request to ensure that the `BlobRequest` is valid:

1. Verify the merkle proof in the requested blob metadata by calling `VerifyProofUsing`.
2. Verify the KZG commitment by using`lowDegreeProof` to verify that `BlobCommitments` in the `BlobHeader` commits to a polynomial of degree equal to the commitments length.
3. Verify the KZG commitment for each blob chunk which was previously encoded into the blob during the dispersal phase.
