package controller

import (
	"encoding/json"
	"fmt"
	"iter"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/dustin/go-humanize"
)

func (c *Controller) ObjectList(bucket, prefix, originalPrefix, delimiter string, recursive, asJson bool) error {
	for l, err := range c.objectList(bucket, prefix, delimiter) {
		if err != nil {
			return err
		}

		if l.Prefix != nil {
			if recursive {
				err := c.ObjectList(bucket, *l.Prefix.Prefix, originalPrefix, delimiter, recursive, asJson)
				if err != nil {
					return err
				}
			} else {
				fmt.Printf("%28s  %s\n", "PREFIX", *l.Prefix.Prefix)
			}
		}

		if l.Object != nil {
			if asJson {
				b, err := json.MarshalIndent(l.Object, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(b))
				continue
			}
			fmt.Printf("%s %8s  %s\n",
				l.Object.LastModified.Local().Format(time.DateTime),
				humanize.IBytes(uint64(*l.Object.Size)),
				strings.TrimPrefix(*l.Object.Key, originalPrefix),
			)
		}
	}

	return nil
}

type ListItem struct {
	Object *types.Object
	Prefix *types.CommonPrefix
}

func (c *Controller) objectList(bucket, prefix, delimiter string) iter.Seq2[ListItem, error] {
	return func(yield func(ListItem, error) bool) {
		paginator := s3.NewListObjectsV2Paginator(c.client, &s3.ListObjectsV2Input{
			Bucket:    aws.String(bucket),
			Delimiter: aws.String(delimiter),
			Prefix:    aws.String(prefix),
			MaxKeys:   aws.Int32(100),
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(c.ctx)
			if err != nil {
				yield(ListItem{}, err)
				return
			}

			for _, p := range page.CommonPrefixes {
				if !yield(ListItem{Prefix: &p}, nil) {
					return
				}
			}

			for _, o := range page.Contents {
				if !yield(ListItem{Object: &o}, nil) {
					return
				}
			}
		}
	}
}
