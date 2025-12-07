package controller

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

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
	Expires           time.Time
}

func (c *Controller) ObjectPut(filePath, dest string, cfg ObjectPutConfig) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		if dest == "" {
			dest = filepath.Base(filePath)
		}

		if strings.HasSuffix(dest, "/") {
			dest = path.Join(dest, filepath.Base(filePath))
		}

		f, err := os.Open(filePath)
		if err != nil {
			return err
		}

		return c.objectPut(f, uint64(info.Size()), dest, cfg)
	}

	// TODO: flatten option which allows storing in the current folder instead of creating the subfolder?
	return filepath.Walk(filePath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(p)
		if err != nil {
			return err
		}

		// Switch to forward slash even when uploading from Windows.
		p = filepath.ToSlash(p)
		filePath = filepath.ToSlash(filepath.Clean(filePath))

		var (
			lastDir       = filepath.Base(filePath)
			trimmedPrefix = strings.TrimPrefix(p, filePath)
			fp            = path.Join(dest, lastDir, trimmedPrefix)
		)

		err = c.objectPut(f, uint64(info.Size()), fp, cfg)
		if err != nil {
			return err
		}

		return nil
	})
}

func (c *Controller) objectPut(body io.Reader, size uint64, key string, cfg ObjectPutConfig) error {
	uploader := manager.NewUploader(c.client, func(u *manager.Uploader) {
		u.Concurrency = cfg.Concurrency
		u.LeavePartsOnError = cfg.LeavePartsOnError
		u.MaxUploadParts = cfg.MaxUploadParts
		u.PartSize = cfg.PartSize
		// u.RequestChecksumCalculation
	})

	pr := progress.NewReader(c.OutWriter, body, size, c.verbosity, key)

	putObjectInput := &s3.PutObjectInput{
		Bucket:  aws.String(cfg.Bucket),
		Key:     aws.String(filepath.ToSlash(filepath.Clean(key))),
		Body:    pr,
		ACL:     types.ObjectCannedACL(cfg.ACL),
		Expires: aws.Time(cfg.Expires),
	}

	if cfg.SSEC.KeyIsSet() {
		putObjectInput.SSECustomerKey = aws.String(cfg.SSEC.Base64Key())
		putObjectInput.SSECustomerKeyMD5 = aws.String(cfg.SSEC.Base64KeyMD5())
		putObjectInput.SSECustomerAlgorithm = aws.String(cfg.SSEC.Algorithm())
	}

	if !cfg.DryRun {
		_, err := uploader.Upload(c.ctx, putObjectInput)
		if err != nil {
			return err
		}
	}
	// Don't put it into a defer after initializing
	// as it would then output the progress even when
	// the upload was abortet due to an error.
	pr.Finish()

	return nil
}
