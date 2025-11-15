package controller

import (
	"encoding/json"
	"fmt"
	"iter"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (c *Controller) BucketMultipartUploadsList(bucket, prefix, delimiter string) error {
	for upload, err := range c.bucketMultipartUploadsList(bucket, prefix, delimiter) {
		if err != nil {
			return err
		}

		b, err := json.MarshalIndent(upload, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(b))
	}

	return nil
}

func (c *Controller) bucketMultipartUploadsList(bucket, prefix, delimiter string) iter.Seq2[types.MultipartUpload, error] {
	return func(yield func(types.MultipartUpload, error) bool) {
		paginator := s3.NewListMultipartUploadsPaginator(c.client, &s3.ListMultipartUploadsInput{
			Bucket:     aws.String(bucket),
			Prefix:     &prefix,
			Delimiter:  &delimiter,
			MaxUploads: aws.Int32(100),
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(c.ctx)
			if err != nil {
				yield(types.MultipartUpload{}, err)
				return
			}

			for _, p := range page.Uploads {
				if !yield(p, nil) {
					return
				}
			}
		}
	}
}

func (c *Controller) BucketMultipartUploadAbort(bucket, key, uploadID string) error {
	_, err := c.client.AbortMultipartUpload(c.ctx, &s3.AbortMultipartUploadInput{
		Bucket:   &bucket,
		Key:      &key,
		UploadId: &uploadID,
	})

	return err
}
