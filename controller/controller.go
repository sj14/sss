package controller

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/logging"
	"github.com/sj14/sss/util/ratelimiter"
	"golang.org/x/time/rate"
)

type Controller struct {
	ctx       context.Context
	client    *s3.Client
	verbosity uint8
}

type ControllerConfig struct {
	Profile   Profile
	Verbosity uint8
	Headers   []string
	Bandwidth uint64
	DryRun    bool
}

type Config struct {
	Profiles map[string]Profile `toml:"profiles"`
}

type Profile struct {
	Endpoint  string `toml:"endpoint"`
	Region    string `toml:"region"`
	AccessKey string `toml:"access_key"`
	SecretKey string `toml:"secret_key"`
	PathStyle bool   `toml:"path_style"`
	Insecure  bool   `toml:"insecure"`
	ReadOnly  bool   `toml:"read_only"`
	SNI       string `toml:"sni"`
}

func New(ctx context.Context, cfg ControllerConfig) (*Controller, error) {
	if cfg.Verbosity > 0 && cfg.Profile.ReadOnly {
		fmt.Println("> read-only mode <")
	}

	if cfg.DryRun {
		// additional protection when dry-run is enabled
		cfg.Profile.ReadOnly = true

		if cfg.Verbosity > 0 {
			fmt.Println("> dry-run mode <")
		}
	}

	clientOptions := []func(o *s3.Options){
		func(o *s3.Options) { o.UsePathStyle = cfg.Profile.PathStyle },
	}

	awsCfg := aws.Config{
		Region:       cfg.Profile.Region,
		BaseEndpoint: &cfg.Profile.Endpoint,
		Credentials: aws.NewCredentialsCache(
			credentials.NewStaticCredentialsProvider(cfg.Profile.AccessKey, cfg.Profile.SecretKey, ""),
		),
	}

	if cfg.Profile.ReadOnly {
		awsCfg.RetryMaxAttempts = 1
	}

	if cfg.Verbosity >= 9 {
		awsCfg.Logger = logging.NewStandardLogger(os.Stdout)
		awsCfg.ClientLogMode = aws.LogRequestWithBody |
			aws.LogResponseWithBody |
			aws.LogRetries |
			aws.LogDeprecatedUsage |
			aws.LogSigning |
			aws.LogRequestEventMessage |
			aws.LogResponseEventMessage
	} else if cfg.Verbosity >= 8 {
		awsCfg.Logger = logging.NewStandardLogger(os.Stdout)
		awsCfg.ClientLogMode = aws.LogRequest |
			aws.LogResponse |
			aws.LogRetries
	}

	clientOptions = append(clientOptions, func(o *s3.Options) {
		baseTransport := http.DefaultTransport.(*http.Transport).Clone()

		if baseTransport.TLSClientConfig == nil {
			baseTransport.TLSClientConfig = &tls.Config{}
		}
		if cfg.Profile.Insecure {
			baseTransport.TLSClientConfig.InsecureSkipVerify = true
		}
		if cfg.Profile.SNI != "" {
			baseTransport.TLSClientConfig.ServerName = cfg.Profile.SNI
		}

		transportWrapper := &TransportWrapper{
			Base:     baseTransport,
			ReadOnly: cfg.Profile.ReadOnly,
			Headers:  cfg.Headers,
		}

		if cfg.Bandwidth > 0 {
			transportWrapper.Limiter = rate.NewLimiter(
				rate.Limit(cfg.Bandwidth),
				128*1024, // add a small burst, otherwise it might fail
			)
		}

		o.HTTPClient = &http.Client{
			Transport: transportWrapper,
		}
	})

	return &Controller{
		ctx:       ctx,
		verbosity: cfg.Verbosity,
		client:    s3.NewFromConfig(awsCfg, clientOptions...),
	}, nil
}

type TransportWrapper struct {
	Base     http.RoundTripper
	ReadOnly bool
	Headers  []string
	Limiter  *rate.Limiter
}

func (t *TransportWrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.ReadOnly {
		switch req.Method {
		case http.MethodDelete, http.MethodPatch, http.MethodPut, http.MethodPost, http.MethodConnect:
			return nil, fmt.Errorf("blocked by read-only mode")
		}
	}

	for _, header := range t.Headers {
		s := strings.Split(header, ":")
		if len(s) != 2 {
			return nil, fmt.Errorf("failed to parse header %q", header)
		}

		req.Header.Set(s[0], s[1])
	}

	if req.Body != nil && t.Limiter != nil {
		req.Body = io.NopCloser(ratelimiter.NewReader(req.Context(), req.Body, t.Limiter))
	}

	resp, err := t.Base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil && t.Limiter != nil {
		resp.Body = io.NopCloser(ratelimiter.NewReader(req.Context(), resp.Body, t.Limiter))
	}

	return resp, nil
}
