package initialize

import (
	"fmt"

	"github.com/spf13/viper"

	"IShare/global"
)

func InitViper() (err error) {
	// rootPath, _ := os.Executable()
	// rootPath = filepath.Dir(rootPath)
	v := viper.New()
	v.SetConfigFile("./config.yml")
	err = v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.Set("root_path", "./")

	global.VP = v
	return err
}
