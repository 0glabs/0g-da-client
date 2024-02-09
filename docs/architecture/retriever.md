# Retriever

When an end user posts a blob of data to ZGDA, the disperser determines which place to store the data and does two things:

1. Directly store the data into the pre-configured s3 bucket.
2. Send the blob request into a queue for the batcher to batch multiple blobs together and send out to ZeroG Storage Node for DA. The batcher will also append the KZG commitment to the batch for later verification use.

A user can request a certain disperser to retrieve the needed blob data. However, if the user doesn't trust any disperser and is worried about the integrity of the data, it can start a retriever service to retrieve the verified data on its own.

1. The user submit the retrieval request to the retriever service with the form of [`BlobRequest`](../data-model.md#blob-request).
2. The retriever client will first fetch the metadata of the blob and verify its merkle proof to guarantee correctness.
3. Then the client will verify the KZG commitments of each chunk in the blob and download the chunk if the check passes.
4. The retriever will finally decode the chunks into the original blob data and send back to the user.
