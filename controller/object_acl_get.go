package controller

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sj14/sss/util"
)

func (c *Controller) ObjectACLGet(bucket, key, version string) error {
	resp, err := c.client.GetObjectAcl(c.ctx, &s3.GetObjectAclInput{
		Bucket:    aws.String(bucket),
		Key:       aws.String(key),
		VersionId: util.NilIfZero(version),
	})
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
