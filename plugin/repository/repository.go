package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/opq-osc/Yui/plugin/meta"
	"github.com/spf13/viper"
)

// 远程仓库 只允许加载签名的插件

const url = "https://raw.githubusercontent.com/mcoo/OPQPlugin/main/repository.json"

//go:generate easyjson repository.go

type Plugin struct {
	meta.PluginMeta
	Permissions []string `json:"Permissions"`
	DownloadUrl string   `json:"DownloadUrl"`
}

//easyjson:json
type ResponseStruct struct {
	ApiVersion int               `json:"ApiVersion"`
	Plugins    map[string]Plugin `json:"Plugins"`
}

func GetPluginList(ctx context.Context) (map[string]Plugin, error) {
	var r *req.Request
	if proxy := viper.GetString("httpProxy"); proxy != "" {
		r = req.C().SetProxyURL(proxy).R().SetContext(ctx)
	} else {
		r = req.R().SetContext(ctx)
	}
	resp, err := r.Get("https://raw.githubusercontent.com/mcoo/OPQPlugin/main/yun.json")
	if err != nil {
		return nil, err
	}
	var data = &ResponseStruct{}
	err = json.Unmarshal(resp.Bytes(), data)
	if err != nil {
		return nil, err
	}
	if meta.PluginApiVersion < data.ApiVersion {
		return nil, fmt.Errorf("api版本过低")
	}
	return data.Plugins, nil
}
