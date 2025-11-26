package controller

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *Controller) BucketPolicyGet(bucket string) error {
	resp, err := c.client.GetBucketPolicy(c.ctx, &s3.GetBucketPolicyInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}

	fmt.Println(*resp.Policy)

	return nil
}

func (c *Controller) BucketPolicyPut(policyPath, bucket string) error {
	pBytes, err := os.ReadFile(policyPath)
	if err != nil {
		return err
	}

	_, err = c.client.PutBucketPolicy(c.ctx, &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucket),
		Policy: aws.String(string(pBytes)),
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) BucketPolicyDelete(bucket string) error {
	_, err := c.client.DeleteBucketPolicy(c.ctx, &s3.DeleteBucketPolicyInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}

	return nil
}
