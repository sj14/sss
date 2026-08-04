package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (c *Controller) BucketPublicAccessBlockGet(bucket string) error {
	resp, err := c.client.GetPublicAccessBlock(c.ctx, &s3.GetPublicAccessBlockInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(resp.PublicAccessBlockConfiguration, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}

func (c *Controller) BucketPublicAccessBlockPut(cfgPath, bucket string) error {
	cfgBytes, err := os.ReadFile(cfgPath)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(bytes.NewBuffer(cfgBytes))
	dec.DisallowUnknownFields()

	var config *types.PublicAccessBlockConfiguration
	if err := dec.Decode(&config); err != nil {
		return fmt.Errorf("failed to unmarshal configuration file: %w", err)
	}

	_, err = c.client.PutPublicAccessBlock(c.ctx, &s3.PutPublicAccessBlockInput{
		Bucket:                         aws.String(bucket),
		PublicAccessBlockConfiguration: config,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) BucketPublicAccessBlockDelete(bucket string) error {
	_, err := c.client.DeletePublicAccessBlock(c.ctx, &s3.DeletePublicAccessBlockInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}

	return nil
}
