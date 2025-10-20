package controller

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *Controller) BucketDelete(bucket string) error {
	if bucket == "" {
		return fmt.Errorf("no bucket name specified")
	}
	_, err := c.client.DeleteBucket(c.ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	return err
}
