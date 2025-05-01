package minio

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	client        *minio.Client
	bucketPublic  string
	bucketPrivate string
}

func NewClient(endpoint string, accessKey string, accessSecret string) (*Client, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint url: %w", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("endpoint scheme must be 'http://' or 'https://'")
	}

	host := strings.TrimPrefix(endpoint, u.Scheme+"://")
	if host == "" {
		return nil, errors.New("invalid endpoint format: missing host")
	}

	c, err := minio.New(host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, accessSecret, ""),
		Secure: u.Scheme == "https",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	return &Client{client: c}, nil
}

// Setup will create the needed bucket(s). The parameter bucketName will
// be appended with '-private' and '-public' to separate them.
//
// The policy for the public bucket will also be set.
func (c *Client) Setup(
	ctx context.Context,
	bucketName string,
	bucketRegion string,
) error {
	var err error

	buckets := []string{bucketName + "-public", bucketName + "-private"}
	for _, bucket := range buckets {
		err = c.makeBucket(ctx, bucket, bucketRegion)
		if err != nil {
			return fmt.Errorf("could not create bucket '%s': %w", bucket, err)
		}
	}

	c.bucketPublic = buckets[0]
	c.bucketPrivate = buckets[1]

	err = c.setBucketPolicy(
		ctx,
		c.bucketPublic,
		fmt.Sprintf(bucketPublicPolicyTemplate, c.bucketPublic),
	)
	if err != nil {
		return fmt.Errorf("could not set public bucket policy: %w", err)
	}

	return nil
}

func (c *Client) makeBucket(ctx context.Context, bucketName string, bucketRegion string) error {
	exists, err := c.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("could not check if bucket exists: %w", err)
	}

	if exists {
		return nil
	}

	return c.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
		Region: bucketRegion,
		// ObjectLocking needs to be set to false,
		// else we cannot delete the objects using this client (for now).
		ObjectLocking: false,
	})
}

const bucketPublicPolicyTemplate = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": "*",
      "Action": [
        "s3:GetObject"
      ],
      "Resource": "arn:aws:s3:::%s/*"
    }
  ]
}`

func (c *Client) setBucketPolicy(ctx context.Context, bucketName string, policy string) error {
	return c.client.SetBucketPolicy(ctx, bucketName, policy)
}
