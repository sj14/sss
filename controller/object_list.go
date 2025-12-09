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

func (c *Controller) ObjectList(bucket, prefix, originalPrefix string, recursive, asJson bool) error {
	for l, err := range c.objectList(bucket, prefix) {
		if err != nil {
			return err
		}

		for _, prefix := range l.CommonPrefixes {
			if recursive {
				err := c.ObjectList(bucket, *prefix.Prefix, originalPrefix, recursive, asJson)
				if err != nil {
					return err
				}
			} else {
				fmt.Fprintf(c.OutWriter, "%28s  %s\n", "PREFIX", *prefix.Prefix)
			}
		}

		for _, object := range l.Contents {
			if asJson {
				b, err := json.Marshal(object)
				if err != nil {
					return err
				}
				fmt.Println(string(b))
				continue
			}
			fmt.Fprintf(c.OutWriter, "%s %8s  %s\n",
				object.LastModified.Local().Format(time.DateTime),
				humanize.IBytes(uint64(*object.Size)),
				strings.TrimPrefix(*object.Key, originalPrefix),
			)
		}
	}

	return nil
}

func (c *Controller) objectList(bucket, prefix string) iter.Seq2[*s3.ListObjectsV2Output, error] {
	return func(yield func(*s3.ListObjectsV2Output, error) bool) {
		paginator := s3.NewListObjectsV2Paginator(c.client, &s3.ListObjectsV2Input{
			Bucket:    aws.String(bucket),
			Prefix:    aws.String(prefix),
			Delimiter: aws.String("/"),
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
