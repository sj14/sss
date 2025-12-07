package controller

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sj14/sss/util"
	"github.com/sj14/sss/util/progress"
)

type ObjectGetConfig struct {
	Bucket            string
	SSEC              util.SSEC
	VersionID         string
	IfMatch           string
	IfModifiedSince   time.Time
	IfNoneMatch       string
	IfUnmodifiedSince time.Time
	Range             string
	PartNumber        int32
	Concurrency       int
	PartSize          int64
	DryRun            bool
}

func (c *Controller) ObjectGet(targetDir, prefix, originalPrefix string, cfg ObjectGetConfig) error {
	if prefix == "" {
		return errors.New("missing key")
	}

	// only get single object
	if !strings.HasSuffix(prefix, "/") {
		fp := path.Join(targetDir, path.Base(prefix))
		return c.objectGet(fp, prefix, cfg)
	}

	// allow downloading the whole bucket
	if prefix == "/" {
		prefix = ""
		originalPrefix = ""
	}

	// recursive get
	for l, err := range c.objectList(cfg.Bucket, prefix) {
		if err != nil {
			return err
		}

		if l.Prefix != nil {
			err := c.ObjectGet(targetDir, *l.Prefix.Prefix, originalPrefix, cfg)
			if err != nil {
				return err
			}
			continue
		}

		lastDir := path.Base(path.Dir(originalPrefix))
		trimmedPrefix := strings.TrimPrefix(*l.Object.Key, originalPrefix)

		fp := path.Join(targetDir, lastDir, trimmedPrefix)

		err = c.objectGet(fp, *l.Object.Key, cfg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) objectGet(targetPath, objectKey string, cfg ObjectGetConfig) error {
	headObjectInput := &s3.HeadObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(objectKey),
	}

	getObjectInput := &s3.GetObjectInput{
		Bucket:            aws.String(cfg.Bucket),
		Key:               aws.String(objectKey),
		VersionId:         util.NilIfZero(cfg.VersionID),
		IfMatch:           util.NilIfZero(cfg.IfMatch),
		IfModifiedSince:   util.NilIfZero(cfg.IfModifiedSince),
		IfNoneMatch:       util.NilIfZero(cfg.IfNoneMatch),
		IfUnmodifiedSince: util.NilIfZero(cfg.IfUnmodifiedSince),
		PartNumber:        util.NilIfZero(cfg.PartNumber),
	}

	if cfg.SSEC.KeyIsSet() {
		headObjectInput.SSECustomerKeyMD5 = aws.String(cfg.SSEC.Base64KeyMD5())
		headObjectInput.SSECustomerKey = aws.String(cfg.SSEC.Base64Key())
		headObjectInput.SSECustomerAlgorithm = aws.String(cfg.SSEC.Algorithm())

		getObjectInput.SSECustomerKeyMD5 = aws.String(cfg.SSEC.Base64KeyMD5())
		getObjectInput.SSECustomerKey = aws.String(cfg.SSEC.Base64Key())
		getObjectInput.SSECustomerAlgorithm = aws.String(cfg.SSEC.Algorithm())
	}

	// Range requests are like "bytes=100-200".
	// It's easy to miss the "bytes=" part, add it when the flag value starts with a digit.
	if cfg.Range != "" && unicode.IsDigit(rune(cfg.Range[0])) {
		cfg.Range = fmt.Sprintf("bytes=%v", cfg.Range)
		getObjectInput.Range = &cfg.Range
	}

	var total uint64 = 0
	if cfg.Range == "" {
		headResp, err := c.client.HeadObject(c.ctx, headObjectInput)
		if err != nil {
			log.Printf("head object: %v\n", err)
		} else {
			total = uint64(*headResp.ContentLength)
		}
	}

	// TODO: represent download ranges
	if cfg.DryRun {
		var file = &os.File{}
		pw := progress.NewWriter(c.OutWriter, file, total, c.verbosity, targetPath)
		pw.Finish()
		return nil
	}

	// create the output dir
	if err := os.MkdirAll(filepath.Dir(targetPath), 0775); err != nil {
		return err
	}

	// create the output file
	file, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer file.Close()

	downloader := manager.NewDownloader(c.client, func(d *manager.Downloader) {
		d.Concurrency = cfg.Concurrency
		d.PartSize = cfg.PartSize
	})

	pw := progress.NewWriter(c.OutWriter, file, total, c.verbosity, targetPath)

	_, err = downloader.Download(c.ctx, pw, getObjectInput)
	if err != nil {
		return err
	}

	pw.Finish()

	return nil
}
