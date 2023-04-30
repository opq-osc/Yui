//go:build tinygo.wasm

package S

import (
	"Yui/proto"
	"context"
	"github.com/opq-osc/OPQBot/v2/events"
)

func LogInfo(ctx context.Context, msg string) {
	proto.NewApi().Log(ctx, &proto.LogReq{
		LogType: proto.LogType_Info,
		Msg:     msg,
	})
}
func LogError(ctx context.Context, msg string) {
	proto.NewApi().Log(ctx, &proto.LogReq{
		LogType: proto.LogType_Error,
		Msg:     msg,
	})
}
func LogWarn(ctx context.Context, msg string) {
	proto.NewApi().Log(ctx, &proto.LogReq{
		LogType: proto.LogType_Warn,
		Msg:     msg,
	})
}
func LogDebug(ctx context.Context, msg string) {
	proto.NewApi().Log(ctx, &proto.LogReq{
		LogType: proto.LogType_Debug,
		Msg:     msg,
	})
}
func HttpGet(ctx context.Context, url string, header map[string]string) (*proto.HttpRes, error) {
	return proto.NewApi().Http(ctx, &proto.HttpReq{
		Url:     url,
		Method:  proto.HttpMethod_GET,
		Header:  header,
		Content: nil,
	})
}
func HttpPost(ctx context.Context, url string, header map[string]string, body []byte) (*proto.HttpRes, error) {
	return proto.NewApi().Http(ctx, &proto.HttpReq{
		Url:     url,
		Method:  proto.HttpMethod_POST,
		Header:  header,
		Content: body,
	})
}
func SendGroupTextMsg(ctx context.Context, toUin, botUin int64, text string) (*proto.SendReply, error) {
	return proto.NewApi().SendGroupMsg(ctx, &proto.MsgReq{
		ToUin:  toUin,
		Msg:    &proto.MsgReq_TextMsg{TextMsg: &proto.TextMsg{Text: text}},
		BotUin: botUin,
	})
}
func ParserEvent(rawMessage []byte) (events.IEvent, error) {
	var event = &events.EventStruct{}
	err := event.UnmarshalJSON(rawMessage)
	if err != nil {
		return nil, err
	}
	return event, nil
}
