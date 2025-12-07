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
		resp, err := c.client.ListObjectVersions(c.ctx, &s3.ListObjectVersionsInput{
			Bucket:  &bucket,
			MaxKeys: aws.Int32(1),
		})
		if err != nil {
			return fmt.Errorf("failed to check if bucket is empty: %w", err)
		}
		if len(resp.Versions) > 0 {
			return fmt.Errorf("bucket not empty, use 'force' flag to delete the bucket")
		}
	}
	_, err := c.client.DeleteBucket(c.ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	return err
}
