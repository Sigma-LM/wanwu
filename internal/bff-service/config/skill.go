package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/UnicomAI/wanwu/pkg/util"
)

const (
	akConfigDir = "configs/microservice/bff-service/configs/agent-skills"
)

type SkillsConfig struct {
	SkillId       string `json:"skillId" mapstructure:"skillId"`
	Name          string `json:"name" mapstructure:"name"`
	Avatar        string `json:"avatar" mapstructure:"avatar"`
	Author        string `json:"author" mapstructure:"author"`
	Desc          string `json:"desc" mapstructure:"desc"`
	MdPath        string `json:"mdPath" mapstructure:"mdPath"`
	SkillMarkdown []byte `json:"-" mapstructure:"-"`
}

func (stf *SkillsConfig) AgentSkillZipToBytes(skillsId string) ([]byte, error) {
	return util.DirToBytes(filepath.Join(akConfigDir, skillsId))
}

// --- internal ---
func (stf *SkillsConfig) load() error {
	markdownPath := filepath.Join(akConfigDir, stf.MdPath)
	b, err := os.ReadFile(markdownPath)
	if err != nil {
		return fmt.Errorf("load skill %v makrdown path %v err: %v", stf.SkillId, markdownPath, err)
	}
	stf.SkillMarkdown = b
	return nil
}
