package componentContainer

import (
	"html/template"
	"strings"

	"github.com/speedianet/control/src/domain/entity"
)

func ContainerTaggedSummary(containerEntity entity.Container) string {
	summaryTemplate := `
	<div class="p-1.5">
		<div><strong>{{.Hostname}}</strong> <small>({{.Id}})</small></div>
		<div class="flex flex-row justify-normal space-x-2">
			<div class="bg-control-300 flex flex-row items-center justify-evenly rounded-md p-1 text-xs">
				<span class="p-1">imageAddress</span>
				<span class="bg-control-100 max-w-48 truncate rounded-md px-2 py-1">{{.ImageAddress}}</span>
			</div>
			<div class="bg-control-300 flex flex-row items-center justify-evenly rounded-md p-1 text-xs">
				<span class="p-1">accountId</span>
				<span class="bg-control-100 max-w-48 truncate rounded-md px-2 py-1">{{.AccountId}}</span>
			</div>
			<div class="bg-control-300 flex flex-row items-center justify-evenly rounded-md p-1 text-xs">
				<span class="p-1">profileId</span>
				<span class="bg-control-100 max-w-48 truncate rounded-md px-2 py-1">{{.ProfileId}}</span>
			</div>
		</div>
	</div>
	`

	templatePtr, err := template.New("containerTaggedSummary").Parse(summaryTemplate)
	if err != nil {
		return ""
	}

	var templateStrBuilder strings.Builder
	err = templatePtr.Execute(&templateStrBuilder, containerEntity)
	if err != nil {
		return ""
	}

	return templateStrBuilder.String()
}
