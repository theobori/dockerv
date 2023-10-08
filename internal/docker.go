package dockerv

type ActionKind int

const (
	Import ActionKind = 0
	Export ActionKind = 1
	Copy   ActionKind = 2
	Move   ActionKind = 3
)

type DockerVolumeConfig struct {
	Source      string
	Destination string
	Kind        ActionKind
	Recursive   bool
}

type DockerVolume struct {
	config DockerVolumeConfig
}

func NewDockerVolume(config *DockerVolumeConfig) *DockerVolume {
	return &DockerVolume{
		*config,
	}
}

func NewDefaultDockerVolume() *DockerVolume {
	return &DockerVolume{
		DockerVolumeConfig{
			Kind: Export,
		},
	}
}

func (dv *DockerVolume) SetConfig(config *DockerVolumeConfig) {
	dv.config = *config
}

func (dv *DockerVolume) Config() DockerVolumeConfig {
	return dv.config
}
