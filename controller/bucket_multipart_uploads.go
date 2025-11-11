package controller

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (c *Controller) BucketMultipartUploadsList(bucket string) error {
	uploads, err := c.bucketMultipartUploadsList(bucket)
	if err != nil {
		return err
	}

	for _, upload := range uploads {
		b, err := json.MarshalIndent(upload, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(b))
	}
	return nil
}

func (c *Controller) bucketMultipartUploadsList(bucket string) ([]types.MultipartUpload, error) {
	paginator := s3.NewListMultipartUploadsPaginator(c.client, &s3.ListMultipartUploadsInput{
		Bucket: aws.String(bucket),
	})

	var result []types.MultipartUpload
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(c.ctx)
		if err != nil {
			return nil, err
		}

		result = append(result, page.Uploads...)
	}

	return result, nil
}

func (c *Controller) BucketMultipartUploadAbort(bucket, key, uploadID string) error {
	_, err := c.client.AbortMultipartUpload(c.ctx, &s3.AbortMultipartUploadInput{
		Bucket:   &bucket,
		Key:      &key,
		UploadId: &uploadID,
	})

	return err
}
