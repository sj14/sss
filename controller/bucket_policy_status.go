package controller

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *Controller) BucketPolicyStatusGet(bucket string) error {
	resp, err := c.client.GetBucketPolicyStatus(c.ctx, &s3.GetBucketPolicyStatusInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(resp.PolicyStatus, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}
