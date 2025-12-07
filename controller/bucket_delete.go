package controller

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *Controller) BucketDelete(bucket string, flagForce bool) error {
	if bucket == "" {
		return fmt.Errorf("no bucket name specified")
	}
	if !flagForce {
		return fmt.Errorf("force flag required")
	}
	_, err := c.client.DeleteBucket(c.ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	return err
}
