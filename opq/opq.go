package opq

import (
	"github.com/charmbracelet/log"
	"github.com/opq-osc/OPQBot/v2"
	"github.com/spf13/viper"
)

var C *OPQBot.Core

var Url = viper.GetString("opqUrl")

func init() {
	var err error
	C, err = OPQBot.NewCore(Url)
	if err != nil {
		log.Fatal(err)
	}
}
