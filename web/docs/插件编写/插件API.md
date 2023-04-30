# 插件API

## Log
打印日志信息
```go
proto.NewApi().Log(ctx, &proto.LogReq{
    LogType: proto.LogType_Info,
    Msg:     "来自插件的消息欧~",
})
```

## Http
对外发送 HTTP 请求
```go
proto.NewApi().Http(ctx, &proto.HttpReq{
    Url:     url,
    Method:  proto.HttpMethod_GET,
    Header:  header,
    Content: nil,
})
```

## SendGroupMsg
发送群消息
```go
_, _ = api.SendGroupMsg(ctx, &proto.MsgReq{
    ToUin:  msg.FromUin,
    Msg:    &proto.MsgReq_TextMsg{TextMsg: &proto.TextMsg{Text: "周期任务测试"}},
    BotUin: msg.SelfId,
})
```

## SendFriendMsg
发送好友消息

## SendPrivateMsg
发送私聊消息

## Upload
上传文件到tx服务器

## RegisterCronJob
注册周期任务

## RemoveCronJob
移除周期任务

## RemoteCall
调用其他插件的功能