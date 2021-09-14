package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/44smkn/ebspv-eraser/pkg/aws"
	"github.com/44smkn/ebspv-eraser/pkg/prompt"
	"github.com/44smkn/ebspv-eraser/pkg/volume"
)

const (
	ExitCodeOK = 0

	// Errors start at 10
	ExitCodeError = 10 + iota
	ExitCodeAWSClientInitilizeError
	ExitCodeListVolumeError
	ExitCodeCancelDeleteProcess
	ExitCodeNoAvailableEBSError
	ExitCodeNoSelectedDeletionVolumeError
)

var (
	kubernetesCluster = flag.String("cluster", "", "target kubernetes cluster")
)

func main() {
	flag.Parse()
	os.Exit(run())
}

func run() int {
	ctx := context.Background()
	cloud, err := aws.NewCloud(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize aws client: %s\n", err.Error())
		return ExitCodeAWSClientInitilizeError
	}
	ve := volume.NewVolumeEraser(cloud)

	availbleVolumes, err := ve.ListAvailablePersistentVolumeEBS(ctx, *kubernetesCluster)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch available volume list: %s\n", err.Error())
		return ExitCodeListVolumeError
	}
	if len(availbleVolumes) == 0 {
		fmt.Fprintf(os.Stderr, "Available volumes in %s is not exists\n", *kubernetesCluster)
		return ExitCodeNoAvailableEBSError
	}

	selectedVolumes := prompt.EBSMultiSelectPrompt(availbleVolumes, *kubernetesCluster)
	if len(selectedVolumes) == 0 {
		fmt.Fprintln(os.Stderr, "You don't choose volumes, so process aborted")
		return ExitCodeNoSelectedDeletionVolumeError
	}
	if delete := prompt.DeleteVolumesConfirm(selectedVolumes); !delete {
		fmt.Fprintln(os.Stderr, "Process aborted")
		return ExitCodeCancelDeleteProcess
	}

	err = ve.DeleteEBSVolumes(ctx, selectedVolumes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete ebs volumes: %s\n", err.Error())
		return ExitCodeListVolumeError
	}

	return ExitCodeOK
}
