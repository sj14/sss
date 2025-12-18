package controller

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/logging"
	smithymiddleware "github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/dustin/go-humanize"
	"github.com/sj14/sss/util"
	"github.com/sj14/sss/util/ratelimiter"
	"golang.org/x/time/rate"
)

type Controller struct {
	ctx       context.Context
	OutWriter io.Writer
	ErrWriter io.Writer
	client    *s3.Client
	verbosity uint8
}

type ControllerConfig struct {
	OutWriter io.Writer
	ErrWriter io.Writer
	Profile   Profile
	Verbosity uint8
	Headers   map[string]string
	Params    map[string]string
	DryRun    bool
	BuildInfo util.BuildInfo
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
	Network   string `toml:"network"`
	Bandwidth string `toml:"bandwidth"`
}

func New(ctx context.Context, cfg ControllerConfig) (*Controller, error) {
	if cfg.Verbosity > 0 && cfg.Profile.ReadOnly {
		fmt.Fprintln(cfg.OutWriter, "> read-only mode <")
	}

	if cfg.DryRun {
		// additional protection when dry-run is enabled
		cfg.Profile.ReadOnly = true

		if cfg.Verbosity > 0 {
			fmt.Fprintln(cfg.OutWriter, "> dry-run mode <")
		}
	}

	clientOptions := []func(o *s3.Options){
		func(o *s3.Options) { o.UsePathStyle = cfg.Profile.PathStyle },
		func(o *s3.Options) {
			o.APIOptions = append(o.APIOptions,
				middleware.AddUserAgentKeyValue("sss", cfg.BuildInfo.Version),
			)
		},
		func(o *s3.Options) {
			if len(cfg.Headers) == 0 {
				return
			}
			o.APIOptions = append(o.APIOptions,
				func(stack *smithymiddleware.Stack) error {
					return stack.Serialize.Insert(
						&AddHeadersMiddleware{
							Headers: cfg.Headers,
						},
						"OperationSerializer",
						smithymiddleware.Before,
					)
				},
			)
		},
		func(o *s3.Options) {
			if len(cfg.Params) == 0 {
				return
			}
			o.APIOptions = append(o.APIOptions,
				func(stack *smithymiddleware.Stack) error {
					return stack.Serialize.Insert(
						&AddParamsMiddleware{
							Params: cfg.Params,
						},
						"OperationSerializer",
						smithymiddleware.Before,
					)
				},
			)
		},
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

		baseTransport.DialContext = func(ctx context.Context, _, addr string) (net.Conn, error) {
			dialer := net.Dialer{}
			return dialer.DialContext(ctx, cfg.Profile.Network, addr)
		}

		transportWrapper := &TransportWrapper{
			Base:     baseTransport,
			ReadOnly: cfg.Profile.ReadOnly,
		}

		if cfg.Profile.Bandwidth != "" {
			bandwidth, err := humanize.ParseBytes(cfg.Profile.Bandwidth)
			if err != nil {
				log.Fatalln(err)
			}

			transportWrapper.Limiter = rate.NewLimiter(
				rate.Limit(bandwidth),
				64*1024, // add a small burst, otherwise it might fail
			)
		}

		o.HTTPClient = &http.Client{
			Transport: transportWrapper,
		}
	})

	return &Controller{
		ctx:       ctx,
		OutWriter: cfg.OutWriter,
		ErrWriter: cfg.ErrWriter,
		verbosity: cfg.Verbosity,
		client:    s3.NewFromConfig(awsCfg, clientOptions...),
	}, nil
}

type TransportWrapper struct {
	Base     http.RoundTripper
	ReadOnly bool
	Limiter  *rate.Limiter
}

func (t *TransportWrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.ReadOnly {
		switch req.Method {
		case http.MethodHead, http.MethodGet, http.MethodOptions, http.MethodTrace:
		default:
			return nil, fmt.Errorf("blocked by read-only mode")
		}
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

type AddHeadersMiddleware struct {
	Headers map[string]string
}

func (m *AddHeadersMiddleware) ID() string {
	return "AddCustomHTTPHeaders"
}

func (m *AddHeadersMiddleware) HandleSerialize(
	ctx context.Context,
	in smithymiddleware.SerializeInput,
	next smithymiddleware.SerializeHandler,
) (
	smithymiddleware.SerializeOutput,
	smithymiddleware.Metadata,
	error,
) {
	if req, ok := in.Request.(*smithyhttp.Request); ok {
		for k, v := range m.Headers {
			req.Header.Set(k, v)
		}
	}
	return next.HandleSerialize(ctx, in)
}

type AddParamsMiddleware struct {
	Params map[string]string
}

func (m *AddParamsMiddleware) ID() string {
	return "AddCustomURLParameter"
}

func (m *AddParamsMiddleware) HandleSerialize(
	ctx context.Context,
	in smithymiddleware.SerializeInput,
	next smithymiddleware.SerializeHandler,
) (
	smithymiddleware.SerializeOutput,
	smithymiddleware.Metadata,
	error,
) {
	if req, ok := in.Request.(*smithyhttp.Request); ok {
		q := req.URL.Query()
		for k, v := range m.Params {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
	return next.HandleSerialize(ctx, in)
}
