package controller

import (
	"fmt"
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
}

func (c *Controller) ObjectGet(targetDir, delimiter string, cfg ObjectGetConfig) error {
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
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(cfg.ObjectKey),
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
		util.SetIfNotZero(&getObjectInput.Range, cfg.Range)
	}

	util.SetIfNotZero(&getObjectInput.VersionId, cfg.VersionID)
	util.SetIfNotZero(&getObjectInput.IfMatch, cfg.IfMatch)
	util.SetIfNotZero(&getObjectInput.IfModifiedSince, cfg.IfModifiedSince)
	util.SetIfNotZero(&getObjectInput.IfNoneMatch, cfg.IfNoneMatch)
	util.SetIfNotZero(&getObjectInput.IfUnmodifiedSince, cfg.IfUnmodifiedSince)
	util.SetIfNotZero(&getObjectInput.PartNumber, cfg.PartNumber)

	headResp, err := c.client.HeadObject(c.ctx, headObjectInput)
	if err != nil {
		return err
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
		// d.Concurrency =
	})

	// TODO: represent download ranges
	pw := progress.NewWriter(file, uint64(*headResp.ContentLength), c.verbosity)
	defer pw.Finish()

	_, err = downloader.Download(c.ctx, pw, getObjectInput)
	if err != nil {
		return err
	}

	return nil
}
