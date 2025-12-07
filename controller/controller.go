package controller

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
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
	Profiles map[string]Profile
}

type Profile struct {
	Endpoint  string `yaml:"endpoint"`
	Region    string `yaml:"region"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	PathStyle bool   `yaml:"path_style"`
	Insecure  bool   `yaml:"insecure"`
	ReadOnly  bool   `yaml:"read_only"`
}

func New(ctx context.Context, verbosity uint8, cfg Profile) (*Controller, error) {
	clientOptions := []func(o *s3.Options){
		func(o *s3.Options) { o.UsePathStyle = cfg.PathStyle },
	}

	awsCfg := aws.Config{
		Region:       cfg.Region,
		BaseEndpoint: &cfg.Endpoint,
		Credentials: aws.NewCredentialsCache(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		),
	}

	if verbosity >= 9 {
		awsCfg.Logger = logging.NewStandardLogger(os.Stdout)
		awsCfg.ClientLogMode = aws.LogRequestWithBody |
			aws.LogResponseWithBody |
			aws.LogRetries |
			aws.LogDeprecatedUsage |
			aws.LogSigning |
			aws.LogRequestEventMessage |
			aws.LogResponseEventMessage
	} else if verbosity >= 8 {
		awsCfg.Logger = logging.NewStandardLogger(os.Stdout)
		awsCfg.ClientLogMode = aws.LogRequest |
			aws.LogResponse |
			aws.LogRetries
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
		verbosity: verbosity,
		client:    s3.NewFromConfig(awsCfg, clientOptions...),
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
