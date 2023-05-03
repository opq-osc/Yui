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
	_ "github.com/opq-osc/Yui/config"
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
)

// 允许自动加载插件

type AutoLoadPluginCfg struct {
	PluginName   string `mapstructure:"pluginname"`
	PluginSha256 string `mapstructure:"pluginsha256"`
}

func GetAutoLoadPlugins() []AutoLoadPluginCfg {
	var autoLoadCfg []AutoLoadPluginCfg
	viper.UnmarshalKey("plugin.autoload", &autoLoadCfg)
	return autoLoadCfg
}
func RemoveAutoLoadPlugin(name string) error {
	var plugins []AutoLoadPluginCfg
	viper.UnmarshalKey("plugin.autoload", &plugins)
	for i, v := range plugins {
		if v.PluginName == name {
			plugins = append(plugins[:i], plugins[i+1:]...)
		}
	}
	viper.Set("plugin.autoload", plugins)
	return viper.WriteConfig()
}
func AddAutoLoadPlugin(name string) error {
	var plugins []AutoLoadPluginCfg
	viper.UnmarshalKey("plugin.autoload", &plugins)
	for _, v := range plugins {
		if v.PluginName == name {
			return fmt.Errorf("插件已经添加到自动加载列表")
		}
	}
	data, err := os.ReadFile(filepath.Join("plugins", name+".opq"))
	if err != nil {
		return err
	}
	sha := sha256.New()
	sha.Write(data)
	hash := hex.EncodeToString(sha.Sum(nil))
	plugins = append(plugins, AutoLoadPluginCfg{
		PluginName:   name,
		PluginSha256: hash,
	})
	viper.Set("plugin.autoload", plugins)
	return viper.WriteConfig()
}

// LoadPluginWithSha256 自动加载时使用
func (m *Manager) LoadPluginWithSha256(ctx context.Context, pluginInfo AutoLoadPluginCfg) error {
	f, err := os.ReadFile(filepath.Join("plugins", pluginInfo.PluginName+".opq"))
	if err != nil {
		return err
	}
	sha := sha256.New()
	sha.Write(f)
	if !strings.EqualFold(hex.EncodeToString(sha.Sum(nil)), pluginInfo.PluginSha256) {
		return fmt.Errorf("插件疑似被修改，无法自动载入")
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
	sha = sha256.New()
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

func init() {
	// auto load plugin
	plugins := GetAutoLoadPlugins()
	for _, v := range plugins {
		err := M.LoadPluginWithSha256(context.Background(), v)
		if err != nil {
			log.Error("自动加载时遇到错误", "err", err)
		}
	}
}
