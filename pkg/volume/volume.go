package volume

import (
	"context"
	"fmt"
	"os"

	"github.com/44smkn/ebspv-eraser/pkg/aws"
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

const (
	kubernetesNamespaceTagKey = "kubernetes.io/created-for/pvc/namespace"
	KubernetesPVCNameTagKey   = "kubernetes.io/created-for/pvc/name"
)

type VolumeEraser interface {
	ListAvailablePersistentVolumeEBS(ctx context.Context, cluster string) ([]EBSVolume, error)
	DeleteEBSVolumes(ctx context.Context, volumes []EBSVolume) error
}

type defaultVolumeEraser struct {
	cloud aws.Cloud
}

type EBSVolume struct {
	ID                  string
	State               string
	KubernetesNamespace string
	KubernetesPVCName   string
}

func NewVolumeEraser(cloud aws.Cloud) VolumeEraser {
	return &defaultVolumeEraser{
		cloud: cloud,
	}
}

// You can ask the compiler to check that the type T implements the interface I by attempting an assignment
// see: https://golang.org/doc/faq#guarantee_satisfies_interface
var _ VolumeEraser = &defaultVolumeEraser{}

func (e *defaultVolumeEraser) ListAvailablePersistentVolumeEBS(ctx context.Context, cluster string) ([]EBSVolume, error) {
	clusterTagKey := fmt.Sprintf("tag:kubernetes.io/cluster/%s", cluster)
	params := &ec2.DescribeVolumesInput{
		Filters: []ec2types.Filter{
			ec2Filter("status", "available"),
			ec2Filter(clusterTagKey, "owned"),
		},
	}

	list, err := e.cloud.EC2().ListVolumesAsList(ctx, params)
	if err != nil {
		return nil, err
	}

	volumes := make([]EBSVolume, 0, len(list))
	for _, v := range list {
		e := EBSVolume{
			ID:                  *v.VolumeId,
			State:               string(v.State),
			KubernetesNamespace: lookUpTag(v.Tags, kubernetesNamespaceTagKey),
			KubernetesPVCName:   lookUpTag(v.Tags, KubernetesPVCNameTagKey),
		}
		volumes = append(volumes, e)
	}
	return volumes, nil
}

func lookUpTag(tags []ec2types.Tag, key string) string {
	for _, tag := range tags {
		if *tag.Key == key {
			return *tag.Value
		}
	}
	return ""
}

func ec2Filter(name, value string) ec2types.Filter {
	return ec2types.Filter{
		Name:   awssdk.String(name),
		Values: []string{value},
	}
}

func (e *defaultVolumeEraser) DeleteEBSVolumes(ctx context.Context, volumes []EBSVolume) error {
	for _, v := range volumes {
		params := &ec2.DeleteVolumeInput{
			VolumeId: awssdk.String(v.ID),
		}
		_, err := e.cloud.EC2().DeleteVolume(ctx, params)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "%s has been deleted\n", v.ID)
	}
	return nil
}
