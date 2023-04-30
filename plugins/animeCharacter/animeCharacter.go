//go:build tinygo.wasm

package main

import (
	"Yui/plugin/S"
	"Yui/plugins/animeCharacter/config"
	"Yui/proto"
	"context"
	"fmt"
	"github.com/knqyf263/go-plugin/types/known/emptypb"
	"github.com/tidwall/gjson"
	"math/rand"
	"time"
)

var dataArray []gjson.Result

type animeCharacter struct {
}

func (s animeCharacter) OnRemoteCallEvent(ctx context.Context, req *proto.RemoteCallReq) (*proto.RemoteCallReply, error) {
	call := &config.RemoteCallS{}
	err := call.UnmarshalJSON(req.Args)
	if err != nil {
		return nil, err
	}
	S.LogInfo(ctx, call.Test)
	return &proto.RemoteCallReply{Data: []byte(dataArray[rand.Intn(len(dataArray))].String())}, nil
}

func (s animeCharacter) Init(ctx context.Context, empty *emptypb.Empty) (*proto.InitReply, error) {
	S.LogInfo(ctx, "初始化插件")
	rand.Seed(time.Now().Unix())
	return &proto.InitReply{
		Ok:      true,
		Message: "",
	}, nil
}

func (s animeCharacter) OnGroupMsg(ctx context.Context, msg *proto.CommonMsg) (*emptypb.Empty, error) {
	event, err := S.ParserEvent(msg.RawMessage)
	if err != nil {
		S.LogError(ctx, err.Error())
	}
	if event.ParseGroupMsg().AtBot() {
		S.SendGroupTextMsg(ctx, event.ParseGroupMsg().GetGroupUin(), event.GetCurrentQQ(), "at我干什么？")
	}

	if msg.SelfId != msg.SenderUin {
		if msg.Message == "来个名字" {
			if dataArray == nil {
				resp, err := S.HttpGet(ctx, "https://animechan.vercel.app/api/available/character", nil)
				if err != nil {
					return nil, err
				}
				S.LogDebug(ctx, string(resp.Content))
				data := gjson.ParseBytes(resp.Content)
				if data.IsArray() {
					dataArray = data.Array()
				} else {
					S.SendGroupTextMsg(ctx, msg.FromUin, msg.SelfId, "获取资源失败")
				}

			}
			S.SendGroupTextMsg(ctx, msg.FromUin, msg.SelfId, fmt.Sprintf("%s这个名字怎么样呢？", dataArray[rand.Intn(len(dataArray))].String()))

		}

	}
	return &emptypb.Empty{}, nil
}

func (s animeCharacter) OnFriendMsg(ctx context.Context, msg *proto.CommonMsg) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s animeCharacter) OnPrivateMsg(ctx context.Context, msg *proto.CommonMsg) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s animeCharacter) Unload(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s animeCharacter) OnCronEvent(ctx context.Context, req *proto.CronEventReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func main() {
	proto.RegisterEvent(&animeCharacter{})
}
