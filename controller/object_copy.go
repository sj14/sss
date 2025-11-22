package controller

import (
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sj14/sss/util"
)

type ObjectCopyConfig struct {
	SrcBucket string
	SrcKey    string
	DstBucket string
	DstKey    string
	SSEC      util.SSEC
}

func (c *Controller) ObjectCopy(cfg ObjectCopyConfig) error {
	if cfg.DstBucket == "" {
		cfg.DstBucket = cfg.SrcBucket
	}
	if cfg.DstKey == "" {
		cfg.DstKey = cfg.SrcKey
	}

	input := &s3.CopyObjectInput{
		Bucket:     aws.String(cfg.DstBucket),
		CopySource: aws.String(path.Join(cfg.SrcBucket, cfg.SrcKey)),
		Key:        aws.String(cfg.DstKey),
	}

	if cfg.SSEC.KeyIsSet() {
		input.SSECustomerAlgorithm = aws.String(cfg.SSEC.Algorithm())
		input.SSECustomerKey = aws.String(cfg.SSEC.Base64Key())
		input.SSECustomerKeyMD5 = aws.String(cfg.SSEC.Base64KeyMD5())
	}

	_, err := c.client.CopyObject(c.ctx, input)
	if err != nil {
		return err
	}

	return nil
}
