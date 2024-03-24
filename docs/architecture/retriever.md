# Retriever

The Retriever is a service for retrieving chunks corresponding to a blob from the 0G DA operator nodes and reconstructing the original blob from the chunks. This is a client-side library that the users are supposed to operationalize.

When an end user posts a blob of data to 0G DA, the disperser determines which place to store the data and does two things:

1. Directly store the data into the pre-configured s3 bucket.
2. Send the blob request into a queue for the batcher to batch multiple blobs together and send out to 0G Storage Node for DA. The batcher will also append the KZG commitment to the batch for later verification use.

Note: Users generally have two ways to retrieve a blob from 0G DA:

1. Retrieve from the Disperser that the user initially used for dispersal: the API is `Disperser.RetrieveBlob()` as defined in `api/proto/disperser/disperser.proto`
2. Retrieve directly from the 0G DA Nodes, which is supported by this Retriever.

The `Disperser.RetrieveBlob()` is generally faster and cheaper as the Disperser manages the blobs that it has processed, whereas the Retriever.RetrieveBlob() removes the need to trust the Disperser, with the downside of more cost and performance.

1. The user submit the retrieval request to the retriever service with the form of [`BlobRequest`](../data-model.md#blob-request).
2. The retriever client will first fetch the metadata of the blob and verify its merkle proof to guarantee correctness.
3. Then the client will verify the KZG commitments of each chunk in the blob and download the chunk if the check passes.
4. The retriever will finally decode the chunks into the original blob data and send back to the user.
