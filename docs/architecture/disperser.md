# Disperser

Requesters that want to store data on 0G Storage make requests to the disperser with the form of [`DisperseBlobRequest`](../data-model.md#disperser).

They specify the data they want to store on 0G Storage. The disperser takes each `DisperseBlobRequest` and stores it into the s3 bucket. The data that the disperser stores into the s3 bucket contains two parts: blob metadata and full blob data.

### Blob Data

The [key](../data-model.md#blob-key) of the blob data is calculated by a certain hash function. This key is used for users to retrieve the blob directly from a disperser. The disperser uploads the blob data with the object key of the blob hash to a certain s3 bucket (defined by the disperser service).

In pseudocode:

```go
blobHash := getBlobHash(blob)
metadataHash := getMetadataHash(requestedAt, blob.RequestHeader.SecurityParams)
metadataKey.BlobHash = blobHash
metadataKey.MetadataHash = metadataHash

err = s.s3Client.UploadObject(ctx, s.bucketName, blobObjectKey(blobHash), blob.Data)
```

### Metadata

The [metadata](../data-model.md#blob-metadata) of a blob is constructed and stored into a table (defined by the disperser service) in aws dynamodb which is a nosql database. The update of the metadata in the dynamodb is monitored by the Batcher service to do further process.

Note that the disperser is not responsible for encoding/decoding of the blob.

