package controller

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (c *Controller) BucketPartsList(bucket, key, uploadID string) error {
	parts, err := c.bucketPartsList(bucket, key, uploadID)
	if err != nil {
		return err
	}

	for _, part := range parts {
		b, err := json.MarshalIndent(part, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(b))
	}

	return nil
}

func (c *Controller) bucketPartsList(bucket, key, uploadID string) ([]types.Part, error) {
	paginator := s3.NewListPartsPaginator(c.client, &s3.ListPartsInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadID),
	})

	var result []types.Part
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(c.ctx)
		if err != nil {
			return nil, err
		}

		result = append(result, page.Parts...)
	}

	return result, nil
}
