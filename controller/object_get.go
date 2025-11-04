package controller

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	"unicode"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sj14/sss/progress"
	"github.com/sj14/sss/util"
)

type ObjectGetConfig struct {
	Recursive         bool
	Bucket            string
	ObjectKey         string
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
}

func (c *Controller) ObjectGet(targetDir, delimiter string, cfg ObjectGetConfig) error {
	if !cfg.Recursive {
		return c.objectGet(targetDir, cfg)
	}

	objects, _, err := c.objectList(cfg.Bucket, cfg.ObjectKey, delimiter)
	if err != nil {
		return err
	}

	if len(objects) == 0 {
		return fmt.Errorf("no objects found")
	}

	for _, object := range objects {
		cfg.ObjectKey = *object.Key
		fmt.Printf("loading %s\n", *object.Key)
		err = c.objectGet(targetDir, cfg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) objectGet(targetDir string, cfg ObjectGetConfig) error {
	headObjectInput := &s3.HeadObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(cfg.ObjectKey),
	}

	getObjectInput := &s3.GetObjectInput{
		Bucket:            aws.String(cfg.Bucket),
		Key:               aws.String(cfg.ObjectKey),
		VersionId:         util.IfNotZero(cfg.VersionID),
		IfMatch:           util.IfNotZero(cfg.IfMatch),
		IfModifiedSince:   util.IfNotZero(cfg.IfModifiedSince),
		IfNoneMatch:       util.IfNotZero(cfg.IfNoneMatch),
		IfUnmodifiedSince: util.IfNotZero(cfg.IfUnmodifiedSince),
		PartNumber:        util.IfNotZero(cfg.PartNumber),
	}

	if cfg.SSEC.KeyIsSet() {
		headObjectInput.SSECustomerKeyMD5 = aws.String(cfg.SSEC.Base64KeyMD5())
		headObjectInput.SSECustomerKey = aws.String(cfg.SSEC.Base64Key())
		headObjectInput.SSECustomerAlgorithm = aws.String(cfg.SSEC.Algorithm())

		getObjectInput.SSECustomerKeyMD5 = aws.String(cfg.SSEC.Base64KeyMD5())
		getObjectInput.SSECustomerKey = aws.String(cfg.SSEC.Base64Key())
		getObjectInput.SSECustomerAlgorithm = aws.String(cfg.SSEC.Algorithm())
	}

	// range requests are like "bytes=100-200".
	// It's easy to miss the "bytes=" part, add it when the flag value starts with a digit
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

	// create the output dir
	filePath := filepath.Join(targetDir, cfg.ObjectKey)
	if err := os.MkdirAll(filepath.Dir(filePath), 0775); err != nil {
		return err
	}

	// create the output file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	downloader := manager.NewDownloader(c.client, func(d *manager.Downloader) {
		d.Concurrency = cfg.Concurrency
		d.PartSize = cfg.PartSize
	})

	// TODO: represent download ranges
	pw := progress.NewWriter(file, total, c.verbosity)
	defer pw.Finish()

	_, err = downloader.Download(c.ctx, pw, getObjectInput)
	if err != nil {
		return err
	}

	return nil
}
