package env

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Env struct {
	MinioEndpoint     string        `env:"MINIO_ENDPOINT"`
	MinioAccessKey    string        `env:"MINIO_ACCESS_KEY"`
	MinioAccessSecret string        `env:"MINIO_ACCESS_SECRET"`
	MinioBucketName   string        `env:"MINIO_BUCKET_NAME"   envDefault:"minyls"`
	MinioRegion       string        `env:"MINIO_REGION"        envDefault:"us-east-1"`
	MinioLinkExpiry   time.Duration `env:"MINIO_LINK_EXPIRY"   envDefault:"168h"`
	YOURLSEndpoint    string        `env:"YOURLS_ENDPOINT"`
	YOURLSSignature   string        `env:"YOURLS_SIGNATURE"`
	YOURLSTitle       string        `env:"YOURLS_TITLE"        envDefault:"shortened using minyls"`
}

func Load() (*Env, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("could not get executable: %w", err)
	}

	// ignore the error since we do not actually care
	// if the env file was loaded or not, it simply eases
	// up the process for the user
	_ = godotenv.Load(filepath.Join(filepath.Dir(exe), ".minyls.env"))

	e := &Env{}
	err = env.ParseWithOptions(e, env.Options{
		RequiredIfNoDef: true,
		Prefix:          "MINYLS_",
	})
	if err != nil {
		return nil, fmt.Errorf("could not parse env: %w", err)
	}

	return e, nil
}
