package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/opq-osc/Yui/plugin/meta"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"os"
	"path/filepath"
	"strings"
)

// 远程仓库 只允许加载签名的插件

const url = "https://api.github.com/repos/opq-osc/Yui-plugins/releases/latest"

//go:generate easyjson repository.go

type Plugin struct {
	meta.PluginMeta
	Permissions []string `json:"Permissions"`
	DownloadUrl string   `json:"DownloadUrl"`
}

func DownloadPlugin(ctx context.Context, plugin Plugin) error {
	var r *req.Request
	if proxy := viper.GetString("httpProxy"); proxy != "" {
		r = req.C().SetProxyURL(proxy).R().SetContext(ctx)
	} else {
		r = req.R().SetContext(ctx)
	}
	resp, err := r.Get(plugin.DownloadUrl)
	if err != nil {
		return err
	}
	sha := sha256.New()
	sha.Write(resp.Bytes())
	if !strings.EqualFold(hex.EncodeToString(sha.Sum(nil)), plugin.Sha256) {
		return fmt.Errorf("sha256比对失败")
	}
	return os.WriteFile(filepath.Join("plugins", plugin.PluginName+".opq"), resp.Bytes(), 0777)
}
func GetPluginList(ctx context.Context) ([]Plugin, error) {
	var r *req.Request
	if proxy := viper.GetString("httpProxy"); proxy != "" {
		r = req.C().SetProxyURL(proxy).R().SetContext(ctx)
	} else {
		r = req.R().SetContext(ctx)
	}
	resp, err := r.Get(url)
	if err != nil {
		return nil, err
	}
	release := gjson.ParseBytes(resp.Bytes())
	// 下载meta
	assetsDownload := map[string]string{}
	for _, v := range release.Get("assets").Array() {
		assetsDownload[v.Get("name").String()] = v.Get("browser_download_url").String()
	}
	resp, err = r.Get(assetsDownload["meta.json"])
	if err != nil {
		return nil, err
	}
	var data = []Plugin{}
	err = json.Unmarshal(resp.Bytes(), &data)
	if err != nil {
		return nil, err
	}
	for i, _ := range data {
		data[i].DownloadUrl = assetsDownload[data[i].PluginName+".opq"]
	}
	return data, nil
}
