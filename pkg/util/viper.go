package util

import (
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig(in string, cfg interface{}) error {
	v := viper.New()
	v.SetConfigFile(in)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()
	v.AllowEmptyEnv(true)
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	return v.Unmarshal(cfg)
}
