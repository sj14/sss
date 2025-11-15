package controller

import (
	"encoding/json"
	"fmt"
	"iter"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (c *Controller) BucketPartsList(bucket, key, uploadID string) error {
	for part, err := range c.bucketPartsList(bucket, key, uploadID) {
		if err != nil {
			return err
		}

		b, err := json.MarshalIndent(part, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(b))
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
