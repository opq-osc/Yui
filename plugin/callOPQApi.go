package plugin

import (
	"Yui/opq"
	"Yui/plugin/meta"
	"Yui/proto"
	"context"
	"fmt"
	"github.com/opq-osc/OPQBot/v2/apiBuilder"
)

func (p *Plugin) SendFriendMsg(ctx context.Context, request *proto.MsgReq) (*proto.SendReply, error) {
	return nil, nil
}
func (p *Plugin) SendPrivateMsg(ctx context.Context, request *proto.MsgReq) (*proto.SendReply, error) {
	return nil, nil
}
func (p *Plugin) SendGroupMsg(ctx context.Context, request *proto.MsgReq) (*proto.SendReply, error) {
	if p.Meta.Permissions&meta.SendGroupMsgPermission == 0 {
		return nil, fmt.Errorf("[%s]: %s", p.Meta.PluginName, "插件未获取发送群聊消息权限！")
	}
	groupMsgApi := apiBuilder.New(opq.C.ApiUrl.String(), request.BotUin).SendMsg().GroupMsg().ToUin(request.ToUin)
	switch v := request.Msg.(type) {
	case *proto.MsgReq_TextMsg:
		groupMsgApi.TextMsg(v.TextMsg.GetText())
	case *proto.MsgReq_PicMsg:
		var files []*apiBuilder.File
		for _, v := range v.PicMsg.GetFile() {
			files = append(files, &apiBuilder.File{
				FileId:    v.GetFileId(),
				FileToken: v.GetFileToken(),
				FileSize:  int(v.GetFileSize()),
				FileMd5:   v.GetFileMd5(),
			})
		}
		groupMsgApi.PicMsg(files...)
	case *proto.MsgReq_XmlMsg:
		groupMsgApi.XmlMsg(v.XmlMsg.GetXml())
	case *proto.MsgReq_JsonMsg:
		groupMsgApi.JsonMsg(v.JsonMsg.GetJson())
	}

	resp, err := groupMsgApi.DoAndResponse(ctx)
	if err != nil {
		return nil, err
	}
	code, msg := resp.Result()
	return &proto.SendReply{
		Ret:    int32(code),
		ErrMsg: &msg,
	}, nil
}

func (p *Plugin) Upload(ctx context.Context, uploadReq *proto.UploadReq) (*proto.UploadReply, error) {
	if p.Meta.Permissions&meta.UploadPermission == 0 {
		return nil, fmt.Errorf("[%s]: %s", p.Meta.PluginName, "插件未获取上传文件权限！")
	}
	api := apiBuilder.New(opq.C.ApiUrl.String(), uploadReq.BotUin).Upload()
	switch v := uploadReq.File.(type) {
	case *proto.UploadReq_Path:
		api.GroupPic().SetFilePath(v.Path)
	case *proto.UploadReq_Url:
		api.SetFileUrlPath(v.Url)
	case *proto.UploadReq_Base64Buf:
		api.SetBase64Buf(v.Base64Buf)
	}
	switch uploadReq.GetUploadId() {
	case proto.UploadId_Friend:
		api.FriendPic()
	case proto.UploadId_Group:
		api.GroupPic()
	}
	file, err := api.DoUpload(ctx)
	if err != nil {
		return nil, err
	}
	return &proto.UploadReply{
		Ret:    0,
		ErrMsg: nil,
		File: &proto.File{
			FileId:    file.FileId,
			FileMd5:   file.FileMd5,
			FileSize:  int32(file.FileSize),
			FileToken: file.FileToken,
		},
	}, nil
}
