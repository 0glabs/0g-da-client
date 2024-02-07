# Disperser API

## Table of Contents

* [Service](disperser.md#service)
  * [Disperser](disperser.md#disperser-Disperser)
* [Data Structure](disperser.md#data-structure)
  * [BatchHeader](disperser.md#disperser-BatchHeader)
  * [BatchMetadata](disperser.md#disperser-BatchMetadata)
  * [BlobHeader](disperser.md#disperser-BlobHeader)
  * [BlobInfo](disperser.md#disperser-BlobInfo)
  * [BlobQuorumParam](disperser.md#disperser-BlobQuorumParam)
  * [BlobStatusReply](disperser.md#disperser-BlobStatusReply)
  * [BlobStatusRequest](disperser.md#disperser-BlobStatusRequest)
  * [BlobVerificationProof](disperser.md#disperser-BlobVerificationProof)
  * [DisperseBlobReply](disperser.md#disperser-DisperseBlobReply)
  * [DisperseBlobRequest](disperser.md#disperser-DisperseBlobRequest)
  * [RetrieveBlobReply](disperser.md#disperser-RetrieveBlobReply)
  * [RetrieveBlobRequest](disperser.md#disperser-RetrieveBlobRequest)
  * [SecurityParams](disperser.md#disperser-SecurityParams)
  * [BlobStatus](disperser.md#disperser-BlobStatus)
* [Scalar Value Types](disperser.md#scalar-value-types)

[Top](disperser.md#top)

## Service

### Disperser

Disperser defines the public APIs for dispersing blobs.

| Method Name   | Request Type                                                      | Response Type                                                 | Description                                                                                                                                                                                                                                                                                                                       |
| ------------- | ----------------------------------------------------------------- | ------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| DisperseBlob  | [DisperseBlobRequest](disperser.md#disperser-DisperseBlobRequest) | [DisperseBlobReply](disperser.md#disperser-DisperseBlobReply) | This API accepts blob to disperse from clients. This executes the dispersal async, i.e. it returns once the request is accepted. The client could use GetBlobStatus() API to poll the the processing status of the blob.                                                                                                          |
| GetBlobStatus | [BlobStatusRequest](disperser.md#disperser-BlobStatusRequest)     | [BlobStatusReply](disperser.md#disperser-BlobStatusReply)     | This API is meant to be polled for the blob status.                                                                                                                                                                                                                                                                               |
| RetrieveBlob  | [RetrieveBlobRequest](disperser.md#disperser-RetrieveBlobRequest) | [RetrieveBlobReply](disperser.md#disperser-RetrieveBlobReply) | This retrieves the requested blob from the Disperser's backend. This is a more efficient way to retrieve blobs than directly retrieving from the DA Nodes (see detail about this approach in api/proto/retriever/retriever.proto). The blob should have been initially dispersed via this Disperser service for this API to work. |

## Data Structure

### BatchHeader

| Field                       | Type                          | Label | Description                                                                                                                                                                       |
| --------------------------- | ----------------------------- | ----- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| batch\_root                 | [bytes](disperser.md#bytes)   |       | The root of the merkle tree with the hashes of blob headers as leaves.                                                                                                            |
| quorum\_numbers             | [bytes](disperser.md#bytes)   |       | All quorums associated with blobs in this batch.                                                                                                                                  |
| quorum\_signed\_percentages | [bytes](disperser.md#bytes)   |       | The percentage of stake that has signed for this batch. The quorum\_signed\_percentages\[i] is percentage for the quorum\_numbers\[i].                                            |
| reference\_block\_number    | [uint32](disperser.md#uint32) |       | The Ethereum block number at which the batch was created. The Disperser will encode and disperse the blobs based on the onchain info (e.g. operator stakes) at this block number. |

### BatchMetadata



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| batch_header | [BatchHeader](#disperser-BatchHeader) |  |  |
| signatory_record_hash | [bytes](#bytes) |  | The hash of all public keys of the operators that did not sign the batch. |
| fee | [bytes](#bytes) |  | The gas fee of confirming this batch. It&#39;s the bytes representation of a big.Int value. |
| confirmation_block_number | [uint32](#uint32) |  | The Ethereum block number at which the batch is confirmed onchain. |
| batch_header_hash | [bytes](#bytes) |  | This is the hash of the ReducedBatchHeader defined onchain, see: https://github.com/0glabs/0g-data-avail/blob/master/contracts/src/interfaces/IZGDAServiceManager.sol#L43 The is the message that the operators will sign their signatures on. |






<a name="disperser-BlobHeader"></a>

### BlobHeader

| Field                | Type                                                      | Label    | Description                                                  |
| -------------------- | --------------------------------------------------------- | -------- | ------------------------------------------------------------ |
| commitment           | [bytes](disperser.md#bytes)                               |          | KZG commitment to the blob.                                  |
| data\_length         | [uint32](disperser.md#uint32)                             |          | The length of the blob in symbols (each symbol is 31 bytes). |
| blob\_quorum\_params | [BlobQuorumParam](disperser.md#disperser-BlobQuorumParam) | repeated | The params of the quorums that this blob participates in.    |

### BlobInfo

BlobInfo contains information needed to confirm the blob against the ZGDA contracts

| Field                     | Type                                                                  | Label | Description |
| ------------------------- | --------------------------------------------------------------------- | ----- | ----------- |
| blob\_header              | [BlobHeader](disperser.md#disperser-BlobHeader)                       |       |             |
| blob\_verification\_proof | [BlobVerificationProof](disperser.md#disperser-BlobVerificationProof) |       |             |

### BlobQuorumParam



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| quorum_number | [uint32](#uint32) |  | The ID of the quorum. |
| adversary_threshold_percentage | [uint32](#uint32) |  | Same as SecurityParams.adversary_threshold. |
| quorum_threshold_percentage | [uint32](#uint32) |  | Same as SecurityParams.quorum_threshold. |
| quantization_param | [uint32](#uint32) |  | This determines the nominal number of chunks for the blob, which is nominal_num_chunks = quantization_param * num_operators. A chunk is the smallest unit that&#39;s distributed to DA Nodes, corresponding to a set of evaluations of the polynomial (representing the blob) and a KZG multiproof. See more details in data model of ZGDA: https://github.com/0glabs/0g-data-avail/blob/master/docs/spec/data-model.md |
| encoded_length | [uint64](#uint64) |  | The length of the blob after encoding (in number of symbols). |






<a name="disperser-BlobStatusReply"></a>

### BlobStatusReply

| Field  | Type                                            | Label | Description                                                                      |
| ------ | ----------------------------------------------- | ----- | -------------------------------------------------------------------------------- |
| status | [BlobStatus](disperser.md#disperser-BlobStatus) |       | The status of the blob.                                                          |
| info   | [BlobInfo](disperser.md#disperser-BlobInfo)     |       | The blob info needed for clients to confirm the blob against the ZGDA contracts. |

### BlobStatusRequest

BlobStatusRequest is used to query the status of a blob.

| Field       | Type                        | Label | Description |
| ----------- | --------------------------- | ----- | ----------- |
| request\_id | [bytes](disperser.md#bytes) |       |             |

### BlobVerificationProof

| Field            | Type                                                  | Label | Description                                                                                                                                                                                                                                                                                                                                       |
| ---------------- | ----------------------------------------------------- | ----- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| batch\_id        | [uint32](disperser.md#uint32)                         |       | batch\_id is an incremental ID assigned to a batch by ZGDAServiceManager                                                                                                                                                                                                                                                                          |
| blob\_index      | [uint32](disperser.md#uint32)                         |       | The index of the blob in the batch (which is logically an ordered list of blobs).                                                                                                                                                                                                                                                                 |
| batch\_metadata  | [BatchMetadata](disperser.md#disperser-BatchMetadata) |       |                                                                                                                                                                                                                                                                                                                                                   |
| inclusion\_proof | [bytes](disperser.md#bytes)                           |       | inclusion\_proof is a merkle proof for a blob header's inclusion in a batch                                                                                                                                                                                                                                                                       |
| quorum\_indexes  | [bytes](disperser.md#bytes)                           |       | indexes of quorums in BatchHeader.quorum\_numbers that match the quorums in BlobHeader.blob\_quorum\_params Ex. BlobHeader.blob\_quorum\_params = \[ { quorum\_number = 0, ... }, { quorum\_number = 3, ... }, { quorum\_number = 5, ... }, ] BatchHeader.quorum\_numbers = \[0, 5, 3] => 0x000503 Then, quorum\_indexes = \[0, 2, 1] => 0x000201 |

### DisperseBlobReply

| Field       | Type                                            | Label | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| ----------- | ----------------------------------------------- | ----- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| result      | [BlobStatus](disperser.md#disperser-BlobStatus) |       | The status of the blob associated with the request\_id.                                                                                                                                                                                                                                                                                                                                                                                                           |
| request\_id | [bytes](disperser.md#bytes)                     |       | The request ID generated by the disperser. Once a request is accepted (although not processed), a unique request ID will be generated. Two different DisperseBlobRequests (determined by the hash of the DisperseBlobRequest) will have different IDs, and the same DisperseBlobRequest sent repeatedly at different times will also have different IDs. The client should use this ID to query the processing status of the request (via the GetBlobStatus API). |

### DisperseBlobRequest

| Field            | Type                                                    | Label    | Description                                                                                                                                                                                                                                                                                                           |
| ---------------- | ------------------------------------------------------- | -------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| data             | [bytes](disperser.md#bytes)                             |          | The data to be dispersed. The size of data must be <= 512KiB.                                                                                                                                                                                                                                                         |
| security\_params | [SecurityParams](disperser.md#disperser-SecurityParams) | repeated | Security parameters allowing clients to customize the safety (via adversary threshold) and liveness (via quorum threshold). Clients can define one SecurityParams per quorum, and specify multiple quorums. The disperser will ensure that the encoded blobs for each quorum are all processed within the same batch. |

### RetrieveBlobReply

RetrieveBlobReply contains the retrieved blob data

| Field | Type                        | Label | Description |
| ----- | --------------------------- | ----- | ----------- |
| data  | [bytes](disperser.md#bytes) |       |             |

### RetrieveBlobRequest

RetrieveBlobRequest contains parameters to retrieve the blob.

| Field               | Type                          | Label | Description |
| ------------------- | ----------------------------- | ----- | ----------- |
| batch\_header\_hash | [bytes](disperser.md#bytes)   |       |             |
| blob\_index         | [uint32](disperser.md#uint32) |       |             |

### SecurityParams

SecurityParams contains the security parameters for a given quorum.

| Field                | Type                          | Label | Description                                                                                                   |
| -------------------- | ----------------------------- | ----- | ------------------------------------------------------------------------------------------------------------- |
| quorum\_id           | [uint32](disperser.md#uint32) |       | The ID of the quorum. The quorum must be already registered on EigenLayer. The ID must be in range \[0, 255]. |
| adversary\_threshold | [uint32](disperser.md#uint32) |       | The max percentage of stake within the quorum that can be held by or delegated to adversarial operators.      |

Clients use this to customize the trust assumption (safety).

Requires: 1 <= adversary\_threshold < 100 | | quorum\_threshold | [uint32](disperser.md#uint32) | | The min percentage of stake that must attest in order to consider the dispersal is successful.

Clients use this to customize liveness requirement. The higher this number, the more operators may need to be up for attesting the blob, so the chance the dispersal request to fail may be higher (liveness for dispersal).

Requires: 1 <= quorum\_threshld <= 100 quorum\_threshld > adversary\_threshold.

Note: The adversary_threshold and quorum_threshold will directly influence the cost of encoding for the blob to be dispersed, roughly by a factor of 100 / (quorum_threshold - adversary_threshold). See the spec for more details: https://github.com/0glabs/0g-data-avail/blob/master/docs/spec/protocol-modules/storage/overview.md |





 


<a name="disperser-BlobStatus"></a>

### BlobStatus

| Name                     | Number | Description                                                                                                                         |
| ------------------------ | ------ | ----------------------------------------------------------------------------------------------------------------------------------- |
| UNKNOWN                  | 0      |                                                                                                                                     |
| PROCESSING               | 1      | PROCESSING means that the blob is currently being processed by the disperser                                                        |
| CONFIRMED                | 2      | CONFIRMED means that the blob has been dispersed to DA Nodes and the dispersed batch containing the blob has been confirmed onchain |
| FAILED                   | 3      | FAILED means that the blob has failed permanently (for reasons other than insufficient signatures, which is a separate state)       |
| FINALIZED                | 4      | FINALIZED means that the block containing the blob's confirmation transaction has been finalized on Ethereum                        |
| INSUFFICIENT\_SIGNATURES | 5      | INSUFFICIENT\_SIGNATURES means that the quorum threshold for the blob was not met for at least one quorum.                          |

## Scalar Value Types

| .proto Type | Notes                                                                                                                                           | C++    | Java       | Python      | Go      | C#         | PHP            | Ruby                           |
| ----------- | ----------------------------------------------------------------------------------------------------------------------------------------------- | ------ | ---------- | ----------- | ------- | ---------- | -------------- | ------------------------------ |
| double      |                                                                                                                                                 | double | double     | float       | float64 | double     | float          | Float                          |
| float       |                                                                                                                                                 | float  | float      | float       | float32 | float      | float          | Float                          |
| int32       | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32  | int        | int         | int32   | int        | integer        | Bignum or Fixnum (as required) |
| int64       | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64  | long       | int/long    | int64   | long       | integer/string | Bignum                         |
| uint32      | Uses variable-length encoding.                                                                                                                  | uint32 | int        | int/long    | uint32  | uint       | integer        | Bignum or Fixnum (as required) |
| uint64      | Uses variable-length encoding.                                                                                                                  | uint64 | long       | int/long    | uint64  | ulong      | integer/string | Bignum or Fixnum (as required) |
| sint32      | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s.                            | int32  | int        | int         | int32   | int        | integer        | Bignum or Fixnum (as required) |
| sint64      | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s.                            | int64  | long       | int/long    | int64   | long       | integer/string | Bignum                         |
| fixed32     | Always four bytes. More efficient than uint32 if values are often greater than 2^28.                                                            | uint32 | int        | int         | uint32  | uint       | integer        | Bignum or Fixnum (as required) |
| fixed64     | Always eight bytes. More efficient than uint64 if values are often greater than 2^56.                                                           | uint64 | long       | int/long    | uint64  | ulong      | integer/string | Bignum                         |
| sfixed32    | Always four bytes.                                                                                                                              | int32  | int        | int         | int32   | int        | integer        | Bignum or Fixnum (as required) |
| sfixed64    | Always eight bytes.                                                                                                                             | int64  | long       | int/long    | int64   | long       | integer/string | Bignum                         |
| bool        |                                                                                                                                                 | bool   | boolean    | boolean     | bool    | bool       | boolean        | TrueClass/FalseClass           |
| string      | A string must always contain UTF-8 encoded or 7-bit ASCII text.                                                                                 | string | String     | str/unicode | string  | string     | string         | String (UTF-8)                 |
| bytes       | May contain any arbitrary sequence of bytes.                                                                                                    | string | ByteString | str         | \[]byte | ByteString | string         | String (ASCII-8BIT)            |
