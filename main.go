package main

import (
	_ "Yui/config"
	"Yui/cron"
	"Yui/opq"
	"Yui/plugin"
	"context"
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

var (
	version = "v0.0.1"
	commit  = "xxxxxxxxxxxxxxx"
)

func main() {
	log.Infof("Yui 欢迎使用 Version:%s (commit:%s)", version, commit)
	if v := viper.Get("logLevel"); v != nil {
		log.SetLevel(log.Level(v.(int)))
	}
	err := opq.C.ListenAndWait(context.Background())
	if err != nil {
		log.Error(err)
	}
	plugin.M.CloseAllPlugin(context.Background())
	cron.C.Stop()
}
