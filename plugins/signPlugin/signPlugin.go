//go:build tinygo.wasm

package main

import (
	"context"
	"fmt"
	"github.com/knqyf263/go-plugin/types/known/emptypb"
	"github.com/opq-osc/Yui/plugin/S"
	anime "github.com/opq-osc/Yui/plugins/animeCharacter/config"
	"github.com/opq-osc/Yui/plugins/signPlugin/config"
	"github.com/opq-osc/Yui/proto"
	"strconv"
	"strings"
	"time"
)

type SignPlugin struct {
}

func (s SignPlugin) OnRemoteCallEvent(ctx context.Context, req *proto.RemoteCallReq) (*proto.RemoteCallReply, error) {
	//TODO implement me
	panic("implement me")
}

func RemoveUser(id int64) error {
	for i, v := range config.C.Users {
		if v == id {
			config.C.Users = append(config.C.Users[:i], config.C.Users[i+1:]...)
		}
	}
	return config.SaveConfig()
}
func AddUser(id int64) error {
	isAdded := false
	for _, v := range config.C.Users {
		if v == id {
			isAdded = true
			break
		}
	}
	if isAdded {
		return fmt.Errorf("已经添加了！")
	}
	config.C.Users = append(config.C.Users, id)
	return config.SaveConfig()
}
func GetList() (result []string) {
	for _, v := range config.C.Users {
		result = append(result, strconv.FormatInt(v, 10))
	}
	return result
}

func (s SignPlugin) Init(ctx context.Context, empty *emptypb.Empty) (*proto.InitReply, error) {
	S.LogInfo(ctx, "初始化插件")
	err := config.ReadConfig()
	if err != nil {
		S.LogError(ctx, err.Error())
	}
	_, err = proto.NewApi().RegisterCronJob(ctx, &proto.CronJob{
		Spec: "0 12 * * *",
		Id:   "sign",
	})
	if err != nil {
		return nil, err
	}
	return &proto.InitReply{
		Ok:      true,
		Message: "",
	}, nil
}

func (s SignPlugin) OnGroupMsg(ctx context.Context, msg *proto.CommonMsg) (*emptypb.Empty, error) {
	if msg.SelfId != msg.SenderUin {
		if msg.Message == "远程call" {
			data, _ := anime.RemoteCallS{Test: "测试"}.MarshalJSON()
			r, err := proto.NewApi().RemoteCall(ctx, &proto.RemoteCallReq{
				DstPluginName: "animeCharacter",
				CallPath:      "randomName",
				Args:          data,
			})
			if err != nil {
				S.SendGroupTextMsg(ctx, msg.FromUin, msg.SelfId, err.Error())
			} else {
				S.SendGroupTextMsg(ctx, msg.FromUin, msg.SelfId, "调用结果："+string(r.Data))
			}
		}
		if msg.Message == "删除我" {
			if err := RemoveUser(msg.SenderUin); err != nil {
				S.SendGroupTextMsg(ctx, msg.FromUin, msg.SelfId, err.Error())
			} else {
				S.SendGroupTextMsg(ctx, msg.FromUin, msg.SelfId, "删除成功")
			}
		}
		if msg.Message == "添加我" {
			if err := AddUser(msg.SenderUin); err != nil {
				S.SendGroupTextMsg(ctx, msg.FromUin, msg.SelfId, err.Error())
			} else {
				S.SendGroupTextMsg(ctx, msg.FromUin, msg.SelfId, "添加成功")
			}
		}
		if msg.Message == "自动签到列表" {
			S.SendGroupTextMsg(ctx, msg.FromUin, msg.SelfId, "自动签到列表:\n"+strings.Join(GetList(), "\n"))

		}

	}
	return &emptypb.Empty{}, nil
}

func (s SignPlugin) OnFriendMsg(ctx context.Context, msg *proto.CommonMsg) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s SignPlugin) OnPrivateMsg(ctx context.Context, msg *proto.CommonMsg) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s SignPlugin) Unload(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	api := proto.NewApi()
	api.RemoveCronJob(ctx, &proto.CronJob{
		Id: "sign",
	})
	return &emptypb.Empty{}, nil
}

func (s SignPlugin) OnCronEvent(ctx context.Context, req *proto.CronEventReq) (*emptypb.Empty, error) {
	if req.Id == "sign" {
		for _, v := range config.C.Users {
			S.SendGroupTextMsg(ctx, 856337734, config.C.BotUin, "签到"+strconv.FormatInt(v, 10))
			time.Sleep(time.Second * 2)
		}

	}
	return &emptypb.Empty{}, nil
}

func main() {
	proto.RegisterEvent(&SignPlugin{})
}
