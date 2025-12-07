package controller

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/dustin/go-humanize"
	"golang.org/x/sync/errgroup"
)

type ObjectDeleteConfig struct {
	Bucket           string
	Force            bool
	Concurrency      int
	DryRun           bool
	BypassGovernance bool
	VersionID        string
}

func (c *Controller) ObjectDelete(prefix string, cfg ObjectDeleteConfig) error {
	if prefix == "" {
		return errors.New("missing key")
	}
	if prefix == "/" && !cfg.Force && !cfg.DryRun {
		return errors.New("use -force flag to empty the whole bucket")
	}

	// only delete single object
	if !strings.HasSuffix(prefix, "/") {
		resp, err := c.client.HeadObject(c.ctx, &s3.HeadObjectInput{
			Bucket: aws.String(cfg.Bucket),
			Key:    aws.String(prefix),
		})
		if err != nil {
			fmt.Fprintf(c.OutWriter, "failed to head object, continuing: %v\n", err)
		} else {
			fmt.Fprintf(c.OutWriter, "deleting %s (%s)\n", prefix, humanize.IBytes(uint64(*resp.ContentLength)))
		}
		return c.objectDelete(cfg.DryRun, cfg.BypassGovernance, cfg.Bucket, prefix, cfg.VersionID)
	}

	// allow deleting the whole bucket
	if prefix == "/" {
		prefix = ""
	}

	// recrusive deletion
	eg, _ := errgroup.WithContext(c.ctx)
	eg.SetLimit(cfg.Concurrency)

	for l, err := range c.objectList(cfg.Bucket, prefix) {
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
				fmt.Fprintf(c.OutWriter, "deleting %s (%s)\n", *l.Object.Key, humanize.IBytes(uint64(*l.Object.Size)))
				err := c.objectDelete(cfg.DryRun, cfg.BypassGovernance, cfg.Bucket, *l.Object.Key, cfg.VersionID)
				if err != nil {
					return err
				}
				return nil
			})
		}
	}

	return eg.Wait()
}

func (c *Controller) objectDelete(dryRun, bypassGovernanceRetention bool, bucket, key, versionID string) error {
	if dryRun {
		return nil
	}

	input := &s3.DeleteObjectInput{
		Bucket:                    aws.String(bucket),
		Key:                       aws.String(key),
		BypassGovernanceRetention: &bypassGovernanceRetention,
	}

	if versionID != "" {
		input.VersionId = &versionID
	}

	_, err := c.client.DeleteObject(c.ctx, input)
	if err != nil {
		return err
	}

	return nil
}
