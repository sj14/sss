package controller

import (
	"encoding/json"
	"fmt"
	"iter"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/sync/errgroup"
)

func (c *Controller) BucketMultipartUploadsList(bucket, prefix, originalPrefix string, recursive, asJson bool) error {
	for upload, err := range c.bucketMultipartUploadsList(bucket, prefix, "/") {
		if err != nil {
			return err
		}

		for _, prefix := range upload.CommonPrefixes {
			if recursive {
				err := c.BucketMultipartUploadsList(bucket, *prefix.Prefix, originalPrefix, recursive, asJson)
				if err != nil {
					return err
				}
			} else {
				fmt.Fprintf(c.OutWriter, "%28s  %s\n", "PREFIX", *prefix.Prefix)
			}
		}

		for _, ul := range upload.Uploads {
			if asJson {
				b, err := json.MarshalIndent(ul, "", "  ")
				if err != nil {
					return err
				}

				fmt.Println(string(b))
				continue
			}

			fmt.Fprintf(c.OutWriter, "%s  %s  %s\n",
				ul.Initiated.Local().Format(time.DateTime),
				*ul.UploadId,
				strings.TrimPrefix(*ul.Key, originalPrefix),
			)
		}
	}

	return nil
}

func (c *Controller) bucketMultipartUploadsList(bucket, prefix, delimiter string) iter.Seq2[*s3.ListMultipartUploadsOutput, error] {
	return func(yield func(*s3.ListMultipartUploadsOutput, error) bool) {
		paginator := s3.NewListMultipartUploadsPaginator(c.client, &s3.ListMultipartUploadsInput{
			Bucket:     aws.String(bucket),
			Prefix:     aws.String(prefix),
			Delimiter:  aws.String(delimiter),
			MaxUploads: aws.Int32(100),
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(c.ctx)
			if !yield(page, err) {
				return
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

func (c *Controller) BucketMultipartUploadAbortAll(bucket string, dryRun bool, concurrency int) error {
	eg, _ := errgroup.WithContext(c.ctx)
	eg.SetLimit(concurrency)

	for resp, err := range c.bucketMultipartUploadsList(bucket, "", "") {
		if err != nil {
			return err
		}

		// No prefixes to handle as we don't set a delimiter,

		for _, upload := range resp.Uploads {
			fmt.Fprintf(c.OutWriter, "deleting %s (%s)\n", *upload.Key, *upload.UploadId)

			if !dryRun {
				continue
			}

			eg.Go(func() error {
				return c.BucketMultipartUploadAbort(bucket, *upload.Key, *upload.UploadId)
			})
		}
	}

	return eg.Wait()
}
