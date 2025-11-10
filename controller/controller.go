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
	ReadOnly  bool
}

func New(ctx context.Context, cfg Config) (*Controller, error) {
	loadOptions := []func(*config.LoadOptions) error{}

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

	if cfg.Profile != "" {
		loadOptions = append(loadOptions, config.WithSharedConfigProfile(cfg.Profile))
	}
	if cfg.Region != "" {
		loadOptions = append(loadOptions, config.WithRegion(cfg.Region))
	}
	if cfg.Endpoint != "" {
		loadOptions = append(loadOptions, config.WithBaseEndpoint(cfg.Endpoint))
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

	clientOptions = append(clientOptions, func(o *s3.Options) {
		transport := &Transport{
			Base:     http.DefaultTransport,
			Insecure: cfg.Insecure,
			ReadOnly: cfg.ReadOnly,
		}

		o.HTTPClient = &http.Client{
			Transport: transport,
		}
	})

	return &Controller{
		ctx:       ctx,
		verbosity: cfg.Verbosity,
		client:    s3.NewFromConfig(awsConfig, clientOptions...),
	}, nil
}

type Transport struct {
	Base     http.RoundTripper
	ReadOnly bool
	Insecure bool
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.ReadOnly {
		switch req.Method {
		case http.MethodDelete, http.MethodPatch, http.MethodPut, http.MethodPost:
			// TODO: disable retries
			return nil, fmt.Errorf("read-only mode: blocked %s %s", req.Method, req.URL)
		}
	}

	if t.Base == nil {
		t.Base = http.DefaultTransport
	}

	if t.Insecure {
		if tr, ok := t.Base.(*http.Transport); ok {
			cloned := tr.Clone()
			if cloned.TLSClientConfig == nil {
				cloned.TLSClientConfig = &tls.Config{}
			}
			cloned.TLSClientConfig.InsecureSkipVerify = true
			t.Base = cloned
		}
	}

	return t.Base.RoundTrip(req)
}
