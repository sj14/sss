package controller

import (
	"fmt"
	"log"

	"github.com/dustin/go-humanize"
)

func (c *Controller) BucketSize(bucket, prefix string) error {
	var (
		sizeCurrent    uint64
		sizeVersioned  uint64
		sizeMultiparts uint64

		countCurrent    uint64
		countVersioned  uint64
		countMultiparts uint64
	)

	for item, err := range c.objectVersions(bucket, prefix, "") {
		if err != nil {
			return err
		}

		if item.Versions == nil {
			continue
		}

		if *item.Versions.IsLatest {
			sizeCurrent += uint64(*item.Versions.Size)
			countCurrent++
			continue
		}

		sizeVersioned += uint64(*item.Versions.Size)
		countVersioned++
	}

	fmt.Printf("current: %v (%d)", humanize.IBytes(sizeCurrent), countCurrent)
	fmt.Printf(" | versions: %v (%d)", humanize.IBytes(sizeVersioned), countVersioned)

	for upload, err := range c.bucketMultipartUploadsList(bucket, prefix) {
		if err != nil {
			return err
		}

		for part, err := range c.bucketPartsList(bucket, *upload.Key, *upload.UploadId) {
			if err != nil {
				log.Println(err)
				continue
			}

			sizeMultiparts += uint64(*part.Size)
			countMultiparts++
		}
	}

	var (
		totalByte  = humanize.IBytes(sizeCurrent + sizeVersioned + sizeMultiparts)
		totalCount = countCurrent + countVersioned + countMultiparts
	)

	fmt.Printf(" | multiparts: %v (%d)", humanize.IBytes(sizeMultiparts), countMultiparts)
	fmt.Printf(" | total: %v (%d)\n", totalByte, totalCount)

	return nil
}
