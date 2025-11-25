package main

import (
	"context"
	"fmt"
	"log"
	"maps"
	"os"
	"slices"

	"github.com/alecthomas/kong"
	"github.com/sj14/sss/controller"
)

// //////// Common

type ArgObject struct {
	Object string `arg:"" name:"object"`
}

type ArgUploadID struct {
	UploadID string `arg:"" name:"upload-id"`
}

type FlagPrefix struct {
	Prefix string `arg:"" name:"prefix" optional:""`
}

type FlagRecursive struct {
	Recursive bool `name:"recursive" short:"r"`
}

type FlagJson struct {
	AsJson bool `name:"json" short:"j"`
}

type FlagConcurrency struct {
	Concurrency int `name:"concurrency" default:"5"`
}

type FlagDryRun struct {
	DryRun bool `name:"dry-run"`
}

type FlagForce struct {
	Force bool `name:"force" short:"f"`
}

////////////

type CLI struct {
	Profiles Profiles `cmd:"" name:"profiles"`
	Buckets  Buckets  `cmd:"" name:"buckets"`
	Bucket   Bucket   `cmd:"" name:"bucket"`
	Config   string   // flag
	Profile  string   `default:"default"`
}

type Profiles struct{}

type Buckets struct{}

type Bucket struct {
	BucketArg BucketArg `arg:"" name:"bucket"`
}

type BucketArg struct {
	BucketName string `arg:"" name:"bucket"`

	BucketCreate BucketCreate `cmd:"" name:"mb"`
	BucketHead   BucketHead   `cmd:"" name:"hb"`
	BucketRemove BucketRemove `cmd:"" name:"rb"`
	BucketList   BucketList   `cmd:"" name:"ls" aliases:"list"`
	BucketTag    BucketTag    `cmd:"" name:"tag"`
	Multiparts   Multiparts   `cmd:"" name:"multiparts"`
	ObjectCopy   ObjectCopy   `cmd:"" name:"cp"`
	ObjectPut    ObjectPut    `cmd:"" name:"put"`
	ObjectDelete ObjectDelete `cmd:"" name:"rm"`
	ObjectGet    ObjectGet    `cmd:"" name:"get"`
}

type ObjectGet struct {
	ObjectPath      string `arg:"" name:"object"`
	DestinationPath string `arg:"" name:"destination" optional:""`
	FlagDryRun
	FlagConcurrency
}

type ObjectDelete struct {
	ObjectPath string `arg:"" name:"path"`
	FlagConcurrency
	FlagDryRun
	FlagForce
}

type ObjectPut struct {
	Filepath     string `arg:"" name:"path"`
	Destinaton   string `arg:"" name:"destination" optional:""`
	FlagPartSize int64  `name:"part-size"`
	FlagConcurrency
	FlagDryRun
}

type ObjectCopy struct {
	SrcObject string `arg:"" name:"src-object"`
	DstBucket string `arg:"" name:"dst-bucket"`
	DstObject string `arg:"" name:"dst-object"`
}

type BucketCreate struct {
	ObjectLock bool `name:"object-lock"`
}

type BucketHead struct{}

type BucketRemove struct{}

type BucketList struct {
	FlagPrefix
	FlagRecursive
	FlagJson
}

type BucketTag struct {
	Get struct{} `cmd:"" name:"get"`
}

type Multiparts struct {
	List struct {
		FlagPrefix
		FlagJson
	} `cmd:"" name:"ls"`
	Remove struct {
		ArgObject
		ArgUploadID
	} `cmd:"" name:"rm"`
	Parts struct {
		List struct {
			ArgObject
			ArgUploadID
			FlagJson
		} `cmd:"" name:"ls"`
	} `cmd:"" name:"parts"`
}

func main() {
	cli := CLI{}

	kctx := kong.Parse(&cli)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := controller.LoadConfig(cli.Config)
	if err != nil {
		log.Fatalln(err)
	}

	profile, ok := config.Profiles[cli.Profile]
	if !ok && cli.Profile != "default" {
		fmt.Printf("profile %q not found, available profiles:\n", cli.Profile)

		keys := slices.Collect(maps.Keys(config.Profiles))
		slices.Sort(keys)

		for _, key := range keys {
			fmt.Println(key)
		}

		os.Exit(1)
	}

	// util.SetIfNotZero(&profile.Endpoint, cfg.Profile.Endpoint)
	// util.SetIfNotZero(&profile.Region, cfg.Profile.Region)
	// util.SetIfNotZero(&profile.PathStyle, cfg.Profile.PathStyle)
	// util.SetIfNotZero(&profile.AccessKey, cfg.Profile.AccessKey)
	// util.SetIfNotZero(&profile.SecretKey, cfg.Profile.SecretKey)
	// util.SetIfNotZero(&profile.Insecure, cfg.Profile.Insecure)
	// util.SetIfNotZero(&profile.ReadOnly, cfg.Profile.ReadOnly)
	// util.SetIfNotZero(&profile.SNI, cfg.Profile.SNI)

	ctrl, err := controller.New(ctx, os.Stdout, os.Stderr, controller.ControllerConfig{
		Profile:   profile,
		Verbosity: 1,
	})
	if err != nil {
		log.Fatalln(err)
	}

	switch kctx.Command() {
	case "profiles":
		keys := slices.Collect(maps.Keys(config.Profiles))
		slices.Sort(keys)
		for _, key := range keys {
			fmt.Println(key)
		}
	case "buckets":
		err = ctrl.BucketList("")
	case "bucket <bucket> mb":
		err = ctrl.BucketCreate(
			cli.Bucket.BucketArg.BucketName,
			cli.Bucket.BucketArg.BucketCreate.ObjectLock,
		)
	case "bucket <bucket> hb":
		err = ctrl.BucketHead(cli.Bucket.BucketArg.BucketName)
	case "bucket <bucket> rb":
		err = ctrl.BucketDelete(cli.Bucket.BucketArg.BucketName)
	case "bucket <bucket> ls":
		err = ctrl.ObjectList(
			cli.Bucket.BucketArg.BucketName,
			cli.Bucket.BucketArg.BucketList.Prefix,
			cli.Bucket.BucketArg.BucketList.Prefix,
			cli.Bucket.BucketArg.BucketList.Recursive,
			cli.Bucket.BucketArg.BucketList.AsJson,
		)

	case "bucket <bucket> tag get":
		err = ctrl.BucketTagging(cli.Bucket.BucketArg.BucketName)
	case "bucket <bucket> multiparts ls":
		err = ctrl.BucketMultipartUploadsList(
			cli.Bucket.BucketArg.BucketName,
			cli.Bucket.BucketArg.Multiparts.List.Prefix,
			cli.Bucket.BucketArg.Multiparts.List.AsJson,
		)
	case "bucket <bucket> multiparts rm <object> <upload-id>":
		err = ctrl.BucketMultipartUploadAbort(
			cli.Bucket.BucketArg.BucketName,
			cli.Bucket.BucketArg.Multiparts.Remove.Object,
			cli.Bucket.BucketArg.Multiparts.Remove.UploadID,
		)
	case "bucket <bucket> multiparts parts ls <object> <upload-id>":
		err = ctrl.BucketPartsList(
			cli.Bucket.BucketArg.BucketName,
			cli.Bucket.BucketArg.Multiparts.Parts.List.Object,
			cli.Bucket.BucketArg.Multiparts.Parts.List.UploadID,
			cli.Bucket.BucketArg.Multiparts.Parts.List.AsJson,
		)
	case "bucket <bucket> cp <src-object> <dst-bucket> <dst-object>":
		err = ctrl.ObjectCopy(controller.ObjectCopyConfig{
			SrcBucket: cli.Bucket.BucketArg.BucketName,
			SrcKey:    cli.Bucket.BucketArg.ObjectCopy.SrcObject,
			DstBucket: cli.Bucket.BucketArg.ObjectCopy.DstBucket,
			DstKey:    cli.Bucket.BucketArg.ObjectCopy.DstObject,
		})
	case "bucket <bucket> put <path>", "bucket <bucket> put <path> <destination>":
		err = ctrl.ObjectPut(
			cli.Bucket.BucketArg.ObjectPut.Filepath,
			cli.Bucket.BucketArg.ObjectPut.Destinaton,
			controller.ObjectPutConfig{
				Bucket:      cli.Bucket.BucketArg.BucketName,
				Concurrency: cli.Bucket.BucketArg.ObjectPut.Concurrency,
				DryRun:      cli.Bucket.BucketArg.ObjectPut.DryRun,
				// SSEC: ,
				// LeavePartsOnError: ,
				// MaxUploadParts: ,
				// PartSize: ,
				// ACL: ,
			},
		)
	case "bucket <bucket> rm <path>":
		err = ctrl.ObjectDelete(
			cli.Bucket.BucketArg.ObjectDelete.ObjectPath,
			controller.ObjectDeleteConfig{
				Bucket:      cli.Bucket.BucketArg.BucketName,
				Force:       cli.Bucket.BucketArg.ObjectDelete.Force,
				Concurrency: cli.Bucket.BucketArg.ObjectDelete.Concurrency,
				DryRun:      cli.Bucket.BucketArg.ObjectDelete.DryRun,
			})
	case "bucket <bucket> get <object>", "bucket <bucket> get <object> <destination>":
		err = ctrl.ObjectGet(
			cli.Bucket.BucketArg.ObjectGet.DestinationPath,
			cli.Bucket.BucketArg.ObjectGet.ObjectPath,
			cli.Bucket.BucketArg.ObjectGet.ObjectPath,
			controller.ObjectGetConfig{
				Bucket:      cli.Bucket.BucketArg.BucketName,
				Concurrency: cli.Bucket.BucketArg.ObjectGet.Concurrency,
				DryRun:      cli.Bucket.BucketArg.ObjectGet.DryRun,
			})
	default:
		panic(kctx.Command())
	}

	if err != nil {
		log.Fatalln(err)
	}
}
