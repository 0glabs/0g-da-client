# Disperser API

## Table of Contents

- [Service](disperser.md#service)
  - [Disperser](disperser.md#disperser)
- [Data Structure](disperser.md#data-structure)
  - [BatchHeader](disperser.md#batchheader)
  - [BatchMetadata](disperser.md#batchmetadata)
  - [BlobHeader](disperser.md#blobheader)
  - [BlobInfo](disperser.md#blobinfo)
  - [BlobQuorumParam](disperser.md#blobquorumparam)
  - [BlobStatusReply](disperser.md#blobstatusreply)
  - [BlobStatusRequest](disperser.md#blobstatusrequest)
  - [BlobVerificationProof](disperser.md#blobverificationproof)
  - [DisperseBlobReply](disperser.md#disperseblobreply)
  - [DisperseBlobRequest](disperser.md#disperseblobrequest)
  - [RetrieveBlobRequest](disperser.md#retrieveblobrequest)
  - [RetrieveBlobReply](disperser.md#retrieveblobreply)
  - [SecurityParams](disperser.md#securityparams)
  - [BlobStatus](disperser.md#blobstatus)
- [Scalar Value Types](disperser.md#scalar-value-types)

[Top](disperser.md#top)

## Service

### Disperser

Disperser defines the public APIs for dispersing blobs.

<table><thead><tr><th width="172">Method Name</th><th>Request Type</th><th>Response Type</th><th>Description</th></tr></thead><tbody><tr><td>DisperseBlob</td><td><a href="disperser.md#disperseblobrequest">DisperseBlobRequest</a></td><td><a href="disperser.md#disperseblobreply">DisperseBlobReply</a></td><td>This API accepts blob to disperse from clients. This executes the dispersal async, i.e. it returns once the request is accepted. The client could use GetBlobStatus() API to poll the the processing status of the blob.</td></tr><tr><td>GetBlobStatus</td><td><a href="disperser.md#blobstatusrequest">BlobStatusRequest</a></td><td><a href="disperser.md#blobstatusreply">BlobStatusReply</a></td><td>This API is meant to be polled for the blob status.</td></tr><tr><td>RetrieveBlob</td><td><a href="disperser.md#retrieveblobrequest">RetrieveBlobRequest</a></td><td><a href="disperser.md#retrieveblobreply">RetrieveBlobReply</a></td><td>This retrieves the requested blob from the Disperser's backend. This is a more efficient way to retrieve blobs than directly retrieving from the DA Nodes (see detail about this approach in api/proto/retriever/retriever.proto). The blob should have been initially dispersed via this Disperser service for this API to work.</td></tr></tbody></table>

## Data Structure

### BatchHeader

<table><thead><tr><th>Field</th><th>Type</th><th width="134">Label</th><th>Description</th></tr></thead><tbody><tr><td>batch_root</td><td>bytes</td><td></td><td>The root of the merkle tree with the hashes of blob headers as leaves.</td></tr><tr><td>quorum_numbers</td><td>bytes</td><td></td><td>All quorums associated with blobs in this batch.</td></tr><tr><td>quorum_signed_percentages</td><td>bytes</td><td></td><td>The percentage of stake that has signed for this batch. The quorum_signed_percentages[i] is percentage for the quorum_numbers[i].</td></tr><tr><td>reference_block_number</td><td>uint32</td><td></td><td>The Ethereum block number at which the batch was created. The Disperser will encode and disperse the blobs based on the onchain info (e.g. operator stakes) at this block number.</td></tr></tbody></table>

### BatchMetadata

<table><thead><tr><th>Field</th><th>Type</th><th width="84">Label</th><th>Description</th></tr></thead><tbody><tr><td>batch_header</td><td><a href="disperser.md#batchheader">BatchHeader</a></td><td></td><td></td></tr><tr><td>signatory_record_hash</td><td>bytes</td><td></td><td>The hash of all public keys of the operators that did not sign the batch.</td></tr><tr><td>fee</td><td>bytes</td><td></td><td>The gas fee of confirming this batch. It's the bytes representation of a big.Int value.</td></tr><tr><td>confirmation_block_number</td><td>uint32</td><td></td><td>The Ethereum block number at which the batch is confirmed onchain.</td></tr><tr><td>batch_header_hash</td><td>bytes</td><td></td><td>This is the hash of the ReducedBatchHeader defined onchain, see: https://github.com/0glabs/0g-data-avail/blob/master/contracts/src/interfaces/IZGDAServiceManager.sol#L43 The is the message that the operators will sign their signatures on.</td></tr></tbody></table>

### BlobHeader

<table><thead><tr><th>Field</th><th>Type</th><th width="137">Label</th><th>Description</th></tr></thead><tbody><tr><td>commitment</td><td>bytes</td><td></td><td>KZG commitment to the blob.</td></tr><tr><td>data_length</td><td>uint32</td><td></td><td>The length of the blob in symbols (each symbol is 31 bytes).</td></tr><tr><td>blob_quorum_params</td><td><a href="disperser.md#blobquorumparam">BlobQuorumParam</a></td><td>repeated</td><td>The params of the quorums that this blob participates in.</td></tr></tbody></table>

### BlobInfo

BlobInfo contains information needed to confirm the blob against the 0G DA contracts

<table><thead><tr><th>Field</th><th>Type</th><th width="132">Label</th><th>Description</th></tr></thead><tbody><tr><td>blob_header</td><td><a href="disperser.md#blobheader">BlobHeader</a></td><td></td><td></td></tr><tr><td>blob_verification_proof</td><td><a href="disperser.md#blobverificationproof">BlobVerificationProof</a></td><td></td><td></td></tr></tbody></table>

### BlobQuorumParam

<table><thead><tr><th width="214">Field</th><th width="95">Type</th><th width="95">Label</th><th>Description</th></tr></thead><tbody><tr><td>quorum_number</td><td>uint32</td><td></td><td>The ID of the quorum.</td></tr><tr><td>adversary_threshold_percentage</td><td>uint32</td><td></td><td>Same as SecurityParams.adversary_threshold.</td></tr><tr><td>quorum_threshold_percentage</td><td>uint32</td><td></td><td>Same as SecurityParams.quorum_threshold.</td></tr><tr><td>quantization_param</td><td>uint32</td><td></td><td>This determines the nominal number of chunks for the blob, which is <code>nominal_num_chunks = quantization_param * num_operators</code>. A chunk is the smallest unit that's distributed to DA Nodes, corresponding to a set of evaluations of the polynomial (representing the blob) and a KZG multiproof.</td></tr><tr><td>encoded_length</td><td>uint64</td><td></td><td>The length of the blob after encoding (in number of symbols).</td></tr></tbody></table>

### BlobStatusReply

<table><thead><tr><th width="171">Field</th><th width="157">Type</th><th width="138">Label</th><th>Description</th></tr></thead><tbody><tr><td>status</td><td><a href="disperser.md#blobstatus">BlobStatus</a></td><td></td><td>The status of the blob.</td></tr><tr><td>info</td><td><a href="disperser.md#blobinfo">BlobInfo</a></td><td></td><td>The blob info needed for clients to confirm the blob against the 0G DA contracts.</td></tr></tbody></table>

### BlobStatusRequest

BlobStatusRequest is used to query the status of a blob.

<table><thead><tr><th>Field</th><th width="184">Type</th><th width="130">Label</th><th>Description</th></tr></thead><tbody><tr><td>request_id</td><td>bytes</td><td></td><td></td></tr></tbody></table>

### BlobVerificationProof

<table><thead><tr><th width="202">Field</th><th width="159">Type</th><th width="89">Label</th><th>Description</th></tr></thead><tbody><tr><td>batch_id</td><td>uint32</td><td></td><td>batch_id is an incremental ID assigned to a batch by ZGDAServiceManager</td></tr><tr><td>blob_index</td><td>uint32</td><td></td><td>The index of the blob in the batch (which is logically an ordered list of blobs).</td></tr><tr><td>batch_metadata</td><td><a href="disperser.md#batchmetadata">BatchMetadata</a></td><td></td><td></td></tr><tr><td>inclusion_proof</td><td>bytes</td><td></td><td>inclusion_proof is a merkle proof for a blob header's inclusion in a batch</td></tr><tr><td>quorum_indexes</td><td>bytes</td><td></td><td>indexes of quorums in BatchHeader.quorum_numbers that match the quorums in BlobHeader.blob_quorum_params Ex. BlobHeader.blob_quorum_params = [ { quorum_number = 0, ... }, { quorum_number = 3, ... }, { quorum_number = 5, ... }, ] BatchHeader.quorum_numbers = [0, 5, 3] => 0x000503 Then, quorum_indexes = [0, 2, 1] => 0x000201</td></tr></tbody></table>

### DisperseBlobReply

<table><thead><tr><th width="156">Field</th><th width="149">Type</th><th width="127">Label</th><th>Description</th></tr></thead><tbody><tr><td>result</td><td><a href="disperser.md#blobstatus">BlobStatus</a></td><td></td><td>The status of the blob associated with the request_id.</td></tr><tr><td>request_id</td><td>bytes</td><td></td><td>The request ID generated by the disperser. Once a request is accepted (although not processed), a unique request ID will be generated. Two different DisperseBlobRequests (determined by the hash of the DisperseBlobRequest) will have different IDs, and the same DisperseBlobRequest sent repeatedly at different times will also have different IDs. The client should use this ID to query the processing status of the request (via the GetBlobStatus API).</td></tr></tbody></table>

### DisperseBlobRequest

<table><thead><tr><th width="183">Field</th><th>Type</th><th width="135">Label</th><th>Description</th></tr></thead><tbody><tr><td>data</td><td>bytes</td><td></td><td>The data to be dispersed. The size of data must be &#x3C;= 512KiB.</td></tr><tr><td>security_params</td><td><a href="disperser.md#securityparams">SecurityParams</a></td><td>repeated</td><td>Security parameters allowing clients to customize the safety (via adversary threshold) and liveness (via quorum threshold). Clients can define one SecurityParams per quorum, and specify multiple quorums. The disperser will ensure that the encoded blobs for each quorum are all processed within the same batch.</td></tr></tbody></table>

### RetrieveBlobRequest

RetrieveBlobRequest contains parameters to retrieve the blob.

| Field             | Type   | Label | Description |
| ----------------- | ------ | ----- | ----------- |
| batch_header_hash | bytes  |       |             |
| blob_index        | uint32 |       |             |

### RetrieveBlobReply

RetrieveBlobReply contains the retrieved blob data

| Field | Type  | Label | Description |
| ----- | ----- | ----- | ----------- |
| data  | bytes |       |             |

### SecurityParams

SecurityParams contains the security parameters for a given quorum.

<table><thead><tr><th width="218">Field</th><th width="106">Type</th><th width="109">Label</th><th>Description</th></tr></thead><tbody><tr><td>quorum_id</td><td>uint32</td><td></td><td>The ID of the quorum. The quorum must be already registered on EigenLayer. The ID must be in range [0, 255].</td></tr><tr><td>adversary_threshold</td><td>uint32</td><td></td><td>The max percentage of stake within the quorum that can be held by or delegated to adversarial operators.</td></tr></tbody></table>

### BlobStatus

<table><thead><tr><th>Name</th><th width="165.33333333333331">Number</th><th>Description</th></tr></thead><tbody><tr><td>UNKNOWN</td><td>0</td><td></td></tr><tr><td>PROCESSING</td><td>1</td><td>PROCESSING means that the blob is currently being processed by the disperser</td></tr><tr><td>CONFIRMED</td><td>2</td><td>CONFIRMED means that the blob has been dispersed to DA Nodes and the dispersed batch containing the blob has been confirmed onchain</td></tr><tr><td>FAILED</td><td>3</td><td>FAILED means that the blob has failed permanently (for reasons other than insufficient signatures, which is a separate state)</td></tr><tr><td>FINALIZED</td><td>4</td><td>FINALIZED means that the block containing the blob's confirmation transaction has been finalized on Ethereum</td></tr><tr><td>INSUFFICIENT_SIGNATURES</td><td>5</td><td>INSUFFICIENT_SIGNATURES means that the quorum threshold for the blob was not met for at least one quorum.</td></tr></tbody></table>

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
