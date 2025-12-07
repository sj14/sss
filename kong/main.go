package main

import (
	"context"
	"fmt"
	"log"
	"maps"
	"os"
	"slices"
	"time"

	"github.com/alecthomas/kong"
	"github.com/dustin/go-humanize"
	"github.com/sj14/sss/controller"
	"github.com/sj14/sss/util"
)

// //////// Common

type ArgPath struct {
	Path string `arg:"" name:"path"`
}

type ArgPathOptional struct {
	Path string `arg:"" name:"path" optional:""`
}

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

type FlagVersion struct {
	Version string `name:"versin"`
}

////////////

type CLI struct {
	// Commands
	Profiles Profiles `cmd:"" name:"profiles"`
	Buckets  Buckets  `cmd:"" name:"buckets"`
	Bucket   Bucket   `cmd:"" name:"bucket"`

	// Flags
	Config    string   `name:"config"`
	Profile   string   `name:"profile" default:"default"`
	Endpoint  string   `name:"endpoint"`
	Region    string   `name:"region"`
	PathStyle bool     `name:"path-style"`
	AccessKey string   `name:"access-key"`
	SecretKey string   `name:"secret-key"`
	Insecure  bool     `name:"insecure"`
	ReadOnly  bool     `name:"read-only"`
	Bandwidth string   `name:"bandwidth"`
	Header    []string `name:"header"`
	SNI       string   `name:"sni"`
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

	BucketCreate     BucketCreate     `cmd:"" group:"bucket" name:"mb"`
	BucketHead       BucketHead       `cmd:"" group:"bucket" name:"hb"`
	BucketRemove     BucketRemove     `cmd:"" group:"bucket" name:"rb"`
	BucketPolicy     BucketPolicy     `cmd:"" group:"bucket" name:"policy"`
	BucketCors       BucketCors       `cmd:"" group:"bucket" name:"cors"`
	BucketTag        BucketTag        `cmd:"" group:"bucket" name:"tag"`
	BucketLifecycle  BucketLifecycle  `cmd:"" group:"bucket" name:"lifecycle"`
	BucketVersioning BucketVersioning `cmd:"" group:"bucket" name:"versioning"`
	ObjectLock       ObjectLock       `cmd:"" group:"bucket" name:"object-lock"`
	BucketSize       BucketSize       `cmd:"" group:"bucket" name:"size"`
	Multiparts       Multiparts       `cmd:"" group:"multiparts" name:"multiparts"`
	ObjectList       ObjectList       `cmd:"" group:"object" name:"ls" aliases:"list"`
	ObjectCopy       ObjectCopy       `cmd:"" group:"object" name:"cp"`
	ObjectPut        ObjectPut        `cmd:"" group:"object" name:"put"`
	ObjectDelete     ObjectDelete     `cmd:"" group:"object" name:"rm"`
	ObjectGet        ObjectGet        `cmd:"" group:"object" name:"get"`
	ObcectHead       ObjectHead       `cmd:"" group:"object" name:"head"`
	ObjectPresign    ObjectPresign    `cmd:"" group:"object" name:"presign"`
	ObjectACL        ObjectACL        `cmd:"" group:"object" name:"acl"`
	ObjectVersions   ObjectVersions   `cmd:"" group:"object" name:"versions"`
}

type ObjectVersions struct {
	ArgPathOptional
	FlagJson
}

func (s ObjectVersions) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.ObjectVersions(
		cli.Bucket.BucketArg.BucketName,
		s.ArgPathOptional.Path,
		s.FlagJson.AsJson,
	)
}

type ObjectACL struct {
	ObjectACLGet ObjectACLGet `cmd:"" name:"get"`
}

type ObjectACLGet struct {
	ArgObject
	FlagVersion
}

func (s ObjectACLGet) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.ObjectACLGet(
		cli.Bucket.BucketArg.BucketName,
		s.ArgObject.Object,
		s.FlagVersion.Version,
	)
}

type BucketSize struct {
	ArgPathOptional
}

func (s BucketSize) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketSize(
		cli.Bucket.BucketArg.BucketName,
		s.ArgPathOptional.Path,
	)
}

type BucketVersioning struct {
	BucketVersioningGet BucketVersioningGet `cmd:"" name:"get"`
	BucketVersioningPut BucketVersioningPut `cmd:"" name:"put"`
}

type BucketVersioningGet struct{}

func (s BucketVersioningGet) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketVersioningGet(
		cli.Bucket.BucketArg.BucketName,
	)
}

type BucketVersioningPut struct {
	ArgObject
}

func (s BucketVersioningPut) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketVersioningPut(
		s.ArgObject.Object,
		cli.Bucket.BucketArg.BucketName,
	)
}

type BucketLifecycle struct {
	BucketLifecycleGet    BucketLifecycleGet    `cmd:"" name:"get"`
	BucketLifecyclePut    BucketLifecyclePut    `cmd:"" name:"put"`
	BucketLifecycleDelete BucketLifecycleDelete `cmd:"" name:"rm"`
}

type BucketLifecycleGet struct{}

func (s BucketLifecycleGet) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketLifecycleGet(
		cli.Bucket.BucketArg.BucketName,
	)
}

type BucketLifecyclePut struct {
	ArgPath
}

func (s BucketLifecyclePut) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketLifecyclePut(
		s.ArgPath.Path,
		cli.Bucket.BucketArg.BucketName,
	)
}

type BucketLifecycleDelete struct{}

func (s BucketLifecycleDelete) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketLifecycleDelete(
		cli.Bucket.BucketArg.BucketName,
	)
}

type ObjectLock struct {
	ObjectLockGet ObjectLockGet `cmd:"" name:"get"`
	ObjectLockPut ObjectLockPut `cmd:"" name:"put"`
}

type ObjectLockGet struct{}

func (s ObjectLockGet) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketObjectLockGet(
		cli.Bucket.BucketArg.BucketName,
	)
}

type ObjectLockPut struct {
	ArgPath
}

func (s ObjectLockPut) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketObjectLockPut(
		s.ArgPath.Path,
		cli.Bucket.BucketArg.BucketName,
	)
}

type BucketCors struct {
	BucketCorsGet    BucketCorsGet    `cmd:"" name:"get"`
	BucketCorsPut    BucketCorsPut    `cmd:"" name:"put"`
	BucketCorsRemove BucketCorsRemove `cmd:"" name:"rm"`
}

type BucketCorsRemove struct{}

func (s BucketCorsRemove) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketCORSDelete(
		cli.Bucket.BucketArg.BucketName,
	)
}

type BucketCorsPut struct {
	ArgPath
}

func (s BucketCorsPut) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketCORSPut(
		s.ArgPath.Path,
		cli.Bucket.BucketArg.BucketName,
	)
}

type BucketCorsGet struct{}

func (s BucketCorsGet) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketCORSGet(
		cli.Bucket.BucketArg.BucketName,
	)
}

type BucketPolicy struct {
	BucketPolicyGet    BucketPolicyGet    `cmd:"" name:"get"`
	BucketPolicyPut    BucketPolicyPut    `cmd:"" name:"put"`
	BucketPolicyRemove BucketPolicyRemove `cmd:"" name:"rm"`
}

type BucketPolicyPut struct {
	ArgPath
}

func (s BucketPolicyPut) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketPolicyPut(
		s.ArgPath.Path,
		cli.Bucket.BucketArg.BucketName,
	)
}

type BucketPolicyGet struct{}

func (s BucketPolicyGet) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketPolicyGet(
		cli.Bucket.BucketArg.BucketName,
	)
}

type BucketPolicyRemove struct{}

func (s BucketPolicyRemove) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketPolicyDelete(
		cli.Bucket.BucketArg.BucketName,
	)
}

type ObjectPresign struct {
	PresignGet    PresignGet    `cmd:"" name:"get"`
	PresignPut    PresignPut    `cmd:"" name:"put"`
	FlagExpiresIn time.Duration `name:"epxires-in"` // flag
}

type PresignPut struct {
	ArgObject
}

func (s PresignPut) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.ObjectPresignPut(
		cli.Bucket.BucketArg.ObjectPresign.FlagExpiresIn,
		s.ArgObject.Object,
		controller.ObjectPutConfig{
			Bucket: cli.Bucket.BucketArg.BucketName,
		},
	)
}

type PresignGet struct {
	ArgObject
}

func (s PresignGet) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.ObjectPresignGet(
		s.ArgObject.Object,
		controller.ObjectGetConfig{
			Bucket: cli.Bucket.BucketArg.BucketName,
		},
		cli.Bucket.BucketArg.ObjectPresign.FlagExpiresIn,
	)
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

type BucketRemove struct {
	FlagForce
}

func (s BucketRemove) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketDelete(
		cli.Bucket.BucketArg.BucketName,
		s.FlagForce.Force,
	)
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

	util.SetIfNotZero(&profile.Endpoint, cli.Endpoint)
	util.SetIfNotZero(&profile.Region, cli.Region)
	util.SetIfNotZero(&profile.PathStyle, cli.PathStyle)
	util.SetIfNotZero(&profile.AccessKey, cli.AccessKey)
	util.SetIfNotZero(&profile.SecretKey, cli.SecretKey)
	util.SetIfNotZero(&profile.Insecure, cli.Insecure)
	util.SetIfNotZero(&profile.ReadOnly, cli.ReadOnly)
	util.SetIfNotZero(&profile.SNI, cli.SNI)

	var bandwidth uint64
	if cli.Bandwidth != "" {
		bandwidth, err = humanize.ParseBytes(cli.Bandwidth)
		if err != nil {
			log.Fatalln(err)
		}
	}

	ctrl, err := controller.New(
		ctx,
		os.Stdout,
		os.Stderr,
		controller.ControllerConfig{
			Profile:   profile,
			Verbosity: 1,
			Headers:   cli.Header,
			Bandwidth: bandwidth,
			// DryRun: ,
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
