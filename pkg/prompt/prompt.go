package prompt

import (
	"fmt"
	"strings"

	"github.com/44smkn/ebspv-eraser/pkg/volume"
	"github.com/AlecAivazis/survey/v2"
)

var customMultiSelectQuestionTemplate = `
{{- define "option"}}
    {{- if eq .SelectedIndex .CurrentIndex }}{{color .Config.Icons.SelectFocus.Format }}{{ .Config.Icons.SelectFocus.Text }}{{color "reset"}}{{else}} {{end}}
    {{- if index .Checked .CurrentOpt.Index }}{{color .Config.Icons.MarkedOption.Format }} {{ .Config.Icons.MarkedOption.Text }} {{else}}{{color .Config.Icons.UnmarkedOption.Format }} {{ .Config.Icons.UnmarkedOption.Text }} {{end}}
    {{- color "reset"}}
    {{- " "}}{{- .CurrentOpt.Value}}
{{end}}
{{- if .ShowHelp }}{{- color .Config.Icons.Help.Format }}{{ .Config.Icons.Help.Text }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color .Config.Icons.Question.Format }}{{ .Config.Icons.Question.Text }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }}{{ .FilterMessage }}{{color "reset"}}
{{- if .ShowAnswer}}{{"\n"}}
{{- else }}
	{{- "  "}}{{"\n"}}{{- color "cyan"}}[Use arrows to move, space to select, <right> to all, <left> to none, type to filter{{- if and .Help (not .ShowHelp)}}, {{ .Config.HelpInput }} for more help{{end}}]{{color "reset"}}{{"\n"}}{{"\n"}}{{- color "default+hb"}}       volumeID                namespace/pvc{{color "reset"}}
  {{- "\n"}}
  {{- range $ix, $option := .PageEntries}}
    {{- template "option" $.IterateOption $ix $option}}
  {{- end}}
{{- end}}`

func init() {
	survey.MultiSelectQuestionTemplate = customMultiSelectQuestionTemplate
}

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
