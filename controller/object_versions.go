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

func (c *Controller) ObjectVersions(bucket, prefix, delimiter string) error {
	versions, prefixes, err := c.objectVersions(bucket, prefix, delimiter)
	if err != nil {
		return err
	}

	for _, cp := range prefixes {
		fmt.Printf("%28s  %s\n", "PREFIX", *cp.Prefix)
	}

	for _, ver := range versions {
		fmt.Printf("%s %8s  %s  %s\n",
			ver.LastModified.Local().Format(time.DateTime),
			*ver.VersionId,
			humanize.Bytes(uint64(*ver.Size)),
			strings.TrimPrefix(*ver.Key, prefix),
		)
	}

	return nil
}

func (c *Controller) objectVersions(bucket, prefix, delimiter string) ([]types.ObjectVersion, []types.CommonPrefix, error) {
	if bucket == "" {
		return nil, nil, fmt.Errorf("missing bucket")
	}

	paginator := s3.NewListObjectVersionsPaginator(c.client, &s3.ListObjectVersionsInput{
		Bucket:    aws.String(bucket),
		Delimiter: aws.String(delimiter),
		Prefix:    aws.String(prefix),
	})

	var versions []types.ObjectVersion
	var prefixes []types.CommonPrefix
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(c.ctx)
		if err != nil {
			return nil, nil, err
		}

		prefixes = append(prefixes, page.CommonPrefixes...)
		versions = append(versions, page.Versions...)
	}

	return versions, prefixes, nil
}
