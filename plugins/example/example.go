//go:build tinygo.wasm

package main

import (
	"Yui/plugin/meta"
	"Yui/proto"
	"context"
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/knqyf263/go-plugin/types/known/emptypb"
	"os"
)

const Permission = meta.SendMsgPermission | meta.HTTPRequestPermission | meta.UploadPermission | meta.GroupAdminPermission | meta.ReceiveAllMsgPermission

type ExamplePlugin struct {
}

func (p ExamplePlugin) OnRemoteCallEvent(ctx context.Context, req *proto.RemoteCallReq) (*proto.RemoteCallReply, error) {
	//TODO implement me
	panic("implement me")
}

var cronJob = map[string]func(){}

func (p ExamplePlugin) OnCronEvent(ctx context.Context, req *proto.CronEventReq) (*emptypb.Empty, error) {
	v, ok := cronJob[req.Id]
	if ok {
		v()
	}
	return &emptypb.Empty{}, nil
}

func (p ExamplePlugin) OnFriendMsg(ctx context.Context, msg *proto.CommonMsg) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (p ExamplePlugin) OnPrivateMsg(ctx context.Context, msg *proto.CommonMsg) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (p ExamplePlugin) OnGroupMsg(ctx context.Context, msg *proto.CommonMsg) (*emptypb.Empty, error) {
	api := proto.NewApi()
	if msg != nil {
		if msg.Message == "http" {
			resp, err := api.Http(ctx, &proto.HttpReq{
				Url:     "https://api.github.com/repositories/332116757/languages",
				Method:  proto.HttpMethod_GET,
				Header:  nil,
				Content: nil,
			})
			if err != nil {
				return nil, err
			}

			_, _ = api.SendGroupMsg(ctx, &proto.MsgReq{
				ToUin:  msg.FromUin,
				Msg:    &proto.MsgReq_TextMsg{TextMsg: &proto.TextMsg{Text: string(resp.Content)}},
				BotUin: msg.SelfId,
			})
		}
		if msg.Message == "pic" {
			png, _ := os.ReadFile("/bq.png")
			pic := base64.StdEncoding.EncodeToString(png)
			file, err := api.Upload(ctx, &proto.UploadReq{
				File: &proto.UploadReq_Base64Buf{
					Base64Buf: pic,
				},
				BotUin:   msg.SelfId,
				UploadId: proto.UploadId_Group,
			})
			if err != nil {
				return nil, err
			}
			_, err = api.SendGroupMsg(ctx, &proto.MsgReq{
				ToUin: msg.FromUin,
				Msg: &proto.MsgReq_PicMsg{PicMsg: &proto.Files{
					File: []*proto.File{
						file.File,
					},
				}},
				BotUin: msg.SelfId,
			})

			if err != nil {
				return nil, err
			}
		}
		if msg.Message == "测试周期任务" {
			id := uuid.New()
			cronJob[id.String()] = func() {
				api := proto.NewApi()
				_, _ = api.SendGroupMsg(ctx, &proto.MsgReq{
					ToUin:  msg.FromUin,
					Msg:    &proto.MsgReq_TextMsg{TextMsg: &proto.TextMsg{Text: "周期任务测试"}},
					BotUin: msg.SelfId,
				})
			}
			api.RegisterCronJob(ctx, &proto.CronJob{Spec: "* * * * *", Id: id.String()})

		}
	}

	return &emptypb.Empty{}, nil
}
func (p ExamplePlugin) Unload(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil

}
func (p ExamplePlugin) Init(ctx context.Context, _ *emptypb.Empty) (*proto.InitReply, error) {
	api := proto.NewApi()
	api.Log(ctx, &proto.LogReq{
		LogType: proto.LogType_Info,
		Msg:     "来自插件的消息欧~",
	})
	api.Log(ctx, &proto.LogReq{
		LogType: proto.LogType_Warn,
		Msg:     "来自插件的消息欧~",
	})
	api.Log(ctx, &proto.LogReq{
		LogType: proto.LogType_Error,
		Msg:     "来自插件的消息欧~",
	})
	err := os.WriteFile("/test.txt", []byte("helloworld"), 0777)
	if err != nil {
		api.Log(ctx, &proto.LogReq{
			LogType: proto.LogType_Error,
			Msg:     err.Error(),
		})
	}

	return &proto.InitReply{
		Ok:      true,
		Message: "Success",
	}, nil
}

func main() {
	proto.RegisterEvent(ExamplePlugin{})
}
