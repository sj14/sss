package controller

import (
	"encoding/json"
	"fmt"
	"iter"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (c *Controller) BucketMultipartUploadsList(bucket, prefix string, asJson bool) error {
	for upload, err := range c.bucketMultipartUploadsList(bucket, prefix) {
		if err != nil {
			return err
		}

		if asJson {
			b, err := json.MarshalIndent(upload, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(b))
			continue
		}

		fmt.Printf("%s  %s  %s\n",
			upload.Initiated.Local().Format(time.DateTime),
			*upload.UploadId,
			*upload.Key,
		)
	}

	return nil
}

func (c *Controller) bucketMultipartUploadsList(bucket, prefix string) iter.Seq2[types.MultipartUpload, error] {
	return func(yield func(types.MultipartUpload, error) bool) {
		paginator := s3.NewListMultipartUploadsPaginator(c.client, &s3.ListMultipartUploadsInput{
			Bucket:     aws.String(bucket),
			Prefix:     &prefix,
			Delimiter:  aws.String("/"),
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
	if key == "" {
		return fmt.Errorf("empty key")
	}
	if uploadID == "" {
		return fmt.Errorf("empty upload ID")
	}

	_, err := c.client.AbortMultipartUpload(c.ctx, &s3.AbortMultipartUploadInput{
		Bucket:   &bucket,
		Key:      &key,
		UploadId: &uploadID,
	})

	return err
}

// TODO:
// - add concurrency
func (c *Controller) BucketMultipartUploadAbortAll(bucket string, dryRun bool) error {
	for upload, err := range c.bucketMultipartUploadsList(bucket, "") {
		if err != nil {
			return err
		}

		fmt.Fprintf(c.OutWriter, "deleting %s (%s)\n", *upload.Key, *upload.UploadId)

		if !dryRun {
			continue
		}

		err := c.BucketMultipartUploadAbort(bucket, *upload.Key, *upload.UploadId)
		if err != nil {
			return err
		}
	}

	return nil
}
