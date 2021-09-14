package prompt

import (
	"fmt"
	"strings"

	"github.com/44smkn/ebspv-eraser/pkg/volume"
	"github.com/AlecAivazis/survey/v2"
)

func EBSMultiSelectPrompt(options []volume.EBSVolume, cluster string) []volume.EBSVolume {

	selectedVolumeTexts := []string{}
	multiselect := &survey.MultiSelect{
		Message: fmt.Sprintf("Select the volumes you want to delete. All volumes is available and managed by %s.", cluster),
		Options: createDescriptionTexts(options),
	}
	survey.AskOne(multiselect, &selectedVolumeTexts)

	selectedVolumes := make([]volume.EBSVolume, 0, len(selectedVolumeTexts))
	for _, t := range selectedVolumeTexts {
		volumeID := strings.Fields(t)[0]
		if obj := lookUpByID(options, volumeID); obj != nil {
			selectedVolumes = append(selectedVolumes, *obj)
		}
	}
	return selectedVolumes
}

func createDescriptionTexts(volumes []volume.EBSVolume) []string {
	descriptions := make([]string, 0, len(volumes))
	for _, v := range volumes {
		description := fmt.Sprintf("%s - %s/%s", v.ID, v.KubernetesNamespace, v.KubernetesPVCName)
		descriptions = append(descriptions, description)
	}
	return descriptions
}

func lookUpByID(volumes []volume.EBSVolume, volumeID string) *volume.EBSVolume {
	for _, v := range volumes {
		if v.ID == volumeID {
			return &v
		}
	}
	return nil
}

func DeleteVolumesConfirm(targets []volume.EBSVolume) bool {
	targetTexts := createDescriptionTexts(targets)
	fmt.Printf("\nVolumes you have specified for deletion: \n")
	for _, v := range targetTexts {
		fmt.Printf("* %s\n", v)
	}
	fmt.Printf("\n")

	delete := false
	confirm := &survey.Confirm{
		Message: "Delete these volumes. Are you sure?",
	}

	survey.AskOne(confirm, &delete)
	return delete
}
