package controller

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sj14/sss/util"
)

type ObjectHeadConfig struct {
	SSEC util.SSEC
}

func (c *Controller) ObjectHead(bucket, key string, cfg ObjectHeadConfig) error {
	headObjectInput := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	if cfg.SSEC.KeyIsSet() {
		headObjectInput.SSECustomerKeyMD5 = aws.String(cfg.SSEC.Base64KeyMD5())
		headObjectInput.SSECustomerKey = aws.String(cfg.SSEC.Base64Key())
		headObjectInput.SSECustomerAlgorithm = aws.String(cfg.SSEC.Algorithm())

	}
	resp, err := c.client.HeadObject(c.ctx, headObjectInput)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}
