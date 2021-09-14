package aws

import (
	"context"

	"github.com/44smkn/ebspv-eraser/pkg/aws/services"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
)

type Cloud interface {
	// S3 provides API to AWS EC2
	EC2() services.EC2
}

func NewCloud(ctx context.Context) (Cloud, error) {
	cfg, err := awscfg.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &defaultCloud{
		ec2: services.NewEC2(cfg),
	}, nil
}

var _ Cloud = &defaultCloud{}

type defaultCloud struct {
	ec2 services.EC2
}

func (c *defaultCloud) EC2() services.EC2 {
	return c.ec2
}
