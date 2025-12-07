package controller

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (c *Controller) BucketCORSGet(bucket string) error {
	resp, err := c.client.GetBucketCors(c.ctx, &s3.GetBucketCorsInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(resp.CORSRules, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}

func (c *Controller) BucketCORSPut(corsPath, bucket string) error {
	lBytes, err := os.ReadFile(corsPath)
	if err != nil {
		return err
	}

	var corsConfig *types.CORSConfiguration
	if err := json.Unmarshal(lBytes, &corsConfig); err != nil {
		return fmt.Errorf("failed to unmarshal lifecycle policy: %w", err)
	}

	_, err = c.client.PutBucketCors(c.ctx, &s3.PutBucketCorsInput{
		Bucket:            aws.String(bucket),
		CORSConfiguration: corsConfig,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) BucketCORSDelete(bucket string) error {
	_, err := c.client.DeleteBucketCors(c.ctx, &s3.DeleteBucketCorsInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}

	return nil
}
