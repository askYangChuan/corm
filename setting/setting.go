package setting

import (
	"github.com/go-ini/ini"
	"log"
)


type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &Database{}


var cfg *ini.File
var localCfg *ini.File

// Setup initialize the configuration instance
func Setup(cfgFile string) {
	var err error
	cfg, err = ini.Load(cfgFile)
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}
	localCfg, _ = ini.Load("conf/app_local.ini")

	mapTo("database", DatabaseSetting)
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	} else if localCfg != nil {
		_ = localCfg.Section(section).MapTo(v)
	}
}
