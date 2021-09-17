# ebs-remove-facilitator

`ebspv-eraser` is the tool that facilinate deletion EBSVolume of Persistent Volume resource remaining even after the EKS Cluster is deleted.

## Demo

This demo retrive EBSVolume list and delete selected volumes.

![schreenshot of demo](https://raw.githubusercontent.com/44smkn/ebspv-eraser/main/.github/images/ebspv-eraser_readme.gif)

## Installation

Download packaged binaries from the [releases page](https://github.com/44smkn/ebspv-eraser/releases).

You will need to have AWS API credentials configured. What works for AWS CLI, should be sufficient. You can use [~/.aws/credentials file](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html) or [environment variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html#envvars-set).

## Usage

```console
$ ebspv-eraser --cluster <CLUSTER_NAME>
```

After executing the command, an interactive screen will appear where you can select the target to delete.
