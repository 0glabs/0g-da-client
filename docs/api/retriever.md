# Retriever API

## Table of Contents

- [Service](retriever.md#service)
  - [Retriever](retriever.md#retriever)
- [Data Structure](retriever.md#data-structure)
  - [BlobReply](retriever.md#blobreply)
  - [BlobRequest](retriever.md#blobrequest)
- [Scaler Value Types](retriever.md#scalar-value-types)

[Top](retriever.md#top)

## Service

### Retriever

| Method Name  | Request Type                            | Response Type                       | Description                                                                                                          |
| ------------ | --------------------------------------- | ----------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| RetrieveBlob | [BlobRequest](retriever.md#blobrequest) | [BlobReply](retriever.md#blobreply) | This fans out request to 0G DA Nodes to retrieve the chunks and returns the reconstructed original blob in response. |

## Data Structure

### BlobRequest

| Field                  | Type   | Label | Description                                                                                                                                                                                                                          |
| ---------------------- | ------ | ----- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| batch_header_hash      | bytes  |       | The hash of the ReducedBatchHeader defined onchain, see: https://github.com/0g-gravity-labs/0g-data-avail/blob/master/contracts/src/interfaces/IZGDAServiceManager.sol#L43 This identifies the batch that this blob belongs to. |
| blob_index             | uint32 |       | Which blob in the batch this is requesting for (note: a batch is logically an ordered list of blobs).                                                                                                                                |
| reference_block_number | uint32 |       | The Ethereum block number at which the batch for this blob was constructed.                                                                                                                                                          |
| quorum_id              | uint32 |       | Which quorum of the blob this is requesting for (note a blob can participate in multiple quorums).                                                                                                                                   |

### BlobReply

| Field | Type  | Label | Description                                                                |
| ----- | ----- | ----- | -------------------------------------------------------------------------- |
| data  | bytes |       | The blob retrieved and reconstructed from the 0G DA Nodes per BlobRequest. |

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
