package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/sj14/sss/util"
)

func (c *Controller) ObjectACLGet(bucket, key, version string) error {
	resp, err := c.client.GetObjectAcl(c.ctx, &s3.GetObjectAclInput{
		Bucket:    aws.String(bucket),
		Key:       aws.String(key),
		VersionId: util.NilIfZero(version),
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

func (c *Controller) ObjectACLPut(aclPath, bucket, key, version string) error {
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

	_, err = c.client.PutObjectAcl(c.ctx, &s3.PutObjectAclInput{
		Bucket:              aws.String(bucket),
		Key:                 aws.String(key),
		VersionId:           util.NilIfZero(version),
		AccessControlPolicy: policy,
	})
	if err != nil {
		return err
	}

	return nil
}
