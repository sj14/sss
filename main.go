package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/sj14/sss/controller"
	"github.com/sj14/sss/util"
	"github.com/urfave/cli/v3"
)

func main() {
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalln(err)
	}
}

func Exec(cmd *cli.Command, fn func(ctrl *controller.Controller) error) error {
	ctrl, err := controller.New(context.Background(), controller.Config{
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
)

func parseSSEC(cmd *cli.Command) util.SSEC {
	return util.NewSSEC(
		cmd.String(flagSSEcAlgo.Name),
		cmd.String(flagSSEcKey.Name),
	)
}

var cmd = &cli.Command{
	Name: "sss",
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
			Name: "buckets",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name: "prefix",
				},
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
				return Exec(cmd, func(ctrl *controller.Controller) error {
					return ctrl.BucketList(cmd.String("prefix"))
				})
			},
		},
		{
			Name: "bucket",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return Exec(cmd, func(ctrl *controller.Controller) error {
					return ctrl.BucketHead(cmd.String(flagBucket.Name))
				})
			},
		},
		{
			Name: "mb",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return Exec(cmd, func(ctrl *controller.Controller) error {
					return ctrl.BucketCreate(cmd.String(flagBucket.Name))
				})
			},
		},
		{
			Name: "rb",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return Exec(cmd, func(ctrl *controller.Controller) error {
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
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.BucketMultipartUploadsList(cmd.String(flagBucket.Name))
						})
					},
				},
			},
		},
		{
			Name: "parts",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name: "key",
				},
				&cli.StringFlag{
					Name: "upload-id",
				},
			},
			Commands: []*cli.Command{
				{
					Name: "ls",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.BucketPartsList(
								cmd.String(flagBucket.Name),
								cmd.String("key"),
								cmd.String("upload-id"),
							)
						})
					},
				},
			},
		},
		{
			Name: "ls",
			Arguments: []cli.Argument{
				&cli.StringArg{
					Name: "prefix",
				},
			},
			Flags: []cli.Flag{
				flagDelimiter,
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return Exec(cmd, func(ctrl *controller.Controller) error {
					return ctrl.ObjectList(
						cmd.String(flagBucket.Name),
						cmd.StringArg("prefix"),
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
				return Exec(cmd, func(ctrl *controller.Controller) error {
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
				&cli.IntFlag{
					Name: "concurrency",
				},
				&cli.IntFlag{
					Name: "leave-parts-on-error",
				},
				&cli.IntFlag{
					Name: "max-parts",
				},
				&cli.Int64Flag{
					Name: "part-size",
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
				return Exec(cmd, func(ctrl *controller.Controller) error {
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
				&cli.StringArg{
					Name: "key",
				},
			},
			Flags: []cli.Flag{
				flagDelimiter,
				flagForce,
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return Exec(cmd, func(ctrl *controller.Controller) error {
					return ctrl.ObjectDelete(
						cmd.String(flagBucket.Name),
						cmd.StringArg("key"),
						cmd.String(flagDelimiter.Name),
						cmd.Bool(flagForce.Name),
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
				&cli.BoolFlag{
					Name: "recursive",
				},
				&cli.StringFlag{
					Name: "version-id",
				},
				&cli.StringFlag{
					Name:  "range",
					Usage: "bytes=beginbyte-endbyte, e.g. 'bytes=0-500' to get the first 500 bytes",
				},
				&cli.Int32Flag{
					Name: "part-number",
				},
				&cli.StringFlag{
					Name: "if-match",
				},
				&cli.StringFlag{
					Name: "if-none-match",
				},
				&cli.TimestampFlag{
					Name: "if-modified-since",
				},
				&cli.TimestampFlag{
					Name: "if-unmodified-since",
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return Exec(cmd, func(ctrl *controller.Controller) error {
					return ctrl.ObjectGet(
						cmd.StringArg("target"),
						cmd.String(flagDelimiter.Name),
						controller.ObjectGetConfig{
							Bucket:            cmd.String(flagBucket.Name),
							ObjectKey:         cmd.StringArg("key"),
							SSEC:              parseSSEC(cmd),
							VersionID:         cmd.String("version-id"),
							Range:             cmd.String("range"),
							PartNumber:        cmd.Int32("part-number"),
							IfMatch:           cmd.String("if-match"),
							IfNoneMatch:       cmd.String("if-none-match"),
							IfModifiedSince:   cmd.Timestamp("if-modified-since"),
							IfUnmodifiedSince: cmd.Timestamp("if-unmodified-since"),
						},
					)
				})
			},
		},
		{
			Name: "head",
			Arguments: []cli.Argument{
				&cli.StringArg{
					Name: "key",
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return Exec(cmd, func(ctrl *controller.Controller) error {
					return ctrl.ObjectHead(
						cmd.String(flagBucket.Name),
						cmd.StringArg("key"),
					)
				})
			},
		},
		{
			Name: "presign",
			Arguments: []cli.Argument{
				&cli.StringArg{
					Name: "key",
				},
			},
			Flags: []cli.Flag{
				&cli.DurationFlag{
					Name: "expires-in",
				},
				&cli.StringFlag{
					Name:  "method",
					Value: "get",
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return Exec(cmd, func(ctrl *controller.Controller) error {
					return ctrl.ObjectPresign(
						cmd.Duration("expires-in"),
						controller.ObjectPresignConfig{
							Method: cmd.String("method"),
							ObjectGetConfig: controller.ObjectGetConfig{
								Bucket:    cmd.String(flagBucket.Name),
								ObjectKey: cmd.StringArg("key"),
							},
						},
					)
				})
			}},
		{
			Name: "policy",
			Commands: []*cli.Command{
				{
					Name: "get",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
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
						return Exec(cmd, func(ctrl *controller.Controller) error {
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
			Name: "lifecycle",
			Commands: []*cli.Command{
				{
					Name: "get",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return Exec(cmd, func(ctrl *controller.Controller) error {
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
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.BucketLifecyclePut(
								cmd.StringArg("lifecycle"),
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
						return Exec(cmd, func(ctrl *controller.Controller) error {
							return ctrl.BucketVersioningGet(
								cmd.String(flagBucket.Name),
							)
						})
					},
				},
			},
		},
	},
}
