# DA Client API

## Table of Contents

* [Service](disperser.md#service)
  * [Disperser](disperser.md#disperser)
* [Data Structure](disperser.md#data-structure)
  * [BlobHeader](api-1.md#disperser-BlobHeader)
  * [BlobInfo](api-1.md#disperser-BlobInfo)
  * [BlobStatusReply](api-1.md#disperser-BlobStatusReply)
  * [BlobStatusRequest](api-1.md#disperser-BlobStatusRequest)
  * [DisperseBlobReply](api-1.md#disperser-DisperseBlobReply)
  * [DisperseBlobRequest](api-1.md#disperser-DisperseBlobRequest)
  * [RetrieveBlobReply](api-1.md#disperser-RetrieveBlobReply)
  * [RetrieveBlobRequest](api-1.md#disperser-RetrieveBlobRequest)
  * [BlobStatus](api-1.md#disperser-BlobStatus)
  * [Disperser](api-1.md#disperser-Disperser)
* [Scalar Value Types](api-1.md#scalar-value-types)

## Service

### Disperser

Disperser defines the public APIs for dispersing blobs.

| Method Name   | Request Type                                                  | Response Type                                             | Description                                                                                                                                                                                                              |
| ------------- | ------------------------------------------------------------- | --------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| DisperseBlob  | [DisperseBlobRequest](api-1.md#disperser-DisperseBlobRequest) | [DisperseBlobReply](api-1.md#disperser-DisperseBlobReply) | This API accepts blob to disperse from clients. This executes the dispersal async, i.e. it returns once the request is accepted. The client could use GetBlobStatus() API to poll the the processing status of the blob. |
| GetBlobStatus | [BlobStatusRequest](api-1.md#disperser-BlobStatusRequest)     | [BlobStatusReply](api-1.md#disperser-BlobStatusReply)     | This API is meant to be polled for the blob status.                                                                                                                                                                      |
| RetrieveBlob  | [RetrieveBlobRequest](api-1.md#disperser-RetrieveBlobRequest) | [RetrieveBlobReply](api-1.md#disperser-RetrieveBlobReply) | This retrieves the requested blob from the Disperser's backend. The blob should have been initially dispersed via this Disperser service for this API to work.                                                           |

## Data Structure

### BlobHeader

| Field         | Type                      | Label | Description          |
| ------------- | ------------------------- | ----- | -------------------- |
| storage\_root | [bytes](api-1.md#bytes)   |       | The data merkle root |
| epoch         | [uint64](api-1.md#uint64) |       | Signers epoch        |
| quorum\_id    | [uint64](api-1.md#uint64) |       | Signers quorum id    |

### BlobInfo

BlobInfo contains information needed to confirm the blob against the ZGDA contracts

| Field        | Type                                        | Label | Description |
| ------------ | ------------------------------------------- | ----- | ----------- |
| blob\_header | [BlobHeader](api-1.md#disperser-BlobHeader) |       |             |

### BlobStatusReply

| Field  | Type                                        | Label | Description                                                                      |
| ------ | ------------------------------------------- | ----- | -------------------------------------------------------------------------------- |
| status | [BlobStatus](api-1.md#disperser-BlobStatus) |       | The status of the blob.                                                          |
| info   | [BlobInfo](api-1.md#disperser-BlobInfo)     |       | The blob info needed for clients to confirm the blob against the ZGDA contracts. |

### BlobStatusRequest

BlobStatusRequest is used to query the status of a blob.

| Field       | Type                    | Label | Description |
| ----------- | ----------------------- | ----- | ----------- |
| request\_id | [bytes](api-1.md#bytes) |       |             |

### DisperseBlobReply

| Field       | Type                                        | Label | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| ----------- | ------------------------------------------- | ----- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| result      | [BlobStatus](api-1.md#disperser-BlobStatus) |       | The status of the blob associated with the request\_id.                                                                                                                                                                                                                                                                                                                                                                                                           |
| request\_id | [bytes](api-1.md#bytes)                     |       | The request ID generated by the disperser. Once a request is accepted (although not processed), a unique request ID will be generated. Two different DisperseBlobRequests (determined by the hash of the DisperseBlobRequest) will have different IDs, and the same DisperseBlobRequest sent repeatedly at different times will also have different IDs. The client should use this ID to query the processing status of the request (via the GetBlobStatus API). |

### DisperseBlobRequest

| Field | Type                    | Label | Description                                                      |
| ----- | ----------------------- | ----- | ---------------------------------------------------------------- |
| data  | [bytes](api-1.md#bytes) |       | The data to be dispersed. The size of data must be <= 31744 KiB. |

### RetrieveBlobReply

RetrieveBlobReply contains the retrieved blob data

| Field | Type                    | Label | Description |
| ----- | ----------------------- | ----- | ----------- |
| data  | [bytes](api-1.md#bytes) |       |             |

### RetrieveBlobRequest

RetrieveBlobRequest contains parameters to retrieve the blob.

| Field         | Type                      | Label | Description                                          |
| ------------- | ------------------------- | ----- | ---------------------------------------------------- |
| storage\_root | [bytes](api-1.md#bytes)   |       | The storage hash of data                             |
| epoch         | [uint64](api-1.md#uint64) |       | This identifies the epoch that this blob belongs to. |
| quorum\_id    | [uint64](api-1.md#uint64) |       | Which quorum of the blob this is requesting for.     |

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