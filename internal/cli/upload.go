package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/devusSs/minyls/internal/log"
	"github.com/devusSs/minyls/internal/minio"
	"github.com/devusSs/minyls/internal/yourls"
)

func Upload() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := initialize()
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	log.Log().Debug().Str("func", "cli.Upload").Msg("initialized")

	if len(os.Args) != neededUploadArgsLen {
		return fmt.Errorf("expected %d arguments, got %d", neededUploadArgsLen, len(os.Args))
	}

	log.Log().Debug().Str("func", "cli.Upload").Int("len_args", len(os.Args)).Msg("got os args")

	fp, err := getUploadFilePath()
	if err != nil {
		return fmt.Errorf("could not get upload file path: %w", err)
	}

	log.Log().Info().Str("func", "cli.Upload").Str("file_path", fp).Msg("got file path")

	p, err := getUploadPolicy()
	if err != nil {
		return fmt.Errorf("could not get upload policy: %w", err)
	}

	log.Log().Info().Str("func", "cli.Upload").Str("policy", p).Msg("got policy")

	mc, err := minio.NewClient(e.MinioEndpoint, e.MinioAccessKey, e.MinioAccessSecret)
	if err != nil {
		return fmt.Errorf("could not create minio client: %w", err)
	}

	log.Log().
		Debug().
		Str("func", "cli.Upload").
		Str("endpoint", e.MinioEndpoint).
		Str("access_key", e.MinioAccessKey).
		Str("access_secret", e.MinioAccessSecret).
		Msg("created minio client")

	err = mc.Setup(ctx, e.MinioBucketName, e.MinioRegion)
	if err != nil {
		return fmt.Errorf("could not setup minio client: %w", err)
	}

	log.Log().
		Debug().
		Str("func", "cli.Upload").
		Str("bucket_name", e.MinioBucketName).
		Str("region", e.MinioRegion).
		Msg("setup minio client")

	minioLink, err := mc.Upload(ctx, fp, p == "public", e.MinioLinkExpiry)
	if err != nil {
		return fmt.Errorf("could not upload file to minio: %w", err)
	}

	log.Log().
		Info().
		Str("func", "cli.Upload").
		Str("minio_link'", minioLink).
		Msg("got minio presigned url")

	yc := yourls.NewClient(e.YOURLSEndpoint, e.YOURLSSignature)
	link, err := yc.Shorten(ctx, minioLink, e.YOURLSTitle)
	if err != nil {
		return fmt.Errorf("could not shorten url: %w", err)
	}

	log.Log().
		Info().
		Str("func", "cli.Upload").
		Str("yourls_link", link).
		Msg("got shortened yourls link")

	// TODO: storage & clip

	return nil
}

const neededUploadArgsLen = 4

func getUploadFilePath() (string, error) {
	fp := os.Args[2]
	if fp == "" {
		return "", errors.New("empty filepath provided")
	}

	_, err := os.Stat(fp)
	if err != nil {
		return "", fmt.Errorf("file at path '%s' could not be found", fp)
	}

	return fp, nil
}

func getUploadPolicy() (string, error) {
	p := os.Args[3]
	if p != "public" && p != "private" {
		return "", fmt.Errorf("unexpected policy '%s' provided (expected 'public' or 'private')", p)
	}

	return p, nil
}
