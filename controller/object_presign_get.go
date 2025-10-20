package controller

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ObjectPresignConfig struct {
	Method          string
	ObjectGetConfig ObjectGetConfig
}

func (c *Controller) ObjectPresign(expiration time.Duration, cfg ObjectPresignConfig) error {
	switch strings.ToLower(cfg.Method) {
	case "get":
		c.objectPresignGet(expiration, cfg.ObjectGetConfig)
	// case "put":
	// c.objectPresignPut(expiration)
	default:
		return fmt.Errorf("presign method not recognized: %q", cfg.Method)
	}

	return nil
}

func (c *Controller) objectPresignGet(expiration time.Duration, cfg ObjectGetConfig) error {
	presigner := s3.NewPresignClient(c.client)

	req, err := presigner.PresignGetObject(c.ctx, &s3.GetObjectInput{
		Bucket: &cfg.Bucket,
		Key:    &cfg.ObjectKey,
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

// func (c *Controller) objectPresignPut(expiration time.Duration, cfg ObjectPutConfig) error {
// 	presigner := s3.NewPresignClient(c.client)

// 	req, err := presigner.PresignPutObject(c.ctx, &s3.PutObjectInput{
// 		Bucket: &cfg.Bucket,
// 		Key:    &cfg.ObjectKey,
// 		// IfMatch: ,
// 		// SSECustomerAlgorithm: ,
// 	}, s3.WithPresignExpires(expiration),
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Printf("%s\n", req.URL)
// 	return nil
// }
