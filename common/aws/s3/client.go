package s3

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/0glabs/0g-da-client/common"
	commonaws "github.com/0glabs/0g-da-client/common/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var (
	once              sync.Once
	ref               *Client
	ErrObjectNotFound = errors.New("object not found")
)

type Object struct {
	Key  string
	Size int64
}

type Client struct {
	s3Client *s3.Client
	logger   common.Logger
}

func NewClient(cfg commonaws.ClientConfig, logger common.Logger) (*Client, error) {
	var err error
	logger.Info("url", cfg.EndpointURL)
	once.Do(func() {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if cfg.EndpointURL != "" {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           cfg.EndpointURL,
					SigningRegion: cfg.Region,
				}, nil
			}
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})
		options := [](func(*config.LoadOptions) error){
			config.WithRegion(cfg.Region),
			config.WithEndpointResolverWithOptions(customResolver),
			config.WithRetryMode(aws.RetryModeStandard),
		}
		// If access key and secret access key are not provided, use the default credential provider
		if len(cfg.AccessKey) > 0 && len(cfg.SecretAccessKey) > 0 {
			options = append(options, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretAccessKey, "")))
		}
		awsConfig, errCfg := config.LoadDefaultConfig(context.Background(), options...)
		if errCfg != nil {
			err = errCfg
			return
		}
		s3Client := s3.NewFromConfig(awsConfig, func(o *s3.Options) {
			o.UsePathStyle = true
		})
		ref = &Client{s3Client: s3Client, logger: logger}
	})
	return ref, err
}

func (s *Client) DownloadObject(ctx context.Context, bucket string, key string) ([]byte, error) {
	var partMiBs int64 = 10
	downloader := manager.NewDownloader(s.s3Client, func(d *manager.Downloader) {
		d.PartSize = partMiBs * 1024 * 1024 // 10MB per part
		d.Concurrency = 3                   //The number of goroutines to spin up in parallel per call to Upload when sending parts
	})

	buffer := manager.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	if buffer == nil || len(buffer.Bytes()) == 0 {
		return nil, ErrObjectNotFound
	}

	return buffer.Bytes(), nil
}

func (s *Client) keyExists(ctx context.Context, bucket string, key string) (bool, error) {
	_, err := s.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var responseError *awshttp.ResponseError
		if errors.As(err, &responseError) && responseError.ResponseError.HTTPStatusCode() == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *Client) UploadObject(ctx context.Context, bucket string, key string, data []byte) error {
	var partMiBs int64 = 10
	uploaded, _ := s.keyExists(ctx, bucket, key)
	if uploaded {
		s.logger.Info("object already uploaded, skip", "key", key)
		return nil
	}
	uploader := manager.NewUploader(s.s3Client, func(u *manager.Uploader) {
		u.PartSize = partMiBs * 1024 * 1024 // 10MB per part
		u.Concurrency = 3                   //The number of goroutines to spin up in parallel per call to Upload when sending parts
	})

	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Client) DeleteObject(ctx context.Context, bucket string, key string) error {
	_, err := s.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return err
}

func (s *Client) ListObjects(ctx context.Context, bucket string, prefix string) ([]Object, error) {
	output, err := s.s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, err
	}

	objects := make([]Object, 0, len(output.Contents))
	for _, object := range output.Contents {
		objects = append(objects, Object{
			Key:  *object.Key,
			Size: object.Size,
		})
	}
	return objects, nil
}

func (s *Client) CreateBucket(ctx context.Context, name, region string) error {
	_, err := s.s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})

	return err
}
