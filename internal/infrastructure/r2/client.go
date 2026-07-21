package r2

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	S3			*s3.Client
	Presign 	*s3.PresignClient
	Bucket		string
}

func New(ctx context.Context, endpoint, accessKey, secretKey, bucket string) (*Client, error) {
	resolver := aws.EndpointResolverWithOptionsFunc(
        func(service, region string, options ...interface{}) (aws.Endpoint, error) {
            return aws.Endpoint{
                URL: endpoint,
            }, nil
        })

    cfg, err := config.LoadDefaultConfig(ctx,
        config.WithEndpointResolverWithOptions(resolver),
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
        config.WithRegion("auto"),
    )
    if err != nil {
        return nil, err
    }

    s3Client := s3.NewFromConfig(cfg)

    return &Client{
        S3:      s3Client,
        Presign: s3.NewPresignClient(s3Client),
        Bucket:  bucket,
    }, nil
}

func (c *Client) PresignGet(ctx context.Context, key string, ttl time.Duration) (string, error) {
	req, err := c.Presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (c *Client) DeleteObject(ctx context.Context, key string) error {
	_, err := c.S3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(key),
	})
	return err
}