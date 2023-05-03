package plugin

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/knqyf263/go-plugin/types/known/emptypb"
	"github.com/opq-osc/OPQBot/v2/events"
	"github.com/opq-osc/Yui/cron"
	"github.com/opq-osc/Yui/opq"
	"github.com/opq-osc/Yui/plugin/meta"
	"github.com/opq-osc/Yui/proto"
	"github.com/opq-osc/Yui/proto/library/systemInfo"
	systemInfoExport "github.com/opq-osc/Yui/proto/library/systemInfo/export"
	cron2 "github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const publicKey = "04bb752482134fc1ed0d7b1bdd8502fb3a364918e047c3233cabe40c285ad386ce6c94b2ac56dab3f9b56491f8e54b88e376c5f5bfc0bde8b1bd186c1f920178b4"

var M = Manager{
	plugins: sync.Map{},
}

type Manager struct {
	plugins sync.Map
}

type Plugin struct {
	Meta     meta.PluginMeta
	event    PluginEvent
	cronJobs map[string]cron2.EntryID
	lock     sync.Mutex // 保护 event
}

func (p *Plugin) RemoteCall(ctx context.Context, req *proto.RemoteCallReq) (*proto.RemoteCallReply, error) {
	if p.Meta.Permissions&meta.RemoteCallEventPermission == 0 {
		return nil, fmt.Errorf("缺少调用权限")
	}

	dstPlugin, err := M.GetPlugin(req.DstPluginName)
	if err != nil {
		return nil, err
	}
	// 审查被调用插件权限
	if dstPlugin.Meta.Permissions&meta.RemoteCallEventPermission == 0 {
		return nil, fmt.Errorf("被调用插件缺少权限，无法调用")
	}
	var result *proto.RemoteCallReply
	dstPlugin.LockEvent(func() {
		result, err = dstPlugin.event.OnRemoteCallEvent(ctx, req)
	})
	return result, err
}

func (p *Plugin) LockEvent(f func()) {
	p.lock.Lock()
	defer p.lock.Unlock()
	f()
}

func (p *Plugin) RemoveCronJob(ctx context.Context, job *proto.CronJob) (*emptypb.Empty, error) {
	if p.Meta.Permissions&meta.RegisterCronEventPermission == 0 {
		return nil, fmt.Errorf("[%s]: %s", p.Meta.PluginName, "插件未获取注册周期任务权限！")
	}
	if v, ok := p.cronJobs[job.Id]; ok {
		cron.C.Remove(v)
		delete(p.cronJobs, job.Id)
	}
	return &emptypb.Empty{}, nil
}

func (p *Plugin) RegisterCronJob(ctx context.Context, job *proto.CronJob) (*emptypb.Empty, error) {
	if p.Meta.Permissions&meta.RegisterCronEventPermission == 0 {
		return nil, fmt.Errorf("[%s]: %s", p.Meta.PluginName, "插件未获取注册周期任务权限！")
	}
	id, err := cron.C.AddFunc(job.Spec, func() {
		log.Debug("调用周期任务")
		if p != nil && p.event != nil {
			p.LockEvent(func() {
				_, err := p.event.OnCronEvent(ctx, &proto.CronEventReq{Id: job.Id})
				if err != nil {
					log.Error(err)
				}
			})
		}
	})
	if err != nil {
		return nil, err
	}
	log.Debug("添加周期任务")
	p.cronJobs[job.Id] = id
	return &emptypb.Empty{}, nil
}

type PluginEvent interface {
	Close(ctx context.Context) error
	proto.Event
}

func init() {
	opq.C.On(events.EventNameGroupMsg, func(ctx context.Context, event events.IEvent) {
		if event.ParseGroupMsg().GetMsgType() == events.MsgTypeGroupMsg {
			M.OnGroupMsgEvent(ctx, event)
		}
	})
}

func (m *Manager) OnGroupMsgEvent(ctx context.Context, event events.IEvent) {
	if m.OnGroupMsgAdmin(ctx, event) {
		return
	}
	m.plugins.Range(func(key, value any) bool {
		v := value.(*Plugin)
		// 权限检查
		if v.Meta.Permissions&meta.ReceiveGroupMsgPermission != 0 {
			msg := &proto.CommonMsg{
				Time:        event.GetMsgTime(),
				SelfId:      event.GetCurrentQQ(),
				FromUin:     event.ParseGroupMsg().GetGroupUin(),
				SenderUin:   event.ParseGroupMsg().GetSenderUin(),
				Message:     event.ParseGroupMsg().ParseTextMsg().GetTextContent(),
				RawMessage:  event.GetRawBytes(),
				MessageId:   int32(event.GetMsgUid()),
				MessageType: int32(event.GetMsgType()),
			}
			go v.LockEvent(func() {
				newCtx, cancel := context.WithTimeout(ctx, time.Second*10)
				defer cancel()

				_, err := v.event.OnGroupMsg(newCtx, msg)
				//_, err := v.event.OnGroupMsg(ctx, msg)
				if err != nil {
					log.Error(err)
				}
			})

		}
		return true
	})
}
func (m *Manager) CloseAllPlugin(ctx context.Context) {
	m.plugins.Range(func(key, value any) bool {
		p := value.(*Plugin)
		closedCh := make(chan struct{}, 1)
		go func() {
			defer close(closedCh)
			for _, v := range p.cronJobs {
				cron.C.Remove(v)
			}
			p.LockEvent(func() {
				_, err := p.event.Unload(ctx, &emptypb.Empty{})
				if err != nil {
					log.Error(err)
				}
			})

		}()
		select {
		case <-closedCh:
			log.Infof("插件[%s]已卸载", p.Meta.PluginName)
		case <-time.After(time.Second * 2):
			log.Infof("超时强制关闭插件 %s", p.Meta.PluginName)
		}
		return true
	})
}
func (m *Manager) UnloadPlugin(ctx context.Context, pluginName string) error {
	pI, ok := m.plugins.Load(pluginName)
	if !ok {
		return fmt.Errorf("未找到该插件")
	}
	defer m.plugins.Delete(pluginName)
	p := pI.(*Plugin)
	defer p.LockEvent(func() {
		p.event.Close(ctx)
	})
	closedCh := make(chan struct{}, 1)
	go func() {
		defer close(closedCh)
		for _, v := range p.cronJobs {
			cron.C.Remove(v)
		}
		p.LockEvent(func() {
			_, err := p.event.Unload(ctx, &emptypb.Empty{})
			if err != nil {
				log.Error(err)
			}
		})
	}()
	select {
	case <-closedCh:
		log.Infof("插件[%s]已卸载", p.Meta.PluginName)
	case <-time.After(time.Second * 2):
		return fmt.Errorf("超时强制关闭插件")
	}
	return nil
}

func (m *Manager) GetPlugin(pluginName string) (*Plugin, error) {
	p, ok := m.plugins.Load(pluginName)
	if !ok {
		return nil, fmt.Errorf("未找到该插件")
	}
	return p.(*Plugin), nil
}
func (m *Manager) GetAllPlugins() map[string]*Plugin {
	var plugins = map[string]*Plugin{}
	m.plugins.Range(func(key, value any) bool {
		plugins[key.(string)] = value.(*Plugin)
		return true
	})
	return plugins
}
func GetPluginInfo(pluginName string) (*meta.PluginMeta, error) {
	f, err := os.ReadFile(filepath.Join("plugins", pluginName+".opq"))
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(f)
	sign := [3]byte{}
	_, err = reader.Read(sign[:])
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(sign[:], []byte("OPQ")) {
		return nil, fmt.Errorf("非OPQ插件")
	}
	var pluginApiVersion int32 = 0
	err = binary.Read(reader, binary.LittleEndian, &pluginApiVersion)
	if err != nil {
		return nil, err
	}
	if pluginApiVersion < meta.PluginApiVersion {
		return nil, fmt.Errorf("插件API版本过低无法载入")
	}
	var length int32 = 0
	// 获取签名信息
	r := big.NewInt(0)
	s := big.NewInt(0)
	var signFlag int32 = 0
	err = binary.Read(reader, binary.LittleEndian, &signFlag)
	if err != nil {
		return nil, err
	}
	if int(signFlag) == 1 {
		err = binary.Read(reader, binary.LittleEndian, &length)
		if err != nil {
			return nil, err
		}
		rByte := make([]byte, length)
		_, err = reader.Read(rByte)

		if err != nil {
			return nil, err
		}
		err = r.UnmarshalText(rByte)
		if err != nil {
			return nil, err
		}

		err = binary.Read(reader, binary.LittleEndian, &length)
		if err != nil {
			return nil, err
		}
		sByte := make([]byte, length)
		_, err = reader.Read(sByte)
		if err != nil {
			return nil, err
		}
		err = s.UnmarshalText(sByte)
		if err != nil {
			return nil, err
		}
	}

	// 读取头信息
	err = binary.Read(reader, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}
	var headerSha256 []byte
	header := make([]byte, length)
	_, err = reader.Read(header)
	if err != nil {
		return nil, err
	}
	sha := sha256.New()
	sha.Write(header)
	headerSha256 = sha.Sum(nil)
	dec := gob.NewDecoder(bytes.NewBuffer(header))
	pluginMeta := meta.PluginMeta{}
	if err = dec.Decode(&pluginMeta); err != nil {
		return nil, err
	}
	pluginMeta.Sign = false
	if int(signFlag) == 1 {
		p, err := hex.DecodeString(publicKey)
		if err != nil {
			log.Error(err)
		}
		pubKey, err := crypto.UnmarshalPubkey(p)
		if err != nil {
			log.Fatal(err)
		}
		// 公钥验证签名
		pluginMeta.Sign = ecdsa.Verify(pubKey, headerSha256, r, s)
	}
	return &pluginMeta, nil
}

/*
[0-2] OPQ
[3-6] int32 header Len
*/

func (m *Manager) LoadPlugin(ctx context.Context, pluginPath string) error {
	f, err := os.ReadFile(filepath.Join("plugins", pluginPath+".opq"))
	if err != nil {
		return err
	}
	reader := bytes.NewReader(f)
	sign := [3]byte{}
	_, err = reader.Read(sign[:])
	if err != nil {
		return err
	}
	if !bytes.Equal(sign[:], []byte("OPQ")) {
		return fmt.Errorf("非OPQ插件")
	}
	var pluginApiVersion int32 = 0
	err = binary.Read(reader, binary.LittleEndian, &pluginApiVersion)
	if err != nil {
		return err
	}
	if pluginApiVersion < meta.PluginApiVersion {
		return fmt.Errorf("插件API版本过低无法载入")
	}
	var length int32 = 0
	// 获取签名信息
	r := big.NewInt(0)
	s := big.NewInt(0)
	var signFlag int32 = 0
	err = binary.Read(reader, binary.LittleEndian, &signFlag)
	if err != nil {
		return err
	}
	if int(signFlag) == 1 {
		err = binary.Read(reader, binary.LittleEndian, &length)
		if err != nil {
			return err
		}
		rByte := make([]byte, length)
		_, err = reader.Read(rByte)

		if err != nil {
			return err
		}
		err = r.UnmarshalText(rByte)
		if err != nil {
			return err
		}

		err = binary.Read(reader, binary.LittleEndian, &length)
		if err != nil {
			return err
		}
		sByte := make([]byte, length)
		_, err = reader.Read(sByte)
		if err != nil {
			return err
		}
		err = s.UnmarshalText(sByte)
		if err != nil {
			return err
		}
	}

	// 读取头信息
	err = binary.Read(reader, binary.LittleEndian, &length)
	if err != nil {
		return err
	}
	header := make([]byte, length)
	_, err = reader.Read(header)
	if err != nil {
		return err
	}
	sha := sha256.New()
	sha.Write(header)
	headerSha256 := sha.Sum(nil)
	dec := gob.NewDecoder(bytes.NewBuffer(header))
	pluginMeta := meta.PluginMeta{}
	if err = dec.Decode(&pluginMeta); err != nil {
		return err
	}
	if _, ok := m.plugins.Load(pluginMeta.PluginName); ok {
		return fmt.Errorf("插件名已存在，无法重复加载！")
	}
	// 依赖
	for _, v := range pluginMeta.Dependencies {
		if _, ok := m.plugins.Load(v); !ok {
			return fmt.Errorf("缺少依赖项%s,插件无法载入", v)
		}
	}

	pluginMeta.Sign = false
	if int(signFlag) == 1 {
		p, err := hex.DecodeString(publicKey)
		if err != nil {
			log.Error(err)
		}
		pubKey, err := crypto.UnmarshalPubkey(p)
		if err != nil {
			log.Fatal(err)
		}
		// 公钥验证签名
		pluginMeta.Sign = ecdsa.Verify(pubKey, headerSha256, r, s)
	}
	log.Info("载入插件中", "pluginName", pluginMeta.PluginName, "sha256", pluginMeta.Sha256)
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	sha = sha256.New()
	sha.Write(data)
	sha256Result := hex.EncodeToString(sha.Sum(nil))
	log.Info("计算sha256", "result", sha256Result)
	if !strings.EqualFold(sha256Result, pluginMeta.Sha256) {
		return fmt.Errorf("插件疑似被篡改，禁止载入")
	}
	dir := filepath.Join("plugins", pluginMeta.PluginName)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	p := Plugin{Meta: pluginMeta, cronJobs: map[string]cron2.EntryID{}}
	mc := wazero.NewModuleConfig().
		WithStdout(io.Discard). // 丢弃
		WithStderr(io.Discard). // 丢弃
		WithFSConfig(wazero.NewFSConfig().WithDirMount(dir, "/"))
	env := viper.GetStringMapString("plugin.env")
	for k, v := range env {
		mc = mc.WithEnv(k, v)
	}
	rc := wazero.NewRuntimeConfig()

	pluginVm, err := proto.NewEventPlugin(ctx, proto.WazeroModuleConfig(mc), proto.WazeroRuntime(func(ctx context.Context) (wazero.Runtime, error) {
		r := wazero.NewRuntimeWithConfig(ctx, rc)
		if _, err := wasi_snapshot_preview1.Instantiate(ctx, r); err != nil {
			return nil, err
		}
		if err = systemInfoExport.Instantiate(ctx, r, &systemInfo.SystemInfo{PluginInfo: pluginMeta}); err != nil {
			return nil, err
		}
		return r, nil
	}))

	if err != nil {
		return err
	}

	p.event, err = pluginVm.LoadWithBytes(ctx, data, &p)
	if err != nil {
		return err
	}
	reply, err := p.event.Init(ctx, &emptypb.Empty{})
	if err != nil {
		p.event.Close(ctx)
		return err
	}
	if !reply.Ok {
		p.event.Close(ctx)
		return fmt.Errorf(reply.Message)
	}
	m.plugins.Store(pluginMeta.PluginName, &p)
	return nil
}

func (p *Plugin) Log(ctx context.Context, req *proto.LogReq) (*emptypb.Empty, error) {
	data := fmt.Sprintf("[%s]: %s", p.Meta.PluginName, req.GetMsg())
	switch req.LogType {
	case proto.LogType_Info:
		log.Info(data)
	case proto.LogType_Debug:
		log.Debug(data)
	case proto.LogType_Error:
		log.Error(data)
	case proto.LogType_Warn:
		log.Warn(data)
	}
	return &emptypb.Empty{}, nil
}
