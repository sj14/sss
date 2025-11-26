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

func (s Profiles) Run(cli CLI, ctrl *controller.Controller, config controller.Config) error {
	keys := slices.Collect(maps.Keys(config.Profiles))
	slices.Sort(keys)
	for _, key := range keys {
		fmt.Println(key)
	}

	return nil
}

type Buckets struct{}

func (s Buckets) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketList("")
}

type Bucket struct {
	BucketArg BucketArg `arg:"" name:"bucket"`
}

type BucketArg struct {
	BucketName string `arg:"" name:"bucket"`

	BucketCreate BucketCreate `cmd:"" name:"mb"`
	BucketHead   BucketHead   `cmd:"" name:"hb"`
	BucketRemove BucketRemove `cmd:"" name:"rb"`
	BucketTag    BucketTag    `cmd:"" name:"tag"`
	Multiparts   Multiparts   `cmd:"" name:"multiparts"`
	ObjectList   ObjectList   `cmd:"" name:"ls" aliases:"list"`
	ObjectCopy   ObjectCopy   `cmd:"" name:"cp"`
	ObjectPut    ObjectPut    `cmd:"" name:"put"`
	ObjectDelete ObjectDelete `cmd:"" name:"rm"`
	ObjectGet    ObjectGet    `cmd:"" name:"get"`
	ObcectHead   ObjectHead   `cmd:"" name:"head"`
}

type ObjectHead struct {
	ArgObject
	DestinationPath string `arg:"" name:"destination" optional:""`
	FlagDryRun
	FlagConcurrency
}

func (s ObjectHead) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.ObjectHead(
		cli.Bucket.BucketArg.BucketName,
		s.ArgObject.Object,
	)
}

type ObjectGet struct {
	ArgObject
	DestinationPath string `arg:"" name:"destination" optional:""`
	FlagDryRun
	FlagConcurrency
}

func (s ObjectGet) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.ObjectGet(
		cli.Bucket.BucketArg.ObjectGet.DestinationPath,
		s.ArgObject.Object,
		s.ArgObject.Object,
		controller.ObjectGetConfig{
			Bucket:      cli.Bucket.BucketArg.BucketName,
			Concurrency: cli.Bucket.BucketArg.ObjectGet.FlagConcurrency.Concurrency,
			DryRun:      cli.Bucket.BucketArg.ObjectGet.FlagDryRun.DryRun,
		})
}

type ObjectDelete struct {
	ArgObject
	FlagConcurrency
	FlagDryRun
	FlagForce
}

func (s ObjectDelete) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.ObjectDelete(
		cli.Bucket.BucketArg.ObjectDelete.Object,
		controller.ObjectDeleteConfig{
			Bucket:      cli.Bucket.BucketArg.BucketName,
			Force:       cli.Bucket.BucketArg.ObjectDelete.FlagForce.Force,
			Concurrency: cli.Bucket.BucketArg.ObjectDelete.FlagConcurrency.Concurrency,
			DryRun:      cli.Bucket.BucketArg.ObjectDelete.FlagDryRun.DryRun,
		})

}

type ObjectPut struct {
	Filepath     string `arg:"" name:"path"`
	Destinaton   string `arg:"" name:"destination" optional:""`
	FlagPartSize int64  `name:"part-size"`
	FlagConcurrency
	FlagDryRun
}

func (s ObjectPut) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.ObjectPut(
		cli.Bucket.BucketArg.ObjectPut.Filepath,
		cli.Bucket.BucketArg.ObjectPut.Destinaton,
		controller.ObjectPutConfig{
			Bucket:      cli.Bucket.BucketArg.BucketName,
			Concurrency: cli.Bucket.BucketArg.ObjectPut.FlagConcurrency.Concurrency,
			DryRun:      cli.Bucket.BucketArg.ObjectPut.FlagDryRun.DryRun,
			// SSEC: ,
			// LeavePartsOnError: ,
			// MaxUploadParts: ,
			// PartSize: ,
			// ACL: ,
		},
	)
}

type ObjectCopy struct {
	SrcObject string `arg:"" name:"src-object"`
	DstBucket string `arg:"" name:"dst-bucket"`
	DstObject string `arg:"" name:"dst-object"`
}

func (s ObjectCopy) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.ObjectCopy(controller.ObjectCopyConfig{
		SrcBucket: cli.Bucket.BucketArg.BucketName,
		SrcKey:    cli.Bucket.BucketArg.ObjectCopy.SrcObject,
		DstBucket: cli.Bucket.BucketArg.ObjectCopy.DstBucket,
		DstKey:    cli.Bucket.BucketArg.ObjectCopy.DstObject,
	})
}

type BucketCreate struct {
	ObjectLock bool `name:"object-lock"`
}

func (s BucketCreate) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketCreate(
		cli.Bucket.BucketArg.BucketName,
		cli.Bucket.BucketArg.BucketCreate.ObjectLock,
	)
}

type BucketHead struct{}

func (s BucketHead) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketHead(cli.Bucket.BucketArg.BucketName)
}

type BucketRemove struct{}

func (s BucketRemove) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketDelete(cli.Bucket.BucketArg.BucketName)
}

type ObjectList struct {
	FlagPrefix
	FlagRecursive
	FlagJson
}

func (s ObjectList) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.ObjectList(
		cli.Bucket.BucketArg.BucketName,
		cli.Bucket.BucketArg.ObjectList.FlagPrefix.Prefix,
		cli.Bucket.BucketArg.ObjectList.FlagPrefix.Prefix,
		cli.Bucket.BucketArg.ObjectList.FlagRecursive.Recursive,
		cli.Bucket.BucketArg.ObjectList.FlagJson.AsJson,
	)
}

type BucketTag struct {
	BucketTagGet BucketTagGet `cmd:"" name:"get"`
}

type BucketTagGet struct{}

func (s BucketTagGet) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketTagging(cli.Bucket.BucketArg.BucketName)
}

type Multiparts struct {
	MultipartRemove MultipartRemove `cmd:"" name:"rm"`
	MultipartList   MultipartList   `cmd:"" name:"ls"`
	MultipartParts  MultipartParts  `cmd:"" name:"parts"`
}

type MultipartList struct {
	FlagPrefix
	FlagJson
}

func (s MultipartList) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketMultipartUploadsList(
		cli.Bucket.BucketArg.BucketName,
		s.FlagPrefix.Prefix,
		s.FlagJson.AsJson,
	)
}

type MultipartRemove struct {
	ArgObject
	ArgUploadID
}

func (s MultipartRemove) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketMultipartUploadAbort(
		cli.Bucket.BucketArg.BucketName,
		cli.Bucket.BucketArg.Multiparts.MultipartRemove.ArgObject.Object,
		cli.Bucket.BucketArg.Multiparts.MultipartRemove.ArgUploadID.UploadID,
	)

}

type MultipartParts struct {
	PartsList PartsList `cmd:"" name:"ls"`
}

type PartsList struct {
	ArgObject
	ArgUploadID
	FlagJson
}

func (s PartsList) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketPartsList(
		cli.Bucket.BucketArg.BucketName,
		cli.Bucket.BucketArg.Multiparts.MultipartParts.PartsList.Object,
		cli.Bucket.BucketArg.Multiparts.MultipartParts.PartsList.UploadID,
		cli.Bucket.BucketArg.Multiparts.MultipartParts.PartsList.AsJson,
	)
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

	err = kctx.Run(cli, ctrl, config)
	if err != nil {
		log.Fatalln(err)
	}

	if err != nil {
		log.Fatalln(err)
	}
}
