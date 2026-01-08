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

		for _, version := range item.Versions {
			if *version.IsLatest {
				sizeCurrent += uint64(*version.Size)
				countCurrent++
				continue
			}

			sizeVersioned += uint64(*version.Size)
			countVersioned++
		}
	}

	fmt.Fprintf(c.OutWriter, "current: %v (%d)", humanize.IBytes(sizeCurrent), countCurrent)
	fmt.Fprintf(c.OutWriter, " | versions: %v (%d)", humanize.IBytes(sizeVersioned), countVersioned)

	for uploads, err := range c.multipartUploadsList(bucket, prefix, "") {
		if err != nil {
			return err
		}

		for _, upload := range uploads.Uploads {
			for part, err := range c.partsList(bucket, *upload.Key, *upload.UploadId) {
				if err != nil {
					log.Println(err)
					continue
				}

				sizeMultiparts += uint64(*part.Size)
				countMultiparts++
			}
		}
	}

	var (
		totalByte  = humanize.IBytes(sizeCurrent + sizeVersioned + sizeMultiparts)
		totalCount = countCurrent + countVersioned + countMultiparts
	)

	fmt.Fprintf(c.OutWriter, " | multiparts: %v (%d)", humanize.IBytes(sizeMultiparts), countMultiparts)
	fmt.Fprintf(c.OutWriter, " | total: %v (%d)\n", totalByte, totalCount)

	return nil
}
