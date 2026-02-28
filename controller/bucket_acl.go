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

func (c *Controller) BucketACLGet(bucket string) error {
	resp, err := c.client.GetBucketAcl(c.ctx, &s3.GetBucketAclInput{
		Bucket: aws.String(bucket),
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

func (c *Controller) BucketACLPut(aclPath, bucket string) error {
	lBytes, err := os.ReadFile(aclPath)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(bytes.NewBuffer(lBytes))
	dec.DisallowUnknownFields()

	var policy *types.AccessControlPolicy
	if err := dec.Decode(&policy); err != nil {
		return fmt.Errorf("failed to unmarshal configuration file: %w", err)
	}

	_, err = c.client.PutBucketAcl(c.ctx, &s3.PutBucketAclInput{
		Bucket:              aws.String(bucket),
		AccessControlPolicy: policy,
	})
	if err != nil {
		return err
	}

	return nil
}
