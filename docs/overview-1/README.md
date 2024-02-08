# Security Guarantee

The overall security guarantee provided by ZGDA is actually a composite of many smaller guarantees, each with its own intricacies. For the purpose of exposition, in this section we break down the ZGDA protocol into a set of two modules, which divide roughly along the lines of the guarantees for which the modules are responsible.

## Storage

The main guarantee supported by the storage module concerns the off-chain conditions which mirror the on-chain conditions of the storage module. In particular, the storage module is responsible for upholding the following guarantee:

* Honest custody of complete blob: When the minimal adversarial threshold assumptions of a blob are met for any quorum, then on-chain acceptance of a blob implies a full blob is held by honest DA nodes of that quorum for the designated period.

The Storage module is largely implemented by the DA nodes, with an untrusted supporting role by the Disperser.&#x20;

## Retrieval

The ZGDA retrievers expect for blobs to correspond to evaluations of a polynomial of a certain degree. The blob payload delivered to the node contains a KZG polynomial commitment identifying the polynomial, as well as a separate commitment allowing the node to verify its degree. The payload also contains KZG reveal proofs allowing the node to verify that its received data corresponds to evaluations of the polynomial at a specific evaluation index. This document describes in detail how the node performs all verifications, including calculating the evaluation indices.

The operator node will then perform the following checks for each quorum they are a part of to ensure that the `StoreChunksRequest` is valid:

1. Confirm Degree of KZG Commitment: Use the `lowDegreeProof` to verify that `kzgCommitment` commits to a polynomial of degree less than or equal to `numSysE`\*`degreeE`.
2. Verify frame for the quorum against KZG Commitment: Use the `headerHash` to determine the proper indices for the chunks held by the operator, and then use the multireveal proofs contained in the chunks to verify each chunk against `kzgCommitment` commit.
