package util

import (
	"fmt"
	"path"

	"github.com/spf13/viper"
)

func InitViper(dir, file, FileType string) *viper.Viper {
	config := viper.New()
	config.AddConfigPath(dir)
	config.SetConfigName(file)
	config.SetConfigType(FileType)

	if err := config.ReadInConfig(); err != nil {
		panic(fmt.Errorf("解析配置文件%s出错:%s", path.Join(dir, file)+"."+FileType, err))
	}

	return config
}
