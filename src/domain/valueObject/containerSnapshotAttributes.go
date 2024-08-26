package valueObject

type ContainerSnapshotAttributes struct {
	PortBindings  []PortBinding           `json:"portBindings"`
	Envs          []ContainerEnv          `json:"envs"`
	RestartPolicy *ContainerRestartPolicy `json:"restartPolicy"`
	Entrypoint    *ContainerEntrypoint    `json:"entrypoint"`
}

func NewContainerSnapshotAttributes(
	portBindings []PortBinding,
	envs []ContainerEnv,
	restartPolicy *ContainerRestartPolicy,
	entrypoint *ContainerEntrypoint,
) ContainerSnapshotAttributes {
	return ContainerSnapshotAttributes{
		PortBindings:  portBindings,
		Envs:          envs,
		RestartPolicy: restartPolicy,
		Entrypoint:    entrypoint,
	}
}
