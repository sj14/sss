package controller

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/dustin/go-humanize"
)

func (c *Controller) ObjectList(bucket, prefix, delimiter string) error {
	objects, prefixes, err := c.objectList(bucket, prefix, delimiter)
	if err != nil {
		return err
	}

	for _, cp := range prefixes {
		fmt.Printf("%28s  %s\n", "PREFIX", *cp.Prefix)
	}

	for _, obj := range objects {
		fmt.Printf("%s %8s  %s\n",
			obj.LastModified.Local().Format(time.DateTime),
			humanize.Bytes(uint64(*obj.Size)),
			strings.TrimPrefix(*obj.Key, prefix),
		)
	}

	return nil
}

func (c *Controller) objectList(bucket, prefix, delimiter string) ([]types.Object, []types.CommonPrefix, error) {
	if bucket == "" {
		return nil, nil, fmt.Errorf("missing bucket")
	}

	paginator := s3.NewListObjectsV2Paginator(c.client, &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Delimiter: aws.String(delimiter),
		Prefix:    aws.String(prefix),
	})

	var objects []types.Object
	var prefixes []types.CommonPrefix
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(c.ctx)
		if err != nil {
			return nil, nil, err
		}

		prefixes = append(prefixes, page.CommonPrefixes...)
		objects = append(objects, page.Contents...)
	}

	return objects, prefixes, nil
}
