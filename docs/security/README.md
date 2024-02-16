# Security Guarantee

The overall security guarantee provided by ZGDA is actually a composite of two guarantees, one at blob dispersal phase and another one at blob retrieval phase.

## Dispersal

The main guarantee at the dispersal phase is implemented in the ZeroG Storage module. In particular, the storage module is responsible for upholding the following guarantee:

* ZeroG Storage Contract: receive merkle root from the batcher and emit on-chain events
* ZeroG Storage Node: receive full batch data and verify the data with the submitted on-chain merkle root (by listening on corresponding events).

The merkle root is constructed from the multiple blobs in the batch by the batcher. Its purpose is to verify that certain blob is in the batch. The storage node is responsible for verifying the correctness of the full batch data using the root.

## Retrieval

The ZGDA retrievers expect for blobs to correspond to evaluations of a polynomial of a certain degree. The blob payload delivered to the receiver contains a KZG polynomial commitment identifying the polynomial, as well as a separate commitment allowing the retriever to verify its degree.

The receiver will perform the following checks for each retrieval request to ensure that the `BlobRequest` is valid:

1. Verify the merkle proof in the requested blob metadata by calling `VerifyProofUsing`.
2. Verify the KZG commitment by using`lowDegreeProof` to verify that `BlobCommitments` in the `BlobHeader` commits to a polynomial of degree equal to the commitments length.
3. Verify the KZG commitment for each blob chunk which was previously encoded into the blob during the dispersal phase.
