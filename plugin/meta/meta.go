package meta

//go:generate go run ./generate/generate.go

type PluginMeta struct {
	PluginName   string     `json:"PluginName"`
	Description  string     `json:"Description"`
	Dependencies []string   `json:"Dependencies"`
	Author       string     `json:"Author"`
	Url          string     `json:"Url"`
	Version      int        `json:"Version"`
	Permissions  Permission `json:"Permissions"`
	Sha256       string     `json:"Sha256"`
	Sign         bool       `json:"Sign"`     // 社区官方签名认证
	SignInfo     string     `json:"SignInfo"` // 签名说明，会打印到插件载入时的信息里
}

type Permission uint64

const (
	SendGroupMsgPermission       Permission = 1 << iota // 发送群消息权限
	SendFriendMsgPermission                             // 发送好友消息权限
	SendPrivateChatMsgPermission                        // 发送私聊权限
	HTTPRequestPermission                               // [⚠️注意] Http联网权限
	UploadPermission                                    // 上传文件权限
	RegisterMiddlePermission                            // 注册中间件权限
	GetClientKeyPermission                              // [⚠️注意] 获取 Client Key 权限
	GetPSeyPermission                                   // [⚠️注意] 获取 PS Key 权限
	GetFriendListPermission                             // 获取好友列表权限
	GetGroupListPermission                              // 获取群列表权限
	GetGroupMemberListPermission                        // 获取群成员列表权限
	GroupRevokePermission                               // 撤回群消息权限
	GroupMemberBanPermission                            // 禁言群成员权限
	GroupMemberRemovePermission                         // 踢出群组成员权限
	ReceiveGroupMsgPermission                           // 接收群消息事件权限
	ReceiveFriendMsgPermission                          // 接收好友消息事件权限
	ReceivePrivateMsgPermission                         // 接收私聊消息事件权限
	RegisterCronEventPermission                         // 添加注册周期定时任务权限
	RemoteCallEventPermission                           // [⚠️注意] 远程调用接口权限
	SystemInfoPermission                                // 获取系统信息权限

	ReceiveAllMsgPermission = ReceiveGroupMsgPermission | ReceiveFriendMsgPermission | ReceivePrivateMsgPermission // 宏：接收消息权限
	SendMsgPermission       = SendGroupMsgPermission | SendFriendMsgPermission | SendPrivateChatMsgPermission      // 宏：发送消息权限

	GroupAdminPermission = GetGroupListPermission | GetGroupMemberListPermission | GroupRevokePermission | GroupMemberBanPermission | GroupMemberRemovePermission // 宏：群管理权限
	AllPermission        = Permission(^uint64(0))                                                                                                                 // 宏：所有权限
)
const PluginApiVersion = 1
