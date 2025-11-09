package controller

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (c *Controller) BucketObjectLockGet(bucket string) error {
	resp, err := c.client.GetObjectLockConfiguration(c.ctx, &s3.GetObjectLockConfigurationInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(resp.ObjectLockConfiguration, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}

func (c *Controller) BucketObjectLockPut(lockConfigPath, bucket string) error {
	lBytes, err := os.ReadFile(lockConfigPath)
	if err != nil {
		return err
	}

	var lockConfiguration *types.ObjectLockConfiguration
	if err := json.Unmarshal(lBytes, &lockConfiguration); err != nil {
		return fmt.Errorf("failed to unmarshal lifecycle policy: %w", err)
	}

	_, err = c.client.PutObjectLockConfiguration(c.ctx, &s3.PutObjectLockConfigurationInput{
		Bucket:                  aws.String(bucket),
		ObjectLockConfiguration: lockConfiguration,
	})
	if err != nil {
		return err
	}

	return nil
}
