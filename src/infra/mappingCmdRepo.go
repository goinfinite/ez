package infra

import (
	"errors"
	"log/slog"
	"slices"
	"strings"
	"text/template"
	"time"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
	infraHelper "github.com/speedianet/control/src/infra/helper"
)

var nginxStreamConfDir string = "/var/nginx/stream.d"
var nginxHttpConfDir string = "/var/nginx/http.d"

type MappingCmdRepo struct {
	persistentDbSvc    *db.PersistentDatabaseService
	mappingQueryRepo   *MappingQueryRepo
	containerQueryRepo *ContainerQueryRepo
}

func NewMappingCmdRepo(persistentDbSvc *db.PersistentDatabaseService) *MappingCmdRepo {
	return &MappingCmdRepo{
		persistentDbSvc:    persistentDbSvc,
		mappingQueryRepo:   NewMappingQueryRepo(persistentDbSvc),
		containerQueryRepo: NewContainerQueryRepo(persistentDbSvc),
	}
}

func (repo *MappingCmdRepo) Create(
	createDto dto.CreateMapping,
) (mappingId valueObject.MappingId, err error) {
	var hostnamePtr *string
	if createDto.Hostname != nil {
		hostnameStr := createDto.Hostname.String()
		hostnamePtr = &hostnameStr
	}

	mappingModel := dbModel.NewMapping(
		0, createDto.AccountId.Uint64(), hostnamePtr,
		uint(createDto.PublicPort.Read()), createDto.Protocol.String(),
		[]dbModel.MappingTarget{}, time.Now(), time.Now(),
	)

	createResult := repo.persistentDbSvc.Handler.Create(&mappingModel)
	if createResult.Error != nil {
		return mappingId, createResult.Error
	}

	// FYI: adding target requires business logic that is laid out in the use case.
	// Although the target information is on the createDto, it is not used here.

	return valueObject.NewMappingId(mappingModel.ID)
}

func (repo *MappingCmdRepo) updateWebServerFile() error {
	mappings, err := repo.mappingQueryRepo.Read()
	if err != nil {
		return err
	}

	if len(mappings) == 0 {
		return nil
	}

	for _, mappingEntity := range mappings {
		mappingFileDir := nginxStreamConfDir
		switch mappingEntity.Protocol.String() {
		case "http", "grpc", "ws":
			mappingFileDir = nginxHttpConfDir
		}

		mappingIdStr := mappingEntity.Id.String()
		mappingFile := mappingFileDir + "/mapping-" + mappingIdStr + ".conf"

		if len(mappingEntity.Targets) == 0 {
			err = infraHelper.RemoveFile(mappingFile)
			if err != nil {
				slog.Error(
					"RemoveMappingFileError",
					slog.String("mappingId", mappingIdStr),
					slog.Any("error", err),
				)
			}
			continue
		}

		fileTemplate := `upstream mapping_{{ .Id }}_backend {
{{- range .Targets }}
	server localhost:{{ .ContainerPrivatePort }};
{{- end }}
}

{{- if eq .Protocol "http" "grpc" "ws" "tcp" "udp" }}
server {
	listen {{ .PublicPort }}{{ if eq .Protocol "udp" }}udp{{end}};
	{{- if eq .Protocol "http" "grpc" "ws" }}
	server_name {{ .Hostname }} www.{{ .Hostname }};

	location / {
		{{- if eq .Protocol "http" "ws" }}
		proxy_pass http://mapping_{{ .Id }}_backend;
		proxy_set_header Host $host;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header X-Forwarded-Proto $scheme;
		{{- end }}
		{{- if eq .Protocol "ws" }}
		proxy_http_version 1.1;
		proxy_set_header Upgrade $http_upgrade;
		proxy_set_header Connection "upgrade";
		{{- end }}
		{{- if eq .Protocol "grpc" }}
		grpc_pass grpc://mapping_{{ .Id }}_backend;
		{{- end }}
	}
	{{- else }}
	proxy_pass mapping_{{ .Id }}_backend;
	{{- end }}
}
{{- end }}
`
		templatePtr, err := template.New("mappingFile").Parse(fileTemplate)
		if err != nil {
			slog.Error(
				"MappingTemplateParsingError",
				slog.String("mappingId", mappingIdStr), slog.Any("error", err),
			)
			continue
		}

		var mappingFileContent strings.Builder
		err = templatePtr.Execute(&mappingFileContent, mappingEntity)
		if err != nil {
			slog.Error(
				"MappingTemplateExecutionError",
				slog.String("mappingId", mappingIdStr), slog.Any("error", err),
			)
			continue
		}

		mappingContentStr := strings.TrimSpace(mappingFileContent.String()) + "\n"

		err = infraHelper.UpdateFile(mappingFile, mappingContentStr, true)
		if err != nil {
			slog.Error(
				"UpdateMappingFileError",
				slog.String("mappingId", mappingIdStr), slog.Any("error", err),
			)
			continue
		}
	}

	sslMappings := []entity.Mapping{}
	validSslProtocolos := []string{"https", "grpcs", "wss"}
	for _, mapping := range mappings {
		if !slices.Contains(validSslProtocolos, mapping.Protocol.String()) {
			continue
		}

		sslMappings = append(sslMappings, mapping)
	}

	fileTemplate := `{{ range . }}
map $ssl_preread_server_name $ssl_{{.PublicPort}}_vhost_upstream_map {
	{{ .Hostname }} mapping_{{ .Id }}_backend;
	www.{{ .Hostname }} mapping_{{ .Id }}_backend;
}

server {
	listen {{ .PublicPort }};
	proxy_pass $ssl_{{ .PublicPort }}_vhost_upstream_map;
	ssl_preread on;
}
{{ end }}`

	if len(sslMappings) == 0 {
		fileTemplate = ``
	}

	templatePtr, err := template.New("sslFile").Parse(fileTemplate)
	if err != nil {
		return errors.New("SslTemplateParsingError: " + err.Error())
	}

	var sslFileContent strings.Builder
	err = templatePtr.Execute(&sslFileContent, sslMappings)
	if err != nil {
		return errors.New("SslTemplateExecutionError: " + err.Error())
	}

	fileContentStr := strings.TrimSpace(sslFileContent.String()) + "\n"

	err = infraHelper.UpdateFile(nginxStreamConfDir+"/ssl.conf", fileContentStr, true)
	if err != nil {
		return errors.New("UpdateSslFileError: " + err.Error())
	}

	_, err = infraHelper.RunCmdWithSubShell("nginx -t && nginx -s reload")
	if err != nil {
		return errors.New("ReloadNginxFailed: " + err.Error())
	}

	return nil
}

func (repo *MappingCmdRepo) CreateTarget(createDto dto.CreateMappingTarget) error {
	containerEntity, err := repo.containerQueryRepo.ReadById(createDto.ContainerId)
	if err != nil {
		return err
	}

	mappingEntity, err := repo.mappingQueryRepo.ReadById(createDto.MappingId)
	if err != nil {
		return err
	}

	containerPrivatePort := uint64(0)
	mappingPublicPortStr := mappingEntity.PublicPort.String()
	for _, portBinding := range containerEntity.PortBindings {
		if portBinding.PublicPort.String() != mappingPublicPortStr {
			continue
		}

		bindingProtocolStr := portBinding.Protocol.String()
		mappingProtocolStr := mappingEntity.Protocol.String()
		if bindingProtocolStr != mappingProtocolStr {
			slog.Error(
				"MappingVsBindingProtocolMismatch",
				slog.String("mappingPublicPort", mappingPublicPortStr),
				slog.String("mappingProtocol", mappingProtocolStr),
				slog.String("bindingProtocol", bindingProtocolStr),
			)
			continue
		}

		containerPrivatePort = portBinding.PrivatePort.Read()
	}
	if containerPrivatePort == 0 {
		return errors.New("ContainerPrivatePortNotFound")
	}

	targetModel := dbModel.NewMappingTarget(
		0,
		createDto.MappingId.Uint64(),
		containerEntity.Id.String(),
		containerEntity.Hostname.String(),
		uint(containerPrivatePort),
	)

	createResult := repo.persistentDbSvc.Handler.Create(&targetModel)
	if createResult.Error != nil {
		return createResult.Error
	}

	return repo.updateWebServerFile()
}

func (repo *MappingCmdRepo) Delete(id valueObject.MappingId) error {
	ormSvc := repo.persistentDbSvc.Handler

	err := ormSvc.Delete(dbModel.MappingTarget{}, "mapping_id = ?", id.Uint64()).Error
	if err != nil {
		return err
	}

	err = ormSvc.Delete(dbModel.Mapping{}, id.Uint64()).Error
	if err != nil {
		return err
	}

	return repo.updateWebServerFile()
}

func (repo *MappingCmdRepo) DeleteTarget(
	mappingId valueObject.MappingId,
	targetId valueObject.MappingTargetId,
) error {
	err := repo.persistentDbSvc.Handler.Delete(
		dbModel.MappingTarget{}, "id = ?", targetId.Uint64(),
	).Error
	if err != nil {
		return err
	}

	return repo.updateWebServerFile()
}
