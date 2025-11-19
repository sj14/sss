package controller

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/dustin/go-humanize"
	"golang.org/x/sync/errgroup"
)

type ObjectDeleteConfig struct {
	Bucket      string
	Delimiter   string
	Force       bool
	Concurrency int
}

func (c *Controller) ObjectDelete(key string, cfg ObjectDeleteConfig) error {
	if key == "" && !cfg.Force {
		return errors.New("use -force flag to empty the whole bucket")
	}

	eg, _ := errgroup.WithContext(c.ctx)
	eg.SetLimit(cfg.Concurrency)

	for l, err := range c.objectList(cfg.Bucket, key, cfg.Delimiter) {
		if err != nil {
			return err
		}

		if l.Prefix != nil {
			err := c.ObjectDelete(*l.Prefix.Prefix, cfg)
			if err != nil {
				return err
			}
		}

		if l.Object != nil {
			eg.Go(func() error {
				fmt.Printf("deleting %s (%s)\n", *l.Object.Key, humanize.IBytes(uint64(*l.Object.Size)))
				err := c.objectDelete(cfg.Bucket, *l.Object.Key)
				if err != nil {
					return err
				}
				return nil
			})
		}
	}

	return eg.Wait()
}

func (c *Controller) objectDelete(bucket, key string) error {
	_, err := c.client.DeleteObject(c.ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}
