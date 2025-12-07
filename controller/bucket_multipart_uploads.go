package controller

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *Controller) BucketMultipartUploadsList(bucket string) error {
	paginator := s3.NewListMultipartUploadsPaginator(c.client, &s3.ListMultipartUploadsInput{
		Bucket: aws.String(bucket),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(c.ctx)
		if err != nil {
			return err
		}
		for _, upload := range page.Uploads {
			b, err := json.MarshalIndent(upload, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(b))
		}
	}

	return nil
}

func (c *Controller) BucketMultipartUploadAbort(bucket, key, uploadID string) error {
	_, err := c.client.AbortMultipartUpload(c.ctx, &s3.AbortMultipartUploadInput{
		Bucket:   &bucket,
		Key:      &key,
		UploadId: &uploadID,
	})

	return err
}
