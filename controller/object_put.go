package controller

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/sj14/sss/progress"
	"github.com/sj14/sss/util"
)

type ObjectPutConfig struct {
	Bucket            string
	SSEC              util.SSEC
	Concurrency       int
	LeavePartsOnError bool
	MaxUploadParts    int32
	PartSize          int64
	ACL               string
}

func (c *Controller) ObjectPut(filePath string, target string, cfg ObjectPutConfig) error {
	// log.Fatalln(filePaths)

	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return c.objectPut(filePath, target, cfg)
	}

	return filepath.Walk(filePath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		fmt.Printf("uploading %s\n", p)
		err = c.objectPut(p, path.Join(target, filepath.Dir(p), path.Base(p)), cfg)
		if err != nil {
			return err
		}

		return nil
	})
}

func (c *Controller) objectPut(filePath, key string, cfg ObjectPutConfig) error {
	if key == "" {
		key = path.Base(filePath)
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

	pr := progress.NewReader(f, uint64(stat.Size()), c.verbosity)
	defer pr.Finish()

	putObjectInput := &s3.PutObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(key),
		Body:   pr,
		ACL:    types.ObjectCannedACL(cfg.ACL),
	}

	if cfg.SSEC.KeyIsSet() {
		putObjectInput.SSECustomerKey = aws.String(cfg.SSEC.Base64Key())
		putObjectInput.SSECustomerKeyMD5 = aws.String(cfg.SSEC.Base64KeyMD5())
		putObjectInput.SSECustomerAlgorithm = aws.String(cfg.SSEC.Algorithm())
	}

	_, err = uploader.Upload(c.ctx, putObjectInput)
	if err != nil {
		return err
	}

	// b, err := json.MarshalIndent(resp, "", "  ")
	// if err != nil {
	// 	return err
	// }

	// fmt.Println(string(b))

	return nil
}
