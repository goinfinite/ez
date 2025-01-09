package infra

import (
	"errors"
	"log/slog"
	"slices"
	"strings"
	"text/template"
	"time"

	"github.com/goinfinite/ez/src/domain/dto"
	"github.com/goinfinite/ez/src/domain/entity"
	"github.com/goinfinite/ez/src/domain/useCase"
	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra/db"
	dbModel "github.com/goinfinite/ez/src/infra/db/model"
	infraHelper "github.com/goinfinite/ez/src/infra/helper"
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
		uint(createDto.PublicPort.Uint16()), createDto.Protocol.String(),
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

	publicPortSslMappingsSliceMap := map[uint16][]entity.Mapping{}
	for _, mapping := range sslMappings {
		publicPortUint := mapping.PublicPort.Uint16()
		if _, exists := publicPortSslMappingsSliceMap[publicPortUint]; !exists {
			publicPortSslMappingsSliceMap[publicPortUint] = []entity.Mapping{}
		}

		publicPortSslMappingsSliceMap[publicPortUint] = append(
			publicPortSslMappingsSliceMap[publicPortUint], mapping,
		)
	}

	fileTemplate := `{{ range $publicPort, $mappings := . }}
map $ssl_preread_server_name $ssl_{{ $publicPort }}_vhost_upstream_map {
{{- range . }}
	{{ .Hostname }} mapping_{{ .Id }}_backend;
	www.{{ .Hostname }} mapping_{{ .Id }}_backend;
{{- end }}
}

server {
	listen {{ $publicPort }};
	proxy_pass $ssl_{{ $publicPort }}_vhost_upstream_map;
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
	err = templatePtr.Execute(&sslFileContent, publicPortSslMappingsSliceMap)
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

func (repo *MappingCmdRepo) CreateTarget(
	createDto dto.CreateMappingTarget,
) (targetId valueObject.MappingTargetId, err error) {
	readContainersRequestDto := dto.ReadContainersRequest{
		Pagination:  useCase.ContainersDefaultPagination,
		ContainerId: []valueObject.ContainerId{createDto.ContainerId},
	}

	readContainersResponseDto, err := repo.containerQueryRepo.Read(readContainersRequestDto)
	if err != nil || len(readContainersResponseDto.Containers) == 0 {
		return targetId, errors.New("ContainerNotFound")
	}
	containerEntity := readContainersResponseDto.Containers[0]

	mappingEntity, err := repo.mappingQueryRepo.ReadById(createDto.MappingId)
	if err != nil {
		return targetId, err
	}

	containerPrivatePort := uint16(0)
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

		containerPrivatePort = portBinding.PrivatePort.Uint16()
	}
	if containerPrivatePort == 0 {
		return targetId, errors.New("ContainerPrivatePortNotFound")
	}

	targetModel := dbModel.NewMappingTarget(
		0,
		createDto.MappingId.Uint64(),
		containerEntity.Id.String(),
		containerEntity.Hostname.String(),
		containerPrivatePort,
	)

	createResult := repo.persistentDbSvc.Handler.Create(&targetModel)
	if createResult.Error != nil {
		return targetId, createResult.Error
	}
	targetId, err = valueObject.NewMappingTargetId(targetModel.ID)
	if err != nil {
		return targetId, errors.New("MappingTargetIdCreationError: " + err.Error())
	}

	err = repo.updateWebServerFile()
	if err != nil {
		deleteTargetDto := dto.NewDeleteMappingTarget(
			createDto.AccountId, createDto.MappingId, targetId,
			createDto.OperatorAccountId, createDto.OperatorIpAddress,
		)
		deleteErr := repo.DeleteTarget(deleteTargetDto)
		if deleteErr != nil {
			return targetId, errors.New("DeleteTargetFailed: " + deleteErr.Error())
		}
		return targetId, errors.New("UpdateWebServerFileFailed: " + err.Error())
	}

	return targetId, nil
}

func (repo *MappingCmdRepo) Delete(deleteDto dto.DeleteMapping) error {
	ormSvc := repo.persistentDbSvc.Handler

	mappingIdUint64 := deleteDto.MappingId.Uint64()

	err := ormSvc.Delete(
		dbModel.MappingTarget{}, "mapping_id = ?", mappingIdUint64,
	).Error
	if err != nil {
		return err
	}

	err = ormSvc.Delete(dbModel.Mapping{}, mappingIdUint64).Error
	if err != nil {
		return err
	}

	return repo.updateWebServerFile()
}

func (repo *MappingCmdRepo) DeleteEmpty() error {
	mappings, err := repo.mappingQueryRepo.Read()
	if err != nil {
		return err
	}

	if len(mappings) == 0 {
		return nil
	}

	nowEpoch := time.Now().Unix()
	for _, mapping := range mappings {
		if len(mapping.Targets) > 0 {
			continue
		}

		isMappingTooRecent := nowEpoch-mapping.CreatedAt.Read() < 60
		if isMappingTooRecent {
			continue
		}

		deleteDto := dto.DeleteMapping{MappingId: mapping.Id}
		err = repo.Delete(deleteDto)
		if err != nil {
			slog.Error(
				"DeleteMappingError",
				slog.String("mappingId", mapping.Id.String()),
				slog.Any("error", err),
			)
			continue
		}
	}

	return nil
}

func (repo *MappingCmdRepo) DeleteTarget(
	deleteDto dto.DeleteMappingTarget,
) error {
	err := repo.persistentDbSvc.Handler.Delete(
		dbModel.MappingTarget{}, "id = ?", deleteDto.TargetId.Uint64(),
	).Error
	if err != nil {
		return err
	}

	return repo.updateWebServerFile()
}

func (repo *MappingCmdRepo) DeleteTargetsByContainerId(
	containerId valueObject.ContainerId,
) error {
	err := repo.persistentDbSvc.Handler.Delete(
		dbModel.MappingTarget{}, "container_id = ?", containerId.String(),
	).Error
	if err != nil {
		return err
	}

	return repo.updateWebServerFile()
}
