package infra

import (
	"errors"
	"strings"
	"text/template"

	"github.com/goinfinite/ez/src/domain/valueObject"
	"github.com/goinfinite/ez/src/infra/db"
	dbModel "github.com/goinfinite/ez/src/infra/db/model"
	infraHelper "github.com/goinfinite/ez/src/infra/helper"
)

const ProxyWebServerFile string = "/var/nginx/http.d/container-proxy.conf"

type ContainerProxyCmdRepo struct {
	persistentDbSvc    *db.PersistentDatabaseService
	containerQueryRepo *ContainerQueryRepo
}

func NewContainerProxyCmdRepo(
	persistentDbSvc *db.PersistentDatabaseService,
) *ContainerProxyCmdRepo {
	return &ContainerProxyCmdRepo{
		persistentDbSvc:    persistentDbSvc,
		containerQueryRepo: NewContainerQueryRepo(persistentDbSvc),
	}
}

func (repo *ContainerProxyCmdRepo) updateWebServerFile() error {
	proxyModels := []dbModel.ContainerProxy{}
	err := repo.persistentDbSvc.Handler.Find(&proxyModels).Error
	if err != nil {
		return err
	}

	serverHostname, err := infraHelper.ReadServerHostname()
	if err != nil {
		return errors.New("InvalidServerHostname: " + err.Error())
	}

	//cspell:disable
	fileTemplate := `server {
	listen 1618 ssl;
	server_name ` + serverHostname.String() + `;

	ssl_certificate /var/infinite/pki/ez.crt;
	ssl_certificate_key /var/infinite/pki/ez.key;
	{{ range . }}
	location /{{ .ContainerId }}/ {
		sub_filter_once off;
		sub_filter_types application/javascript;
		sub_filter '"/_/"' '"/{{ .ContainerId }}/_/"';
		sub_filter 'src="/_/' 'src="/{{ .ContainerId }}/_/';
		sub_filter 'href="/_/' 'href="/{{ .ContainerId }}/_/';
		sub_filter 'src=/_/' 'src=/{{ .ContainerId }}/_/';
		sub_filter 'href=/_/' 'href=/{{ .ContainerId }}/_/';

		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header X-Forwarded-Proto $scheme;
		proxy_set_header X-Proxied-By "Infinite Ez";
		proxy_set_header X-Container-Id "{{ .ContainerId }}";
		proxy_pass https://localhost:{{ .ContainerPrivatePort }}/;
	}
	{{- end }}
}
{{ range . }}
server {
	listen 1618 ssl;
	server_name {{ .ContainerHostname}} os.{{ .ContainerHostname}};

	ssl_certificate /var/infinite/pki/ez.crt;
	ssl_certificate_key /var/infinite/pki/ez.key;

	location / {
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header X-Forwarded-Proto $scheme;
		proxy_set_header X-Proxied-By "Infinite Ez";
		proxy_set_header X-Container-Id "{{ .ContainerId }}";
		proxy_pass https://localhost:{{ .ContainerPrivatePort }};
	}
}
{{ end }}`
	//cspell:enable

	if len(proxyModels) == 0 {
		fileTemplate = ``
	}

	templatePtr, err := template.New("webServerFile").Parse(fileTemplate)
	if err != nil {
		return errors.New("TemplateParsingError: " + err.Error())
	}

	var webServerFileContent strings.Builder
	err = templatePtr.Execute(&webServerFileContent, proxyModels)
	if err != nil {
		return errors.New("TemplateExecutionError: " + err.Error())
	}

	err = infraHelper.UpdateFile(ProxyWebServerFile, webServerFileContent.String(), true)
	if err != nil {
		return err
	}

	_, err = infraHelper.RunCmdWithSubShell("nginx -t && nginx -s reload")
	if err != nil {
		return errors.New("ReloadNginxFailed: " + err.Error())
	}

	return nil
}

func (repo *ContainerProxyCmdRepo) Create(containerId valueObject.ContainerId) error {
	containerEntity, err := repo.containerQueryRepo.ReadById(containerId)
	if err != nil {
		return err
	}

	containerPrivatePort := uint16(0)
	for _, portBinding := range containerEntity.PortBindings {
		if portBinding.ContainerPort.String() != "1618" {
			continue
		}

		containerPrivatePort = portBinding.PrivatePort.Uint16()
	}
	if containerPrivatePort == 0 {
		return errors.New("InfiniteOsPrivatePortNotFound")
	}

	proxyModel := dbModel.NewContainerProxy(
		0,
		containerEntity.Id.String(),
		containerEntity.Hostname.String(),
		uint(containerPrivatePort),
	)

	createResult := repo.persistentDbSvc.Handler.Create(&proxyModel)
	if createResult.Error != nil {
		return createResult.Error
	}

	return repo.updateWebServerFile()
}

func (repo *ContainerProxyCmdRepo) Delete(containerId valueObject.ContainerId) error {
	err := repo.persistentDbSvc.Handler.Delete(
		dbModel.ContainerProxy{},
		"container_id = ?", containerId.String(),
	).Error
	if err != nil {
		return err
	}

	return repo.updateWebServerFile()
}
