package main

import (
	"context"
	"github.com/charmbracelet/log"
	_ "github.com/opq-osc/Yui/config"
	"github.com/opq-osc/Yui/cron"
	"github.com/opq-osc/Yui/opq"
	"github.com/opq-osc/Yui/plugin"
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
