package controller

import (
	"fmt"
	"log"

	"github.com/dustin/go-humanize"
)

func (c *Controller) BucketSize(bucket, prefix, delimiter string) error {
	var (
		sizeCurrent    uint64
		sizeVersioned  uint64
		sizeMultiparts uint64

		countCurrent    uint64
		countVersioned  uint64
		countMultiparts uint64
	)

	for item, err := range c.objectVersions(bucket, prefix, delimiter) {
		if err != nil {
			return err
		}

		if item.Prefix != nil {
			fmt.Printf("%28s  %s\n", "PREFIX", *item.Prefix.Prefix)
		}

		if *item.Versions.IsLatest {
			sizeCurrent += uint64(*item.Versions.Size)
			countCurrent++
			continue
		}
		sizeVersioned += uint64(*item.Versions.Size)
		countVersioned++
	}

	fmt.Printf("current: %v (%d)", humanize.Bytes(sizeCurrent), countCurrent)
	fmt.Printf(" | versions: %v (%d)", humanize.Bytes(sizeVersioned), countVersioned)

	uploads, err := c.bucketMultipartUploadsList(bucket)
	if err != nil {
		return err
	}

	for _, upload := range uploads {
		parts, err := c.bucketPartsList(bucket, *upload.Key, *upload.UploadId)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, part := range parts {
			sizeMultiparts += uint64(*part.Size)
			countMultiparts++
		}
	}

	var (
		totalByte  = humanize.Bytes(sizeCurrent + sizeVersioned + sizeMultiparts)
		totalCount = countCurrent + countVersioned + countMultiparts
	)

	fmt.Printf(" | multiparts: %v (%d)", humanize.Bytes(sizeMultiparts), countMultiparts)
	fmt.Printf(" | total: %v (%d)\n", totalByte, totalCount)

	return nil
}
