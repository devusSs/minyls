package minio

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
)

// Upload uploads the specified file to either
// the public or private bucket and creates a share link
// with the specified expiry.
func (c *Client) Upload(
	ctx context.Context,
	filePath string,
	public bool,
	expiry time.Duration,
) (string, error) {
	if c.bucketPublic == "" || c.bucketPrivate == "" {
		return "", errors.New("buckets not setup, run Setup() first")
	}

	bucketName := c.bucketPrivate
	if public {
		bucketName = c.bucketPublic
	}

	fn, err := randomizeFileName(filePath)
	if err != nil {
		return "", fmt.Errorf("could not randomize file name: %w", err)
	}

	ct, err := findContentType(filePath)
	if err != nil {
		return "", fmt.Errorf("could not find content type: %w", err)
	}

	info, err := c.client.FPutObject(
		ctx,
		bucketName,
		fn,
		filePath,
		minio.PutObjectOptions{ContentType: ct},
	)
	if err != nil {
		return "", fmt.Errorf("could not fput object: %w", err)
	}

	// public object do not need a presigned url
	// do we want one tho for content disposition purposes?
	if public {
		link := filepath.Join(c.client.EndpointURL().String(), bucketName, info.Key)
		return link, nil
	}

	// no content disposition so far
	link, err := c.client.PresignedGetObject(ctx, bucketName, info.Key, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("could not egt presigned url: %w", err)
	}

	return link.String(), nil
}
