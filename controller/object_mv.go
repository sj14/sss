package controller

import "fmt"

// TODO: move whole prefix
func (c *Controller) ObjectMove(cfg ObjectCopyConfig) error {
	err := c.ObjectCopy(cfg)
	if err != nil {
		return fmt.Errorf("copy object: %w", err)
	}

	err = c.objectDelete(false, false, cfg.SrcBucket, cfg.SrcKey, "")
	if err != nil {
		return fmt.Errorf("delete object: %w", err)
	}

	return nil
}
