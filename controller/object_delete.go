package controller

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/dustin/go-humanize"
)

func (c *Controller) ObjectDelete(bucket, key, delimiter string, force bool) error {
	if key == "" && !force {
		return errors.New("use -force flag to empty the whole bucket")
	}

	objects, prefixes, err := c.objectList(bucket, key, delimiter)
	if err != nil {
		return err
	}

	for _, prefix := range prefixes {
		err := c.ObjectDelete(bucket, *prefix.Prefix, delimiter, force)
		if err != nil {
			return err
		}
	}

	for _, object := range objects {
		fmt.Printf("deleting %s (%s)\n", *object.Key, humanize.Bytes(uint64(*object.Size)))
		err := c.objectDelete(bucket, *object.Key)
		if err != nil {
			return err
		}
	}

	return nil
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
