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

	objects, prefixes, err := c.objectList(cfg.Bucket, key, cfg.Delimiter)
	if err != nil {
		return err
	}

	for _, prefix := range prefixes {
		err := c.ObjectDelete(*prefix.Prefix, cfg)
		if err != nil {
			return err
		}
	}

	eg, _ := errgroup.WithContext(c.ctx)
	eg.SetLimit(cfg.Concurrency)

	for _, object := range objects {
		eg.Go(func() error {
			fmt.Printf("deleting %s (%s)\n", *object.Key, humanize.Bytes(uint64(*object.Size)))
			err := c.objectDelete(cfg.Bucket, *object.Key)
			if err != nil {
				return err
			}
			return nil
		})
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
