package controller

import (
	"encoding/json"
	"fmt"
	"iter"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/dustin/go-humanize"
)

func (c *Controller) ObjectVersions(bucket, prefix, originalPrefix string, recursive, asJson bool) error {
	for resp, err := range c.objectVersions(bucket, prefix, "/") {
		if err != nil {
			return err
		}

		for _, prefix := range resp.CommonPrefixes {
			if recursive {
				err := c.ObjectVersions(bucket, *prefix.Prefix, originalPrefix, recursive, asJson)
				if err != nil {
					return err
				}
			} else {
				fmt.Printf("%61s  %s\n", "PREFIX", *prefix.Prefix)
			}
		}

		for _, v := range resp.Versions {
			if asJson {
				b, err := json.MarshalIndent(v, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(b))
				continue
			}
			fmt.Printf("%s  %s %8s  %s\n",
				v.LastModified.Local().Format(time.DateTime),
				*v.VersionId,
				humanize.IBytes(uint64(*v.Size)),
				strings.TrimPrefix(*v.Key, originalPrefix),
			)
		}
	}

	return nil
}

func (c *Controller) objectVersions(bucket, prefix, delimiter string) iter.Seq2[*s3.ListObjectVersionsOutput, error] {
	return func(yield func(*s3.ListObjectVersionsOutput, error) bool) {
		paginator := s3.NewListObjectVersionsPaginator(c.client, &s3.ListObjectVersionsInput{
			Bucket:    aws.String(bucket),
			Delimiter: aws.String(delimiter),
			Prefix:    aws.String(prefix),
			MaxKeys:   aws.Int32(100),
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(c.ctx)
			if !yield(page, err) {
				return
			}
		}
	}
}
