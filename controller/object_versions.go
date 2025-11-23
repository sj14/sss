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

func (c *Controller) ObjectVersions(bucket, prefix string, asJson bool) error {
	for v, err := range c.objectVersions(bucket, prefix) {
		if err != nil {
			return err
		}

		if v.Prefix != nil {
			fmt.Printf("%61s  %s\n", "PREFIX", *v.Prefix.Prefix)
		}

		if v.Versions != nil {
			if asJson {
				b, err := json.MarshalIndent(v.Versions, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(b))
				continue
			}
			fmt.Printf("%s  %s %8s  %s\n",
				v.Versions.LastModified.Local().Format(time.DateTime),
				*v.Versions.VersionId,
				humanize.IBytes(uint64(*v.Versions.Size)),
				strings.TrimPrefix(*v.Versions.Key, prefix),
			)
		}
	}

	return nil
}

type VersionsItem struct {
	Versions *types.ObjectVersion
	Prefix   *types.CommonPrefix
}

func (c *Controller) objectVersions(bucket, prefix string) iter.Seq2[VersionsItem, error] {
	return func(yield func(VersionsItem, error) bool) {
		paginator := s3.NewListObjectVersionsPaginator(c.client, &s3.ListObjectVersionsInput{
			Bucket:    aws.String(bucket),
			Delimiter: aws.String("/"),
			Prefix:    aws.String(prefix),
			MaxKeys:   aws.Int32(100),
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(c.ctx)
			if err != nil {
				yield(VersionsItem{}, err)
				return
			}

			for _, p := range page.CommonPrefixes {
				if !yield(VersionsItem{Prefix: &p}, nil) {
					return
				}
			}

			for _, v := range page.Versions {
				if !yield(VersionsItem{Versions: &v}, nil) {
					return
				}
			}
		}
	}
}
