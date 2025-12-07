package controller

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// TODO: prefix doesn't seem to work
func (c *Controller) BucketList(prefix string) error {
	paginator := s3.NewListBucketsPaginator(c.client, &s3.ListBucketsInput{
		Prefix: aws.String(prefix),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(c.ctx)
		if err != nil {
			return err
		}
		for _, bucket := range page.Buckets {
			fmt.Fprintf(c.OutWriter, "%s %s\n", bucket.CreationDate.Local().Format(time.DateTime), *bucket.Name)
		}
	}

	return nil
}
