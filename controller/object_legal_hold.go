package controller

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/sj14/sss/util"
)

func (c *Controller) ObjectGetLegalHold(bucket, key, versionID string) error {
	input := &s3.GetObjectLegalHoldInput{
		Bucket:    aws.String(bucket),
		Key:       aws.String(key),
		VersionId: util.NilIfZero(versionID),
	}

	resp, err := c.client.GetObjectLegalHold(c.ctx, input)
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

func (c *Controller) ObjectPutLegalHold(bucket, key, status, versionID string) error {
	status = strings.ToUpper(status)

	if status != "ON" && status != "OFF" {
		return fmt.Errorf("allowed status: 'on|off', got %q", status)
	}

	input := &s3.PutObjectLegalHoldInput{
		Bucket:    aws.String(bucket),
		Key:       aws.String(key),
		VersionId: util.NilIfZero(versionID),
		LegalHold: &types.ObjectLockLegalHold{Status: types.ObjectLockLegalHoldStatus(status)},
	}

	_, err := c.client.PutObjectLegalHold(c.ctx, input)
	if err != nil {
		return err
	}

	return nil
}
