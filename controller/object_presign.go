package controller

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *Controller) ObjectPresignGet(objectKey string, cfg ObjectGetConfig, expiration time.Duration) error {
	presigner := s3.NewPresignClient(c.client)

	req, err := presigner.PresignGetObject(c.ctx, &s3.GetObjectInput{
		Bucket: &cfg.Bucket,
		Key:    &objectKey,
		// IfMatch: ,
		// SSECustomerAlgorithm: ,
	}, s3.WithPresignExpires(expiration),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", req.URL)
	return nil
}

func (c *Controller) ObjectPresignPut(expiration time.Duration, key string, cfg ObjectPutConfig) error {
	presigner := s3.NewPresignClient(c.client)

	req, err := presigner.PresignPutObject(c.ctx, &s3.PutObjectInput{
		Bucket: &cfg.Bucket,
		Key:    &key,
		// IfMatch: ,
		// SSECustomerAlgorithm: ,
	}, s3.WithPresignExpires(expiration),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", req.URL)
	return nil
}
