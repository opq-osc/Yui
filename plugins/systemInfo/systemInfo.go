//go:build tinygo.wasm

package main

import (
	"Yui/plugin/S"
	"Yui/proto"
	"Yui/proto/library/systemInfo/export"
	"context"
	"fmt"

	"github.com/knqyf263/go-plugin/types/known/emptypb"
)

type SystemInfoPlugin struct {
}

func (s SystemInfoPlugin) OnRemoteCallEvent(ctx context.Context, req *proto.RemoteCallReq) (*proto.RemoteCallReply, error) {
	//TODO implement me
	panic("implement me")
}

func (s SystemInfoPlugin) Init(ctx context.Context, empty *emptypb.Empty) (*proto.InitReply, error) {
	S.LogInfo(ctx, "初始化插件")
	return &proto.InitReply{
		Ok:      true,
		Message: "",
	}, nil
}

func (s SystemInfoPlugin) OnGroupMsg(ctx context.Context, msg *proto.CommonMsg) (*emptypb.Empty, error) {
	if msg.SelfId != msg.SenderUin {
		if msg.Message == "系统信息" {
			systemInfo := export.NewSystemInfo()
			cpuInfoBytes, err := systemInfo.CpuInfo(ctx, &emptypb.Empty{})
			if err != nil {
				S.SendGroupTextMsg(ctx, msg.FromUin, msg.SelfId, err.Error())
				return &emptypb.Empty{}, nil
			}
			var cpuInfo = &export.CpuInfo{}
			err = cpuInfo.UnmarshalJSON(cpuInfoBytes.Data)

			memInfoBytes, err := systemInfo.MemInfo(ctx, &emptypb.Empty{})
			if err != nil {
				S.SendGroupTextMsg(ctx, msg.FromUin, msg.SelfId, err.Error())
				return &emptypb.Empty{}, nil
			}
			var memInfo = &export.MemInfo{}
			err = memInfo.UnmarshalJSON(memInfoBytes.Data)
			if err != nil {
				S.SendGroupTextMsg(ctx, msg.FromUin, msg.SelfId, err.Error())
				return &emptypb.Empty{}, nil
			}
			S.SendGroupTextMsg(ctx, msg.FromUin, msg.SelfId, fmt.Sprintf("CPU: %v\nMem: virtual:%v/%v",
				cpuInfo.CPU[0].Family+" "+cpuInfo.CPU[0].Model+" "+cpuInfo.CPU[0].ModelName,
				memInfo.VirtualMemory.Used, memInfo.VirtualMemory.Total,
			))
		}

	}
	return &emptypb.Empty{}, nil
}

func (s SystemInfoPlugin) OnFriendMsg(ctx context.Context, msg *proto.CommonMsg) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s SystemInfoPlugin) OnPrivateMsg(ctx context.Context, msg *proto.CommonMsg) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s SystemInfoPlugin) Unload(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s SystemInfoPlugin) OnCronEvent(ctx context.Context, req *proto.CronEventReq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func main() {
	proto.RegisterEvent(&SystemInfoPlugin{})
}
