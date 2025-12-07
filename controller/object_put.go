package controller

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/sj14/sss/util"
	"github.com/sj14/sss/util/progress"
)

type ObjectPutConfig struct {
	Bucket            string
	SSEC              util.SSEC
	Concurrency       int
	LeavePartsOnError bool
	MaxUploadParts    int32
	PartSize          int64
	ACL               string
	DryRun            bool
}

func (c *Controller) ObjectPut(filePath, dest string, cfg ObjectPutConfig) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return c.objectPut(filePath, path.Join(dest, filepath.Base(filePath)), cfg)
	}

	return filepath.Walk(filePath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Switch to forward slash even when uploading from Windows.
		p = filepath.ToSlash(p)
		filePath = filepath.ToSlash(filepath.Clean(filePath))

		var (
			lastDir       = filepath.Base(filePath)
			trimmedPrefix = strings.TrimPrefix(p, filePath)
			fp            = path.Join(dest, lastDir, trimmedPrefix)
		)

		err = c.objectPut(p, fp, cfg)
		if err != nil {
			return err
		}

		return nil
	})
}

func (c *Controller) objectPut(filePath, key string, cfg ObjectPutConfig) error {
	if key == "" {
		key = filepath.Base(filePath)
	}
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	uploader := manager.NewUploader(c.client, func(u *manager.Uploader) {
		u.Concurrency = cfg.Concurrency
		u.LeavePartsOnError = cfg.LeavePartsOnError
		u.MaxUploadParts = cfg.MaxUploadParts
		u.PartSize = cfg.PartSize
		// u.RequestChecksumCalculation
	})

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	pr := progress.NewReader(c.OutWriter, f, uint64(stat.Size()), c.verbosity, key)

	putObjectInput := &s3.PutObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(filepath.ToSlash(filepath.Clean(key))),
		Body:   pr,
		ACL:    types.ObjectCannedACL(cfg.ACL),
	}

	if cfg.SSEC.KeyIsSet() {
		putObjectInput.SSECustomerKey = aws.String(cfg.SSEC.Base64Key())
		putObjectInput.SSECustomerKeyMD5 = aws.String(cfg.SSEC.Base64KeyMD5())
		putObjectInput.SSECustomerAlgorithm = aws.String(cfg.SSEC.Algorithm())
	}

	if !cfg.DryRun {
		_, err = uploader.Upload(c.ctx, putObjectInput)
		if err != nil {
			return err
		}
	}
	// don't put it into a defer after initializing
	// as it would then output the progress even when
	// the upload was abortet due to an error
	pr.Finish()

	return nil
}
