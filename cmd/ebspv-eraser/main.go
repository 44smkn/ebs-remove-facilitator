package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/44smkn/ebspv-eraser/pkg/aws"
	"github.com/44smkn/ebspv-eraser/pkg/build"
	"github.com/44smkn/ebspv-eraser/pkg/prompt"
	"github.com/44smkn/ebspv-eraser/pkg/volume"
	"gopkg.in/alecthomas/kingpin.v2"
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
	kubernetesCluster = kingpin.Flag("cluster", "kubernetes cluster using EBS Volume of your delete target").Short('c').String()
	printVersion      = kingpin.Flag("version", "prints the build information for the executable").Short('v').Bool()
)

func main() {
	kingpin.Parse()
	if *printVersion {
		fmt.Fprintf(os.Stdout, FormatVersion(build.Version, build.Date))
		os.Exit(ExitCodeOK)
	}
	os.Exit(run())
}

func FormatVersion(version, buildDate string) string {
	version = strings.TrimPrefix(version, "v")
	var dateStr string
	if buildDate != "" {
		dateStr = fmt.Sprintf(" (%s)", buildDate)
	}
	return fmt.Sprintf("ebspv-eraser version %s%s\n", version, dateStr)
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
		fmt.Fprintln(os.Stderr, "You don't choose volumes, so aborted")
		return ExitCodeNoSelectedDeletionVolumeError
	}
	if delete := prompt.DeleteVolumesConfirm(selectedVolumes); !delete {
		fmt.Fprintln(os.Stderr, "aborted")
		return ExitCodeCancelDeleteProcess
	}

	err = ve.DeleteEBSVolumes(ctx, selectedVolumes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete ebs volumes: %s\n", err.Error())
		return ExitCodeListVolumeError
	}

	return ExitCodeOK
}
