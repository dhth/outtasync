package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/dhth/outtasync/internal/types"

	"github.com/aws/aws-sdk-go-v2/config"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

const (
	describeDriftSleepMillis = 3000
)

type CFClient struct {
	Client *cf.Client
	Err    error
}

func GetAWSConfig(source types.ConfigSource) (aws.Config, error) {
	var cfg aws.Config
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	switch source.Kind {
	case types.Env:
		cfg, err = config.LoadDefaultConfig(ctx)
	case types.SharedProfile:
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithSharedConfigProfile(source.Value))
	case types.AssumeRole:
		cfg, err = config.LoadDefaultConfig(ctx)
		if err != nil {
			return cfg, err
		}
		stsSvc := sts.NewFromConfig(cfg)
		creds := stscreds.NewAssumeRoleProvider(stsSvc, source.Value)
		cfg.Credentials = aws.NewCredentialsCache(creds)
	}
	return cfg, err
}
