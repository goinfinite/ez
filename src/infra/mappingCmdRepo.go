package infra

import (
	"errors"
	"os"
	"strings"

	"github.com/speedianet/control/src/domain/dto"
	"github.com/speedianet/control/src/domain/entity"
	"github.com/speedianet/control/src/domain/valueObject"
	"github.com/speedianet/control/src/infra/db"
	dbModel "github.com/speedianet/control/src/infra/db/model"
	infraHelper "github.com/speedianet/control/src/infra/helper"
)

var nginxStreamConfDir = "/var/nginx/stream.d"
var nginxHttpConfDir = "/var/nginx/http.d"

type MappingCmdRepo struct {
	dbSvc     *db.DatabaseService
	queryRepo *MappingQueryRepo
}

func NewMappingCmdRepo(dbSvc *db.DatabaseService) *MappingCmdRepo {
	return &MappingCmdRepo{
		dbSvc:     dbSvc,
		queryRepo: NewMappingQueryRepo(dbSvc),
	}
}

func (repo MappingCmdRepo) Add(mappingDto dto.AddMapping) (valueObject.MappingId, error) {
	var mappingId valueObject.MappingId

	mappingModel := dbModel.Mapping{}.AddDtoToModel(mappingDto)

	createResult := repo.dbSvc.Orm.Create(&mappingModel)
	if createResult.Error != nil {
		return mappingId, createResult.Error
	}

	return valueObject.NewMappingId(mappingModel.ID)
}

func (repo MappingCmdRepo) getHttpsPreReadBlock() (string, error) {
	httpsProtocol, _ := valueObject.NewNetworkProtocol("https")

	allHttpsMappings, err := repo.queryRepo.FindAll(
		nil,
		nil,
		&httpsProtocol,
	)
	if err != nil {
		return "", err
	}

	if len(allHttpsMappings) == 0 {
		return "", nil
	}

	type hostUpstream struct {
		hostname string
		upstream string
	}

	portHostUpstreamMap := map[string][]hostUpstream{}
	for _, mapping := range allHttpsMappings {
		hostPort := mapping.PublicPort.String()
		hostname := "default"
		if mapping.Hostname != nil {
			hostname = mapping.Hostname.String()
		}

		upstreamName := "mapping_" + mapping.Id.String() + "_backend"
		portHostUpstreamMap[hostPort] = append(
			portHostUpstreamMap[hostPort],
			hostUpstream{
				hostname: hostname,
				upstream: upstreamName,
			},
		)
	}

	preReadBlock := ""
	for hostPort, hostUpstreams := range portHostUpstreamMap {
		if len(hostUpstreams) == 0 {
			continue
		}

		varName := "https_" + hostPort + "_container_name"
		preReadBlock += "map $ssl_preread_server_name $" + varName + " {\n"

		for _, hostUpstream := range hostUpstreams {
			preReadBlock += "\t" + hostUpstream.hostname + " " + hostUpstream.upstream + ";\n"
		}

		preReadBlock += "}\n"
	}

	return preReadBlock, nil
}

func (repo MappingCmdRepo) nginxConfigFactory(
	mappingEntity entity.Mapping,
) (string, error) {
	containerIdTargetEntityMap := map[valueObject.ContainerId]entity.MappingTarget{}
	for _, mappingTarget := range mappingEntity.Targets {
		containerIdTargetEntityMap[mappingTarget.ContainerId] = mappingTarget
	}

	containerQueryRepo := NewContainerQueryRepo(repo.dbSvc)
	containerEntities, err := containerQueryRepo.Get()
	if err != nil {
		return "", err
	}

	if len(containerEntities) == 0 {
		return "", errors.New("NoContainersFound")
	}

	containerIdContainerEntityMap := map[valueObject.ContainerId]entity.Container{}
	for _, containerEntity := range containerEntities {
		containerIdContainerEntityMap[containerEntity.Id] = containerEntity
	}

	serversList := ""
	for containerId, containerEntity := range containerIdContainerEntityMap {
		_, isContainerTarget := containerIdTargetEntityMap[containerId]
		if !isContainerTarget {
			continue
		}

		for _, containerPortBinding := range containerEntity.PortBindings {
			privatePort := containerPortBinding.PrivatePort
			if privatePort == nil {
				continue
			}

			if containerPortBinding.PublicPort != mappingEntity.PublicPort {
				continue
			}

			serversList += "server localhost:" + privatePort.String() + ";\n"
		}
	}
	serversList = strings.TrimSpace(serversList)

	if serversList == "" {
		return "", errors.New("UpstreamServersListEmpty")
	}

	upstreamName := "mapping_" + mappingEntity.Id.String() + "_backend"
	upstreamBlock := `
upstream ` + upstreamName + ` {
	` + serversList + `
}
`

	hostPort := mappingEntity.PublicPort.String()

	serverNameLine := ``
	if mappingEntity.Hostname != nil {
		serverNameLine = "server_name " + mappingEntity.Hostname.String() + ";\n"
	}

	httpNginxConf := `
server {
	listen      ` + hostPort + `;
	` + serverNameLine + `

	location / {
		proxy_pass http://` + upstreamName + `;
	}
}
`

	httpsPreReadBlock := ""
	if mappingEntity.Protocol.String() == "https" {
		httpsPreReadBlock, err = repo.getHttpsPreReadBlock()
		if err != nil {
			return "", err
		}

		err = infraHelper.UpdateFile(
			nginxStreamConfDir+"/https_pre_read.conf",
			httpsPreReadBlock,
			true,
		)
		if err != nil {
			return "", errors.New("UpdateNginxPreReadBlockFailed: " + err.Error())
		}
	}

	httpsConf := `
server {
	listen      ` + hostPort + `;
	proxy_pass  $https_` + hostPort + `_container_name;
	ssl_preread on;
}
`

	tcpNginxConf := `
server {
	listen      ` + hostPort + `;
	proxy_pass ` + upstreamName + `;
}
`

	udpNginxConf := `
server {
	listen      ` + hostPort + ` udp;
	proxy_pass ` + upstreamName + `;
}
`

	nginxConf := ""
	switch mappingEntity.Protocol.String() {
	case "http", "grpc", "ws":
		nginxConf = httpNginxConf
	case "https", "grpcs", "wss":
		nginxConf = httpsConf
	case "udp":
		nginxConf = udpNginxConf
	default:
		nginxConf = tcpNginxConf
	}

	return strings.TrimSpace(upstreamBlock+nginxConf) + "\n", nil
}

func (repo MappingCmdRepo) getNginxConfDirByProtocol(
	protocol valueObject.NetworkProtocol,
) string {
	switch protocol.String() {
	case "http", "grpc", "ws":
		return nginxHttpConfDir
	}

	return nginxStreamConfDir
}

func (repo MappingCmdRepo) updateMappingFile(mappingId valueObject.MappingId) error {
	mappingEntity, err := repo.queryRepo.GetById(mappingId)
	if err != nil {
		return err
	}

	if len(mappingEntity.Targets) == 0 {
		return errors.New("MappingHasNoTarget")
	}

	nginxConf, err := repo.nginxConfigFactory(mappingEntity)
	if err != nil {
		return err
	}

	nginxConfDir := repo.getNginxConfDirByProtocol(mappingEntity.Protocol)
	err = infraHelper.UpdateFile(
		nginxConfDir+"/mapping_"+mappingId.String()+".conf",
		nginxConf,
		true,
	)
	if err != nil {
		return errors.New("UpdateNginxConfFailed: " + err.Error())
	}

	_, err = infraHelper.RunCmd("nginx", "-s", "reload")
	if err != nil {
		return errors.New("ReloadNginxFailed: " + err.Error())
	}

	return nil
}

func (repo MappingCmdRepo) AddTarget(addDto dto.AddMappingTarget) error {
	targetModel := dbModel.MappingTarget{}.AddDtoToModel(addDto)

	createResult := repo.dbSvc.Orm.Create(&targetModel)
	if createResult.Error != nil {
		return createResult.Error
	}

	return repo.updateMappingFile(addDto.MappingId)
}

func (repo MappingCmdRepo) deleteMappingFile(mappingId valueObject.MappingId) error {
	mappingEntity, err := repo.queryRepo.GetById(mappingId)
	if err != nil {
		return err
	}

	nginxConfDir := repo.getNginxConfDirByProtocol(mappingEntity.Protocol)
	err = os.Remove(nginxConfDir + "/mapping_" + mappingId.String() + ".conf")
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	_, err = infraHelper.RunCmd("nginx", "-s", "reload")
	if err != nil {
		return errors.New("ReloadNginxFailed: " + err.Error())
	}

	return nil
}

func (repo MappingCmdRepo) Delete(id valueObject.MappingId) error {
	ormSvc := repo.dbSvc.Orm

	err := ormSvc.Delete(dbModel.MappingTarget{}, "mapping_id = ?", id.Get()).Error
	if err != nil {
		return err
	}

	err = repo.deleteMappingFile(id)
	if err != nil {
		return err
	}

	return ormSvc.Delete(dbModel.Mapping{}, id.Get()).Error
}

func (repo MappingCmdRepo) DeleteTarget(id valueObject.MappingTargetId) error {
	targetEntity, err := repo.queryRepo.GetTargetById(id)
	if err != nil {
		return err
	}

	err = repo.dbSvc.Orm.Delete(dbModel.MappingTarget{}, id.Get()).Error
	if err != nil {
		return err
	}

	mappingEntity, err := repo.queryRepo.GetById(targetEntity.MappingId)
	if err != nil {
		return err
	}

	if len(mappingEntity.Targets) < 1 {
		return repo.deleteMappingFile(targetEntity.MappingId)
	}

	return repo.updateMappingFile(targetEntity.MappingId)
}
