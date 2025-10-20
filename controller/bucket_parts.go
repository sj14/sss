package controller

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *Controller) BucketPartsList(bucket, key, uploadID string) error {
	paginator := s3.NewListPartsPaginator(c.client, &s3.ListPartsInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadID),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(c.ctx)
		if err != nil {
			return err
		}
		for _, part := range page.Parts {
			b, err := json.MarshalIndent(part, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(b))
		}
	}

	return nil
}
