package controller

import (
	"fmt"

	"github.com/dustin/go-humanize"
)

func (c *Controller) BucketSize(bucket, prefix, delimiter string) error {
	versions, prefixes, err := c.objectVersions(bucket, prefix, delimiter)
	if err != nil {
		return err
	}

	for _, cp := range prefixes {
		fmt.Printf("%28s  %s\n", "PREFIX", *cp.Prefix)
	}

	var (
		sizeCurrent   uint64 = 0
		sizeVersioned uint64 = 0
	)

	for _, ver := range versions {
		if *ver.IsLatest {
			sizeCurrent += uint64(*ver.Size)
			continue
		}
		sizeVersioned += uint64(*ver.Size)
	}

	fmt.Printf("current:  %v\n", humanize.Bytes(sizeCurrent))
	fmt.Printf("versions: %v\n", humanize.Bytes(sizeVersioned))
	fmt.Printf("total:    %v\n", humanize.Bytes(sizeCurrent+sizeVersioned))

	return nil
}
