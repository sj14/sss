package controller

import (
	"fmt"
)

type BucketCleanupConfig struct {
	Bucket           string
	Concurrency      int
	Force            bool
	DryRun           bool
	Multiparts       bool
	ObjectVersion    bool
	BypassGovernance bool
}

func (c *Controller) BucketCleanup(cfg BucketCleanupConfig) error {
	if !cfg.Force && !cfg.DryRun {
		return fmt.Errorf("--force flag required")
	}

	if !cfg.ObjectVersion && !cfg.Multiparts {
		return fmt.Errorf("at least one of --all-object-versions or --all-multiparts needs to be set")
	}

	if cfg.ObjectVersion {
		fmt.Fprintln(c.OutWriter, "> deleting all objects <")

		for versions, err := range c.objectVersions(cfg.Bucket, "", "") {
			if err != nil {
				return err
			}

			for _, version := range versions.Versions {
				err := c.ObjectDelete("/", ObjectDeleteConfig{
					Bucket: cfg.Bucket,
					// Delimiter:
					Force:            cfg.Force,
					Concurrency:      cfg.Concurrency,
					DryRun:           cfg.DryRun,
					BypassGovernance: cfg.BypassGovernance,
					VersionID:        *version.VersionId,
				})
				if err != nil {
					return err
				}
			}
		}
	}

	if cfg.Multiparts {
		fmt.Fprintln(c.OutWriter, "> deleting all multipart uploads <")
		err := c.BucketMultipartUploadAbortAll(
			cfg.Bucket,
			cfg.DryRun,
			cfg.Concurrency,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
