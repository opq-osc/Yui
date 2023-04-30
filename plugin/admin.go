package plugin

import (
	"context"
	"fmt"
	"github.com/opq-osc/OPQBot/v2/events"
	"github.com/opq-osc/Yui/plugin/S"
	"github.com/opq-osc/Yui/plugin/meta"
	"github.com/opq-osc/Yui/plugin/repository"
	"github.com/opq-osc/Yui/session"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"strings"
)

// 内置的插件管理功能

func IsAdmin(uin int64) bool {
	for _, v := range cast.ToSlice(viper.Get("admin")) {
		if cast.ToInt64(v) == uin {
			return true
		}
	}
	return false
}

func (m *Manager) OnGroupMsgAdmin(ctx context.Context, event events.IEvent) bool {
	s := session.S.SessionStart(event.ParseGroupMsg().GetSenderUin())
	//log.Debug(string(event.GetRawBytes()))
	path, _ := s.GetString("loadPlugin")
	if path != "" {
		if event.ParseGroupMsg().ParseTextMsg().GetTextContent() == "是" {
			err := M.LoadPlugin(ctx, path)
			if err != nil {
				S.GetApi(event).SendMsg().GroupMsg().ToUin(event.ParseGroupMsg().GetGroupUin()).TextMsg(err.Error()).Do(ctx)
			} else {
				S.GetApi(event).SendMsg().GroupMsg().ToUin(event.ParseGroupMsg().GetGroupUin()).TextMsg("加载成功").Do(ctx)

			}
		}
		s.Delete("loadPlugin")
		return true
	}
	if strings.HasPrefix(event.ParseGroupMsg().ParseTextMsg().GetTextContent(), ".admin") && IsAdmin(event.ParseGroupMsg().GetSenderUin()) {
		cmd := strings.Split(event.ParseGroupMsg().ParseTextMsg().GetTextContent(), " ")
		if len(cmd) <= 1 {
			return false
		}

		switch cmd[1] {
		case "unload":
			err := M.UnloadPlugin(ctx, cmd[2])
			if err != nil {
				S.GetApi(event).SendMsg().GroupMsg().ToUin(event.ParseGroupMsg().GetGroupUin()).TextMsg(err.Error()).Do(ctx)
			} else {
				S.GetApi(event).SendMsg().GroupMsg().ToUin(event.ParseGroupMsg().GetGroupUin()).TextMsg("卸载成功").Do(ctx)
			}
		case "load":
			pluginInfo, err := GetPluginInfo(cmd[2])
			if err != nil {
				S.GetApi(event).SendMsg().GroupMsg().ToUin(event.ParseGroupMsg().GetGroupUin()).TextMsg(err.Error()).Do(ctx)
			} else {
				info := fmt.Sprintf("是否加载插件 [%s] 作者: %s\n说明:%s\n它需要的权限有:\n%s", pluginInfo.PluginName, pluginInfo.Author, pluginInfo.Description, meta.GetPermissions(pluginInfo.Permissions))
				if pluginInfo.Sign {
					info = fmt.Sprintf("✅ %s\n", pluginInfo.SignInfo) + info
				} else {
					info = "⚠️ 插件未知来源\n" + info
				}
				S.GetApi(event).SendMsg().GroupMsg().ToUin(event.ParseGroupMsg().GetGroupUin()).TextMsg(info).Do(ctx)
				s.Set("loadPlugin", cmd[2])
			}
		case "list":
			plugins := M.GetAllPlugins()
			msg := []string{}
			for k, v := range plugins {
				msg = append(msg, fmt.Sprintf("%s 作者:%s 说明:%s 权限:%d 签名:%v", k, v.Meta.Author, v.Meta.Description, v.Meta.Permissions, v.Meta.Sign))
			}
			if len(msg) == 0 {
				msg = []string{"没有插件呢？"}
			}
			S.GetApi(event).SendMsg().GroupMsg().ToUin(event.ParseGroupMsg().GetGroupUin()).TextMsg(strings.Join(msg, "\n")).Do(ctx)
		case "yun":
			switch cmd[2] {
			case "list":
				lists, err := repository.GetPluginList(ctx)
				if err != nil {
					S.GetApi(event).SendMsg().GroupMsg().ToUin(event.ParseGroupMsg().GetGroupUin()).TextMsg(err.Error()).Do(ctx)
					return true
				}
				var msg = []string{"仓库插件："}
				for _, v := range lists {
					msg = append(msg, fmt.Sprintf("%s 作者:%s 说明:%s 权限:%v", v.PluginName, v.Author, v.Description, v.Permissions))
				}
				S.GetApi(event).SendMsg().GroupMsg().ToUin(event.ParseGroupMsg().GetGroupUin()).TextMsg(strings.Join(msg, "\n")).Do(ctx)
			}
		case "permission":
			plugin, err := M.GetPlugin(cmd[2])
			if err != nil {
				S.GetApi(event).SendMsg().GroupMsg().ToUin(event.ParseGroupMsg().GetGroupUin()).TextMsg(err.Error()).Do(ctx)
			} else {
				S.GetApi(event).SendMsg().GroupMsg().ToUin(event.ParseGroupMsg().GetGroupUin()).TextMsg(plugin.Meta.PluginName + "的权限有:\n" + meta.GetPermissions(plugin.Meta.Permissions)).Do(ctx)
			}
		}
		return true
	}
	return false
}
