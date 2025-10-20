package controller

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/logging"
)

type Controller struct {
	ctx       context.Context
	client    *s3.Client
	verbosity uint8
}

type Config struct {
	Profile   string
	Endpoint  string
	Region    string
	PathStyle bool
	AccessKey string
	SecretKey string
	Verbosity uint8
	Insecure  bool
}

func New(ctx context.Context, cfg Config) (*Controller, error) {
	loadOptions := []func(*config.LoadOptions) error{
		config.WithSharedConfigProfile(cfg.Profile),
	}

	if cfg.Verbosity >= 9 {
		loadOptions = append(loadOptions,
			config.WithLogger(logging.NewStandardLogger(os.Stdout)),
			config.WithClientLogMode(
				aws.LogRequestWithBody|
					aws.LogResponseWithBody|
					aws.LogRetries|
					aws.LogDeprecatedUsage|
					aws.LogSigning|
					aws.LogRequestEventMessage|
					aws.LogResponseEventMessage,
			))
	} else if cfg.Verbosity >= 8 {
		loadOptions = append(loadOptions,
			config.WithLogger(logging.NewStandardLogger(os.Stdout)),
			config.WithClientLogMode(
				aws.LogRequest|
					aws.LogResponse|
					aws.LogRetries,
			))
	}

	awsConfig, err := config.LoadDefaultConfig(ctx,
		loadOptions...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	clientOptions := []func(o *s3.Options){
		func(o *s3.Options) { o.UsePathStyle = cfg.PathStyle },
	}

	if cfg.AccessKey != "" || cfg.SecretKey != "" {
		clientOptions = append(clientOptions, func(o *s3.Options) {
			o.Credentials = aws.NewCredentialsCache(
				credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
			)
		})
	}

	if cfg.Endpoint != "" {
		clientOptions = append(clientOptions, func(o *s3.Options) { o.BaseEndpoint = &cfg.Endpoint })
	}

	if cfg.Insecure {
		clientOptions = append(clientOptions, func(o *s3.Options) {
			customTransport := http.DefaultTransport.(*http.Transport).Clone()
			customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

			o.HTTPClient = &http.Client{
				Transport: customTransport,
			}
		})
	}

	if cfg.Region != "" {
		clientOptions = append(clientOptions, func(o *s3.Options) { o.Region = cfg.Region })
	}

	return &Controller{
		ctx:       ctx,
		verbosity: cfg.Verbosity,
		client:    s3.NewFromConfig(awsConfig, clientOptions...),
	}, nil
}
