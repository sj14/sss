package controller

import (
	"encoding/json"
	"fmt"
	"iter"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/dustin/go-humanize"
)

func (c *Controller) BucketPartsList(bucket, key, uploadID string, asJson bool) error {
	if key == "" {
		return fmt.Errorf("empty key")
	}
	if uploadID == "" {
		return fmt.Errorf("empty upload ID")
	}

	for part, err := range c.bucketPartsList(bucket, key, uploadID) {
		if err != nil {
			return err
		}

		if asJson {
			b, err := json.MarshalIndent(part, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(b))
			continue
		}

		fmt.Fprintf(c.OutWriter, "%s  #%d  %8s  %s\n",
			part.LastModified.Local().Format(time.DateTime),
			*part.PartNumber,
			humanize.IBytes(uint64(*part.Size)),
			*part.ETag,
		)
	}

	return nil

}

func (c *Controller) bucketPartsList(bucket, key, uploadID string) iter.Seq2[types.Part, error] {
	return func(yield func(types.Part, error) bool) {
		paginator := s3.NewListPartsPaginator(c.client, &s3.ListPartsInput{
			Bucket:   aws.String(bucket),
			Key:      aws.String(key),
			UploadId: aws.String(uploadID),
			MaxParts: aws.Int32(100),
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(c.ctx)
			if err != nil {
				yield(types.Part{}, err)
				return
			}

			for _, p := range page.Parts {
				if !yield(p, nil) {
					return
				}
			}
		}
	}
}
