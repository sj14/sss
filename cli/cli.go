package cli

import (
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/alecthomas/kong"
	"github.com/sj14/sss/controller"
	"github.com/sj14/sss/util"
)

type CLI struct {
	// Commands
	Profiles Profiles `cmd:"" name:"profiles" help:"List availale profiles."`
	Buckets  Buckets  `cmd:"" name:"buckets" group:"Bucket Commands" help:"List all buckets."`
	Bucket   Bucket   `cmd:"" name:"bucket"   help:"Manage bucket and objects."`
	Version  Version  `cmd:"" name:"version"  help:"Show version information."`

	// Flags
	Config    string            `name:"config"    short:"c"         help:"Path to the config file (default: ~/.config/sss/config.toml)."`
	Profile   string            `name:"profile"   short:"p" default:"default" help:"Profile to use." `
	Verbosity uint8             `name:"verbosity" short:"v" default:"1"       help:"Verbose output"`
	Endpoint  string            `name:"endpoint"                     help:"S3 endpoint URL."`
	Region    string            `name:"region"                       help:"S3 region."`
	PathStyle bool              `name:"path-style"                   help:"Use path style S3 requests."`
	AccessKey string            `name:"access-key"                   help:"S3 access key."`
	SecretKey string            `name:"secret-key"                   help:"S3 secret key."`
	Insecure  bool              `name:"insecure"                     help:"Skip TLS verification."`
	ReadOnly  bool              `name:"read-only"                    help:"Only allow safe HTTP methods."`
	Bandwidth string            `name:"bandwidth"                    help:"Limit bandwith per second, e.g. '1 MiB' (always 64 KiB burst)."`
	Headers   map[string]string `name:"header"                       help:"Set HTTP headers (format: 'key1=val1;key2=val2')."`
	Params    map[string]string `name:"param"                        help:"Set URL parameters (format: 'key1=val1;key2=val2')."`
	SNI       string            `name:"sni"                          help:"TLS Server Name Indication."`
}

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

type ArgPrefix struct {
	Prefix string `arg:"" name:"prefix" optional:""`
}

type FlagRecursive struct {
	Recursive bool `name:"recursive" short:"r"`
}

type FlagJson struct {
	AsJson bool `name:"json" short:"j"`
}

type FlagConcurrency struct {
	Concurrency int `name:"concurrency" short:"C" default:"5"`
}

type FlagDryRun struct {
	DryRun bool `name:"dry-run"`
}

type FlagBypassGovernance struct {
	BypassGovernance bool `name:"bypass-governance" help:"delete even when objects are governance retention locked"`
}

type FlagForce struct {
	Force bool `name:"force" short:"f"`
}

type FlagRange struct {
	Range string `name:"range" help:"'bytes=0-500' to get the first 501 bytes"`
}

type FlagVersionID struct {
	VersionID string `name:"version" help:"Version ID"`
}

type flagsSSEC struct {
	Algo string `name:"sse-c-algorithm" default:"AES256"`
	Key  string `name:"sse-c-key" help:"32 bytes key for AES256"`
}

func isFlagSet(flags []*kong.Flag, name string) bool {
	for _, f := range flags {
		if !f.Set {
			continue
		}
		if f.Name == name {
			return true
		}
	}
	return false
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

type Version struct {
	version string
}

func (s Version) Run(cli CLI, ctrl *controller.Controller, config controller.Config) error {
	fmt.Fprintln(ctrl.OutWriter, s.version)
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

	BucketCreate     BucketCreate     `cmd:"" group:"Bucket Commands" name:"mb"`
	BucketHead       BucketHead       `cmd:"" group:"Bucket Commands" name:"hb"`
	BucketRemove     BucketRemove     `cmd:"" group:"Bucket Commands" name:"rb"`
	BucketPolicy     BucketPolicy     `cmd:"" group:"Bucket Commands" name:"policy"`
	BucketCors       BucketCors       `cmd:"" group:"Bucket Commands" name:"cors"`
	BucketTag        BucketTag        `cmd:"" group:"Bucket Commands" name:"tag"`
	BucketLifecycle  BucketLifecycle  `cmd:"" group:"Bucket Commands" name:"lifecycle"`
	BucketVersioning BucketVersioning `cmd:"" group:"Bucket Commands" name:"versioning"`
	BucketCleanup    BucketCleanup    `cmd:"" group:"Bucket Commands" name:"cleanup"`
	ObjectLock       ObjectLock       `cmd:"" group:"Bucket Commands" name:"object-lock"`
	BucketSize       BucketSize       `cmd:"" group:"Bucket Commands" name:"size"`
	Multiparts       Multiparts       `cmd:"" group:"Multipart Commands" name:"multipart"`
	ObjectList       ObjectList       `cmd:"" group:"Object Commands" name:"ls"`
	ObjectCopy       ObjectCopy       `cmd:"" group:"Object Commands" name:"cp"`
	ObjectPut        ObjectPut        `cmd:"" group:"Object Commands" name:"put"`
	ObjectDelete     ObjectDelete     `cmd:"" group:"Object Commands" name:"rm"`
	ObjectGet        ObjectGet        `cmd:"" group:"Object Commands" name:"get" help:"Download files. Requires HeadObject permission."`
	ObcectHead       ObjectHead       `cmd:"" group:"Object Commands" name:"head"`
	ObjectPresign    ObjectPresign    `cmd:"" group:"Object Commands" name:"presign"`
	ObjectACL        ObjectACL        `cmd:"" group:"Object Commands" name:"acl"`
	ObjectVersions   ObjectVersions   `cmd:"" group:"Object Commands" name:"versions"`
}

type BucketCleanup struct {
	FlagConcurrency
	FlagForce
	FlagDryRun
	FlagObjectsVersions bool `name:"all-object-versions" help:"Removes all object versions from a bucket"`
	FlagMultiparts      bool `name:"all-multiparts"      help:"Removes all multipart uploads from a bucket"`
}

func (s BucketCleanup) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketCleanup(controller.BucketCleanupConfig{
		Bucket:           cli.Bucket.BucketArg.BucketName,
		Concurrency:      s.Concurrency,
		Force:            s.Force,
		DryRun:           s.DryRun,
		Multiparts:       s.FlagMultiparts,
		ObjectVersion:    s.FlagObjectsVersions,
		BypassGovernance: true,
	})
}

type ObjectVersions struct {
	ArgPathOptional
	FlagJson
	FlagRecursive
}

func (s ObjectVersions) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.ObjectVersions(
		cli.Bucket.BucketArg.BucketName,
		s.ArgPathOptional.Path,
		s.ArgPathOptional.Path,
		s.Recursive,
		s.FlagJson.AsJson,
	)
}

type ObjectACL struct {
	ObjectACLGet ObjectACLGet `cmd:"" name:"get"`
}

type ObjectACLGet struct {
	ArgObject
	FlagVersionID
}

func (s ObjectACLGet) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.ObjectACLGet(
		cli.Bucket.BucketArg.BucketName,
		s.ArgObject.Object,
		s.FlagVersionID.VersionID,
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
	ArgPath
}

func (s BucketVersioningPut) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketVersioningPut(
		s.ArgPath.Path,
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
	flagsSSEC
}

func (s ObjectHead) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.ObjectHead(
		cli.Bucket.BucketArg.BucketName,
		s.ArgObject.Object,
		controller.ObjectHeadConfig{
			SSEC: util.NewSSEC(s.flagsSSEC.Algo, s.flagsSSEC.Key),
		},
	)
}

type ObjectGet struct {
	ArgObject
	DestinationPath string `arg:"" name:"destination" optional:""`
	FlagDryRun
	FlagConcurrency
	flagsSSEC
	FlagVersionID
	FlagRange
	FlagForce
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
			SSEC:        util.NewSSEC(s.flagsSSEC.Algo, s.flagsSSEC.Key),
			VersionID:   s.FlagVersionID.VersionID,
			Range:       s.FlagRange.Range,
			Force:       s.Force,
			// PartNumber:        cmd.Int32(flagPartNumber.Name),
			// PartSize:          cmd.Int64(flagPartSize.Name),
			// IfMatch:           cmd.String(flagIfMatch.Name),
			// IfNoneMatch:       cmd.String(flagIfNoneMatch.Name),
			// IfModifiedSince:   cmd.Timestamp(flagIfModifiedSince.Name),
			// IfUnmodifiedSince: cmd.Timestamp(flagIfUnmodifiedSince.Name),
		})
}

type ObjectDelete struct {
	ArgObject
	FlagConcurrency
	FlagDryRun
	FlagForce
	FlagVersionID
}

func (s ObjectDelete) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.ObjectDelete(
		cli.Bucket.BucketArg.ObjectDelete.Object,
		controller.ObjectDeleteConfig{
			Bucket:      cli.Bucket.BucketArg.BucketName,
			Force:       cli.Bucket.BucketArg.ObjectDelete.FlagForce.Force,
			Concurrency: cli.Bucket.BucketArg.ObjectDelete.FlagConcurrency.Concurrency,
			DryRun:      cli.Bucket.BucketArg.ObjectDelete.FlagDryRun.DryRun,
			VersionID:   s.FlagVersionID.VersionID,
			// BypassGovernance: ,
		})

}

type ObjectPut struct {
	Filepath              string `arg:"" name:"path"`
	Destinaton            string `arg:"" name:"destination" optional:""`
	FlagPartSize          int64  `name:"part-size"`
	FlagMaxUploadParts    int32  `name:"max-parts"`
	FlagLeavePartsOnError bool   `name:"leave-error-parts"`
	FlagACL               string `name:"acl"`
	FlagConcurrency
	FlagDryRun
	flagsSSEC
}

func (s ObjectPut) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.ObjectPut(
		cli.Bucket.BucketArg.ObjectPut.Filepath,
		cli.Bucket.BucketArg.ObjectPut.Destinaton,
		controller.ObjectPutConfig{
			Bucket:            cli.Bucket.BucketArg.BucketName,
			Concurrency:       cli.Bucket.BucketArg.ObjectPut.FlagConcurrency.Concurrency,
			DryRun:            cli.Bucket.BucketArg.ObjectPut.FlagDryRun.DryRun,
			SSEC:              util.NewSSEC(s.flagsSSEC.Algo, s.flagsSSEC.Key),
			PartSize:          s.FlagPartSize,
			MaxUploadParts:    s.FlagMaxUploadParts,
			LeavePartsOnError: s.FlagLeavePartsOnError,
			ACL:               s.FlagACL,
		},
	)
}

type ObjectCopy struct {
	SrcObject string `arg:"" name:"src-object"`
	DstBucket string `arg:"" name:"dst-bucket"`
	DstObject string `arg:"" name:"dst-object"`
	flagsSSEC
}

func (s ObjectCopy) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.ObjectCopy(controller.ObjectCopyConfig{
		SrcBucket: cli.Bucket.BucketArg.BucketName,
		SrcKey:    cli.Bucket.BucketArg.ObjectCopy.SrcObject,
		DstBucket: cli.Bucket.BucketArg.ObjectCopy.DstBucket,
		DstKey:    cli.Bucket.BucketArg.ObjectCopy.DstObject,
		SSEC:      util.NewSSEC(s.flagsSSEC.Algo, s.flagsSSEC.Key),
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
	ArgPrefix
	FlagRecursive
	FlagJson
}

func (s ObjectList) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.ObjectList(
		cli.Bucket.BucketArg.BucketName,
		cli.Bucket.BucketArg.ObjectList.ArgPrefix.Prefix,
		cli.Bucket.BucketArg.ObjectList.ArgPrefix.Prefix,
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
	ArgPrefix
	FlagRecursive
	FlagJson
}

func (s MultipartList) Run(cli CLI, ctrl *controller.Controller) error {
	return ctrl.BucketMultipartUploadsList(
		cli.Bucket.BucketArg.BucketName,
		s.ArgPrefix.Prefix,
		s.ArgPrefix.Prefix,
		s.FlagRecursive.Recursive,
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
