package paths

import (
	"cvgo/console/internal/console"
	"path/filepath"
)

// 当前路径为模块目录
type pathForAnywhere struct {
	rootPath     string // 工程根目录
	commandsPath string // console 根目录
}

func NewPathForAnywhere() *pathForAnywhere {
	// 从 runtime.json 获取根路径
	rootPath := console.GetProjectRootPath()
	return &pathForAnywhere{
		rootPath:     rootPath,
		commandsPath: filepath.Join(rootPath, commandsDir),
	}
}

func (p *pathForAnywhere) DockerEnvTpl() string {
	ret := filepath.Join(p.commandsPath, "create", "sampletpl", "docker-compose-env.yml")
	return ret
}

func (p *pathForAnywhere) DockerDirTpl() string {
	ret := filepath.Join(p.commandsPath, "create", "sampletpl", "docker")
	return ret
}

func (p *pathForAnywhere) AppDockerEnv() string {
	ret := filepath.Join(p.rootPath, "app", "docker-compose-env.yml")
	return ret
}

func (p *pathForAnywhere) AppDockerDir() string {
	ret := filepath.Join(p.rootPath, "app", "docker")
	return ret
}
