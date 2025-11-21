package controller

import (
	"errors"
	"fmt"
	"log"

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
	DryRun      bool
}

func (c *Controller) ObjectDelete(key string, cfg ObjectDeleteConfig) (err error) {
	if cfg.DryRun {
		fmt.Println("> dry-run mode <")
	}
	if key == "" && !cfg.Force {
		return errors.New("use -force flag to empty the whole bucket")
	}

	eg, _ := errgroup.WithContext(c.ctx)
	eg.SetLimit(cfg.Concurrency)

	defer func() {
		// make sure to wait even when we return early somwhere
		e := eg.Wait()
		if err == nil {
			err = e
		} else {
			log.Println(e)
		}
	}()

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
				err := c.objectDelete(cfg.DryRun, cfg.Bucket, *l.Object.Key)
				if err != nil {
					return err
				}
				return nil
			})

			exactMatch := key == *l.Object.Key
			if exactMatch {
				// Single file deletiong, mimicing "normal" behaviour.
				// e.g. ls => "file", "file1"
				// Without this check, "file1" would also be deleted
				// when only "file" is requested.
				// As an alternative, add a -recursive flag or similar.
				break
			}
		}
	}

	return eg.Wait()
}

func (c *Controller) objectDelete(dryRun bool, bucket, key string) error {
	if dryRun {
		return nil
	}
	_, err := c.client.DeleteObject(c.ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}
