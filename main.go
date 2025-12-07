package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"slices"

	"github.com/BurntSushi/toml"
	"github.com/dustin/go-humanize"
	"github.com/sj14/sss/controller"
	"github.com/sj14/sss/util"
	docs "github.com/urfave/cli-docs/v3"
	"github.com/urfave/cli/v3"
)

var (
	// will be replaced during the build process
	version = "undefined"
	commit  = "undefined"
	date    = "undefined"
)

func main() {
	cmd := &cli.Command{
		Name:                  "sss",
		Usage:                 "S3 client",
		Version:               fmt.Sprintf("%s %s %s", version, commit, date),
		EnableShellCompletion: true,
		ConfigureShellCompletionCommand: func(c *cli.Command) {
			c.Hidden = false
			c.Before = func(ctx context.Context, c *cli.Command) (context.Context, error) {
				flagBucket.Required = false
				return ctx, nil
			}
		},
		Flags: []cli.Flag{
			// order will appear the same in the help text
			flagConfig,
			flagProfile,
			flagAccessKey,
			flagSecretKey,
			flagEndpoint,
			flagRegion,
			flagPathStyle,
			flagInsecure,
			flagBucket,
			flagReadOnly,
			flagBandwidth,
			flagSNI,
			flagHeaders,
			flagVerbosity,
		},
		Commands: []*cli.Command{
			// order will appear the same in the help text
			cmdDocs,
			cmdListProfiles,
			cmdBucketList,
			cmdBucketHead,
			cmdBucketMake,
			cmdBucketRemove,
			cmdBucketSize,
			cmdBucketPolicy,
			cmdBucketVersioning,
			cmdBucketObjectLock,
			cmdBucketLifecycle,
			cmdBucketCors,
			cmdBucketTag,
			cmdBucketMultiparts,
			cmdBucketParts,
			cmdObjectsList,
			cmdObjectHead,
			cmdObjectGet,
			cmdObjectPut,
			cmdObjectRemove,
			cmdObjectsCopy,
			cmdObjectVersions,
			cmdObjectACL,
			cmdObjectPresign,
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalln(err)
	}
}

func loadConfig(cmd *cli.Command) (controller.Config, error) {
	var config controller.Config

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return config, err
	}

	configPath := cmd.Root().String(flagConfig.Name)

	if configPath == "" {
		configPath = filepath.Join(homeDir, ".config", "sss", "config.toml")
	}

	md, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		return config, err
	}

	if undecoded := md.Undecoded(); len(undecoded) > 0 {
		return config, fmt.Errorf("unknown fields in config: %v", undecoded)
	}

	return config, nil
}

func exec(ctx context.Context, cmd *cli.Command, fn func(ctrl *controller.Controller) error) error {
	config, err := loadConfig(cmd)
	// Do not return an error when the config file does not exist,
	// as the tool should be usable wihout a config file.
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("failed loading config: %w", err)
	}

	profileName := cmd.Root().String(flagProfile.Name)
	profile, ok := config.Profiles[profileName]
	if !ok && profileName != "default" {
		fmt.Printf("profile %q not found, available profiles:\n", profileName)

		keys := slices.Collect(maps.Keys(config.Profiles))
		slices.Sort(keys)

		for _, key := range keys {
			fmt.Println(key)
		}

		os.Exit(1)
	}

	util.SetIfNotZero(&profile.Endpoint, cmd.Root().String(flagEndpoint.Name))
	util.SetIfNotZero(&profile.Region, cmd.Root().String(flagRegion.Name))
	util.SetIfNotZero(&profile.PathStyle, cmd.Root().Bool(flagPathStyle.Name))
	util.SetIfNotZero(&profile.AccessKey, cmd.Root().String(flagAccessKey.Name))
	util.SetIfNotZero(&profile.SecretKey, cmd.Root().String(flagSecretKey.Name))
	util.SetIfNotZero(&profile.Insecure, cmd.Root().Bool(flagInsecure.Name))
	util.SetIfNotZero(&profile.ReadOnly, cmd.Root().Bool(flagReadOnly.Name))
	util.SetIfNotZero(&profile.SNI, cmd.Root().String(flagSNI.Name))

	var bandwidth uint64 = 0
	{
		bwStr := cmd.Root().String(flagBandwidth.Name)
		if bwStr != "" {
			bandwidth, err = humanize.ParseBytes(bwStr)
			if err != nil {
				return err
			}
		}
	}

	ctrl, err := controller.New(ctx, controller.ControllerConfig{
		Verbosity: cmd.Root().Uint8(flagVerbosity.Name),
		Headers:   cmd.Root().StringSlice(flagHeaders.Name),
		Profile:   profile,
		Bandwidth: bandwidth,
	})
	if err != nil {
		return err
	}
	return fn(ctrl)
}

func parseSSEC(cmd *cli.Command) util.SSEC {
	return util.NewSSEC(
		cmd.String(flagSSEcAlgo.Name),
		cmd.String(flagSSEcKey.Name),
	)
}

var (
	argPrefix = &cli.StringArg{
		Name: "prefix",
	}
	argKey = &cli.StringArg{
		Name: "key",
	}
	argConfigPath = &cli.StringArg{
		Name:      "config",
		UsageText: "Path to the config",
	}

	flagConfig = &cli.StringFlag{
		Name:      "config",
		Usage:     "~/.config/sss/config.toml",
		Sources:   cli.EnvVars("SSS_CONFIG"),
		TakesFile: true,
	}
	flagVerbosity = &cli.Uint8Flag{
		Name:    "verbosity",
		Value:   1,
		Sources: cli.EnvVars("SSS_VERBOSITY"),
	}
	flagBandwidth = &cli.StringFlag{
		Name:    "bandwidth",
		Usage:   "Limit the bandwith per second (e.g. '1 MiB'). If set, an initial burst of 128 KiB is added.",
		Sources: cli.EnvVars("SSS_BANDWIDTH"),
	}
	flagProfile = &cli.StringFlag{
		Name:    "profile",
		Value:   "default",
		Sources: cli.EnvVars("SSS_PROFILE"),
	}
	flagPathStyle = &cli.BoolFlag{
		Name:    "path-style",
		Sources: cli.EnvVars("SSS_PATH_STYLE"),
	}
	flagEndpoint = &cli.StringFlag{
		Name:    "endpoint",
		Sources: cli.EnvVars("SSS_ENDPOINT"),
	}
	flagInsecure = &cli.BoolFlag{
		Name:    "insecure",
		Sources: cli.EnvVars("SSS_INSECURE"),
	}
	flagReadOnly = &cli.BoolFlag{
		Name:    "read-only",
		Sources: cli.EnvVars("SSS_READ_ONLY"),
	}
	flagSNI = &cli.StringFlag{
		Name:    "sni",
		Sources: cli.EnvVars("SSS_SNI"),
	}
	flagRegion = &cli.StringFlag{
		Name:    "region",
		Sources: cli.EnvVars("SSS_REGION"),
	}
	flagAccessKey = &cli.StringFlag{
		Name:    "access-key",
		Sources: cli.EnvVars("SSS_ACCESS_KEY"),
	}
	flagSecretKey = &cli.StringFlag{
		Name:    "secret-key",
		Sources: cli.EnvVars("SSS_SECRET_KEY"),
	}
	flagBucket = &cli.StringFlag{
		Name:     "bucket",
		Required: true,
		Sources:  cli.EnvVars("SSS_BUCKET"),
	}
	flagHeaders = &cli.StringSliceFlag{
		Name:  "header",
		Usage: "format: 'key1:val1,key2:val2'",
	}
	flagObjectLock = &cli.BoolFlag{
		Name: "object-lock",
	}
	flagPrefix = &cli.StringFlag{
		Name: "prefix",
	}
	flagSSEcKey = &cli.StringFlag{
		Name:  "sse-c-key",
		Usage: "32 bytes key",
	}
	flagSSEcAlgo = &cli.StringFlag{
		Name:  "sse-c-algorithm",
		Value: "AES256",
	}
	flagDelimiter = &cli.StringFlag{
		Name:  "delimiter",
		Value: "/",
	}
	flagForce = &cli.BoolFlag{
		Name: "force",
	}
	flagConcurrency = &cli.IntFlag{
		Name:  "concurrency",
		Value: 5,
	}
	flagPartSize = &cli.Int64Flag{
		Name: "part-size",
	}
	flagObjectKey = &cli.StringFlag{
		Name: "key",
	}
	flagUploadID = &cli.StringFlag{
		Name: "upload-id",
	}
	flagVersionID = &cli.StringFlag{
		Name: "version-id",
	}
	flagRange = &cli.StringFlag{
		Name:  "range",
		Usage: "bytes=BeginByte-EndByte, e.g. 'bytes=0-500' to get the first 501 bytes",
	}
	flagPartNumber = &cli.Int32Flag{
		Name: "part-number",
	}
	flagIfMatch = &cli.StringFlag{
		Name: "if-match",
	}
	flagIfNoneMatch = &cli.StringFlag{
		Name: "if-none-match",
	}
	flagIfModifiedSince = &cli.TimestampFlag{
		Name: "if-modified-since",
	}
	flagIfUnmodifiedSince = &cli.TimestampFlag{
		Name: "if-unmodified-since",
	}
)

var (
	cmdDocs = &cli.Command{
		Name:   "docs",
		Hidden: true,
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			flagBucket.Required = false
			return ctx, nil
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			s, err := docs.ToMarkdown(cmd.Root())
			if err != nil {
				return err
			}
			return os.WriteFile("DOCS.md", []byte(s), os.ModePerm)
		},
	}
	cmdListProfiles = &cli.Command{
		Name:  "profiles",
		Usage: "Config Profiles",
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			flagBucket.Required = false
			return ctx, nil
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			config, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			keys := slices.Collect(maps.Keys(config.Profiles))
			slices.Sort(keys)

			for _, key := range keys {
				fmt.Println(key)
			}

			return nil
		},
	}
	cmdBucketList = &cli.Command{
		Category: "bucket management",
		Name:     "buckets",
		Usage:    "Bucket List",
		Flags: []cli.Flag{
			flagPrefix,
		},
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			flagBucket.Required = false
			if flagBucket.IsSet() && flagVerbosity.Value > 0 {
				slog.InfoContext(ctx, "ignoring -bucket flag")
				fmt.Println()
			}
			return ctx, nil
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return exec(ctx, cmd, func(ctrl *controller.Controller) error {
				return ctrl.BucketList(cmd.String(flagPrefix.Name))
			})
		},
	}
	cmdBucketHead = &cli.Command{
		Category: "bucket management",
		Name:     "bucket",
		Usage:    "Bucket Head",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return exec(ctx, cmd, func(ctrl *controller.Controller) error {
				return ctrl.BucketHead(cmd.String(flagBucket.Name))
			})
		},
	}
	cmdBucketTag = &cli.Command{
		Category: "bucket management",
		Name:     "tag",
		Usage:    "Bucket Tagging",
		Commands: []*cli.Command{
			{
				Name: "get",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return exec(ctx, cmd, func(ctrl *controller.Controller) error {
						return ctrl.BucketTagging(cmd.String(flagBucket.Name))
					})
				},
			},
		},
	}
	cmdBucketMake = &cli.Command{
		Category: "bucket management",
		Name:     "mb",
		Usage:    "Bucket Create",
		Flags: []cli.Flag{
			flagObjectLock,
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return exec(ctx, cmd, func(ctrl *controller.Controller) error {
				return ctrl.BucketCreate(
					cmd.String(flagBucket.Name),
					cmd.Bool(flagObjectLock.Name),
				)
			})
		},
	}
	cmdBucketRemove = &cli.Command{
		Category: "bucket management",
		Name:     "rb",
		Usage:    "Bucket Remove",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return exec(ctx, cmd, func(ctrl *controller.Controller) error {
				return ctrl.BucketDelete(cmd.String(flagBucket.Name))
			})
		},
	}
	cmdBucketMultiparts = &cli.Command{
		Category: "multipart management",
		Name:     "multiparts",
		Usage:    "Multipart Uploads",
		Commands: []*cli.Command{
			{
				Name: "ls",
				Flags: []cli.Flag{
					flagPrefix,
					flagDelimiter,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return exec(ctx, cmd, func(ctrl *controller.Controller) error {
						return ctrl.BucketMultipartUploadsList(
							cmd.String(flagBucket.Name),
							cmd.String(flagPrefix.Name),
							cmd.String(flagDelimiter.Name),
						)
					})
				},
			},
			{
				Name: "rm",
				Flags: []cli.Flag{
					flagObjectKey,
					flagUploadID,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return exec(ctx, cmd, func(ctrl *controller.Controller) error {
						return ctrl.BucketMultipartUploadAbort(
							cmd.String(flagBucket.Name),
							cmd.String(flagObjectKey.Name),
							cmd.String(flagUploadID.Name),
						)
					})
				},
			},
		},
	}
	cmdBucketParts = &cli.Command{
		Category: "multipart management",
		Name:     "parts",
		Usage:    "Multipart Parts",
		Flags: []cli.Flag{
			flagObjectKey,
			flagUploadID,
		},
		Commands: []*cli.Command{
			{
				Name: "ls",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return exec(ctx, cmd, func(ctrl *controller.Controller) error {
						return ctrl.BucketPartsList(
							cmd.String(flagBucket.Name),
							cmd.String(flagObjectKey.Name),
							cmd.String(flagUploadID.Name),
						)
					})
				},
			},
		},
	}
	cmdObjectsList = &cli.Command{
		Category: "object management",
		Name:     "ls",
		Usage:    "Object List",
		Arguments: []cli.Argument{
			argPrefix,
		},
		Flags: []cli.Flag{
			flagDelimiter,
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return exec(ctx, cmd, func(ctrl *controller.Controller) error {
				return ctrl.ObjectList(
					cmd.String(flagBucket.Name),
					cmd.StringArg(argPrefix.Name),
					cmd.String(flagDelimiter.Name),
				)
			})
		},
	}
	cmdObjectsCopy = &cli.Command{
		Category: "object management",
		Name:     "cp",
		Usage:    "Object Server Side Copy",
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			flagBucket.Required = false
			if flagBucket.IsSet() && flagVerbosity.Value > 0 {
				slog.InfoContext(ctx, "ignoring -bucket flag")
				fmt.Println()
			}
			return ctx, nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "src-bucket",
				Usage:    "Source bucket",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "src-key",
				Usage:    "Source key",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "dst-bucket",
				Usage:    "Destinaton bucket",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "dst-key",
				Usage: "Destination key. When empty, the src-key will be used",
			},
			flagSSEcKey,
			flagSSEcAlgo,
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return exec(ctx, cmd, func(ctrl *controller.Controller) error {
				return ctrl.ObjectCopy(
					controller.ObjectCopyConfig{
						SrcBucket: cmd.String("src-bucket"),
						SrcKey:    cmd.String("src-key"),
						DstBucket: cmd.String("dst-bucket"),
						DstKey:    cmd.String("dst-key"),
						SSEC:      parseSSEC(cmd),
					})
			})
		},
	}
	cmdObjectPut = &cli.Command{
		Category: "object management",
		Name:     "put",
		Usage:    "Object Upload",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "path",
			},
		},
		Flags: []cli.Flag{
			flagSSEcKey,
			flagSSEcAlgo,
			flagPartSize,
			flagConcurrency,
			&cli.IntFlag{
				Name: "leave-parts-on-error",
			},
			&cli.IntFlag{
				Name: "max-parts",
			},
			&cli.StringFlag{
				Name:  "target",
				Usage: "target key for single file or prefix multiple files",
			},
			&cli.StringFlag{
				Name:  "acl",
				Usage: "e.g. 'public-read'",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return exec(ctx, cmd, func(ctrl *controller.Controller) error {
				return ctrl.ObjectPut(
					cmd.StringArg("path"),
					cmd.String("target"),
					controller.ObjectPutConfig{
						Bucket:            cmd.String(flagBucket.Name),
						SSEC:              parseSSEC(cmd),
						Concurrency:       cmd.Int("concurrency"),
						LeavePartsOnError: cmd.Bool("leave-parts-on-error"),
						MaxUploadParts:    cmd.Int32("max-parts"),
						PartSize:          cmd.Int64("part-size"),
						ACL:               cmd.String("acl"),
					},
				)
			})
		},
	}
	cmdObjectRemove = &cli.Command{
		Category: "object management",
		Name:     "rm",
		Usage:    "Object Remove",
		Arguments: []cli.Argument{
			argKey,
		},
		Flags: []cli.Flag{
			flagDelimiter,
			flagForce,
			flagConcurrency,
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return exec(ctx, cmd, func(ctrl *controller.Controller) error {
				return ctrl.ObjectDelete(
					cmd.StringArg("key"),
					controller.ObjectDeleteConfig{
						Bucket:      cmd.String(flagBucket.Name),
						Delimiter:   cmd.String(flagDelimiter.Name),
						Force:       cmd.Bool(flagForce.Name),
						Concurrency: cmd.Int(flagConcurrency.Name),
					},
				)
			})
		},
	}
	cmdObjectGet = &cli.Command{
		Category: "object management",
		Name:     "get",
		Usage:    "Object Download",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "key",
			},
			&cli.StringArg{
				Name: "target",
			},
		},
		Flags: []cli.Flag{
			flagSSEcKey,
			flagSSEcAlgo,
			flagConcurrency,
			flagPartSize,
			flagVersionID,
			flagRange,
			flagPartNumber,
			flagIfMatch,
			flagIfNoneMatch,
			flagIfModifiedSince,
			flagIfUnmodifiedSince,
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return exec(ctx, cmd, func(ctrl *controller.Controller) error {
				return ctrl.ObjectGet(
					cmd.StringArg("target"),
					cmd.String(flagDelimiter.Name),
					controller.ObjectGetConfig{
						Bucket:            cmd.String(flagBucket.Name),
						ObjectKey:         cmd.StringArg(argKey.Name),
						SSEC:              parseSSEC(cmd),
						VersionID:         cmd.String(flagVersionID.Name),
						Range:             cmd.String(flagRange.Name),
						PartNumber:        cmd.Int32(flagPartNumber.Name),
						Concurrency:       cmd.Int(flagConcurrency.Name),
						PartSize:          cmd.Int64(flagPartSize.Name),
						IfMatch:           cmd.String(flagIfMatch.Name),
						IfNoneMatch:       cmd.String(flagIfNoneMatch.Name),
						IfModifiedSince:   cmd.Timestamp(flagIfModifiedSince.Name),
						IfUnmodifiedSince: cmd.Timestamp(flagIfUnmodifiedSince.Name),
					},
				)
			})
		},
	}
	cmdObjectHead = &cli.Command{
		Category: "object management",
		Name:     "head",
		Usage:    "Object Head",
		Arguments: []cli.Argument{
			argKey,
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return exec(ctx, cmd, func(ctrl *controller.Controller) error {
				return ctrl.ObjectHead(
					cmd.String(flagBucket.Name),
					cmd.StringArg(argKey.Name),
				)
			})
		},
	}
	cmdObjectPresign = &cli.Command{
		Category: "object management",
		Name:     "presign",
		Usage:    "Object pre-signed URL",
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			if flagReadOnly.IsSet() {
				return ctx, errors.New("deactivated due to read only mode")
			}
			return ctx, nil
		},
		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name: "expires-in",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "get",
				Usage: "Presigned URL for a GET request",
				Arguments: []cli.Argument{
					argKey,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return exec(ctx, cmd, func(ctrl *controller.Controller) error {
						return ctrl.ObjectPresignGet(
							cmd.Duration("expires-in"),
							controller.ObjectGetConfig{
								Bucket:    cmd.String(flagBucket.Name),
								ObjectKey: cmd.StringArg(argKey.Name),
							},
						)
					})
				}},
			{
				Name:  "put",
				Usage: "Presigned URL for a PUT request",
				Arguments: []cli.Argument{
					argKey,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return exec(ctx, cmd, func(ctrl *controller.Controller) error {
						return ctrl.ObjectPresignPut(
							cmd.Duration("expires-in"),
							cmd.StringArg(argKey.Name),
							controller.ObjectPutConfig{
								Bucket: cmd.String(flagBucket.Name),
							},
						)
					})
				}},
		},
	}
	cmdBucketPolicy = &cli.Command{
		Category: "bucket management",
		Name:     "policy",
		Usage:    "Bucket Policy",
		Commands: []*cli.Command{
			{
				Name: "get",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return exec(ctx, cmd, func(ctrl *controller.Controller) error {
						return ctrl.BucketPolicyGet(
							cmd.String(flagBucket.Name),
						)
					})
				},
			},
			{
				Name: "put",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name: "policy",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return exec(ctx, cmd, func(ctrl *controller.Controller) error {
						return ctrl.BucketPolicyPut(
							cmd.StringArg("policy"),
							cmd.String(flagBucket.Name),
						)
					})
				},
			},
		},
	}
	cmdBucketCors = &cli.Command{
		Category: "bucket management",
		Name:     "cors",
		Usage:    "Bucket CORS",
		Commands: []*cli.Command{
			{
				Name: "get",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return exec(ctx, cmd, func(ctrl *controller.Controller) error {
						return ctrl.BucketCORSGet(
							cmd.String(flagBucket.Name),
						)
					})
				},
			},
			{
				Name: "put",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name: "cors",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return exec(ctx, cmd, func(ctrl *controller.Controller) error {
						return ctrl.BucketCORSPut(
							cmd.StringArg("cors"),
							cmd.String(flagBucket.Name),
						)
					})
				},
			},
			{
				Name: "rm",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return exec(ctx, cmd, func(ctrl *controller.Controller) error {
						return ctrl.BucketCORSDelete(
							cmd.String(flagBucket.Name),
						)
					})
				},
			},
		},
	}
	cmdBucketObjectLock = &cli.Command{
		Category: "bucket management",
		Name:     "object-lock",
		Usage:    "Bucket Object Locking",
		Commands: []*cli.Command{
			{
				Name: "get",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return exec(ctx, cmd, func(ctrl *controller.Controller) error {
						return ctrl.BucketObjectLockGet(
							cmd.String(flagBucket.Name),
						)
					})
				},
			},
			{
				Name: "put",
				Arguments: []cli.Argument{
					argConfigPath,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return exec(ctx, cmd, func(ctrl *controller.Controller) error {
						return ctrl.BucketObjectLockPut(
							cmd.StringArg(argConfigPath.Name),
							cmd.String(flagBucket.Name),
						)
					})
				},
			},
		},
	}
	cmdBucketLifecycle = &cli.Command{
		Category: "bucket management",
		Name:     "lifecycle",
		Usage:    "Bucket Lifecycle",
		Commands: []*cli.Command{
			{
				Name: "get",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return exec(ctx, cmd, func(ctrl *controller.Controller) error {
						return ctrl.BucketLifecycleGet(
							cmd.String(flagBucket.Name),
						)
					})
				},
			},
			{
				Name: "put",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name: "lifecycle",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return exec(ctx, cmd, func(ctrl *controller.Controller) error {
						return ctrl.BucketLifecyclePut(
							cmd.StringArg("lifecycle"),
							cmd.String(flagBucket.Name),
						)
					})
				},
			},
			{
				Name: "rm",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return exec(ctx, cmd, func(ctrl *controller.Controller) error {
						return ctrl.BucketLifecycleDelete(
							cmd.String(flagBucket.Name),
						)
					})
				},
			},
		},
	}
	cmdBucketVersioning = &cli.Command{
		Category: "bucket management",
		Name:     "versioning",
		Usage:    "Bucket Versioning",
		Commands: []*cli.Command{
			{
				Name: "get",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return exec(ctx, cmd, func(ctrl *controller.Controller) error {
						return ctrl.BucketVersioningGet(
							cmd.String(flagBucket.Name),
						)
					})
				},
			},
			{
				Name: "put",
				Arguments: []cli.Argument{
					argConfigPath,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return exec(ctx, cmd, func(ctrl *controller.Controller) error {
						return ctrl.BucketVersioningPut(
							cmd.StringArg(argConfigPath.Name),
							cmd.String(flagBucket.Name),
						)
					})
				},
			},
		},
	}
	cmdBucketSize = &cli.Command{
		Category: "bucket management",
		Name:     "size",
		Usage:    "Bucket Size",
		Arguments: []cli.Argument{
			argPrefix,
		},
		Flags: []cli.Flag{
			flagDelimiter,
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return exec(ctx, cmd, func(ctrl *controller.Controller) error {
				return ctrl.BucketSize(
					cmd.String(flagBucket.Name),
					cmd.StringArg(argPrefix.Name),
					cmd.String(flagDelimiter.Name),
				)
			})
		},
	}
	cmdObjectACL = &cli.Command{
		Category: "object management",
		Name:     "acl",
		Usage:    "Object ACL",
		Commands: []*cli.Command{
			{
				Name: "get",
				Arguments: []cli.Argument{
					argKey,
				},
				Flags: []cli.Flag{
					flagVersionID,
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return exec(ctx, cmd, func(ctrl *controller.Controller) error {
						return ctrl.ObjectACLGet(
							cmd.String(flagBucket.Name),
							cmd.StringArg(argKey.Name),
							cmd.String(flagVersionID.Name),
						)
					})
				},
			},
		},
	}
	cmdObjectVersions = &cli.Command{
		Category: "object management",
		Name:     "versions",
		Usage:    "Object Versions",
		Arguments: []cli.Argument{
			argPrefix,
		},
		Flags: []cli.Flag{
			flagDelimiter,
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return exec(ctx, cmd, func(ctrl *controller.Controller) error {
				return ctrl.ObjectVersions(
					cmd.String(flagBucket.Name),
					cmd.StringArg(argPrefix.Name),
					cmd.String(flagDelimiter.Name),
				)
			})
		},
	}
)
