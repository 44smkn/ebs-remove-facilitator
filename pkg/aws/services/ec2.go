package services

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

//go:generate mockgen -destination=../../../mocks/aws/services/mock_ec2.go -package=mock_services github.com/44smkn/ebspv-eraser/pkg/aws/services EC2
type EC2 interface {
	// Deletes the specified EBS volume.
	DeleteVolume(ctx context.Context, params *ec2.DeleteVolumeInput, optFns ...func(*ec2.Options)) (*ec2.DeleteVolumeOutput, error)

	// Wrapper to DescribeVolume, which aggregates paged results into list.
	ListVolumesAsList(ctx context.Context, params *ec2.DescribeVolumesInput, optFns ...func(*ec2.Options)) ([]ec2types.Volume, error)
}

func NewEC2(cfg aws.Config, optFns ...func(*ec2.Options)) EC2 {
	return &defaultEC2{
		Client: ec2.NewFromConfig(cfg, optFns...),
	}
}

// default implementation for EC2.
type defaultEC2 struct {
	*ec2.Client
}

func (c *defaultEC2) ListVolumesAsList(ctx context.Context, params *ec2.DescribeVolumesInput, optFns ...func(*ec2.Options)) ([]ec2types.Volume, error) {
	var result []ec2types.Volume
	paginator := ec2.NewDescribeVolumesPaginator(c.Client, params)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		result = append(result, output.Volumes...)
	}
	return result, nil
}
