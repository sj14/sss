package controller

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (c *Controller) BucketLifecycleGet(bucket string) error {
	resp, err := c.client.GetBucketLifecycleConfiguration(c.ctx, &s3.GetBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(resp.Rules, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}

func (c *Controller) BucketLifecyclePut(lifecyclePath, bucket string) error {
	lBytes, err := os.ReadFile(lifecyclePath)
	if err != nil {
		return err
	}

	var lifecycleConfig *types.BucketLifecycleConfiguration
	if err := json.Unmarshal(lBytes, &lifecycleConfig); err != nil {
		return fmt.Errorf("failed to unmarshal lifecycle policy: %w", err)
	}

	_, err = c.client.PutBucketLifecycleConfiguration(c.ctx, &s3.PutBucketLifecycleConfigurationInput{
		Bucket:                 aws.String(bucket),
		LifecycleConfiguration: lifecycleConfig,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) BucketLifecycleDelete(bucket string) error {
	_, err := c.client.DeleteBucketLifecycle(c.ctx, &s3.DeleteBucketLifecycleInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}

	return nil
}
