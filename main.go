package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/sj14/sss/controller"
	"github.com/sj14/sss/util"
	docs "github.com/urfave/cli-docs/v3"
	"github.com/urfave/cli/v3"
)

func main() {
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalln(err)
	}
}

func Exec(ctx context.Context, cmd *cli.Command, fn func(ctrl *controller.Controller) error) error {
	ctrl, err := controller.New(ctx, controller.Config{
		Profile:   cmd.Root().String(flagProfile.Name),
		Endpoint:  cmd.Root().String(flagEndpoint.Name),
		Region:    cmd.Root().String(flagRegion.Name),
		PathStyle: cmd.Root().Bool(flagPathStyle.Name),
		AccessKey: cmd.Root().String(flagAccessKey.Name),
		SecretKey: cmd.Root().String(flagSecretKey.Name),
		Verbosity: cmd.Root().Uint8(flagVerbosity.Name),
		Insecure:  cmd.Root().Bool(flagInsecure.Name),
	})
	if err != nil {
		return err
	}
	return fn(ctrl)
}

var (
	argPrefix = &cli.StringArg{
		Name: "prefix",
	}
	argKey = &cli.StringArg{
		Name: "key",
	}

	flagVerbosity = &cli.Uint8Flag{
		Name:  "verbosity",
		Value: 1,
	}
	flagProfile = &cli.StringFlag{
		Name:  "profile",
		Value: "default",
	}
	flagPathStyle = &cli.BoolFlag{
		Name: "path-style",
	}
	flagEndpoint = &cli.StringFlag{
		Name: "endpoint",
	}
	flagInsecure = &cli.BoolFlag{
		Name: "insecure",
	}
	flagRegion = &cli.StringFlag{
		Name: "region",
	}
	flagAccessKey = &cli.StringFlag{
		Name: "access-key",
	}
	flagSecretKey = &cli.StringFlag{
		Name: "secret-key",
	}
	flagBucket = &cli.StringFlag{
		Name:     "bucket",
		Required: true,
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

func parseSSEC(cmd *cli.Command) util.SSEC {
	return util.NewSSEC(
		cmd.String(flagSSEcAlgo.Name),
		cmd.String(flagSSEcKey.Name),
	)
}

var cmd = &cli.Command{
	Name:  "sss",
	Usage: "S3 client",
	Flags: []cli.Flag{
		flagEndpoint,
		flagInsecure,
		flagRegion,
		flagPathStyle,
		flagProfile,
		flagBucket,
		flagSecretKey,
		flagAccessKey,
		flagVerbosity,
	},
	Commands: []*cli.Command{
		{
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
		},
		{
			Name: "buckets",
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
				return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
					return ctrl.BucketList(cmd.String("prefix"))
				})
			},
		},
		{
			Name: "bucket",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
					return ctrl.BucketHead(cmd.String(flagBucket.Name))
				})
			},
		},
		{
			Name: "mb",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
					return ctrl.BucketCreate(cmd.String(flagBucket.Name))
				})
			},
		},
		{
			Name: "rb",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
					return ctrl.BucketDelete(cmd.String(flagBucket.Name))
				})
			},
		},
		{
			Name: "multiparts",
			Commands: []*cli.Command{
				{
					Name: "ls",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
							return ctrl.BucketMultipartUploadsList(cmd.String(flagBucket.Name))
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
						return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
							return ctrl.BucketMultipartUploadAbort(
								cmd.String(flagBucket.Name),
								cmd.String(flagObjectKey.Name),
								cmd.String(flagUploadID.Name),
							)
						})
					},
				},
			},
		},
		{
			Name: "parts",
			Flags: []cli.Flag{
				flagObjectKey,
				flagUploadID,
			},
			Commands: []*cli.Command{
				{
					Name: "ls",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
							return ctrl.BucketPartsList(
								cmd.String(flagBucket.Name),
								cmd.String(flagObjectKey.Name),
								cmd.String(flagUploadID.Name),
							)
						})
					},
				},
			},
		},
		{
			Name: "ls",
			Arguments: []cli.Argument{
				argPrefix,
			},
			Flags: []cli.Flag{
				flagDelimiter,
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
					return ctrl.ObjectList(
						cmd.String(flagBucket.Name),
						cmd.StringArg(argPrefix.Name),
						cmd.String(flagDelimiter.Name),
					)
				})
			},
		},
		{
			Name:        "cp",
			Description: "server side object copy",
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
					Required: true,
				},
				&cli.StringFlag{
					Name:     "src-key",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "dst-bucket",
					Required: true,
				},
				&cli.StringFlag{
					Name:  "dst-key",
					Usage: "when empty, the src-key will be used",
				},
				flagSSEcKey,
				flagSSEcAlgo,
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
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
		},
		{
			Name: "put",
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
				return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
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
		},
		{
			Name: "rm",
			Arguments: []cli.Argument{
				argKey,
			},
			Flags: []cli.Flag{
				flagDelimiter,
				flagForce,
				flagConcurrency,
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
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
		},
		{
			Name: "get",
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
				return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
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
		},
		{
			Name: "head",
			Arguments: []cli.Argument{
				argKey,
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
					return ctrl.ObjectHead(
						cmd.String(flagBucket.Name),
						cmd.StringArg(argKey.Name),
					)
				})
			},
		},
		{
			Name: "presign",
			Flags: []cli.Flag{
				&cli.DurationFlag{
					Name: "expires-in",
				},
			},
			Commands: []*cli.Command{
				{
					Name: "get",
					Arguments: []cli.Argument{
						argKey,
					},
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
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
					Name: "put",
					Arguments: []cli.Argument{
						argKey,
					},
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
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
		},
		{
			Name: "policy",
			Commands: []*cli.Command{
				{
					Name: "get",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
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
						return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
							return ctrl.BucketPolicyPut(
								cmd.StringArg("policy"),
								cmd.String(flagBucket.Name),
							)
						})
					},
				},
			},
		},
		{
			Name: "cors",
			Commands: []*cli.Command{
				{
					Name: "get",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
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
						return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
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
						return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
							return ctrl.BucketCORSDelete(
								cmd.String(flagBucket.Name),
							)
						})
					},
				},
			},
		},
		{
			Name: "object-lock",
			Commands: []*cli.Command{
				{
					Name: "get",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
							return ctrl.BucketObjectLockGet(
								cmd.String(flagBucket.Name),
							)
						})
					},
				},
				{
					Name: "put",
					Arguments: []cli.Argument{
						&cli.StringArg{
							Name: "config",
						},
					},
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
							return ctrl.BucketObjectLockPut(
								cmd.StringArg("config"),
								cmd.String(flagBucket.Name),
							)
						})
					},
				},
			},
		},
		{
			Name: "lifecycle",
			Commands: []*cli.Command{
				{
					Name: "get",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
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
						return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
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
						return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
							return ctrl.BucketLifecycleDelete(
								cmd.String(flagBucket.Name),
							)
						})
					},
				},
			},
		},
		{
			Name: "versioning",
			Commands: []*cli.Command{
				{
					Name: "get",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
							return ctrl.BucketVersioningGet(
								cmd.String(flagBucket.Name),
							)
						})
					},
				},
			},
		},
		{
			Name: "acl",
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
						return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
							return ctrl.ObjectACLGet(
								cmd.String(flagBucket.Name),
								cmd.StringArg(argKey.Name),
								cmd.String(flagVersionID.Name),
							)
						})
					},
				},
			},
		},
		{
			Name: "versions",
			Arguments: []cli.Argument{
				argPrefix,
			},
			Flags: []cli.Flag{
				flagDelimiter,
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return Exec(ctx, cmd, func(ctrl *controller.Controller) error {
					return ctrl.ObjectVersions(
						cmd.String(flagBucket.Name),
						cmd.StringArg(argPrefix.Name),
						cmd.String(flagDelimiter.Name),
					)
				})
			},
		},
	},
}
