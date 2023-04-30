package plugin

import (
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/imroc/req/v3"
	"github.com/opq-osc/Yui/plugin/meta"
	"github.com/opq-osc/Yui/proto"
	"github.com/spf13/viper"
)

func (p *Plugin) Http(ctx context.Context, request *proto.HttpReq) (*proto.HttpRes, error) {
	if p.Meta.Permissions&meta.HTTPRequestPermission == 0 {
		return nil, fmt.Errorf("[%s]: %s", p.Meta.PluginName, "插件未获取HTTP联网权限！")
	}
	log.Debug("插件发起HTTP连接", "plugin", p.Meta.PluginName, "method", request.Method, "url", request.Url)
	var r *req.Request
	if proxy := viper.GetString("httpProxy"); proxy != "" {
		r = req.C().SetProxyURL(proxy).R().SetContext(ctx)
	} else {
		r = req.R().SetContext(ctx)
	}

	switch request.GetMethod() {
	case proto.HttpMethod_GET:
		r.Method = "GET"
		if request.Header != nil {
			r.SetHeaders(request.Header)
		}

		resp, err := r.Get(request.Url)
		if err != nil {
			return nil, err
		}
		log.Debug("响应", "statusCode", resp.StatusCode, "data", resp.String())
		return &proto.HttpRes{
			Header:     resp.HeaderToString(),
			StatusCode: int32(resp.StatusCode),
			Content:    resp.Bytes(),
		}, nil
	case proto.HttpMethod_POST:
		r.Method = "POST"
		r.SetBody(request.Content)
		if request.Header != nil {
			r.SetHeaders(request.Header)
		}
		resp, err := r.Post(request.Url)
		if err != nil {
			return nil, err
		}
		log.Debug("响应", "statusCode", resp.StatusCode, "data", resp.String())
		return &proto.HttpRes{
			Header:     resp.HeaderToString(),
			StatusCode: int32(resp.StatusCode),
			Content:    resp.Bytes(),
		}, nil
	}
	return nil, fmt.Errorf("未知类型")
}
