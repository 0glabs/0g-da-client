# Introduction

The ZeroGDA (ZGDA) system is a scalable Data Availability (DA) service which interacts with [ZeroG Storage Node](https://github.com/zero-gravity-labs/zerog-storage-rust) to store and retrieve data. In informal terms, DA is a guarantee that a given piece of data is available to anyone who wishes to retrieve it. ZGDA is focused on providing DA with both high security and throughput.

At a high level, a DA system is one which accepts blobs of data via some interface and then makes them available to retrievers through another interface.

Two important aspects of a DA system are:

1. Security: The security of a DA system constitutes the set of conditions which are sufficient to ensure that all data blobs certified by the system as available are indeed available for honest retrievers to download.
2. Throughput: The throughput of a DA system is the rate at which the system is able to accept blobs of data, typically measured in bytes/second.

## Essential Properties of ZGDA

### ZGDA Security

The security of ZGDA is guaranteed at retrieval phase. When an end user posts a blob of data to ZGDA, the disperser determines which place to store the data and does two things:

1. Directly store the data into the pre-configured s3 bucket.
2. Send the blob request into a queue for the batcher to batch multiple blobs together and send out to ZeroG Storage Node for DA. The batcher will also append the KZG commitment to the batch for later verification use.

When an end user wants to retrieve a blob of data from ZGDA, he can directly call the disperser to download the blob from the s3 provided that he trusts the disperser. Otherwise, he can start his own retriever service, and retrieve the blob from ZeroG Storage Node. Here comes with the security guarantee. The retriever service will verify the KZG commitments to check the authentication of the data and send the authenticated data back to the user.

In this way, ZGDA provides not only security guarantee to the data but also efficiency for the user to quickly retrieve the data from the disperser.

### ZGDA Throughput

\[TODO: Complete this section]
