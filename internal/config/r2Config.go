package config

import (
	"fmt"
	"os"
	"velo-server/internal/model"
)

func Load() (*model.Config, error) {
	cfg := &model.Config{
		R2AccountID:       os.Getenv("R2_ACCOUNT_ID"),
		R2AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
		R2SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
		R2BucketName:      os.Getenv("R2_BUCKET_NAME"),
		R2Endpoint:        os.Getenv("R2_ENDPOINT"),
	}

	if cfg.R2AccessKeyID == "" ||
		cfg.R2SecretAccessKey == "" ||
		cfg.R2BucketName == "" ||
		cfg.R2Endpoint == "" {
		return nil, fmt.Errorf("missing required R2 environment variables")
	}

	return cfg, nil
}
