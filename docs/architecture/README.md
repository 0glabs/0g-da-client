# Architecture

0G system consists of a data availability layer (0G DA) on top of a decentralized storage system (0G Storage). There is a separate consensus network that is part of both the 0G DA and the 0G Storage.&#x20;

* For 0G Storage, the consensus is responsible for determining the ordering of the uploaded data blocks, realizing the storage mining verification and the corresponding incentive mechanism through smart contract.&#x20;
* For 0G DA, the consensus is in charge of guaranteeing the data availability property of each data block via data availability sampling. In other words, each validator does data availability sampling independently, and once the majority of the validators reach the consensus of the successful sampling results, the data will be treated as available. The data availability sampling is mainly to verify that the specific data block is not maliciously withheld by the client and is indeed ingested into the 0G Storage where it is stored in a reliable and persistent way.

As is shown in Figure 1, data is made available on 0G DA through the following flow:

#### Blob Dispersal

1. A user submits a data blob to a Disperser service.
2. The Disperser encodes the data in accordance with the encoding requirements, constructs the appropriate metadata, and directly stores the blob as well as the metadata to the s3 bucket.
3. A Batcher listens to the metadata updates on s3 and encodes the blob into an encoded blob and stores it in memory.&#x20;
4. The Batcher then packs multiple blobs into one batch, together with the KZG commitments of each blob. It also constructs a merkle tree of the batch which is used to verify a certain blob is in a batch for data integrity.
5. The Batcher sends the merkle root to the on-chain 0G Storage contract. It also sends the full batch to a 0G Storage node.
6. The 0G Storage node will listen to the on-chain event, fetch the merkle root and verify the batch data is aligned with the merkle root. If yes, it stores the batch together with its metadata into 0G Storage kv.

#### Blob Retrieval

1.

<figure><img src="../../.gitbook/assets/image (1).png" alt=""><figcaption><p>Figure 1. Architecture Overview </p></figcaption></figure>

